# Python SDK (`sdkpy`) Roadmap & Detailed Task Breakdown

This document provides a comprehensive technical blueprint and task tracker for implementing the Python SDK (`sdkpy`) and its AST-based code generator (`forgez/forgepy`).

---

## Core Design Principles

1. **Zero Runtime Dependencies (`dependencies = []`)**:
   - The runtime library strictly relies on the **Python Standard Library** (Python $\ge 3.10$).
   - Avoids dependency hell, version conflicts, and heavy download footprints.
   - Works immediately in lightweight, constrained, or secure environments (AWS Lambda, Airflow, Spark, Cloud Functions, air-gapped systems).

2. **Standard Library Alignments**:
   - **HTTP Client**: Built-in `urllib.request` / `http.client` (no required `httpx` or `requests`).
   - **Data Schemas**: Built-in `typing.TypedDict` and `dataclasses` (no `pydantic`).
   - **Database**: Built-in `sqlite3` (no third-party WASM or wrapper libraries).
   - **Crypto / Encodings**: Built-in `hashlib`, `base64`, `secrets` (no external cryptography packages).
   - **Type System**: Python $\ge 3.10$ native type hints (no `typing-extensions`).
   - **Inference**: Server-side model evaluation via `tfzModelPredict` API + local `sqlite3` queries (browser TFLite was dropped; no `tflite_runtime` needed).

3. **Pluggable Fetcher Adapter (Optional Extensibility)**:
   - Mirroring `sdkweb`'s `fetcher` hook, `I6SdkConfig` supports an optional `fetcher` callable.
   - Power users already running `httpx`, `aiohttp`, or `requests` in their applications can supply a custom adapter if desired, without making it a required dependency for everyone.

4. **1-to-1 Conceptual Symmetry with `sdkweb`**:
   - Reuses identical service architecture (`Services`, `Invoker`, `Auth`, `Tokens`, `Apis`, `SdkFuncs`, `Ingest`, `Infer`, `Dash`).
   - Reuses identical API specs (`apispecz.Api`) and schema ASTs (`typeastz.Typedefs`) via `forgepy`.

---

## Architecture Overview

```text
Backend Go APIs (apiuserz, apidatasetz, apiingestz, apipipez, apitfz)
                                  │
                                  ▼
                     apicollectionz.Apis() & SdkFuncs()
                                  │
                                  ▼
                 forgez/forgepy (Go Code Generator)
    ├── forgepyschema  (Maps Go typeastz -> Python typing.TypedDict & Standard Types)
    ├── forgepyapi     (Generates HTTP service methods from apispecz.Api)
    ├── forgepysdkfunc (Generates high-level workflow functions from moldsdkfunc)
    └── forgepybundle  (Bundles & formats output Python package structure)
                                  │
                   [go:generate / build.sh codegen]
                                  ▼
                   sdkpy/src/i6sdk/sdkapis/
    ├── types.py       (Generated request, response, and model TypedDicts)
    ├── apis.py        (Generated API client methods and service classes)
    └── sdkfuncs.py    (Generated high-level functional wrappers)
                                  │
                                  ▼
                 sdkpy Runtime (Zero External Dependencies)
    ├── config.py      (I6SdkConfig, endpoint resolution, env var overrides)
    ├── invoker.py     (HTTP transport via urllib.request, path interpolation, auth injection)
    ├── services.py    (Services container managing clients, lifecycle, and config)
    ├── auth/          (Device code flow, PKCE via hashlib/secrets, token refresh)
    ├── tokens/        (Pluggable token storage: Memory, File ~/.config/i6/tokens.json)
    ├── ingest/        (Ingestion client: batch, streaming, events)
    ├── infer/         (Inference client: sqlite queries + server-side tfzModelPredict)
    ├── dash/          (Dashboard & analytics queries)
    ├── internal/      (Internal utils: downloader, fs cache, sqlitez, bitmapz)
    └── exceptions.py  (I6Error, I6ApiError with status & code matching sdkweb)
```

---

## Phase 1: Code Generator (`forgez/forgepy`)

### 1.1 Project Structure & CLI Scaffolding
- [x] **Directory Setup**: Create package directories:
  - `forgez/forgepy/forgepyschema/`
  - `forgez/forgepy/forgepyapi/`
  - `forgez/forgepy/forgepysdkfunc/`
  - `forgez/forgepy/forgepybundle/`
- [x] **CLI Subcommand**: Register `forgepy` generation commands under `forgez/cmdforgez/` for standalone invocation during development and debugging.

