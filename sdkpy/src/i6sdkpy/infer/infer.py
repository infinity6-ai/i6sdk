"""Inference client running local SQLite queries and server-side model evaluation."""

from __future__ import annotations

import datetime
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
    from ..services import I6Client

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
    session_id: str
    user_id: str
    limit: Optional[int]
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
        input_skus = (req.get("input") or {}).get("skus") or []
        for s in input_skus:
            if isinstance(s, dict):
                if s.get("sku_id"):
                    cart_skus.add(str(s["sku_id"]))
            elif isinstance(s, str):
                cart_skus.add(s)

    selected: List[Dict[str, Any]] = []
    seen_skus = set()

    for item in raw_results:
        item_skus = item.get("skus") or []
        if not item_skus and item.get("sku_id"):
            item_skus = [{"sku_id": item["sku_id"]}]

        for sku_obj in item_skus:
            sku_id = sku_obj.get("sku_id") if isinstance(sku_obj, dict) else str(sku_obj)
            if sku_id and sku_id not in cart_skus and sku_id not in seen_skus:
                seen_skus.add(sku_id)
                res_entry: Dict[str, Any] = {
                    "sku_id": sku_id,
                    "rating": item.get("rating", 0.0),
                }
                if item.get("enriched") is not None:
                    res_entry["enriched"] = item["enriched"]
                selected.append(res_entry)
                break

    return {"results": selected}


