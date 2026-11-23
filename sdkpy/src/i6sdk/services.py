"""Central service container managing all SDK client instances."""

from __future__ import annotations

from typing import Any, Dict, Optional, Union

from .auth.auth import Auth
from .config import DEFAULT_CONFIG, I6SdkConfig
from .dash.dash import Dash
from .infer.infer import Infer
from .ingest.ingest import Ingest
from .invoker import Invoker
from .sdkapis.apis import Apis
from .sdkapis.sdkfuncs import SdkFuncs
from .tokens.tokens import Tokens


class Services:
    """Service container orchestrating configuration and sub-clients."""

    def __init__(self, config: Optional[I6SdkConfig] = None) -> None:
        self._config: I6SdkConfig = config or DEFAULT_CONFIG
        self._auth: Optional[Auth] = None
        self._tokens: Optional[Tokens] = None
        self._apis: Optional[Apis] = None
        self._sdk_funcs: Optional[SdkFuncs] = None

    def config(self) -> I6SdkConfig:
        """Returns the active SDK configuration."""
        return self._config

    def set_config(self, config: Union[I6SdkConfig, Dict[str, Any]]) -> None:
        """Updates the active SDK configuration."""
        if isinstance(config, I6SdkConfig):
            self._config = config
        elif isinstance(config, dict):
            for k, v in config.items():
                if hasattr(self._config, k):
                    setattr(self._config, k, v)

    setConfig = set_config

    def invoker(self) -> Invoker:
        """Returns an Invoker instance for executing HTTP requests."""
        return Invoker(self)

    def auth(self) -> Auth:
        """Returns the Authentication client."""
        if self._auth is None:
            self._auth = Auth(self)
        return self._auth

    def tokens(self) -> Tokens:
        """Returns the Token management client."""
        if self._tokens is None:
            self._tokens = Tokens(self)
        return self._tokens

    def ingest(self) -> Ingest:
        """Returns the Ingest client."""
        return Ingest(self)

    def infer(self) -> Infer:
        """Returns the Infer client."""
        return Infer(self)

    def dash(self) -> Dash:
        """Returns the Dash analytics client."""
        return Dash(self)

    def apis(self) -> Apis:
        """Returns the generated low-level APIs client."""
        if self._apis is None:
            self._apis = Apis(self)
        return self._apis

    def sdk_funcs(self) -> SdkFuncs:
        """Returns the generated high-level SDK workflow functions client."""
        if self._sdk_funcs is None:
            self._sdk_funcs = SdkFuncs(self)
        return self._sdk_funcs

    sdkFuncs = sdk_funcs