### 1.2 Type & Schema Generator (`forgepyschema`)
- [x] **Go AST Parser Extension**:
  - Implement `forgepyschema.NewSdk() *schemamold.Parser` with known type mappings:
    - `fsz.FileStat` $\to$ Standard Python dictionary or `FileStat(TypedDict)`.
    - `bigquery.NullString` $\to$ `Optional[str]`.
    - `bigquery.NullInt64` $\to$ `Optional[int]`.
    - `bigquery.NullFloat64` $\to$ `Optional[float]`.
    - `time.Time` $\to$ `str` (ISO-8601 string representation).
- [x] **Primitive Type Mapping (`typeastz.Typ` $\to$ Python Standard Types)**:
  - `typeastz.TypString` $\to$ `str`
  - `typeastz.TypInt`, `typeastz.TypInt64`, `typeastz.TypUint` $\to$ `int`
  - `typeastz.TypFloat`, `typeastz.TypFloat64` $\to$ `float`
  - `typeastz.TypBool` $\to$ `bool`
  - `typeastz.TypArray` $\to$ `list[T]`
  - `typeastz.TypMap` $\to$ `dict[str, T]`
  - `typeastz.TypAny` $\to$ `Any`
  - Nullable / Optional fields $\to$ `Optional[T]` (with `None` default)
- [x] **Data Model Code Generation (Zero-Dependency `TypedDict`)**:
  - Generate standard Python `typing.TypedDict` structures (`total=False` for flexible payload compatibility).
  - Include docstrings on schemas based on Go AST field comments.
  - Include safe forward reference handling (`from __future__ import annotations`).
- [x] **Module Typedef Output (`WriteModuleTypedefs`)**:
  - Assemble all extracted types into `types.py`.
  - Include common base types (`ServiceOptions`, `Resp[T]`, `i6JsonEmbedded`).
  - Export `__all__` list with all generated model type names.
- [x] **Unit Tests**:
  - `forgepyschema_unit_test.go`: Verify mapping of primitives, slices, nested structs, maps, pointers, and BigQuery nullables.

### 1.3 API Method Generator (`forgepyapi`)
- [x] **Service & Method Generation**:
  - Inspect `apispecz.Api` list from `apicollectionz.Apis()`.
  - Group endpoints into an `Apis` class matching the `sdkweb` structure.
- [x] **Endpoint Invocation Logic**:
  - Extract route labels (e.g. `servname` for service-specific base URL routing like `tfinfer` and `userz`).
  - Interpolate URL path parameters (e.g. `/api/v1/store/{store_id}` $\to$ dynamically populated from request).
  - Query parameter serialization (filtering `None` values, encoding via standard `urllib.parse`).
  - Header mapping (Content-Type: application/json, custom headers, Authorization bearer).
  - Body payload serialization (via standard `json.dumps`).
- [x] **Method Signatures & Typing**:
  - Produce type-hinted methods using standard library annotations:
    ```python
    def get_dataset(
        self,
        req: DatasetGetReq,
        opts: Optional[ServiceOptions] = None,
    ) -> Resp[DatasetGetResp]:
        """Fetch dataset details.

        :param req: DatasetGetReq request payload.
        :param opts: Optional per-request execution overrides.
        :return: Resp containing parsed DatasetGetResp.
        """
    ```
- [x] **Unit Tests**:
  - `forgepyapi_unit_test.go`: Verify generation of URL formatting, query strings, headers, and docstrings.

### 1.4 High-Level SDK Functions (`forgepysdkfunc`)
- [x] **Parse Functions from `moldsdkfunc`**:
  - Process `moldsdkfunc.SdkFunc` specifications (complex multi-step workflows like dataset warmup, inference, batch ingestion).
  - Generate idiomatic Python wrapper classes and helper methods in `sdkfuncs.py`.
- [x] **Unit Tests**:
  - `forgepysdkfunc_unit_test.go`: Verify generation of higher-level functional flows.

### 1.5 Python Bundler & Formatter (`forgepybundle`)
- [x] **Bundle Assembly**:
  - Combine `forgepyschema`, `forgepyapi`, and `forgepysdkfunc` outputs.
  - Generate clean standard library imports (`from typing import ...`, `from ..invoker import Invoker, Resp, ServiceOptions`).
  - Output files:
    - `sdkapis/__init__.py`: Package re-exports.
    - `sdkapis/types.py`: Generated data schemas.
    - `sdkapis/apis.py`: Generated low-level API endpoints.
    - `sdkapis/sdkfuncs.py`: High-level wrapper workflows.
- [x] **PEP 8 Formatting**:
  - Ensure generated code strictly complies with PEP 8 standards out of the box.
