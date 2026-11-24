"""Configuration dataclass and default settings for i6 SDK."""

from __future__ import annotations

import os
import pathlib
from dataclasses import dataclass, field
from typing import TYPE_CHECKING, Callable, Dict, Optional, Protocol, Union

if TYPE_CHECKING:
    from .tokens.storage import TokenStorage

DEFAULT_DECSUITE_URL = "https://decsuite.sandbox.rs.infinity6.ai"
DEFAULT_TFINFER_URL = "https://tfinfer.sandbox.lb.infinity6.ai"
DEFAULT_USERZ_URL = "https://userz.sandbox.rs.infinity6.ai"


class Fetcher(Protocol):
    """Protocol for custom HTTP fetchers intercepting SDK requests."""

    def __call__(
        self,
        url: str,
        method: str,
        headers: Dict[str, str],
        body: Optional[bytes],
        timeout: float,
    ) -> tuple[int, Dict[str, str], str]: ...


@dataclass
class I6SdkConfig:
    """Configuration options for Infinity6 SDK client."""

    serv_urls: Dict[str, str] = field(
        default_factory=lambda: {
            "decsuite": DEFAULT_DECSUITE_URL,
            "tfinfer": DEFAULT_TFINFER_URL,
            "userz": DEFAULT_USERZ_URL,
        }
    )
    userz_url: Optional[str] = None
    app: str = "default"
    timeout: float = 30.0
    storage: Optional[Union[str, pathlib.Path, TokenStorage]] = None
    cache_dir: Optional[Union[str, pathlib.Path]] = None
    fetcher: Optional[Fetcher] = None


DEFAULT_CONFIG = I6SdkConfig()


def load_config_from_env() -> I6SdkConfig:
    """Creates an I6SdkConfig populated from environment variables."""
    serv_urls = {
        "decsuite": os.environ.get("I6_SDK_DECSUITE_URL", DEFAULT_DECSUITE_URL),
        "tfinfer": os.environ.get("I6_SDK_TFINFER_URL", DEFAULT_TFINFER_URL),
        "userz": os.environ.get("I6_SDK_USERZ_URL", DEFAULT_USERZ_URL),
    }

    cache_dir = os.environ.get("I6_SDK_CACHE_DIR")

    return I6SdkConfig(
        serv_urls=serv_urls,
        cache_dir=cache_dir,
    )
