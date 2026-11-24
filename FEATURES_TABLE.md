# Infinity6 SDK Feature Parity Matrix

This document provides a comprehensive feature-by-feature comparison across the three Infinity6 SDK implementations: **Go (`sdkgo`)**, **Web / JavaScript (`sdkweb`)**, and **Python (`sdkpy`)**.

---

## Feature Matrix

| Feature Area | Feature / Capability | `sdkgo` | `sdkweb` | `sdkpy` | Nuances & Behavioral Differences |
| :--- | :--- | :---: | :---: | :---: | :--- |
| **Architecture** | **Client Container Class** | ⚠️ Partial (`Client`) | ✅ Yes (`I6Sdk` / `Services`) | ✅ Yes (`I6Client` / `Services`) | `sdkgo` is a minimal metadata container. `sdkweb` wraps `Services` in `I6Sdk`. `sdkpy` aliases `I6Client = Services`. |
| | **Code Generation** | ❌ None | ✅ `forgejs` | ✅ `forgepy` | `forgejs` generates bundle `sdkapis.js`; `forgepy` generates `sdkapis/` module with `apis.py`, `sdkfuncs.py`, `types.py`. No Go generator exists. |
| | **Typed Client API Methods** | ❌ None | ✅ Yes (`apis()`) | ✅ Yes (`apis()`) | `sdkweb` uses camelCase (e.g. `ingestzGetUrl`). `sdkpy` uses snake_case (e.g. `ingestz_get_url`). `sdkgo` only exposes raw server-side route definitions via `Apis()`. |
| | **Typed SdkFunc Workflow Methods** | ❌ None | ✅ Yes (`funcs()`) | ✅ Yes (`sdk_funcs()` / `sdkFuncs`) | `sdkweb` exposes `funcs()`. `sdkpy` exposes `sdk_funcs()` / `sdkFuncs`, but has no `funcs()` method. `sdkgo` only returns raw `[]sdkfuncz.Func` structs. |
| **Authentication** | **Browser PKCE Redirect Flow** | ❌ None | ✅ Yes (`auth.login`, `auth.handleCallback`) | ❌ None | `sdkweb` generates RFC 7636 PKCE `code_verifier`/`code_challenge`, persists state in `sessionStorage`, redirects to `/api/userz/auth`, and verifies CSRF `state`. |
| | **OAuth 2.0 Device Code Flow** | ❌ None | ❌ None | ✅ Yes (`device_code_start`, `device_code_poll`) | Designed for CLI/headless environments. `sdkpy` starts flow and polls `/api/userz/auth/device-token` until authorization completes or times out. |
| | **Direct Email / Password Sign-in** | ❌ None | ❌ None | ✅ Yes (`auth.sign_in`, `auth.sign_up`) | `sdkpy` can authenticate directly via password API endpoints. `sdkweb` delegates user authentication to the interactive `userz` web UI. |
| | **Profile Retrieval (`me` / `user`)** | ❌ None | ✅ Yes (`auth.user`) | ✅ Yes (`auth.me`) | `sdkweb` sends bearer token and parses payload. `sdkpy` handles 401 silently returning `None` or raises `I6AuthError`. |
| | **Enforced Login (`requireUser`)** | ❌ None | ✅ Yes (`auth.requireUser`) | ❌ None | `sdkweb` automatically initiates `login()` redirect if the user is unauthenticated. |
| | **Safe Redirect Validation** | ❌ None | ✅ Yes (`_isSafeRedirect`) | ❌ None | `sdkweb` guards against open redirect attacks by verifying URLs against current origin and configured callback URLs. |
| | **Checkout Ticket Redemption** | ❌ None | ✅ Yes (`redeemCheckout`) | ❌ None | `sdkweb` handles Stripe checkout ticket tokens (`checkout_code`) via `/api/decsuite/auth/checkoutredeem` inside `handleCallback`. |
| | **Session Revocation / Logout** | ❌ None | ✅ Yes (`auth.logout`) | ✅ Yes (`auth.sign_out`) | `sdkweb` calls `/api/userz/auth/logout` and navigates browser to `backurl`. `sdkpy` calls API and clears local token storage. |
| **Token Management** | **In-flight Deduplication** | ❌ None | ✅ Yes (`_refreshPromise`, `_datasetTokenPromises`) | ❌ None | `sdkweb` deduplicates simultaneous concurrent refresh calls and concurrent dataset token requests to prevent duplicate network calls. |
| | **Silent Automatic Token Renewal** | ❌ None | ✅ Yes (`getValidAccessToken`) | ⚠️ Partial | `sdkweb` automatically triggers a token refresh before making API requests if the access token is near expiry. `sdkpy` checks expiration timestamp in `get_token()`, but does not automatically auto-refresh. |
| | **Storage Adapters** | ❌ None | ✅ `localStorage`, `sessionStorage`, custom | ✅ `MemoryTokenStorage`, `FileTokenStorage` | `sdkweb` uses Web Storage API or custom adapters. `sdkpy` persists to local JSON files on disk or memory. |
| | **Dataset Access Token Retrieval** | ❌ None | ✅ Yes (`getValidDatasetAccessToken`) | ✅ Yes (`get_valid_dataset_access_token`) | Both derive scoped dataset tokens with grant bitmasks via `/api/userz/dataset-access-token/{dataset}`. |
| **Ingestion** | **File Upload Execution** | ❌ None | ✅ Yes (`ingest.upload`) | ✅ Yes (`ingest.upload`) | `sdkweb` requires a browser `File` instance. `sdkpy` accepts `pathlib.Path`, file path `str`, `bytes`, or `bytearray`. |
| | **Upload Progress Monitoring** | ❌ None | ✅ Yes (`onProgress`) | ❌ None | `sdkweb` uses `XMLHttpRequest.upload` progress events to calculate bytes sent and percent completed. `sdkpy` uses synchronous `urllib.request`. |
| | **Upload Cancellation** | ❌ None | ✅ Yes (`AbortSignal`) | ❌ None | `sdkweb` accepts standard DOM `AbortSignal` to cancel in-flight uploads. |
| | **Signed Upload URL Request** | ❌ None | ⚠️ Private helper (`#getUrl`) | ✅ Public method (`get_url`) | `sdkpy` exposes `get_url()` as a reusable public API returning `{ingest_id, url}`. `sdkweb` wraps it inside `upload()`. |
| **Inference Engine** | **Named SQLite Query Execution** | ❌ None | ✅ Yes (`infer.run`) | ✅ Yes (`infer.run`) | Both execute input query (`model-{name}-input`), call server prediction (`tfz_model_predict`), then run output query (`model-{name}-output`). |
| | **Telemetry Event Tracking** | ❌ None | ✅ Yes (`recraw` & `recshow`) | ❌ None | `sdkweb` automatically logs `recraw` event on `run()` and `recshow` on `selectResults()` using `ingestzEventRelevanceFashionEvent`. `sdkpy` logs no events. |
| | **Required Identifiers** | N/A | `sessionId` & `userId` required | Neither required | In `sdkweb`, `sessionId` and `userId` are mandatory for telemetry tracking. In `sdkpy`, they are omitted from `InferQueryOptions`. |
| | **Background Model Revalidation** | ❌ None | ✅ Yes | ❌ None | When a model file is cached, `sdkweb` returns and computes immediately while triggering background cache revalidation. `sdkpy` revalidates synchronously. |
| | **Bulk Preloading / Warmup** | ❌ None | ✅ Yes (`warmupDataset`) | ❌ None | `sdkweb` provides `I6Sdk.warmupDataset()` to preload multiple model files in parallel. |
| | **Model Metadata Fetch Schema** | N/A | ✅ Matches Server API | ⚠️ Outdated Schema | `sdkweb` sends POST body `{model_version, model_partitions}` and parses `urls`. `sdkpy` sends `query: partitions` and expects `files.sqlite.bin.url` (divergence from server). |
| **Dashboard Analytics** | **Dash Query Execution** | ❌ None | ✅ Yes (`dash.query`) | ✅ Yes (`dash.query`) | Both resolve signed SQLite URL via `dsz_signed_url` and execute named SQLite query `dash-{name}`. |
| **Asset Caching & FS** | **Storage Backend** | ❌ None | ✅ OPFS (Origin Private File System) | ✅ Local POSIX Filesystem | `sdkweb` uses `navigator.storage.getDirectory()`. `sdkpy` stores files in `~/.cache/i6/sdkpy/fs/` (configurable via `cache_dir`). |
| | **ETag Cache Invalidation** | ❌ None | ✅ Yes (`checkedat.txt`, `etag.txt`) | ✅ Yes (matches OPFS schema) | Both avoid re-downloading files if the remote ETag has not changed and check interval is within 5 minutes. |
| **SQLite Runtime** | **Engine** | ❌ None | ✅ WebAssembly (`sql.js`) | ✅ Native C (`sqlite3`) | `sdkweb` dynamically loads `sql.js` wasm. `sdkpy` uses Python's built-in standard library `sqlite3` in read-only mode (`?mode=ro`). |
| **Transport & Invoker** | **Underlying HTTP Client** | Standard `net/http` (internal) | Browser Fetch API | Python `urllib.request` | Both `sdkweb` and `sdkpy` support pluggable custom `fetcher` functions. |
| | **Pre-Request Mutation / URL Signing**| ❌ None | ✅ Yes (`prepareRequest`) | ❌ None | `sdkweb` provides `prepareRequest` hook for client-side Ed25519/HMAC URL signing to avoid bearer token exposure in the browser. |
| | **Synthetic Auth Error Responses** | ❌ None | ✅ Yes | ❌ None | `sdkweb` synthesizes 401 or 403 `Response` objects if required token types are missing before dispatching network calls. |
| | **Structured Error Model** | Go standard `error` | `apiErrorFromResponse` / `Resp.error` | Typed exceptions (`I6ApiError`, etc.) | `sdkpy` maps HTTP statuses to custom exception hierarchies with status and error codes. |
| **Bitmasks & Grants** | **Bitmask Operations** | `sdkcorez/grants` & `commonz/bitmapz` | ✅ `I6Sdk.bitmapz` | ⚠️ Partial (`internal.bitmapz`) | `sdkweb` exports `bitmapz` with `join`, `containsAll`, `containsAny`, `split`, `set`, `serialize`, `deserialize`. `sdkpy` lacks binary serialization. |
| | **Grants Enum / Constants** | `sdkcorez/grants.Grant*` | ✅ `I6Sdk.Grants` | ❌ None | `sdkweb` exports `Grants` (`System: 1`, `Admin: 2`, `Write: 4`, `Read: 8`, `All: 15`). `sdkpy` does not define `Grants`. |
| **Event System** | **EventTarget Binding** | ❌ None | ✅ Yes (`sdk.bind(evt, fn)`) | ❌ None | `sdkweb` incorporates an internal `EventTarget` dispatcher on `Services` for custom event subscription. |
| **Developer Tooling** | **Interactive Web UI Tools** | ❌ None | ✅ Yes (`src/site/*`) | ❌ None | `sdkweb` includes developer console (`console.html`), token management dashboard (`tokens/`), and interactive inference demo (`infer-demo/`). |