- [x] **Unit Tests**:
  - `forgepybundle_unit_test.go`: End-to-end bundling and generation validation test.

---

## Phase 2: Python SDK Runtime Modules (`sdkpy`)

### 2.1 Package Infrastructure & Tooling
- [x] **Scaffold Go Module**:
  - File: `sdkpy/go.mod`
  - File: `sdkpy/main.go` (containing `//go:generate go run . generate`)
  - File: `sdkpy/cmdsdkpy/cmdsdkpy.go` (Cobra root command)
  - File: `sdkpy/cmdsdkpy/cmdgenerate.go` (Runs `forgepybundle.Write` outputting to `src/i6sdk/sdkapis/`)
- [x] **Python Package Manifest**:
  - File: `sdkpy/pyproject.toml`
  - Configure `hatchling` or `flit_core` build backend.
  - Set `requires-python = ">=3.10"`.
  - Set `dependencies = []` (Zero external runtime dependencies).
  - Configure optional dev groups: `[project.optional-dependencies] dev = ["ruff", "mypy"]`.
- [x] **PEP 561 Type Marker**:
  - File: `sdkpy/src/i6sdk/py.typed`
  - Empty marker file signaling type annotations are bundled for downstream type checkers.

### 2.2 Core Modules

#### `sdkpy/src/i6sdk/exceptions.py`
- [x] **Implement Exception Hierarchy**:
  - `I6Error(Exception)`: Root base exception for all SDK errors.
  - `I6ApiError(I6Error)`: Raised when server returns non-2xx status. Includes:
    - `status: int` (HTTP status code)
    - `code: str` (matches `sdkweb`: `"unauthenticated"`, `"forbidden"`, `"not_found"`, `"unavailable"`, `"bad_request"`)
    - `message: str` (error message or HTTP text)
    - `cause: Any` (raw payload or upstream error cause)
  - `I6AuthError(I6Error)`: Raised for authentication failures, expired tokens, or PKCE state mismatches.
  - `I6ConnectionError(I6Error)`: Raised on network transport issues, DNS failures, or connection timeouts.
  - `I6ValidationError(I6Error)`: Raised on missing or malformed client-side arguments (maps to `TypeError` in `sdkweb`).

#### `sdkpy/src/i6sdk/config.py`
- [x] **Implement Configuration Dataclass & Defaults**:
  - `I6SdkConfig` dataclass:
    - `base_url: str` (defaults to `"https://decsuite.sandbox.rs.infinity6.ai"`)
    - `serv_urls: dict[str, str]` (service overrides, e.g. `{"tfinfer": "...", "userz": "..."}`)
    - `userz_url: Optional[str]` (direct auth service override)
    - `app: str` (client app identifier, default `"default"`)
    - `timeout: float` (HTTP timeout in seconds, default `30.0`)
    - `fetcher: Optional[Callable]` (pluggable custom HTTP handler)
  - `DEFAULT_CONFIG: I6SdkConfig`: Default instance matching sandbox environment.
  - `load_config_from_env()`: Helper parsing `I6_SDK_DECSUITE_URL`, `I6_SDK_TFINFER_URL`, `I6_SDK_USERZ_URL`.

