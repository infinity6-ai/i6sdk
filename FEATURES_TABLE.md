# Infinity6 SDK Feature Parity Matrix

This document provides a comprehensive feature-by-feature comparison across the three Infinity6 SDK implementations: **Go (`sdkgo`)**, **Web / JavaScript (`sdkweb`)**, and **Python (`sdkpy`)**.

Specific language implementations, method signatures, behaviors, and nuances are documented directly within each language's column.

---

## Feature Matrix

| Feature Area | Feature / Capability | `sdkgo` | `sdkweb` | `sdkpy` |
| :--- | :--- | :--- | :--- | :--- |
| **Architecture** | **Client Container Class** | ⚠️ **Partial (`Client`)**<br>• Minimal metadata wrapper (`New()`, `Version()`, `Config()`).<br>• Does not instantiate, configure, or hold sub-service clients.<br>• Configuration is limited to server URLs and `http.Client`. | ✅ **Full (`I6Sdk` / `Services`)**<br>• `I6Sdk` wraps internal `Services` container.<br>• Manages config lifecycle, global event dispatching, and singleton services (`auth`, `tokens`, etc.).<br>• Re-exported to `window.I6Sdk` for browser usage. | ✅ **Full (`I6Client` / `Services`)**<br>• `Services` acts as central service container.<br>• Aliased as `I6Client = Services`.<br>• Factory `create_client()` loads env vars (`I6_SDK_*`) and overrides. |
| | **Code Generation (`forgez`)** | ❌ **None**<br>• No Go generator exists (`forgego` is not implemented in `forgez`).<br>• No client-side call stubs are generated. | ✅ **`forgejs`**<br>• Generates single consolidated bundle `src/main/sdkapis/sdkapis.js`.<br>• Emits JSDoc typedefs, `Apis` class, and `SdkFuncs` class. | ✅ **`forgepy`**<br>• Generates modular Python package `sdkapis/` with `apis.py`, `sdkfuncs.py`, and `types.py`.<br>• Uses `TypedDict` and type annotations. |
| | **Typed Client API Methods** | ❌ **None**<br>• `c.Apis()` returns raw server route definitions (`[]apiz.Api`).<br>• Cannot be used to invoke endpoints over HTTP. | ✅ **`sdk.apis().<methodName>()`**<br>• CamelCase methods (e.g. `ingestzGetUrl`).<br>• Returns `Promise<Resp<T>>`.<br>• Does not throw on non-2xx HTTP status. | ✅ **`client.apis().<method_name>()`**<br>• Snake_case methods (e.g. `ingestz_get_url`).<br>• Returns `Resp[T]`.<br>• `.parse()` raises `I6ApiError` on non-2xx HTTP status. |
| | **Typed SdkFunc Workflow Methods** | ❌ **None**<br>• `c.SdkFuncs()` returns raw Go server interface definitions (`[]sdkfuncz.Func`).<br>• Cannot execute workflows. | ✅ **`sdk.funcs().<methodName>()`**<br>• High-level workflows exposed via `sdk.funcs()` and `services.funcs()`.<br>• CamelCase methods (e.g. `productRelevanceFashionFbt`). | ✅ **`client.sdk_funcs().<method_name>()`**<br>• Exposed via `client.sdk_funcs()` or alias `client.sdkFuncs()`.<br>• Snake_case methods (e.g. `product_relevance_fashion_fbt`).<br>• ⚠️ `client.funcs()` alias is missing. |
| **Authentication** | **Browser PKCE Redirect Flow** | ❌ **None** | ✅ **Supported**<br>• `auth.login(callbackUrl)` generates S256 PKCE challenge, stores verifier in `sessionStorage`, and navigates to `/api/userz/auth`.<br>• `auth.handleCallback()` validates CSRF `state`, exchanges code via `POST /api/userz/auth/exchange`, and redirects to app URL. | ❌ **None**<br>• No browser navigation or session-scoped redirect flow. |
| | **OAuth 2.0 Device Code Flow** | ❌ **None** | ❌ **None** | ✅ **Supported**<br>• `auth.device_code_start(client_id)` requests device code via `user_auth_device_code`.<br>• `auth.device_code_poll()` polls `/api/userz/auth/device-token` with configurable timeout/interval until confirmed/denied. |
| | **Direct Email / Password Auth** | ❌ **None** | ❌ **None**<br>• Browser users authenticate via interactive `i6userz` web portal. | ✅ **Supported**<br>• Direct non-interactive login via `auth.sign_in(email, password)` and account creation via `auth.sign_up(email, password)`. |
| | **User Profile Retrieval** | ❌ **None** | ✅ **`auth.user()`**<br>• Calls `GET /api/userz/auth/me` with Bearer token.<br>• Performs automatic silent token renewal before calling.<br>• Throws on HTTP error. | ✅ **`auth.me()`**<br>• Calls `apis().user_auth_me()`.<br>• Returns `None` if unauthenticated or on 401.<br>• Raises `I6AuthError` on other errors. |
| | **Enforced Login (`requireUser`)** | ❌ **None** | ✅ **`auth.requireUser()`**<br>• Fetches user profile; if null, automatically redirects the browser to `auth.login()`. | ❌ **None** |
| | **Safe Redirect Validation** | ❌ **None** | ✅ **`_isSafeRedirect(url)`**<br>• Prevents open redirect attacks by verifying URLs against `window.location.origin`, `config.callbackUrl`, or `config.userExchangeCallbackUrl`. | ❌ **None** |
| | **Checkout Ticket Redemption** | ❌ **None** | ✅ **`auth.handleCallback()`**<br>• Detects `checkout_code` param, calls `tokens.redeemCheckout()`, exchanges ticket via `/api/decsuite/auth/checkoutredeem`, and replaces `{{DATASET}}` in callback URL. | ❌ **None** |
| | **Session Revocation / Logout** | ❌ **None** | ✅ **`auth.logout(backurl)`**<br>• Calls `POST /api/userz/auth/logout` to revoke token family on server, clears tokens/sessionStorage, and redirects to `backurl`. | ✅ **`auth.sign_out()`**<br>• Calls `apis().user_auth_logout()` and clears local token storage. No redirect logic. |
| **Token Management** | **In-flight Deduplication** | ❌ **None** | ✅ **Supported**<br>• Deduplicates concurrent `refresh()` calls via `_refreshPromise`.<br>• Deduplicates concurrent dataset token requests via `_datasetTokenPromises` map. | ❌ **None**<br>• Concurrent refresh calls can trigger token family reuse detection on server, revoking the session. |
| | **Silent Automatic Token Renewal** | ❌ **None** | ✅ **Automatic**<br>• `getValidAccessToken(skewSeconds)` checks expiration threshold (default 30s) and automatically invokes `refresh()` before returning. Invoker uses this transparently. | ⚠️ **Manual**<br>• `get_token()` checks `expires_at > now_ms` but returns `None` if expired rather than refreshing.<br>• Does not auto-refresh during API calls. |
| | **Storage Adapters** | ❌ **None** | ✅ **Pluggable Web Storage**<br>• Defaults to global `localStorage` and `sessionStorage`.<br>• Accepts custom `I6Storage` objects or factory functions via `config.storage`. | ✅ **File & Memory Storage**<br>• `MemoryTokenStorage` (in-memory dict).<br>• `FileTokenStorage` (persists JSON to disk with auto-directory creation). |
| | **Dataset Token Derivation** | ❌ **None** | ✅ **`getValidDatasetAccessToken(dataset)`**<br>• Checks cache, auto-renews user token, requests all grants (`0x0F`) via `POST /api/userz/dataset-access-token/{dataset}`, and caches result. | ✅ **`get_valid_dataset_access_token(dataset)`**<br>• Checks cache, requests token via `POST /api/userz/dataset-access-token/{dataset}` using `get_token()`. Fails if access token expired. |
| **Ingestion** | **File Upload Payloads** | ❌ **None** | ✅ **Browser `File` Only**<br>• `ingest.upload({ file })` strictly requires browser `File` instance (`file instanceof File`). | ✅ **Flexible Payloads**<br>• `ingest.upload(file=...)` accepts `pathlib.Path`, file path string, raw `bytes`, or `bytearray`. |
| | **Upload Progress Monitoring** | ❌ **None** | ✅ **Supported**<br>• Provides real-time upload progress events via `onProgress({ loaded, total, percent })` powered by `XMLHttpRequest.upload`. | ❌ **None**<br>• Uses synchronous blocking `urllib.request.urlopen()`. |
| | **Upload Cancellation** | ❌ None | ✅ **Supported**<br>• Accepts standard DOM `AbortSignal` (`signal` option), rejecting promise with `AbortError`. | ❌ **None**<br>• In-flight uploads cannot be aborted. |
| | **Signed Upload URL Request** | ❌ None | ⚠️ **Private**<br>• Handled via private `#getUrl()` helper inside `upload()`. | ✅ **Public Method**<br>• `ingest.get_url(dataset, table, partitions)` is a public API returning `{ingest_id, url}`. |
| **Inference Engine** | **Named SQLite Execution (`run`)** | ❌ **None** | ✅ **Supported**<br>• Runs named query `model-{name}-input`.<br>• Calls `apis().tfzModelPredict()`.<br>• Runs output query `model-{name}-output`.<br>• Assigns `rec_id` via `crypto.randomUUID()`. | ✅ **Supported**<br>• Runs named query `model-{name}-input`.<br>• Calls `apis().tfz_model_predict()`.<br>• Runs output query `model-{name}-output`.<br>• Assigns `rec_id` via `uuid.uuid4()`. |
| | **Telemetry Event Tracking** | ❌ **None** | ✅ **Automatic**<br>• `run()` emits `recraw` event with session ID, user ID, and results.<br>• `selectResults()` emits `recshow` event with shown SKUs.<br>• Calls `ingestzEventRelevanceFashionEvent`. | ❌ **None**<br>• Neither `recraw` nor `recshow` event is emitted.<br>• `session_id` and `user_id` are omitted from `InferQueryOptions`. |
| | **Required Run Identifiers** | N/A | **`sessionId` & `userId` Required**<br>• Throws `TypeError` if `sessionId` or `userId` is omitted (needed for event attribution). | **Neither Required**<br>• `InferQueryOptions` takes only `dataset`, `model_name`, `partitions`, `params`. |
| | **Background Model Revalidation** | ❌ **None** | ✅ **Non-blocking**<br>• If model file is cached, `run()` executes against it immediately and triggers background revalidation (`this.update()`). | ❌ **Blocking**<br>• `update()` revalidates remote ETag synchronously before executing query. |
| | **Bulk Preload / Warmup** | ❌ **None** | ✅ **`warmupDataset(opts)`**<br>• Top-level helper on `I6Sdk` to preload multiple model SQLite bundles in parallel on page load. | ❌ **None** |
| | **Model Metadata Fetch Schema** | N/A | ✅ **Matches Server**<br>• Calls `dszModelGet` via POST body `{model_version, model_partitions}` and parses `resp.urls`. | ⚠️ **Schema Bug**<br>• Calls `dsz_model_get` with `query: partitions` and expects `resp.files['sqlite.bin']`. Fails against live server. |
| **Dashboard Analytics** | **Dash Query Execution** | ❌ **None** | ✅ **`dash.query(opts)`**<br>• Resolves signed SQLite URL via `dszSignedUrl` (`i6tb://results...`).<br>• Downloads asset, opens in SQLite, and runs named JSON query `dash-{name}`. | ✅ **`dash.query(opts)`**<br>• Resolves signed SQLite URL via `dsz_signed_url`.<br>• Downloads asset, opens in SQLite, and runs named JSON query `dash_name`. |
| **Asset Caching & FS** | **Storage Backend** | ❌ **None** | ✅ **OPFS (Origin Private File System)**<br>• Operates under browser `navigator.storage.getDirectory()` path `i6/fs`. | ✅ **POSIX Disk Filesystem**<br>• Operates under `~/.cache/i6/sdkpy/fs` (configurable via `cache_dir`). |
| | **ETag Cache Invalidation** | ❌ **None** | ✅ **Supported**<br>• Manages `{name}/version.txt`, `{name}/checkedat.txt`, `{name}/v{version}/etag.txt`, and `blob.bin`.<br>• Skips check if checked within 300s. | ✅ **Supported**<br>• Mirrors OPFS file structure and 300s check interval on local disk. |
| **SQLite Runtime** | **Engine** | ❌ **None** | ✅ **WebAssembly (`sql.js`)**<br>• Loaded lazily from window global (`initSqlJs`) or jsDelivr CDN. | ✅ **Native C (`sqlite3`)**<br>• Python standard library `sqlite3` using URI read-only connection (`?mode=ro`). |
| **Transport & Invoker** | **Underlying HTTP Client** | Standard `net/http` (internal only) | **Browser Fetch API**<br>• Pluggable `fetcher` on container and per-request.<br>• Supports `credentials: "include"`. | **Python `urllib.request`**<br>• Pluggable `fetcher` callback on `I6SdkConfig`.<br>• Configurable socket timeout. |
| | **Pre-Request Mutation / URL Signing** | ❌ **None** | ✅ **`prepareRequest(opts)`**<br>• Async hook executed before dispatch.<br>• Enables client-side URL signing (e.g. Ed25519 `sig` parameter) to avoid browser bearer tokens. | ❌ **None** |
| | **Synthetic Auth Error Responses** | ❌ **None** | ✅ **Supported**<br>• Synthesizes 401 or 403 `Response` objects client-side if required token types are missing before network dispatch. | ❌ **None**<br>• Dispatches request without authorization header, relying on server rejection. |
| | **Structured Error Model** | Standard Go `error` | **`apiErrorFromResponse` / `Resp.error`**<br>• Returns `{status, code, message, cause}`.<br>• `parse()` throws standard `Error`. | **Typed Exception Hierarchy**<br>• `parse()` raises `I6ApiError` (subclass of `I6Error`), with `status`, `code`, and `cause`.<br>• Specific `I6AuthError`, `I6ConnectionError`, `I6ValidationError`. |
| **Bitmasks & Grants** | **Bitmask Operations** | `sdkcorez/grants` & `commonz/bitmapz` | ✅ **`I6Sdk.bitmapz`**<br>• `join`, `containsAll`, `containsAny`, `split`, `set`.<br>• Binary serialization: `serialize(v)` and `deserialize(bytes)`. | ⚠️ **`internal.bitmapz`**<br>• `join`, `contains_all`, `contains_any`, `split`, `set_bits`.<br>• ❌ Binary `serialize`/`deserialize` missing.<br>• Not exported at package root. |
| | **Grants Enum / Constants** | `sdkcorez/grants` (`GrantSystem`, `GrantAdmin`, `GrantWrite`, `GrantRead`, `GrantAll`) | ✅ **`I6Sdk.Grants`**<br>• Frozen constants: `System: 1`, `Admin: 2`, `Write: 4`, `Read: 8`, `All: 15`. Exported at root. | ❌ **None**<br>• No `Grants` constants defined. Developers must use raw numeric flags. |
| **Event System** | **EventTarget Binding** | ❌ **None** | ✅ **`sdk.bind(evt, fn)`**<br>• Internal `EventTarget` dispatcher on `Services`.<br>• Allows subscribing to SDK events. | ❌ **None** |
| **Developer Tooling** | **Interactive Web UI Tools** | ❌ **None** | ✅ **Included**<br>• Interactive web console (`console.html`).<br>• Token management UI (`src/site/tokens/`).<br>• Visual inference demo (`src/site/infer-demo/`). | ❌ **None** |

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
