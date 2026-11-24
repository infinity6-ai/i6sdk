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
DATASET_TOKENS_KEY = "i6_dataset_access_tokens"


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
        self._storage.remove_item(DATASET_TOKENS_KEY)

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
        val = self._storage.get_item(DATASET_TOKENS_KEY)
        if not val:
            return {}
        try:
            return json.loads(val)
        except Exception:
            return {}

    def get_dataset_token(self, dataset: str) -> Optional[DatasetToken]:
        """Returns the cached dataset token for the specified dataset."""
        m = self.get_dataset_tokens_map()
        return m.get(dataset)

    def set_dataset_token(self, dataset: str, token: DatasetToken) -> None:
        """Stores a dataset token."""
        m = self.get_dataset_tokens_map()
        m[dataset] = token
        self._storage.set_item(DATASET_TOKENS_KEY, json.dumps(m))

    def remove_dataset_token(self, dataset: str) -> None:
        """Removes the cached token for the specified dataset."""
        m = self.get_dataset_tokens_map()
        if dataset in m:
            del m[dataset]
            self._storage.set_item(DATASET_TOKENS_KEY, json.dumps(m))

    def is_dataset_token_expired(self, token: Optional[DatasetToken], skew_seconds: float = 30.0) -> bool:
        """Checks whether the given dataset access token is expired or will expire within the given skew window."""
        if not token or not token.get("access_token"):
            return True
        expires_at = token.get("expires_at", 0)
        now_ms = int(time.time() * 1000)
        return now_ms >= (expires_at - int(skew_seconds * 1000))

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

