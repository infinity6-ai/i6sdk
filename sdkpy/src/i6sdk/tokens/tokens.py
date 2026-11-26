"""Token management client for i6 SDK."""

from __future__ import annotations

import json
import threading
import time
import urllib.parse
import urllib.request
from typing import TYPE_CHECKING, Any, Dict, Optional, TypedDict

if TYPE_CHECKING:
    from ..services import Services

from .storage import MemoryTokenStorage, TokenStorage


class UserToken(TypedDict, total=False):
    """User authentication token payload."""

    user_access_token: str
    user_refresh_token: str
    user_id: str
    expires_at: int
    refresh_token_expires_at: Optional[int]


class DatasetToken(TypedDict, total=False):
    """Dataset access token payload."""

    dataset: str
    access_token: str
    expires_at: int
    grants: int


USER_TOKEN_KEY = "i6_user_user_token"
DATASET_TOKEN_PREFIX = "i6_ds:"
DATASET_REFRESH_PREFIX = "i6_rf:"

# Deprecated monolithic aliases retained for compatibility
DATASET_TOKENS_KEY = "i6_dataset_access_tokens"
DATASET_REFRESH_TOKENS_KEY = "i6_dataset_refresh_tokens"


def _dataset_token_key(dataset: str) -> str:
    return f"{DATASET_TOKEN_PREFIX}{dataset}"


def _dataset_refresh_key(dataset: str) -> str:
    return f"{DATASET_REFRESH_PREFIX}{dataset}"


