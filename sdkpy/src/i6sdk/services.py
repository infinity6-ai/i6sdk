"""Central service container managing all SDK client instances."""

from __future__ import annotations

from typing import Any, Dict, Optional, Union

import pathlib

from .auth.auth import Auth
from .config import DEFAULT_CONFIG, I6SdkConfig
from .dash.dash import Dash
from .exceptions import I6ValidationError
from .grants import Grants
from .infer.infer import Infer
from .ingest.ingest import Ingest
from .internal import bitmapz
from .internal.downloader import Downloader
from .internal.fs import Fs
from .invoker import Invoker
from .sdkapis.apis import Apis
from .sdkapis.sdkfuncs import SdkFuncs
from .tokens.storage import SqliteTokenStorage
from .tokens.tokens import Tokens


class Services:
    """Service container orchestrating configuration and sub-clients."""

    bitmapz = bitmapz
    Grants = Grants

    def __init__(self, config: Optional[I6SdkConfig] = None) -> None:
        self._config: I6SdkConfig = config or DEFAULT_CONFIG
        self._auth: Optional[Auth] = None
        self._tokens: Optional[Tokens] = None
        self._apis: Optional[Apis] = None
        self._funcs: Optional[SdkFuncs] = None
        self._fs: Optional[Fs] = None
        self._downloader: Optional[Downloader] = None
        self._infer: Optional[Infer] = None
        self._dash: Optional[Dash] = None
        self._ingest: Optional[Ingest] = None

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
        # Invalidate cached services that depend on config
        self._fs = None
        self._downloader = None
        self._infer = None
        self._dash = None

    setConfig = set_config

    def fs(self) -> Fs:
        """Returns the shared local filesystem cache manager."""
        if self._fs is None:
            self._fs = Fs(base_dir=self._config.cache_dir)
        return self._fs

    def downloader(self) -> Downloader:
        """Returns the shared file downloader."""
        if self._downloader is None:
            self._downloader = Downloader(self.fs())
        return self._downloader

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
            storage = self._config.storage
            if isinstance(storage, (str, pathlib.Path)):
                storage = SqliteTokenStorage(storage)
            self._tokens = Tokens(self, storage=storage)
        return self._tokens

    def ingest(self) -> Ingest:
        """Returns the Ingest client."""
        if self._ingest is None:
            self._ingest = Ingest(self)
        return self._ingest

    def infer(self) -> Infer:
        """Returns the Infer client."""
        if self._infer is None:
            self._infer = Infer(self, fs=self.fs(), downloader=self.downloader())
        return self._infer

    def dash(self) -> Dash:
        """Returns the Dash analytics client."""
        if self._dash is None:
            self._dash = Dash(self, fs=self.fs(), downloader=self.downloader())
        return self._dash

    def apis(self) -> Apis:
        """Returns the generated low-level APIs client."""
        if self._apis is None:
            self._apis = Apis(self)
        return self._apis

    def funcs(self) -> SdkFuncs:
        """Returns the generated high-level SDK workflow functions client."""
        if self._funcs is None:
            self._funcs = SdkFuncs(self)
        return self._funcs

