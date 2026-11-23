"""Authentication client for i6 SDK using Python Standard Library."""

from __future__ import annotations

import base64
import hashlib
import secrets
import time
from typing import Any, Dict, Optional, TypedDict

from ..exceptions import I6AuthError
from ..tokens.tokens import UserToken


def generate_code_verifier(length: int = 64) -> str:
    """Generates an RFC 7636 PKCE code verifier."""
    return secrets.token_urlsafe(length)[:length]


def generate_code_challenge(verifier: str) -> str:
    """Generates an RFC 7636 S256 PKCE code challenge from a verifier."""
    digest = hashlib.sha256(verifier.encode("ascii")).digest()
    return base64.urlsafe_b64encode(digest).decode("ascii").rstrip("=")


class AuthPayload(TypedDict, total=False):
    """User profile information returned by /me endpoint."""

    user_id: str
    email_redacted: Optional[str]
    providers: Optional[list[str]]
    created_at: Optional[int]


class Auth:
    """Authentication client handling login, sign-up, device flow, and token rotation."""

    def __init__(self, services: Any) -> None:
        self._services = services

    def userz_base_url(self) -> str:
        """Resolves the userz authentication base URL."""
        config = self._services.config()
        if serv_urls := config.serv_urls:
            if "userz" in serv_urls:
                return serv_urls["userz"]
        if config.userz_url:
            return config.userz_url
        return config.base_url

    def me(self) -> Dict[str, Any]:
        """Fetches current authenticated user profile."""
        resp = self._services.apis().user_auth_me()
        if not resp.ok:
            raise I6AuthError(f"Failed to fetch user profile ({resp.status}): {resp.text}")
        return resp.parse() or {}

    def sign_in(self, email: str, password: str) -> UserToken:
        """Authenticates with email and password."""
        resp = self._services.apis().user_auth_password_signin(
            req={
                "payload": {
                    "email": email,
                    "password": password,
                    "app": self._services.config().app,
                }
            }
        )
        if not resp.ok:
            raise I6AuthError(f"Sign-in failed ({resp.status}): {resp.text}")
        data = resp.parse()
        token: UserToken = {
            "user_access_token": data.get("user_access_token", ""),
            "user_refresh_token": data.get("user_refresh_token", ""),
            "user_id": data.get("user_id", ""),
            "expires_at": int(time.time() * 1000) + int(data.get("expires_in", 3600) * 1000),
        }
        self._services.tokens().set_user_token(token)
        return token

    def sign_up(self, email: str, password: str) -> UserToken:
        """Registers a new user account with email and password."""
        resp = self._services.apis().user_auth_password_signup(
            req={
                "payload": {
                    "email": email,
                    "password": password,
                    "app": self._services.config().app,
                }
            }
        )
        if not resp.ok:
            raise I6AuthError(f"Sign-up failed ({resp.status}): {resp.text}")
        data = resp.parse()
        token: UserToken = {
            "user_access_token": data.get("user_access_token", ""),
            "user_refresh_token": data.get("user_refresh_token", ""),
            "user_id": data.get("user_id", ""),
            "expires_at": int(time.time() * 1000) + int(data.get("expires_in", 3600) * 1000),
        }
        self._services.tokens().set_user_token(token)
        return token

    def sign_out(self) -> None:
        """Signs out the user and clears stored credentials."""
        try:
            self._services.apis().user_auth_logout()
        except Exception:
            pass
        self._services.tokens().remove_user_token()

    def refresh(self) -> Optional[UserToken]:
        """Refreshes the user access token using the stored refresh token."""
        ut = self._services.tokens().get_user_token()
        if not ut or not ut.get("user_refresh_token"):
            return None

        resp = self._services.apis().user_auth_refresh(
            req={
                "payload": {
                    "user_refresh_token": ut["user_refresh_token"],
                }
            }
        )
        if not resp.ok:
            self._services.tokens().remove_user_token()
            raise I6AuthError(f"Token refresh failed ({resp.status}): {resp.text}")

        data = resp.parse()
        refreshed: UserToken = {
            "user_access_token": data.get("user_access_token", ut.get("user_access_token", "")),
            "user_refresh_token": data.get("user_refresh_token", ut["user_refresh_token"]),
            "user_id": data.get("user_id", ut.get("user_id", "")),
            "expires_at": int(time.time() * 1000) + int(data.get("expires_in", 3600) * 1000),
        }
        self._services.tokens().set_user_token(refreshed)
        return refreshed

    def device_code_start(self, client_id: Optional[str] = None) -> Dict[str, Any]:
        """Initiates an OAuth device authorization flow."""
        resp = self._services.apis().user_auth_device_code(
            req={
                "payload": {
                    "client_id": client_id or self._services.config().app,
                }
            }
        )
        if not resp.ok:
            raise I6AuthError(f"Device code initiation failed ({resp.status}): {resp.text}")
        return resp.parse() or {}

    def device_code_poll(
        self,
        device_code: str,
        interval: float = 5.0,
        timeout: float = 300.0,
    ) -> UserToken:
        """Polls for completion of the device authorization flow."""
        start_time = time.time()
        while time.time() - start_time < timeout:
            resp = self._services.apis().user_auth_device_token(
                req={
                    "payload": {
                        "device_code": device_code,
                    }
                }
            )
            if resp.ok:
                data = resp.parse()
                token: UserToken = {
                    "user_access_token": data.get("user_access_token", ""),
                    "user_refresh_token": data.get("user_refresh_token", ""),
                    "user_id": data.get("user_id", ""),
                    "expires_at": int(time.time() * 1000) + int(data.get("expires_in", 3600) * 1000),
                }
                self._services.tokens().set_user_token(token)
                return token

            # Check if pending or actual error
            err_code = (resp.error or {}).get("code", "")
            if err_code not in ("authorization_pending", "bad_request"):
                raise I6AuthError(f"Device token polling failed: {resp.text}")

            time.sleep(interval)

        raise I6AuthError(f"Device authorization flow timed out after {timeout} seconds.")
