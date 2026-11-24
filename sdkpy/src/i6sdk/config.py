"""Configuration dataclass and default settings for i6 SDK."""

from __future__ import annotations

import os
from dataclasses import dataclass, field
from typing import Any, Callable, Dict, Optional

DEFAULT_DECSUITE_URL = "https://decsuite.sandbox.rs.infinity6.ai"
DEFAULT_TFINFER_URL = "https://tfinfer.sandbox.lb.infinity6.ai"
DEFAULT_USERZ_URL = "https://userz.sandbox.rs.infinity6.ai"


@dataclass
class I6SdkConfig:
    """Configuration options for Infinity6 SDK client."""

    serv_urls: Dict[str, str] = field(
        default_factory=lambda: {
            "decsuite": DEFAULT_DECSUITE_URL,
            "tfinfer": DEFAULT_TFINFER_URL,
            "userz": DEFAULT_USERZ_URL
        }
    )
    userz_url: Optional[str] = None
    app: str = "default"
    timeout: float = 30.0
    fetcher: Optional[Callable[..., Any]] = None


DEFAULT_CONFIG = I6SdkConfig()


def load_config_from_env() -> I6SdkConfig:
    """Creates an I6SdkConfig populated from environment variables."""
    serv_urls = {
        "decsuite": os.environ.get("I6_SDK_DECSUITE_URL", DEFAULT_DECSUITE_URL),
        "tfinfer": os.environ.get("I6_SDK_TFINFER_URL", DEFAULT_TFINFER_URL),
        "userz": os.environ.get("I6_SDK_USERZ_URL", DEFAULT_USERZ_URL),
    }

    return I6SdkConfig(
        serv_urls=serv_urls,
    )