class Tokens:
    """Manages storage, retrieval, and lifecycle of user and dataset access tokens."""

    def __init__(self, services: Services, storage: Optional[TokenStorage] = None) -> None:
        self._services = services
        self._storage: TokenStorage = storage or MemoryTokenStorage()
        self._refresh_lock = threading.Lock()
        self._refreshing_thread: Optional[int] = None
        self._dataset_lock = threading.Lock()

    @property
    def storage(self) -> TokenStorage:
        return self._storage

    @storage.setter
    def storage(self, value: TokenStorage) -> None:
        self._storage = value

    def get_user_token(self) -> Optional[UserToken]:
        """Returns the stored user token if available."""
        val = self._storage.get_item(USER_TOKEN_KEY)
        if not val:
            return None
        try:
            return json.loads(val)
        except Exception:
            return None

    def set_user_token(self, token: UserToken) -> None:
        """Stores the user token."""
        self._storage.set_item(USER_TOKEN_KEY, json.dumps(token))

    def remove_user_token(self) -> None:
        """Removes the stored user token."""
        self._storage.remove_item(USER_TOKEN_KEY)

    def clear_tokens(self) -> None:
        """Clears all stored user and dataset tokens."""
        self.remove_user_token()
        self._storage.delete_prefix(DATASET_TOKEN_PREFIX)
        self._storage.delete_prefix(DATASET_REFRESH_PREFIX)

    def get_dataset_refresh_tokens_map(self) -> Dict[str, str]:
        """Returns all stored dataset refresh tokens."""
        raw_map = self._storage.list_prefix(DATASET_REFRESH_PREFIX)
        return {
            k[len(DATASET_REFRESH_PREFIX):]: v
            for k, v in raw_map.items()
        }

    def get_dataset_refresh_token(self, dataset: str) -> Optional[str]:
        """Returns the stored dataset refresh token (Ed25519 private key/keypair) for dataset."""
        if not dataset:
            return None
        return self._storage.get_item(_dataset_refresh_key(dataset))

    def set_dataset_refresh_token(self, dataset: str, token: Union[str, Dict[str, Any]]) -> None:
        """Stores a dataset refresh token for machine robot authentication."""
        if not dataset:
            raise ValueError("dataset is required")
        raw_key = token
        if isinstance(token, dict):
            raw_key = token.get("priv_key") or token.get("private_key") or json.dumps(token)
        self._storage.set_item(_dataset_refresh_key(dataset), str(raw_key))

    def remove_dataset_refresh_token(self, dataset: str) -> None:
        """Removes the stored dataset refresh token for the specified dataset."""
        if not dataset:
            return
        self._storage.remove_item(_dataset_refresh_key(dataset))

    def is_access_token_expired(
        self, token: Optional[UserToken] = None, skew_seconds: float = 30.0
    ) -> bool:
        """Checks whether the user access token is missing, expired, or expiring within skew_seconds."""
        ut = token if token is not None else self.get_user_token()
        if not ut or not ut.get("user_access_token"):
            return True
        expires_at = ut.get("expires_at", 0)
        if expires_at == 0:
            return False
        now_ms = int(time.time() * 1000)
        return now_ms >= (expires_at - int(skew_seconds * 1000))

    def is_refresh_token_expired(self, token: Optional[UserToken] = None) -> bool:
        """Checks whether the user refresh token is missing or expired."""
        ut = token if token is not None else self.get_user_token()
        if not ut or not ut.get("user_refresh_token"):
            return True
        rf_expires_at = ut.get("refresh_token_expires_at")
        if isinstance(rf_expires_at, int) and rf_expires_at > 0:
            return int(time.time() * 1000) >= rf_expires_at
        return False

    def get_valid_access_token(self, skew_seconds: float = 30.0) -> Optional[str]:
        """Retrieves a valid user access token string, silently refreshing if expired."""
        ut = self.get_user_token()
        if not ut:
            return None

        if not self.is_access_token_expired(ut, skew_seconds):
            return ut.get("user_access_token")

        current_thread = threading.get_ident()
        if self._refreshing_thread == current_thread:
            # Prevent recursive refresh on the same thread
            return ut.get("user_access_token")

        # Expired or expiring soon; refresh session with lock to prevent concurrent rotation races
        with self._refresh_lock:
            # Re-check under lock in case another thread already refreshed
            ut = self.get_user_token()
            if not ut:
                return None
            if not self.is_access_token_expired(ut, skew_seconds):
                return ut.get("user_access_token")

            if ut.get("user_refresh_token") and not self.is_refresh_token_expired(ut):
                self._refreshing_thread = current_thread
                try:
                    refreshed = self._services.auth().refresh()
                    if refreshed and (new_token := refreshed.get("user_access_token")):
                        return new_token
                except Exception:
                    return None
                finally:
                    self._refreshing_thread = None

            return None

    def get_token(self) -> Optional[str]:
        """Returns the active access token (user access token)."""
        ut = self.get_user_token()
        if ut and (access_token := ut.get("user_access_token")):
            expires_at = ut.get("expires_at", 0)
            now_ms = int(time.time() * 1000)
            if expires_at == 0 or expires_at > now_ms:
                return access_token
        return None

    def get_dataset_tokens_map(self) -> Dict[str, DatasetToken]:
        """Returns all cached dataset tokens."""
        raw_map = self._storage.list_prefix(DATASET_TOKEN_PREFIX)
        tokens: Dict[str, DatasetToken] = {}
        for k, v in raw_map.items():
            ds = k[len(DATASET_TOKEN_PREFIX):]
            try:
                tokens[ds] = json.loads(v)
            except Exception:
                pass
        return tokens

    def get_dataset_token(self, dataset: str) -> Optional[DatasetToken]:
        """Returns the cached dataset token for the specified dataset."""
        if not dataset:
            return None
        val = self._storage.get_item(_dataset_token_key(dataset))
        if not val:
            return None
        try:
            return json.loads(val)
        except Exception:
            return None

    def set_dataset_token(self, dataset: str, token: DatasetToken) -> None:
        """Stores a dataset token."""
        if not dataset:
            raise ValueError("dataset is required")
        self._storage.set_item(_dataset_token_key(dataset), json.dumps(token))

    def remove_dataset_token(self, dataset: str) -> None:
        """Removes the cached token for the specified dataset."""
        if not dataset:
            return
        self._storage.remove_item(_dataset_token_key(dataset))

    def is_dataset_token_expired(self, token: Optional[DatasetToken], skew_seconds: float = 30.0) -> bool:
        """Checks whether the given dataset access token is expired or will expire within the given skew window."""
        if not token or not token.get("access_token"):
            return True
        expires_at = token.get("expires_at", 0)
        now_ms = int(time.time() * 1000)
        return now_ms >= (expires_at - int(skew_seconds * 1000))

    def refresh_dataset_access_token(
        self,
        dataset: str,
        subsolution: Optional[str] = None,
        subgrants: Optional[int] = None,
        exp_seconds: Optional[int] = None,
        unbound: bool = True,
    ) -> Optional[DatasetToken]:
        """Mints a fresh dataset access token using the stored dataset refresh token (unbound Bearer by default)."""
        rf_token = self.get_dataset_refresh_token(dataset)
        if not rf_token:
            return None

        from ..internal.reqsign import build_refresh_assertion

        payload: Dict[str, Any] = {
            "dataset": dataset,
            "unbound": unbound,
        }
        if subsolution:
            payload["subsolution"] = subsolution
        if subgrants:
            payload["subgrants"] = subgrants
        if exp_seconds:
            payload["exp_seconds"] = exp_seconds

        assertion_sig = build_refresh_assertion(
            priv_key=rf_token,
            dataset=dataset,
            method="POST",
            path="/api/userz/refresh-dataset",
            body_bytes=b"",
        )

        resp = self._services.apis().user_dataset_token_refresh(
            req={
                "headers": {
                    "x-i6-signature": assertion_sig,
                },
                "payload": payload,
            },
            opts={
                "labels": {"token_types": "none"},
            },
        )
        if not resp.ok:
            return None

        data = resp.parse() or {}
        access_tok = data.get("access_token")
        if not access_tok:
            return None

        exp_at = data.get("expires_at", 0)
        if 0 < exp_at < 10_000_000_000:
            exp_at *= 1000

        dt: DatasetToken = {
            "dataset": data.get("dataset", dataset),
            "access_token": access_tok,
            "expires_at": int(exp_at),
            "grants": int(data.get("grants", 0)),
        }
        self.set_dataset_token(dataset, dt)
        return dt

    def get_valid_dataset_access_token(self, dataset: str, skew_seconds: float = 30.0) -> Optional[str]:
        """Retrieves a valid dataset access token string for the target dataset."""
        if not dataset:
            return None

        cached = self.get_dataset_token(dataset)
        if cached and not self.is_dataset_token_expired(cached, skew_seconds):
            return cached.get("access_token")

        with self._dataset_lock:
            # Re-check inside lock
            cached = self.get_dataset_token(dataset)
            if cached and not self.is_dataset_token_expired(cached, skew_seconds):
                return cached.get("access_token")

            # 1. If dataset refresh token is available for machine robots, refresh via Ed25519 assertion
            if self.get_dataset_refresh_token(dataset):
                try:
                    dt = self.refresh_dataset_access_token(dataset)
                    if dt and (tok := dt.get("access_token")):
                        return tok
                except Exception:
                    pass

            # 2. Otherwise fall back to user token
            user_token = self.get_valid_access_token(skew_seconds)
            if not user_token:
                return None

            userz_url = self._services.auth().userz_base_url().rstrip("/")
            req_url = f"{userz_url}/api/userz/dataset-access-token/{urllib.parse.quote(dataset)}"

            config = self._services.config()
            if config.fetcher:
                status, _, text = config.fetcher(
                    req_url,
                    "POST",
                    {
                        "Content-Type": "application/json",
                        "Authorization": f"Bearer {user_token}",
                        "User-Agent": "i6sdk-python/0.0.1",
                    },
                    b"{}",
                    config.timeout,
                )
                if 200 <= status < 300:
                    data = json.loads(text) if text else {}
                    dt: DatasetToken = {
                        "dataset": data.get("dataset", dataset),
                        "access_token": data.get("access_token", ""),
                        "expires_at": int(data.get("expires_at", 0)),
                        "grants": int(data.get("grants", 0)),
                    }
                    self.set_dataset_token(dataset, dt)
                    return dt.get("access_token")
                return None

            req = urllib.request.Request(
                req_url,
                data=b"{}",
                headers={
                    "Content-Type": "application/json",
                    "Authorization": f"Bearer {user_token}",
                    "User-Agent": "i6sdk-python/0.0.1",
                },
                method="POST",
            )
            try:
                with urllib.request.urlopen(req, timeout=config.timeout) as resp:
                    data = json.loads(resp.read().decode("utf-8"))
                    dt: DatasetToken = {
                        "dataset": data.get("dataset", dataset),
                        "access_token": data.get("access_token", ""),
                        "expires_at": int(data.get("expires_at", 0)),
                        "grants": int(data.get("grants", 0)),
                    }
                    self.set_dataset_token(dataset, dt)
                    return dt.get("access_token")
            except Exception:
                return None

