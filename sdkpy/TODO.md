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

3. **Pluggable Fetcher Adapter (Optional Extensibility)**:
   - Mirring `sdkweb`'s `fetcher` hook, `I6SdkConfig` supports an optional `fetcher` callable.
   - Power users already running `httpx`, `aiohttp`, or `requests` in their applications can supply a custom adapter if desired, without making it a required dependency for everyone.

4. **1-to-1 Conceptual Symmetry with `sdkweb`**:
   - Reuses identical service architecture (`Services`, `Invoker`, `Auth`, `Tokens`, `Apis`, `SdkFuncs`).
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
    └── exceptions.py  (I6Error, I6ApiError with status & code matching sdkweb)
```

---

## Phase 1: Code Generator (`forgez/forgepy`)

### 1.1 Project Structure & CLI Scaffolding
- [ ] **Directory Setup**: Create package directories:
  - `forgez/forgepy/forgepyschema/`
  - `forgez/forgepy/forgepyapi/`
  - `forgez/forgepy/forgepysdkfunc/`
  - `forgez/forgepy/forgepybundle/`
- [ ] **CLI Subcommand**: Register `forgepy` generation commands under `forgez/cmdforgez/` for standalone invocation during development and debugging.

### 1.2 Type & Schema Generator (`forgepyschema`)
- [ ] **Go AST Parser Extension**:
  - Implement `forgepyschema.NewSdk() *schemamold.Parser` with known type mappings:
    - `fsz.FileStat` $\to$ Standard Python dictionary or `FileStat(TypedDict)`.
    - `bigquery.NullString` $\to$ `Optional[str]`.
    - `bigquery.NullInt64` $\to$ `Optional[int]`.
    - `bigquery.NullFloat64` $\to$ `Optional[float]`.
    - `time.Time` $\to$ `str` (ISO-8601 string representation).
- [ ] **Primitive Type Mapping (`typeastz.Typ` $\to$ Python Standard Types)**:
  - `typeastz.TypString` $\to$ `str`
  - `typeastz.TypInt`, `typeastz.TypInt64`, `typeastz.TypUint` $\to$ `int`
  - `typeastz.TypFloat`, `typeastz.TypFloat64` $\to$ `float`
  - `typeastz.TypBool` $\to$ `bool`
  - `typeastz.TypArray` $\to$ `list[T]`
  - `typeastz.TypMap` $\to$ `dict[str, T]`
  - `typeastz.TypAny` $\to$ `Any`
  - Nullable / Optional fields $\to$ `Optional[T]` (with `None` default)
- [ ] **Data Model Code Generation (Zero-Dependency `TypedDict`)**:
  - Generate standard Python `typing.TypedDict` structures (`total=False` for flexible payload compatibility).
  - Include docstrings on schemas based on Go AST field comments.
  - Include safe forward reference handling (`from __future__ import annotations`).
- [ ] **Module Typedef Output (`WriteModuleTypedefs`)**:
  - Assemble all extracted types into `types.py`.
  - Include common base types (`ServiceOptions`, `Resp[T]`, `i6JsonEmbedded`).
  - Export `__all__` list with all generated model type names.
- [ ] **Unit Tests**:
  - `forgepyschema_unit_test.go`: Verify mapping of primitives, slices, nested structs, maps, pointers, and BigQuery nullables.

### 1.3 API Method Generator (`forgepyapi`)
- [ ] **Service & Method Generation**:
  - Inspect `apispecz.Api` list from `apicollectionz.Apis()`.
  - Group endpoints into an `Apis` class matching the `sdkweb` structure.
- [ ] **Endpoint Invocation Logic**:
  - Extract route labels (e.g. `servname` for service-specific base URL routing like `tfinfer` and `userz`).
  - Interpolate URL path parameters (e.g. `/api/v1/store/{store_id}` $\to$ dynamically populated from request).
  - Query parameter serialization (filtering `None` values, encoding via standard `urllib.parse`).
  - Header mapping (Content-Type: application/json, custom headers, Authorization bearer).
  - Body payload serialization (via standard `json.dumps`).
- [ ] **Method Signatures & Typing**:
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
- [ ] **Unit Tests**:
  - `forgepyapi_unit_test.go`: Verify generation of URL formatting, query strings, headers, and docstrings.

### 1.4 High-Level SDK Functions (`forgepysdkfunc`)
- [ ] **Parse Functions from `moldsdkfunc`**:
  - Process `moldsdkfunc.SdkFunc` specifications (complex multi-step workflows like dataset warmup, inference, batch ingestion).
  - Generate idiomatic Python wrapper classes and helper methods in `sdkfuncs.py`.
- [ ] **Unit Tests**:
  - `forgepysdkfunc_unit_test.go`: Verify generation of higher-level functional flows.

### 1.5 Python Bundler & Formatter (`forgepybundle`)
- [ ] **Bundle Assembly**:
  - Combine `forgepyschema`, `forgepyapi`, and `forgepysdkfunc` outputs.
  - Generate clean standard library imports (`from typing import ...`, `from ..invoker import Invoker, Resp, ServiceOptions`).
  - Output files:
    - `sdkapis/__init__.py`: Package re-exports.
    - `sdkapis/types.py`: Generated data schemas.
    - `sdkapis/apis.py`: Generated low-level API endpoints.
    - `sdkapis/sdkfuncs.py`: High-level wrapper workflows.
- [ ] **PEP 8 Formatting**:
  - Ensure generated code strictly complies with PEP 8 standards out of the box.
- [ ] **Unit Tests**:
  - `forgepybundle_unit_test.go`: End-to-end bundling and generation validation test.

---

## Phase 2: Python SDK Runtime (`sdkpy`)

### 2.1 Package Setup & Tooling Configuration
- [ ] **Module & Directory Layout**:
  - Create `sdkpy/` Go module (`go.mod`, `main.go`, `cmdsdkpy/`).
  - Create Python source root under `sdkpy/src/i6sdk/`.
- [ ] **Zero-Dependency `pyproject.toml`**:
  - Build backend: `hatchling` or `flit_core` (standard, lightweight PEP 517 backends).
  - Python version target: `>= 3.10`.
  - **Runtime dependencies**:
    ```toml
    dependencies = []
    ```
  - Optional development dependency groups:
    ```toml
    [project.optional-dependencies]
    dev = ["ruff", "mypy"]
    test = ["pytest"]
    ```
- [ ] **Type Marker**:
  - Add `src/i6sdk/py.typed` to signal PEP 561 compliance to type checkers.

### 2.2 Client Core Runtime (`sdkpy/src/i6sdk/`)
- [ ] **Configuration (`config.py`)**:
  - `I6SdkConfig` standard dataclass:
    - `base_url`: Default API gateway URL (`https://decsuite.sandbox.rs.infinity6.ai`).
    - `serv_urls`: Dict mapping service names (e.g. `tfinfer`, `userz`) to service-specific host URLs.
    - `app`: Application identifier.
    - `token`: Static bearer token / API key.
    - `timeout`: Default request timeout in seconds (default `30.0`).
    - `fetcher`: Optional custom HTTP callable `(url, method, headers, body, timeout) -> (status, headers, body)`.
  - Auto-configuration from environment variables:
    - `I6_BASE_URL`, `I6_API_KEY`, `I6_APP`, `I6_TIMEOUT`.
