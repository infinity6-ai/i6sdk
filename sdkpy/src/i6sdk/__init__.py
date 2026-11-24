"""Official Infinity6 Python SDK."""

from __future__ import annotations

import pathlib
from typing import Any, Optional, Union

from .config import DEFAULT_CONFIG, Fetcher, I6SdkConfig, load_config_from_env
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
    serv_urls: Optional[dict[str, str]] = None,
    app: Optional[str] = None,
    timeout: Optional[float] = None,
    storage: Optional[Union[str, pathlib.Path, TokenStorage]] = None,
) -> I6Client:
    """Convenience factory function creating an I6Client."""
    cfg = load_config_from_env()
    if serv_urls is not None:
        cfg.serv_urls.update(serv_urls)
    if app is not None:
        cfg.app = app
    if timeout is not None:
        cfg.timeout = timeout
    if storage is not None:
        cfg.storage = storage
    return I6Client(cfg)


__all__ = [
    "I6Client",
    "Services",
    "I6SdkConfig",
    "DEFAULT_CONFIG",
    "load_config_from_env",
    "create_client",
    "Fetcher",
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
