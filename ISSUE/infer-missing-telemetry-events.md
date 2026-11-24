# Issue: Missing Telemetry Events (`recraw` & `recshow`) in `sdkpy` Inference

**Component:** `sdkpy` (`i6sdk.infer.infer`)  
**Severity:** Medium-High (Silently drops feedback loop and recommendation analytics)  
**Status:** Open  

---

## Summary

In `sdkweb`, the inference lifecycle automatically logs telemetry events:
1. `infer.run()` automatically emits a `recraw` event via `apis().ingestzEventRelevanceFashionEvent(...)` containing the complete generated recommendation payload, model ID, version, `session_id`, `user_id`, and `rec_id`.
2. `infer.selectResults()` automatically emits a `recshow` event containing the list of displayed SKU IDs.

In `sdkpy/src/i6sdk/infer/infer.py`, `rec_id` is generated via `uuid.uuid4()`, but **neither event is ever sent**. Furthermore, `session_id` and `user_id` are not even accepted in `InferQueryOptions`, leaving Python inference unable to report model performance and engagement telemetry.

---

## Root Cause Analysis

### 1. `sdkweb` Implementation
In `sdkweb/src/main/infer/infer.js`:
```javascript
// During run():
const recId = crypto.randomUUID()
const results = { ...result, rec_id: recId }
this.#sendRecEvent(eventOpts, ids, recId, "recraw", results)

// During selectResults():
this.#sendRecEvent(opts, ids, recId, "recshow", shown)
```
Where `#sendRecEvent` invokes:
```javascript
await this._services.apis().ingestzEventRelevanceFashionEvent({
  params: { dataset: opts.dataset },
  payload: {
    id: crypto.randomUUID(),
    store_id: opts.storeId,
    session_id: ids.sessionId,
    user_id: ids.userId,
    created_at: new Date().toISOString().split('.')[0] + 'Z',
    [kind]: {
      rec_id: recId,
      model_id: opts.modelName,
      model_version: opts.modelVersion,
      payload: JSON.stringify(payload),
    },
  },
}, opts.serviceOptions)
```

### 2. `sdkpy` Implementation
In `sdkpy/src/i6sdk/infer/infer.py`:
```python
rec_id = str(uuid.uuid4())
results = dict(result or {})
results["rec_id"] = rec_id
return results  # <--- BUG: Returns immediately without emitting recraw event!
```
And `select_results`:
```python
def select_results(self, opts: Dict[str, Any]) -> Dict[str, Any]:
    callback = opts.get("select_results") or default_select_results
    req = opts.get("req")
    results = opts.get("results")
    return callback(req, results)  # <--- BUG: No recshow event emitted!
```

---

## Impact

When recommendation models are evaluated through `sdkpy` (e.g. in backend batch services, Python microservices, or Jupyter notebooks), the central telemetry and data lake receive no `recraw` or `recshow` events. This skews or eliminates conversion attribution, offline evaluation metrics, and retraining data pipelines.

---

## Proposed Remediation

1. Add `session_id: str` and `user_id: str` to `InferQueryOptions`.
2. Implement `#_send_rec_event` helper in `sdkpy/src/i6sdk/infer/infer.py` using `self._services.apis().ingestz_event_relevance_fashion_event(...)`.
3. Call `_send_rec_event(..., kind="recraw")` asynchronously or synchronously inside `run()`.
4. Call `_send_rec_event(..., kind="recshow")` inside `select_results()`.