- [ ] **Exceptions (`exceptions.py`)**:
  - `I6Error(Exception)`: Base SDK exception.
  - `I6ApiError(I6Error)`: Raised when server returns $\ge 400$. Attributes:
    - `status: int`
    - `code: str` (matches `sdkweb`: `"unauthenticated"`, `"forbidden"`, `"not_found"`, `"unavailable"`, `"bad_request"`)
    - `message: str`
    - `cause: Any`
  - `I6AuthError(I6Error)`: Authentication / PKCE / token errors.
  - `I6ConnectionError(I6Error)`: Network / timeout failures wrapping `urllib.error.URLError`.
  - `I6ValidationError(I6Error)`: Missing or invalid required parameters.
- [ ] **Standard Library HTTP Invoker (`invoker.py`)**:
  - `Invoker` class using standard `urllib.request` and `http.client`:
    - Handles Keep-Alive connections and SSL verification.
    - URL path template interpolation and query parameter encoding via `urllib.parse.urlencode`.
    - Header injection (User-Agent: `i6sdk-python/<version>`, Authorization: Bearer, Content-Type: application/json).
    - Request payload JSON serialization via standard `json.dumps`.
  - `Resp[T]` container matching `sdkweb`'s `Resp`:
    - `status: int`
    - `headers: dict[str, str]`
    - `ok: bool` (`200 <= status < 300`)
    - `text: str`
    - `error: Optional[dict]` (structured error object if `!ok`)
    - `parse() -> T`: Returns parsed JSON data; raises `I6ApiError` if `!ok`.
  - Support pluggable `fetcher` hook for users who supply custom async or third-party transports.

