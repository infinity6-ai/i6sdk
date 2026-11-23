"""Dash client for executing dashboard analytical queries on SQLite files."""

from __future__ import annotations

import pathlib
from typing import Any, Dict, Optional

from ..exceptions import I6ApiError, I6ValidationError
from ..infer.infer import partitions_key
from ..internal.downloader import Downloader
from ..internal.fs import Fs
from ..internal.sqlitez import sqlitez


class Dash:
    """Dashboard analytics client querying pre-aggregated SQLite bundles."""

    def __init__(
        self,
        services: Any,
        fs: Optional[Fs] = None,
        downloader: Optional[Downloader] = None,
    ) -> None:
        self._services = services
        self._fs = fs or Fs()
        self._downloader = downloader or Downloader(self._fs)

    def _sqlite_url(self, opts: Dict[str, Any]) -> str:
        """Finds where the dash sqlite is: signed URL from model API or direct URL."""
        if direct_url := opts.get("sqlite_url"):
            return direct_url

        dataset = opts.get("dataset", "")
        dash_name = opts.get("dash_name", "")
        partitions = opts.get("partitions") or {}

        resp = self._services.apis().dsz_model_get(
            req={
                "params": {
                    "dataset": dataset,
                    "model_name": dash_name,
                },
                "query": partitions,
            },
            opts=opts.get("service_options"),
        )
        if not resp.ok:
            raise I6ApiError(
                message=f"Failed to get dash metadata: {resp.text}",
                status=resp.status,
            )
        data = resp.parse() or {}
        sqlite_info = (data.get("files") or {}).get("sqlite.bin") or {}
        url = sqlite_info.get("url")
        if not url:
            raise LookupError(f"Dash {dash_name} has no sqlite.bin file available.")
        return str(url)

    def query(self, opts: Dict[str, Any]) -> Any:
        """Downloads the dash sqlite, executes named query dash-{dash_name}, and returns result."""
        dash_name = opts.get("dash_name")
        if not dash_name:
            raise I6ValidationError("dash_name is required")

        pkey = partitions_key(opts.get("partitions") or {})
        cache_name = f"{dash_name}-{pkey}-sqlite"
        url = self._sqlite_url(opts)

        sqlite_path = self._downloader.update(name=cache_name, url=url, min_preload_time=300.0)

        with sqlitez.open(sqlite_path) as conn:
            params = opts.get("params")
            query_params = [] if params is None else [params]
            return conn.named_json_query(name=dash_name, params=query_params)
