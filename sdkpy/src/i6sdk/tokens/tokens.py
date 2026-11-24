"""Token management client for i6 SDK."""

from __future__ import annotations

import json
import time
from typing import Any, Dict, Optional, TypedDict

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

    def __init__(self, services: Any, storage: Optional[TokenStorage] = None) -> None:
        self._services = services
        self._storage: TokenStorage = storage or MemoryTokenStorage()

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
