"""Inference client running local SQLite queries and server-side model evaluation."""

from __future__ import annotations

import json
import pathlib
import time
import uuid
from typing import (
    TYPE_CHECKING,
    Any,
    Callable,
    Dict,
    List,
    Optional,
    TypedDict,
    Union,
)

if TYPE_CHECKING:
    from ..invoker import ServiceOptions
    from ..services import Services

from ..exceptions import I6ApiError, I6ValidationError
from ..internal.downloader import Downloader
from ..internal.fs import Fs
from ..internal.sqlitez import sqlitez


class InferQueryOptions(TypedDict, total=False):
    """Options for model inference queries."""

    dataset: str
    model_name: str
    partitions: Optional[Dict[str, str]]
    params: Optional[Dict[str, Any]]
    sqlite_url: Optional[str]
    service_options: Optional[Union[ServiceOptions, Dict[str, Any]]]


def partitions_key(partitions: Dict[str, str]) -> str:
    """Builds a deterministic string key from a partitions dictionary."""
    parts = [f"{k}_{partitions[k]}" for k in sorted(partitions.keys())]
    return "-".join(parts) or "default"


def default_select_results(req: Any, results: Any) -> Dict[str, Any]:
    """Default recommendation filter selecting the top SKUs not already in the cart."""
    if not isinstance(results, dict):
        return {"results": []}

    raw_results = results.get("results") or []
    cart_skus = set()
    if isinstance(req, dict):
        cart_skus = set((req.get("input") or {}).get("skus") or [])

    selected: List[Dict[str, Any]] = []
    seen_skus = set()

    for item in raw_results:
        sku_id = item.get("sku_id")
        if sku_id and sku_id not in cart_skus and sku_id not in seen_skus:
            seen_skus.add(sku_id)
            selected.append(
                {
                    "sku_id": sku_id,
                    "rating": item.get("rating", 0.0),
                    "enriched": item.get("enriched"),
                }
            )

    return {"results": selected}


class Infer:
    """Inference client coordinating local query preparation and remote model prediction."""

    def __init__(
        self,
        services: Services,
        fs: Optional[Fs] = None,
        downloader: Optional[Downloader] = None,
    ) -> None:
        self._services = services
        self._fs = fs or Fs()
        self._downloader = downloader or Downloader(self._fs)

    def _sqlite_url(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> str:
        """Retrieves the signed download URL for the model SQLite database."""
        if direct_url := opts.get("sqlite_url"):
            return direct_url

        dataset = opts.get("dataset", "")
        model_name = opts.get("model_name", "")
        partitions = opts.get("partitions") or {}

        resp = self._services.apis().dsz_model_get(
            req={
                "params": {
                    "dataset": dataset,
                    "model_name": model_name,
                },
                "query": partitions,
            },
            opts=opts.get("service_options"),
        )
        if not resp.ok:
            raise I6ApiError(
                message=f"Failed to get model metadata: {resp.text}",
                status=resp.status,
            )
        data = resp.parse() or {}
        sqlite_info = (data.get("files") or {}).get("sqlite.bin") or {}
        url = sqlite_info.get("url")
        if not url:
            raise LookupError(f"Model {model_name} has no sqlite.bin file available.")
        return str(url)

    def update(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> pathlib.Path:
        """Preloads or revalidates the model SQLite bundle in the local cache."""
        model_name = opts.get("model_name")
        if not model_name:
            raise I6ValidationError("model_name is required")
        pkey = partitions_key(opts.get("partitions") or {})
        cache_name = f"infer-{model_name}-{pkey}-sqlite"
        url = self._sqlite_url(opts)
        return self._downloader.update(name=cache_name, url=url, min_preload_time=300.0)

    def run(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> Dict[str, Any]:
        """Runs inference: local input query -> server-side prediction -> local output query."""
        model_name = opts.get("model_name")
        if not model_name:
            raise I6ValidationError("model_name is required")
        dataset = opts.get("dataset", "")
        partitions = opts.get("partitions") or {}
        params = opts.get("params")
        service_options = opts.get("service_options")

        # 1. Download or get cached SQLite database
        sqlite_path = self.update(opts)

        # 2. Run input query to format model inputs
        query_input = f"model-{model_name}-input"
        query_output = f"model-{model_name}-output"

        with sqlitez.open(sqlite_path) as conn:
            inputs_str = conn.named_string_query(name=query_input, params=[params])
            if not inputs_str:
                raise RuntimeError(f"Input query {query_input!r} returned no inputs")

            # 3. Call server-side predict API (apitfz)
            predict_resp = self._services.apis().tfz_model_predict(
                req={
                    "params": {
                        "dataset": dataset,
                        "model_name": model_name,
                        "version": partitions.get("version", ""),
                    },
                    "payload": {
                        "model_partitions": partitions,
                        "inputs": inputs_str,
                    },
                },
                opts=service_options,
            )
            if not predict_resp.ok:
                raise I6ApiError(
                    message=f"Model prediction API failed: {predict_resp.text}",
                    status=predict_resp.status,
                )

            data = predict_resp.parse() or {}
            raw_outputs = data.get("outputs", "{}")
            model_outputs = json.loads(raw_outputs) if isinstance(raw_outputs, str) else raw_outputs

            # 4. Run output query to format final ranked recommendations
            limit = opts.get("limit", 20)
            result = conn.named_json_query(name=query_output, params=[model_outputs, limit])

        rec_id = str(uuid.uuid4())
        results = dict(result or {})
        results["rec_id"] = rec_id
        return results

    def select_results(self, opts: Dict[str, Any]) -> Dict[str, Any]:
        """Applies selection logic and returns final recommendation payload."""
        callback = opts.get("select_results") or default_select_results
        req = opts.get("req")
        results = opts.get("results")
        return callback(req, results)