class Infer:
    """Inference client coordinating local query preparation and remote model prediction."""

    def __init__(
        self,
        services: I6Client,
        fs: Optional[Fs] = None,
        downloader: Optional[Downloader] = None,
    ) -> None:
        self._services = services
        self._fs = fs or Fs()
        self._downloader = downloader or Downloader(self._fs)

    def _send_rec_event(
        self,
        dataset: str,
        model_name: str,
        model_version: str,
        store_id: str,
        session_id: str,
        user_id: str,
        rec_id: str,
        kind: str,
        payload: Any,
        service_options: Optional[Union[ServiceOptions, Dict[str, Any]]] = None,
    ) -> None:
        """Sends a recommendation telemetry event (recraw or recshow) to the data lake."""
        event_id = str(uuid.uuid4())
        created_at = (
            datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S") + "Z"
        )
        event_payload: Dict[str, Any] = {
            "id": event_id,
            "store_id": str(store_id),
            "session_id": str(session_id),
            "user_id": str(user_id),
            "created_at": created_at,
            kind: {
                "rec_id": str(rec_id),
                "model_id": str(model_name),
                "model_version": str(model_version),
                "payload": json.dumps(payload),
            },
        }

        try:
            self._services.apis().ingestz_event_relevance_fashion_event(
                req={
                    "params": {"dataset": dataset},
                    "payload": event_payload,
                },
                opts=service_options,
            )
        except Exception:
            # Telemetry dispatch should be non-fatal to the inference call
            pass

    def _sqlite_url(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> str:
        """Retrieves the signed download URL for the model SQLite database."""
        if direct_url := opts.get("sqlite_url"):
            return str(direct_url)

        dataset = opts.get("dataset", "")
        model_name = opts.get("model_name", "")
        partitions = opts.get("partitions") or {}
        service_options = opts.get("service_options")

        resp = self._services.apis().dsz_model_get(
            req={
                "params": {
                    "dataset": dataset,
                    "model_name": model_name,
                },
                "payload": {
                    "model_version": partitions.get("version", ""),
                    "model_partitions": partitions,
                },
            },
            opts=service_options,
        )
        if not resp.ok:
            raise I6ApiError(
                message=f"Failed to get model metadata: {resp.text}",
                status=resp.status,
            )
        data = resp.parse() or {}
        urls = data.get("urls") or []
        by_name = {
            f.get("name"): f.get("url")
            for f in urls
            if isinstance(f, dict) and f.get("name") and f.get("url")
        }
        url = by_name.get("sqlite.bin")
        if not url:
            raise LookupError(f"Model {model_name} has no sqlite.bin file available.")
        return str(url)

    def update(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> pathlib.Path:
        """Preloads or revalidates the model SQLite bundle in the local cache."""
        model_name = opts.get("model_name")
        if not model_name:
            raise I6ValidationError("infer: model_name is required")
        pkey = partitions_key(opts.get("partitions") or {})
        cache_name = f"infer-{model_name}-{pkey}-sqlite"
        url = self._sqlite_url(opts)
        return self._downloader.update(name=cache_name, url=url, min_preload_time=300.0)

    def run(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> Dict[str, Any]:
        """Runs inference: local input query -> server-side prediction -> local output query."""
        if not opts:
            raise I6ValidationError("infer: opts is required")

        model_name = opts.get("model_name")
        if not model_name:
            raise I6ValidationError("infer: model_name is required")

        dataset = opts.get("dataset")
        if not dataset:
            raise I6ValidationError("infer: dataset is required")

        session_id = opts.get("session_id")
        user_id = opts.get("user_id")
        if not session_id or not user_id:
            raise I6ValidationError("infer: session_id and user_id are required")

        partitions = opts.get("partitions")
        if not isinstance(partitions, dict):
            raise I6ValidationError("infer: partitions must be a dictionary")

        if opts.get("params") is None:
            raise I6ValidationError("infer: params is required")

        limit = opts.get("limit", 20)
        if limit is not None and (not isinstance(limit, int) or limit <= 0):
            raise I6ValidationError("infer: limit must be a positive integer")

        service_options = opts.get("service_options")
        params = opts.get("params")

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
            result = conn.named_json_query(name=query_output, params=[model_outputs, limit])

        rec_id = str(uuid.uuid4())
        results = dict(result or {})
        results["rec_id"] = rec_id

        # 5. Dispatch recraw telemetry event
        self._send_rec_event(
            dataset=dataset,
            model_name=model_name,
            model_version=str(partitions.get("version", "")),
            store_id=str(partitions.get("store_id", "1")),
            session_id=str(session_id),
            user_id=str(user_id),
            rec_id=rec_id,
            kind="recraw",
            payload=results,
            service_options=service_options,
        )

        return results

    def select_results(self, opts: Dict[str, Any]) -> Dict[str, Any]:
        """Applies selection logic, dispatches recshow telemetry, and returns final recommendation payload."""
        if not opts:
            raise I6ValidationError("infer: opts is required")

        callback = opts.get("select_results") or default_select_results
        if not callable(callback):
            raise TypeError("infer: select_results must be callable")

        rec_id = opts.get("rec_id")
        if not rec_id:
            raise I6ValidationError("infer: rec_id is required")

        session_id = opts.get("session_id")
        user_id = opts.get("user_id")
        if not session_id or not user_id:
            raise I6ValidationError("infer: session_id and user_id are required")

        req = opts.get("req")
        results = opts.get("results")

        selected = callback(req, results)
        if not isinstance(selected, dict) or not isinstance(selected.get("results"), list):
            raise TypeError("infer: select_results must return {'results': [...]}")

        # The event carries only the sku ids
        skus_list = []
        for i, item in enumerate(selected["results"]):
            sku_id = item.get("sku_id") if isinstance(item, dict) else None
            if not sku_id:
                raise TypeError(f"infer: select_results returned an item with no sku_id at results[{i}]")
            skus_list.append({"sku_id": str(sku_id)})

        shown = {"skus": skus_list}

        req_input = (req.get("input") or {}) if isinstance(req, dict) else {}
        dataset = opts.get("dataset") or req_input.get("dataset") or ""
        model_name = opts.get("model_name") or ""
        partitions = opts.get("partitions") or {}
        store_id = str(partitions.get("store_id") or req_input.get("store_id") or "1")
        model_version = str(partitions.get("version") or opts.get("model_version") or "")
        service_options = opts.get("service_options")

        self._send_rec_event(
            dataset=dataset,
            model_name=model_name,
            model_version=model_version,
            store_id=store_id,
            session_id=str(session_id),
            user_id=str(user_id),
            rec_id=str(rec_id),
            kind="recshow",
            payload=shown,
            service_options=service_options,
        )

        return selected

