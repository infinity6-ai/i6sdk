# Issue: `infer.py` Schema Mismatch in `_sqlite_url` (`dsz_model_get`)

**Component:** `sdkpy` (`i6sdk.infer.infer`)  
**Severity:** Critical (Blocks model inference runtime in Python)  
**Status:** Open  

---

## Summary

In `sdkpy/src/i6sdk/infer/infer.py`, the `_sqlite_url()` method makes a request to `apis().dsz_model_get()` with an incorrect request schema and parses an incorrect response structure. As a result, downloading the model SQLite bundle fails at runtime with `LookupError: Model <name> has no sqlite.bin file available.` or an API error from the server.

---

## Root Cause Analysis

### 1. Request Payload Discrepancy
The server API definition in `apidatasetz/dszapimodelget/modelzapisignedget.go` specifies:
* **Method:** `POST`
* **Path:** `/api/model/dataset/{dataset}/signed-get/{model_name}`
* **Path Params:** `dataset`, `model_name`
* **Request Body (`ReqBody`):**
  * `model_version`: string
  * `model_partitions`: map[string]string

However, `sdkpy` constructs the request as:
```python
# sdkpy/src/i6sdk/infer/infer.py
resp = self._services.apis().dsz_model_get(
    req={
        "params": {
            "dataset": dataset,
            "model_name": model_name,
        },
        "query": partitions,  # <--- BUG: sent as query string instead of body payload
    },
    opts=opts.get("service_options"),
)
```

### 2. Response Payload Discrepancy
The server responds with:
```go
type Resp struct {
    Urls []*ModelFile `json:"urls"`
}
// where ModelFile has "name" and "url"
```

However, `sdkpy` attempts to parse a non-existent `files` dictionary:
```python
# sdkpy/src/i6sdk/infer/infer.py
data = resp.parse() or {}
sqlite_info = (data.get("files") or {}).get("sqlite.bin") or {}
url = sqlite_info.get("url")
if not url:
    raise LookupError(f"Model {model_name} has no sqlite.bin file available.")
```

By comparison, `sdkweb/src/main/infer/infer.js` correctly formats and parses this API:
```javascript
const resp = await this._services.apis().dszModelGet({
  params: { dataset, model_name: modelName },
  payload: { model_version: opts.partitions.version, model_partitions: partitions },
}, opts.serviceOptions)

const data = resp.parse()
const byName = Object.fromEntries((data?.urls ?? []).map((file) => [file.name, file.url]))
return { sqlite: byName["sqlite.bin"] }
```

---

## Steps to Reproduce

1. Initialize `I6Client` in Python.
2. Call `client.infer().run(...)` or `client.infer().update(...)` with a valid dataset and model name.
3. Observe `LookupError: Model <name> has no sqlite.bin file available.` or HTTP 400 Bad Request because `model_version` is missing in the request body.

---

## Proposed Remediation

Update `sdkpy/src/i6sdk/infer/infer.py` to match the server schema and `sdkweb`:

```python
def _sqlite_url(self, opts: Union[InferQueryOptions, Dict[str, Any]]) -> str:
    if direct_url := opts.get("sqlite_url"):
        return direct_url

    dataset = opts.get("dataset", "")
    model_name = opts.get("model_name", "")
    partitions = opts.get("partitions") or {}
    model_version = partitions.get("version", "")

    resp = self._services.apis().dsz_model_get(
        req={
            "params": {
                "dataset": dataset,
                "model_name": model_name,
            },
            "payload": {
                "model_version": model_version,
                "model_partitions": partitions,
            },
        },
        opts=opts.get("service_options"),
    )
    if not resp.ok:
        raise I6ApiError(
            message=f"Failed to get model metadata: {resp.text}",
            status=resp.status,
        )
    data = resp.parse() or {}
    urls = data.get("urls") or []
    by_name = {f.get("name"): f.get("url") for f in urls if isinstance(f, dict)}
    url = by_name.get("sqlite.bin")
    if not url:
        raise LookupError(f"Model {model_name} has no sqlite.bin file available in response.")
    return str(url)
```
