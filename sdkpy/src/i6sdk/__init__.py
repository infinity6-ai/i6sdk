"""Official Infinity6 Python SDK."""

from __future__ import annotations

from typing import Optional

from .config import DEFAULT_CONFIG, I6SdkConfig, load_config_from_env
from .exceptions import (
    I6ApiError,
    I6AuthError,
    I6ConnectionError,
    I6Error,
    I6ValidationError,
)
from .invoker import Resp, ServiceOptions
from .services import Services
from .tokens.storage import FileTokenStorage, MemoryTokenStorage, TokenStorage

# Client alias
I6Client = Services


def create_client(
    api_key: Optional[str] = None,
    base_url: Optional[str] = None,
    app: Optional[str] = None,
    timeout: Optional[float] = None,
) -> I6Client:
    """Convenience factory function creating an I6Client."""
    cfg = load_config_from_env()
    if api_key is not None:
        cfg.token = api_key
    if base_url is not None:
        cfg.base_url = base_url
    if app is not None:
        cfg.app = app
    if timeout is not None:
        cfg.timeout = timeout
    return I6Client(cfg)


__all__ = [
    "I6Client",
    "Services",
    "I6SdkConfig",
    "DEFAULT_CONFIG",
    "load_config_from_env",
    "create_client",
    "I6Error",
    "I6ApiError",
    "I6AuthError",
    "I6ConnectionError",
    "I6ValidationError",
    "Resp",
    "ServiceOptions",
    "TokenStorage",
    "MemoryTokenStorage",
    "FileTokenStorage",
]
