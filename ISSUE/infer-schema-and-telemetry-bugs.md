# Issue: `infer.py` Schema Mismatch and Missing Telemetry Events

**Component:** `sdkpy` (`i6sdk.infer.infer`)  
**Severity:** Critical (Blocks model inference runtime and drops telemetry)  
**Status:** Open  

---

## Summary

This issue consolidates two critical defects in the Python `sdkpy` Inference runtime (`i6sdk.infer.infer`):
1. **Model Fetch Schema Bug:** The API request to fetch model SQLite metadata (`dsz_model_get`) uses an outdated request payload format (query params instead of JSON body) and attempts to parse an incorrect response schema (`files["sqlite.bin"]` instead of `urls`). This causes inference to fail entirely against the live server.
2. **Missing Telemetry Events:** The Python inference client does not dispatch `recraw` or `recshow` telemetry events to the data lake, and its `InferQueryOptions` completely omits `session_id` and `user_id`. This causes Python-based inferences to silently drop all recommendation analytics and feedback loops.

---

## Root Cause Analysis

### 1. `dsz_model_get` Schema Mismatch
The backend server (and `sdkweb`) expects:
* **Method:** `POST`
* **Path:** `/api/model/dataset/{dataset}/signed-get/{model_name}`
* **Request Body (`payload`):** `{"model_version": "...", "model_partitions": {...}}`
* **Response Body:** `{"urls": [{"name": "sqlite.bin", "url": "https..."}]}`

`sdkpy` erroneously sends the `partitions` via query parameters and expects a `files` dictionary instead of `urls`.

### 2. Missing `recraw` and `recshow` Events
In `sdkweb`, the `Infer` client enforces the presence of `sessionId` and `userId` and automatically dispatches events via `ingestzEventRelevanceFashionEvent`:
* `run()` dispatches a `recraw` event encapsulating the raw model output and `rec_id`.
* `selectResults()` dispatches a `recshow` event encapsulating the final selected SKU list and the same `rec_id`.

In `sdkpy`, `InferQueryOptions` does not accept `session_id` or `user_id`, and `run()` returns the results directly after generating a `rec_id` without sending any telemetry to the API.

---

## Proposed Remediation

### Fix 1: Update `_sqlite_url` to match the correct API Schema
Update `sdkpy/src/i6sdk/infer/infer.py` to correctly format the POST payload and parse `urls`:
```python
def _sqlite_url(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> str:
    if direct_url := opts.get("sqlite_url"):
        return direct_url

    dataset = opts.get("dataset", "")
    model_name = opts.get("model_name", "")
    partitions = opts.get("partitions") or {}

    resp = self._services.apis().dsz_model_get(
        req={
            "params": {"dataset": dataset, "model_name": model_name},
            "payload": {
                "model_version": partitions.get("version", ""),
                "model_partitions": partitions,
            },
        },
        opts=opts.get("service_options"),
    )
    if not resp.ok:
        raise I6ApiError(message=f"Failed to get model metadata: {resp.text}", status=resp.status)
    
    data = resp.parse() or {}
    urls = data.get("urls") or []
    by_name = {f.get("name"): f.get("url") for f in urls if isinstance(f, dict)}
    url = by_name.get("sqlite.bin")
    if not url:
        raise LookupError(f"Model {model_name} has no sqlite.bin file available.")
    return str(url)
```

### Fix 2: Add Telemetry Dispatching to `run()` and `select_results()`
Update `InferQueryOptions` to include `session_id` and `user_id`. Add a helper method to dispatch `ingestz_event_relevance_fashion_event` and invoke it from `run()` and `select_results()`:

```python
class InferQueryOptions(TypedDict, total=False):
    dataset: str
    model_name: str
    partitions: Optional[Dict[str, str]]
    params: Optional[Dict[str, Any]]
    session_id: str
    user_id: str
    limit: Optional[int]
    sqlite_url: Optional[str]
    service_options: Optional[Union[ServiceOptions, Dict[str, Any]]]

class Infer:
    ...
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
        service_options: Optional[Dict[str, Any]] = None,
    ) -> None:
        created_at = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S") + "Z"
        event_payload: Dict[str, Any] = {
            "id": str(uuid.uuid4()),
            "store_id": store_id,
            "session_id": session_id,
            "user_id": user_id,
            "created_at": created_at,
            kind: {
                "rec_id": rec_id,
                "model_id": model_name,
                "model_version": model_version,
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
        except Exception as e:
            # Telemetry dispatch should fail silently or log, not crash the inference loop
            pass

    def run(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> Dict[str, Any]:
        # Enforce required telemetry fields
        if not opts.get("session_id") or not opts.get("user_id"):
            raise I6ValidationError("infer: session_id and user_id are required")
        
        # ... (execute inference query) ...
        
        rec_id = str(uuid.uuid4())
        results = dict(result or {})
        results["rec_id"] = rec_id
        
        partitions = opts.get("partitions") or {}
        self._send_rec_event(
            dataset=opts.get("dataset", ""),
            model_name=model_name,
            model_version=partitions.get("version", ""),
            store_id=str(partitions.get("store_id", "1")),
            session_id=opts["session_id"],
            user_id=opts["user_id"],
            rec_id=rec_id,
            kind="recraw",
            payload=results,
            service_options=opts.get("service_options"),
        )
        return results
```
This aligns Python's model retrieval and telemetry observability completely with the JavaScript SDK and backend schemas.
