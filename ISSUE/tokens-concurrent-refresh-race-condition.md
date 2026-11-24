# Issue: Missing Concurrency Deduplication in `sdkpy` Token Management

**Component:** `sdkpy` (`i6sdk.tokens.tokens` & `i6sdk.auth.auth`)  
**Severity:** High (Token family revocation on concurrent requests)  
**Status:** Open  

---

## Summary

In `sdkweb`, concurrent requests to refresh user tokens or obtain dataset access tokens are deduplicated using in-flight promise caches (`_refreshPromise` and `_datasetTokenPromises`). In `sdkpy`, there is no deduplication or concurrency control. In multi-threaded or concurrent applications, simultaneous requests when a token expires cause multiple refresh requests with the same refresh token. Because `i6userz` implements single-use refresh token rotation, the second request is detected as a token family replay attack, causing the server to revoke the user's entire session family.

---

## Root Cause Analysis

### 1. In-flight Deduplication in `sdkweb`
In `sdkweb/src/main/tokens/tokens.js`:
```javascript
async refresh() {
  if (this._refreshPromise) {
    return this._refreshPromise;
  }

  this._refreshPromise = (async () => {
    try {
      // execute POST /api/userz/refresh-user
    } finally {
      this._refreshPromise = null;
    }
  })();

  return this._refreshPromise;
}
```
And similarly for dataset tokens:
```javascript
if (this._datasetTokenPromises.has(dataset)) {
  const inFlight = await this._datasetTokenPromises.get(dataset);
  return inFlight?.access_token || null;
}
```

### 2. Missing Deduplication / Thread Safety in `sdkpy`
In `sdkpy/src/i6sdk/auth/auth.py`:
```python
def refresh(self) -> Optional[UserToken]:
    ut = self._services.tokens().get_user_token()
    if not ut or not (refresh_token := ut.get("user_refresh_token")):
        return None

    resp = self._services.apis().user_auth_refresh(...)  # <--- Multiple threads send the same refresh token simultaneously!
    ...
```
When multiple threads in a Python web server (e.g. FastAPI/Gunicorn/ThreadPoolExecutor) invoke SDK APIs concurrently as the access token expires:
1. Thread A and Thread B both read the same `user_refresh_token`.
2. Thread A sends `POST /api/userz/refresh-user`. Server rotates the token family and responds with new tokens.
3. Thread B sends `POST /api/userz/refresh-user` with the now-invalidated prior refresh token.
4. Server detects refresh token reuse, revokes all tokens for that user session, and responds with HTTP 400/401.
5. Thread B catches the error and wipes the user session (`self._services.tokens().remove_user_token()`).

---

## Proposed Remediation

1. Add thread-safe deduplication locking in `Tokens` using `threading.Lock` or future/promise caching:
```python
import threading

class Tokens:
    def __init__(self, services: Services, storage: Optional[TokenStorage] = None) -> None:
        self._services = services
        self._storage: TokenStorage = storage or MemoryTokenStorage()
        self._refresh_lock = threading.Lock()
        self._dataset_locks: Dict[str, threading.Lock] = {}
        self._dataset_locks_lock = threading.Lock()
```
2. When refreshing or obtaining dataset tokens, synchronize across callers or reuse in-flight results so only a single network refresh is dispatched per rotation cycle.
