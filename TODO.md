# Pre-Production SDK Cleanup & Simplification TODOs

Checklist of cleanup tasks to eliminate legacy aliases, redundant overloads, dual-casing fallbacks, and deprecated pre-release artifacts before entering production.

---

## Category 1: Drop CamelCase Method Aliases & Built-in Shadowing in Python (`sdkpy`)

- [ ] **Drop `forgepy` camelCase method alias generator**
  - `forgez/forgepy/forgepyapi/forgepyapi.go`: Remove generation of `%s = %s` camelCase aliases in `Apis` class.
  - `forgez/forgepy/forgepysdkfunc/forgepysdkfunc.go`: Remove generation of `%s = %s` camelCase aliases in `SdkFuncs` class.
- [ ] **Regenerate `sdkpy` APIs and SdkFuncs**
  - `sdkpy/src/i6sdk/sdkapis/apis.py`: Re-generate to strip ~43 redundant camelCase method aliases.
  - `sdkpy/src/i6sdk/sdkapis/sdkfuncs.py`: Re-generate to strip redundant camelCase SdkFunc method aliases.
- [ ] **Remove manual camelCase method aliases in Python service classes**
  - `sdkpy/src/i6sdk/services.py`: Remove `setConfig = set_config`.
  - `sdkpy/src/i6sdk/ingest/ingest.py`: Remove `getUrl = get_url`.
  - `sdkpy/src/i6sdk/infer/infer.py`: Remove `selectResults = select_results`.
- [ ] **Remove built-in shadowing alias in `bitmapz.py`**
  - `sdkpy/src/i6sdk/internal/bitmapz.py`: Remove `set = set_bits` to avoid shadowing Python's built-in `set` type.
- [ ] **Update tests & documentation**
  - `sdkpy/tests/test_apis.py`: Remove test assertions expecting camelCase method aliases.
  - `FEATURES_TABLE.md`: Update documentation to reflect pure `snake_case` in `sdkpy`.

---

## Category 2: Standardize Parameter Casing & Drop Dual-Lookup Fallbacks

- [ ] **Python (`sdkpy`): Enforce canonical `snake_case` options**
  - `sdkpy/src/i6sdk/infer/infer.py`: Drop camelCase option lookups (`sessionId`, `userId`, `modelName`, `sqliteUrl`, `serviceOptions`, `selectResults`, `recId`, `modelVersion`). Keep strictly canonical `snake_case`.
  - `sdkpy/src/i6sdk/invoker.py`: Remove `bearerToken` fallback; retain only `bearer_token`.
  - `sdkpy/src/i6sdk/ingest/ingest.py`: Remove `onProgress` (keep `on_progress`) and `signal` (keep `cancel_event`) in `upload()`.
  - `sdkpy/src/i6sdk/tokens/tokens.py`: Remove `private_key` (keep `priv_key`) and `sub` (keep `user_id`).
- [ ] **Go (`sdkgo`): Standardize partition key on `"store_id"`**
  - `sdkgo/dash/dash.go`: Remove fallback to `opts.Partitions["storeId"]`; standardize strictly on `"store_id"`.
  - `sdkgo/infer/infer.go`: Remove fallback to `partitions["storeId"]`; standardize strictly on `"store_id"`.
- [ ] **Web (`sdkweb`): Enforce canonical configuration & option keys**
  - `sdkweb/src/main/index.js`: Remove `partition` alias in `WarmupModelOptions`; retain `partitions`.
  - `sdkweb/src/main/invoker/invoker.js`: Remove `config.fetch` fallback; retain `config.fetcher`.
  - `sdkweb/src/main/tokens/tokens.js`: Remove `config.localStorage` fallback; retain `config.storage`.
  - `sdkweb/src/main/auth/auth.js`: Remove `config.storageSession` fallback; retain `config.sessionStorage`.

---

## Category 3: Clean Up Class, Type, and Model Aliases

- [ ] **Go (`sdkcorez`): Drop unused request signing type aliases**
  - `sdkcorez/reqsign/reqsign.go`: Remove unused `BuildRequestSignatureData` and `BuildRequestSignatureOptions`.
- [ ] **Python (`sdkpy`): Drop legacy storage & input aliases**
  - `sdkpy/src/i6sdk/tokens/storage.py`: Remove `FileTokenStorage = SqliteTokenStorage`.
  - `sdkpy/src/i6sdk/tokens/__init__.py` and `sdkpy/src/i6sdk/__init__.py`: Remove `FileTokenStorage` from exports.
  - `sdkpy/src/i6sdk/ingest/ingest.py`: Remove `StreamInput = StreamItem`.
- [ ] **Python (`sdkpy`): Standardize `I6Client` naming**
  - `sdkpy/src/i6sdk/services.py`: Rename class `Services` directly to `I6Client` (aligning with Go's `I6Client` and Web's `I6Sdk`).
  - `sdkpy/src/i6sdk/__init__.py`: Export `I6Client` directly.

---

## Category 4: Eliminate 1-Line Functional Wrappers & Redundant Prefixing

- [ ] **Go (`sdkcorez/grants`): Remove 1-line wrapper and stuttering constants**
  - `sdkcorez/grants/grants.go`: Remove `AllGrants()` (keep `All()`).
  - `sdkcorez/grants/grants.go`: Remove `GrantSystem`, `GrantAdmin`, `GrantWrite`, `GrantRead` (keep `System`, `Admin`, `Write`, `Read`).
  - `sdkcorez/grants/grants_unit_test.go`: Update unit tests to verify canonical constants (`System`, `Admin`, etc.).
- [ ] **Python (`sdkpy/grants.py`): Enforce PEP 8 uppercase enum constants**
  - `sdkpy/src/i6sdk/grants.py`: Remove TitleCase aliases (`System`, `Admin`, `Write`, `Read`, `All`); retain canonical `SYSTEM`, `ADMIN`, `WRITE`, `READ`, `ALL`.
  - `sdkpy/tests/`: Update any test references relying on TitleCase `Grants`.

---

## Category 5: Purge Deprecated Legacy Storage Keys & Obsolete Fields

- [ ] **Go (`sdkgo`): Drop deprecated monolithic token storage keys**
  - `sdkgo/tokens/tokens.go`: Remove `DatasetTokensKey` and `DatasetRefreshTokensKey`.
- [ ] **Python (`sdkpy`): Drop deprecated monolithic token storage keys**
  - `sdkpy/src/i6sdk/tokens/tokens.py`: Remove `DATASET_TOKENS_KEY` and `DATASET_REFRESH_TOKENS_KEY`.
- [ ] **Web (`sdkweb`): Remove legacy JWT user fields in documentation & types**
  - `sdkweb/src/main/auth/auth.js`: Remove deprecated legacy fields `sub`, `display`, and `email` from `AuthPayload` typedef.