#### `sdkpy/src/i6sdk/invoker.py`
- [x] **Implement Standard Library HTTP Invoker**:
  - `ServiceOptions` dataclass: per-request overrides (`timeout`, `headers`, `service_name`).
  - `Resp[T]` response container (matching `sdkweb`'s `Resp` class):
    - `status: int`, `headers: dict[str, str]`, `ok: bool` (`200 <= status < 300`), `text: str`.
    - Property `error: Optional[dict]`: returns `{status, code, message, cause}` if `!ok`, `None` otherwise.
    - Method `parse() -> T`: returns deserialized JSON data; raises `I6ApiError` if `!ok`.
  - `Invoker` class:
    - Constructor taking `services: Services`.
    - Method `invoke(opts)`: Executes HTTP call via `urllib.request.Request` and `urllib.request.urlopen`.
    - Handles URL path parameter interpolation (`/api/v1/store/{store_id}`).
    - Handles query parameter encoding via `urllib.parse.urlencode`.
    - Injects headers: `User-Agent: i6sdk-python/<version>`, `Content-Type: application/json`, `Authorization: Bearer <token>`.
    - Serializes JSON payloads using standard `json.dumps`.
    - Method `invoke_result(opts) -> Resp[T]`: Wraps execution into `Resp[T]` with lazy parser.
    - Honors pluggable `fetcher` callable if configured in `I6SdkConfig`.

### 2.3 Token Management Modules

#### `sdkpy/src/i6sdk/tokens/storage.py`
- [x] **Implement Pluggable Storage Providers**:
  - `TokenStorage(abc.ABC)`: Abstract interface (`get_item(key)`, `set_item(key, val)`, `remove_item(key)`).
  - `MemoryTokenStorage`: Thread-safe in-memory cache protected by `threading.Lock`.
  - `FileTokenStorage`: File-backed storage persisting to `~/.config/i6/tokens.json` with secure permissions (`0600`) via `os.open` and `pathlib.Path`.

#### `sdkpy/src/i6sdk/tokens/tokens.py`
- [x] **Implement Token Manager Client**:
  - `Tokens` class:
    - Constructor accepting `services: Services`.
    - Access token caching and automatic expiration tracking.
    - Scoped attenuation method (`attenuate(opts)`).
    - Exchange checkout and token redeem helpers (`redeem_checkout(opts)`).

#### `sdkpy/src/i6sdk/tokens/__init__.py`
- [x] **Package Exports**:
  - Re-export `Tokens`, `TokenStorage`, `MemoryTokenStorage`, `FileTokenStorage`.

### 2.4 Authentication Modules

#### `sdkpy/src/i6sdk/auth/auth.py`
- [x] **Implement Authentication Client**:
  - `Auth` class:
    - Constructor accepting `services: Services`.
    - PKCE code challenge generation using standard library `hashlib.sha256` and `base64.urlsafe_b64encode`.
    - Random state & verifier generation using standard library `secrets.token_urlsafe`.
    - Device code authorization flow (`device_code_start()`, `device_code_poll()`).
    - Password sign-in / sign-up methods (`sign_in()`, `sign_up()`, `sign_out()`).
    - Refresh token rotation (`refresh_token()`).
    - User profile query (`me()`).

#### `sdkpy/src/i6sdk/auth/__init__.py`
- [x] **Package Exports**:
  - Re-export `Auth`.

### 2.5 Internal Utilities (`sdkpy/src/i6sdk/internal/`)

#### `sdkpy/src/i6sdk/internal/fs.py`
- [x] **Implement File Cache System**:
  - Local disk cache management using `pathlib.Path`.
  - Managing cached model and SQLite files by version and etag.
  - Resolving file references to local filesystem paths.

#### `sdkpy/src/i6sdk/internal/downloader.py`
- [x] **Implement Asset Downloader**:
  - Asset downloading with conditional HTTP requests (`If-None-Match` / ETag support) using `urllib.request`.
  - Streaming download with progress callback support.

#### `sdkpy/src/i6sdk/internal/sqlitez.py`
- [x] **Implement SQLite Query Helper**:
  - Wrapper around Python's built-in `sqlite3`.
  - Connection management (`open()`, `close()`).
  - Helper `named_string_query(name, params)`: executes named model queries and returns strings/parameters.
  - Helper `named_json_query(name, params)`: executes named output queries and returns parsed JSON data.

#### `sdkpy/src/i6sdk/internal/bitmapz.py`
- [x] **Implement Roaring / Compact Bitmap Logic**:
  - Pure Python port of `gox/commonz/bitmapz` / `sdkweb/src/main/internal/bitmapz/bitmapz.js`.
  - Decoding byte headers, length calculation, and bitmask checks.

#### `sdkpy/src/i6sdk/internal/__init__.py`
- [x] **Package Exports**:
  - Internal module declarations.

### 2.6 Domain Client Modules

#### `sdkpy/src/i6sdk/ingest/ingest.py`
- [x] **Implement Ingest Client**:
  - `Ingest` class:
    - Constructor accepting `services: Services`.
    - Event batch streaming and upload URL retrieval.
    - Multipart file uploads via standard library `urllib.request`.

#### `sdkpy/src/i6sdk/ingest/__init__.py`
- [x] **Package Exports**:
  - Re-export `Ingest`.

#### `sdkpy/src/i6sdk/infer/infer.py`
- [x] **Implement Inference Client**:
  - `Infer` class:
    - Constructor accepting `services: Services`.
    - Locate and download model SQLite files via `internal/downloader.py` and `internal/fs.py`.
    - Run input query via `internal/sqlitez.py` on local SQLite DB.
    - **Call model prediction on server via `services.apis().tfz_model_predict(...)`** (no local TFLite or LiteRT runtime).
    - Run output query via `internal/sqlitez.py` to produce final recommendation/ranking result.
    - Asynchronously emit recommendation telemetry event (`recraw`).

#### `sdkpy/src/i6sdk/infer/__init__.py`
- [x] **Package Exports**:
  - Re-export `Infer`.

#### `sdkpy/src/i6sdk/dash/dash.py`
- [x] **Implement Analytics & Dashboard Client**:
  - `Dash` class:
    - Constructor accepting `services: Services`.
    - Download and query analytics SQLite bundles.

#### `sdkpy/src/i6sdk/dash/__init__.py`
- [x] **Package Exports**:
  - Re-export `Dash`.

### 2.7 Generated SDK API Modules (`sdkpy/src/i6sdk/sdkapis/`)

#### `sdkpy/src/i6sdk/sdkapis/types.py` *(Generated)*
- [x] **Model & TypedDict Definitions**:
  - Generated by `forgepyschema`.
  - Contains TypedDict schemas for all request payloads, response bodies, and entities.

#### `sdkpy/src/i6sdk/sdkapis/apis.py` *(Generated)*
- [x] **Low-Level API Methods**:
  - Generated by `forgepyapi`.
  - `Apis` class containing methods for every endpoint in `apicollectionz`.

#### `sdkpy/src/i6sdk/sdkapis/sdkfuncs.py` *(Generated)*
- [x] **High-Level Workflow Wrappers**:
  - Generated by `forgepysdkfunc`.
  - `SdkFuncs` class containing multi-step workflow helpers.

#### `sdkpy/src/i6sdk/sdkapis/__init__.py`
- [x] **Package Exports**:
  - Re-export `Apis`, `SdkFuncs`, and public model types.

### 2.8 Service Container & Public Package Entry Point

#### `sdkpy/src/i6sdk/services.py`
- [x] **Implement Central Services Container**:
  - `Services` class (mirrors `sdkweb/src/main/services/services.js`):
    - Constructor accepting `config: I6SdkConfig`.
    - Accessors returning singleton instances:
      - `config()` / `set_config(config)`
      - `invoker() -> Invoker`
      - `auth() -> Auth`
      - `tokens() -> Tokens`
      - `ingest() -> Ingest`
      - `infer() -> Infer`
      - `dash() -> Dash`
      - `apis() -> Apis`
      - `sdk_funcs() -> SdkFuncs`

#### `sdkpy/src/i6sdk/__init__.py`
- [x] **Top-Level Package Interface**:
  - Export primary user-facing symbols:
    - `I6Client` (alias for `Services`)
    - `I6SdkConfig`
    - `I6Error`, `I6ApiError`, `I6AuthError`, `I6ConnectionError`, `I6ValidationError`
  - Convenience factory: `create_client(...)`

---

## Phase 3: Integration, Tooling & Testing

### 3.1 Mono-Repo Build Integration
- [x] **`build.sh` Integration**:
  - Add `sdkpy` to component list in `build.sh`.
  - Support `./build.sh test` running both Go tests and Python unittests.
  - Include `sdkpy` in `./build.sh codegen`.
- [x] **`comp.sh` Integration**:
  - Support `./comp.sh sdkpy run generate`.
  - Support `./comp.sh sdkpy test` running both Go tests and Python unittests.
  - Support `./comp.sh sdkpy codegen`.

### 3.2 Testing Strategy
- [x] **Go Unit Tests (`forgez/forgepy/`)**:
  - Ensure 100% test pass on all schema parsers, AST mappings, and generator outputs.
- [x] **Standard Library Unit Tests (`sdkpy/tests/`)**:
  - Using Python standard library `unittest` (runs anywhere without installing extra packages):
    - `test_config.py`: Environment variable loading, priority resolution.
    - `test_invoker.py`: URL building, header attachment, timeout handling, error status parsing with mock HTTP handlers.
    - `test_auth.py`: PKCE code challenge verification, token storage, refresh cycle.
    - `test_infer.py`: SQLite input/output query execution and mock server prediction flow.
    - `test_apis.py`: Smoke test calling generated API methods against mocked endpoints to ensure correct serialization and deserialization.
    - `test_bitmapz.py`: Bitmask manipulation and decomposition.
    - `test_exceptions.py`: Exception class hierarchy and attributes.
    - `test_tokens.py`: FileTokenStorage (0600 permissions), MemoryTokenStorage, Tokens manager.

### 3.3 Documentation & Examples
- [x] **`sdkpy/README.md`**:
  - Highlight: **Zero external runtime dependencies**.
  - Installation instructions (`pip install .`).
  - Quickstart examples (initialization, authentication, querying a dataset, running predictions).
  - Configuration reference.
- [x] **Sample Scripts**:
  - Add `sdkpy/examples/quickstart.py` demonstrating basic workflow using pure Python standard library.