### 2.3 Authentication & Token Management (`sdkpy/src/i6sdk/auth/` & `tokens/`)
- [ ] **Token Storage (`tokens/storage.py`)**:
  - `TokenStorage` abstract base class (`abc.ABC`): `get(key)`, `set(key, val)`, `delete(key)`.
  - `MemoryTokenStorage`: In-memory thread-safe dictionary using `threading.Lock`.
  - `FileTokenStorage`: File-backed storage in `~/.config/i6/tokens.json` with restricted permissions (`0600`) using `os.open` and `pathlib.Path`.
- [ ] **Auth Client (`auth/auth.py`)**:
  - PKCE OAuth code challenge calculation using standard library `hashlib.sha256` and `base64.urlsafe_b64encode`.
  - Random state/verifier generation using standard library `secrets.token_urlsafe`.
  - Device code authentication flow (`userapiauthdevicecode` & `userapiauthdevicetoken`).
  - Refresh token rotation (`userapiauthrefresh`).
  - Token attenuation / scoped delegation (`userapidatasettokenaccessattenuate`).

### 2.4 Service Container (`services.py`) & Top-Level SDK API
- [ ] **Service Orchestrator (`services.py`)**:
  - `Services` container matching the TypeScript `Services` pattern:
    - Holds reference to shared `I6SdkConfig` and `Invoker`.
    - Lazy-initializes service clients (`auth`, `tokens`, `ingest`, `infer`, `dash`, `apis`).
- [ ] **Top-Level Package Interface (`src/i6sdk/__init__.py`)**:
  - Export primary user-facing entry points:
    ```python
    from i6sdk.services import Services as I6Client
    from i6sdk.config import I6SdkConfig
    from i6sdk.exceptions import I6Error, I6ApiError, I6AuthError
    ```
  - Convenience factory: `client = i6sdk.create_client(api_key="...", ...)`

### 2.5 Generator CLI Command (`sdkpy/cmdsdkpy/`)
- [ ] **Generator Runner (`cmdsdkpy/cmdgenerate.go`)**:
  - Implement Cobra CLI command matching `sdkweb/cmdsdkweb/cmdgenerate.go`.
  - Load all registered APIs (`apicollectionz.Apis()`) and SDK functions (`apicollectionz.SdkFuncs()`).
  - Invoke `forgepybundle.Write(ctx, opts)`.
  - Write output files to `sdkpy/src/i6sdk/sdkapis/`.
- [ ] **`go:generate` Directive**:
  - Place `//go:generate go run . generate` inside `sdkpy/main.go`.

---

## Phase 3: Integration, Tooling & Testing

### 3.1 Mono-Repo Build Integration
- [ ] **`build.sh` Integration**:
  - Add `sdkpy` to component list in `build.sh`.
  - Add `cmd_test_py()` in `build.sh` (running `unittest` or `pytest` under `sdkpy/`).
  - Include `sdkpy` in `./build.sh codegen`.
- [ ] **`comp.sh` Integration**:
  - Support `./comp.sh sdkpy run generate`.
  - Support `./comp.sh sdkpy test`.

### 3.2 Testing Strategy
- [ ] **Go Unit Tests (`forgez/forgepy/`)**:
  - Ensure 100% test pass on all schema parsers, AST mappings, and generator outputs.
- [ ] **Standard Library Unit Tests (`sdkpy/tests/`)**:
  - Using Python standard library `unittest` (runs anywhere without installing extra packages):
    - `test_config.py`: Environment variable loading, priority resolution.
    - `test_invoker.py`: URL building, header attachment, timeout handling, error status parsing with mock HTTP handlers.
    - `test_auth.py`: PKCE code challenge verification, token storage, refresh cycle.
    - `test_apis.py`: Smoke test calling generated API methods against mocked endpoints to ensure correct serialization and deserialization.
- [ ] **Python Static Analysis (Dev Tools)**:
  - `ruff check sdkpy/src/` (linting).
  - `ruff format --check sdkpy/src/` (formatting).
  - `mypy --strict sdkpy/src/` (strict static typing).

### 3.3 Documentation & Examples
- [ ] **`sdkpy/README.md`**:
  - Highlight: **Zero external dependencies**.
  - Installation instructions (`pip install .`).
  - Quickstart examples (initialization, authentication, querying a dataset, running predictions).
  - Configuration reference.
- [ ] **Sample Scripts**:
  - Add `sdkpy/examples/quickstart.py` demonstrating basic workflow using pure Python standard library.
