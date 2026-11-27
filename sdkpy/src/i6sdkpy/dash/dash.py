"""Dash client for executing dashboard analytical queries on SQLite files."""

from __future__ import annotations

import pathlib
from typing import TYPE_CHECKING, Any, Dict, Optional, TypedDict, Union

if TYPE_CHECKING:
    from ..invoker import ServiceOptions
    from ..services import I6Client

from ..exceptions import I6ApiError, I6ValidationError
from ..infer.infer import partitions_key
from ..internal.downloader import Downloader
from ..internal.fs import Fs
from ..internal.sqlitez import sqlitez


class DashQueryOptions(TypedDict, total=False):
    """Execution options for Dash queries."""

    dataset: str
    dash_name: str
    partitions: Optional[Dict[str, str]]
    params: Optional[Dict[str, Any]]
    sqlite_url: Optional[str]
    service_options: Optional[Union[ServiceOptions, Dict[str, Any]]]


class Dash:
    """Dashboard analytics client querying pre-aggregated SQLite bundles."""

    def __init__(
        self,
        services: I6Client,
        fs: Optional[Fs] = None,
        downloader: Optional[Downloader] = None,
    ) -> None:
        self._services = services
        self._fs = fs or Fs()
        self._downloader = downloader or Downloader(self._fs)

    def _sqlite_url(self, opts: Union[DashQueryOptions, Dict[str, Any]]) -> str:
        """Finds where the dash sqlite is: signed URL from dsz_signed_url or direct URL."""
        if direct_url := opts.get("sqlite_url"):
            return direct_url

        dataset = opts.get("dataset", "")
        dash_name = opts.get("dash_name", "")
        partitions = opts.get("partitions") or {}
        store_id = str(partitions.get("store_id", "1"))

        resp = self._services.apis().dsz_signed_url(
            req={
                "params": {"dataset": dataset},
                "payload": {
                    "paths": [
                        {"method": "GET", "path": f"i6tb://results.relevance-fashion-store-dash/store_id={store_id}/sqlite.bin"}
                    ]
                },
            },
            opts=opts.get("service_options"),
        )
        if not resp.ok:
            raise I6ApiError(
                message=f"dash: the api gave no sqlite.bin file: {resp.text}",
                status=resp.status,
                code=(resp.error or {}).get("code", "bad_request"),
            )
        data = resp.parse() or {}
        urls = data.get("urls") or []
        first_url = urls[0].get("url") if urls else None
        if not first_url:
            raise LookupError(f"dash: the api gave no sqlite.bin file for {dash_name}: empty answer")
        return str(first_url)

    def query(self, opts: Union[DashQueryOptions, Dict[str, Any]]) -> Any:
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
