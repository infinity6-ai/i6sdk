"""Configuration dataclass and default settings for i6 SDK."""

from __future__ import annotations

import os
from dataclasses import dataclass, field
from typing import Any, Callable, Dict, Optional

DEFAULT_BASE_URL = "https://decsuite.sandbox.rs.infinity6.ai"


@dataclass
class I6SdkConfig:
    """Configuration options for Infinity6 SDK client."""

    base_url: str = DEFAULT_BASE_URL
    serv_urls: Dict[str, str] = field(
        default_factory=lambda: {"tfinfer": DEFAULT_BASE_URL}
    )
    userz_url: Optional[str] = None
    app: str = "default"
    token: Optional[str] = None
    timeout: float = 30.0
    fetcher: Optional[Callable[..., Any]] = None


DEFAULT_CONFIG = I6SdkConfig()


def load_config_from_env() -> I6SdkConfig:
    """Creates an I6SdkConfig populated from environment variables."""
    base_url = os.environ.get("I6_BASE_URL", DEFAULT_BASE_URL)
    token = os.environ.get("I6_API_KEY") or os.environ.get("I6_TOKEN")
    app = os.environ.get("I6_APP", "default")
    userz_url = os.environ.get("I6_USERZ_URL")

    timeout_val = 30.0
    if timeout_env := os.environ.get("I6_TIMEOUT"):
        try:
            timeout_val = float(timeout_env)
        except ValueError:
            pass

    serv_urls = {"tfinfer": base_url}
    if tfinfer_url := os.environ.get("I6_TFINFER_URL"):
        serv_urls["tfinfer"] = tfinfer_url

    return I6SdkConfig(
        base_url=base_url,
        serv_urls=serv_urls,
        userz_url=userz_url,
        app=app,
        token=token,
        timeout=timeout_val,
    )
