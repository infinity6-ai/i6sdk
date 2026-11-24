# Issue: Missing Upload Progress and Cancellation in `sdkpy` Ingest

**Component:** `sdkpy` (`i6sdk.ingest.ingest`)  
**Severity:** Medium (Degraded UX for large file uploads)  
**Status:** Open  

---

## Summary

In `sdkweb`, `ingest.upload()` supports real-time upload progress reporting via `onProgress({ loaded, total, percent })` and cancellation via standard `AbortSignal`. In `sdkpy`, `upload()` uses blocking `urllib.request.urlopen()` with no progress feedback and no cancellation mechanism, which causes poor observability and inability to abort large file or dataset uploads.

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