---

## Detailed Nuances & Divergences

### 1. Model Get API Divergence (`dsz_model_get`)
* **`sdkweb` (Correct implementation matching server):**
  ```javascript
  const resp = await this._services.apis().dszModelGet({
    params: { dataset, model_name: modelName },
    payload: { model_version: opts.partitions.version, model_partitions: partitions },
  });
  const data = resp.parse();
  const byName = Object.fromEntries((data?.urls ?? []).map((file) => [file.name, file.url]));
  return { sqlite: byName["sqlite.bin"] };
  ```
* **`sdkpy` (Outdated schema):**
  ```python
  resp = self._services.apis().dsz_model_get(
      req={
          "params": {"dataset": dataset, "model_name": model_name},
          "query": partitions,
      }
  )
  data = resp.parse() or {}
  sqlite_info = (data.get("files") or {}).get("sqlite.bin") or {}
  url = sqlite_info.get("url")
  ```
  `sdkpy` passes partitions via query params rather than POST payload, and inspects `data["files"]["sqlite.bin"]` instead of `data["urls"]`.

### 2. Method Naming Discrepancies
* **SdkFuncs Accessor:**
  * In `sdkweb`: `sdk.funcs()` or `services.funcs()`.
  * In `sdkpy`: `client.sdk_funcs()` or `client.sdkFuncs()` (no `funcs()` shorthand).
* **Token Expiration Checks:**
  * In `sdkweb`: `tokens.isAccessTokenExpired()` and `tokens.isRefreshTokenExpired()` accept an explicit `skewSeconds` safety threshold.
  * In `sdkpy`: Expiration check is embedded in `get_token()`, but silent refresh is not automated during standard API invocations.

### 3. Telemetry in Inference
* In `sdkweb`, calling `infer.run()` automatically emits a `recraw` event containing the raw recommendation results, session ID, and user ID to the ingestion pipeline. Subsequently, `infer.selectResults()` emits a `recshow` event containing only the displayed SKU IDs.
* In `sdkpy`, no telemetry events are emitted, and the `Infer` client functions solely as an inference calculator.
