# Issue: Missing Upload Progress and Cancellation in `sdkpy` Ingest

**Component:** `sdkpy` (`i6sdk.ingest.ingest`)  
**Severity:** Medium (Degraded UX for large file uploads)  
**Status:** Resolved  

---

## Summary

In `sdkweb`, `ingest.upload()` supports real-time upload progress reporting via `onProgress({ loaded, total, percent })` and cancellation via standard `AbortSignal`. In `sdkpy`, `upload()` previously used blocking `urllib.request.urlopen()` with whole-file memory buffering, no progress feedback, and no cancellation mechanism.

---

## Resolution

1. Updated `Ingest.upload()` in `sdkpy/src/i6sdk/ingest/ingest.py` to adopt a stream-first architecture accepting `stream: Union[BinaryIO, str, bytes]`.
2. Implemented `ProgressStreamReader` to stream data in 64KB chunks with $O(1)$ memory usage.
3. Added real-time progress callbacks (`on_progress` / `onProgress`) providing `{"loaded": int, "total": int, "percent": float}`.
4. Added cancellation support via `cancel_event` (`threading.Event`) and `signal` alias, aborting in-flight uploads with `I6ValidationError`.
5. Fixed `Ingest.get_url()` API schema bug to send `table`, `partitions`, and `amount: 1` in the POST body `payload` and extract signed URLs from `uploads[0]`.
6. Added comprehensive unit tests in `sdkpy/tests/test_ingest.py`.

---

## Root Cause Analysis

In `sdkweb/src/main/ingest/ingest.js`:
```javascript
// Uploads using XMLHttpRequest with upload progress listener and abort signal
xhr.upload.addEventListener('progress', (e) => {
  const total = e.lengthComputable ? e.total : file.size;
  onProgress({
    loaded: e.loaded,
    total,
    percent: total > 0 ? (e.loaded / total) * 100 : 100,
  });
});
signal?.addEventListener('abort', abort, { once: true });
```

In `sdkpy/src/i6sdk/ingest/ingest.py`:
```python
req = urllib.request.Request(
    url=signed_url,
    data=payload_bytes,
    headers={...},
    method="PUT",
)
with urllib.request.urlopen(req, timeout=120.0) as resp:
    ...
```
* `payload_bytes` is sent in a single monolithic buffer.
* No chunked streaming or progress callback mechanism is available.
* Cannot cancel or interrupt an in-flight upload without terminating the thread or process.

---

## Proposed Remediation

1. Introduce an `on_progress: Optional[Callable[[int, int, float], None]]` parameter to `Ingest.upload()`.
2. Wrap payload streaming using a custom file-like generator or chunked reader that notifies progress after each block write.
3. Accept an optional cancellation token or timeout override.
