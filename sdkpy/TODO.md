# Python SDK (`sdkpy`) Roadmap & Detailed Task Breakdown

This document provides a comprehensive technical blueprint and task tracker for implementing the Python SDK (`sdkpy`) and its AST-based code generator (`forgez/forgepy`).

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
    ├── forgepyschema  (Maps Go typeastz -> Python types / TypedDicts / Dataclasses)
    ├── forgepyapi     (Generates HTTP service methods from apispecz.Api)
    ├── forgepysdkfunc (Generates high-level workflow functions from moldsdkfunc)
    └── forgepybundle  (Bundles & formats output Python package structure)
                                  │
                   [go:generate / build.sh codegen]
                                  ▼
                   sdkpy/src/i6sdk/sdkapis/
    ├── types.py       (Generated request, response, and model type definitions)
    ├── apis.py        (Generated API client methods and service classes)
    └── sdkfuncs.py    (Generated high-level functional wrappers)
                                  │
                                  ▼
                 sdkpy Runtime (Handcrafted Core)
    ├── config.py      (I6SdkConfig, endpoint resolution, env var overrides)
    ├── invoker.py     (HTTP transport, path interpolation, auth injection, error wrapping)
    ├── services.py    (Services container managing clients, lifecycle, and config)
    ├── auth/          (Device code flow, PKCE, token refresh, attenuation)
    ├── tokens/        (Pluggable token storage: Memory, File, Keyring)
    └── exceptions.py  (Structured SDK exception hierarchy)
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
    - `fsz.FileStat` $\to$ File metadata model (`FileStat`).
    - `bigquery.NullString` $\to$ `Optional[str]`.
    - `bigquery.NullInt64` $\to$ `Optional[int]`.
    - `bigquery.NullFloat64` $\to$ `Optional[float]`.
    - `time.Time` $\to$ `datetime` / ISO-8601 `str`.
- [ ] **Primitive Type Mapping (`typeastz.Typ` $\to$ Python)**:
  - `typeastz.TypString` $\to$ `str`
  - `typeastz.TypInt`, `typeastz.TypInt64`, `typeastz.TypUint` $\to$ `int`
  - `typeastz.TypFloat`, `typeastz.TypFloat64` $\to$ `float`
  - `typeastz.TypBool` $\to$ `bool`
  - `typeastz.TypArray` $\to$ `list[T]`
  - `typeastz.TypMap` $\to$ `dict[str, T]`
  - `typeastz.TypAny` $\to$ `Any`
  - Nullable / Optional fields $\to$ `Optional[T]` (with `None` default)
- [ ] **Data Model Code Generation**:
  - Generate Python `TypedDict` structures (with `total=False` for partial/optional payloads or standard `TypedDict` with `NotRequired[T]`).
  - Support alternative mode for `dataclasses` or `pydantic.BaseModel` if configured.
  - Preserve documentation comments from Go AST docstrings as Python docstrings or field comments.
  - Handle circular / forward references safely using `from __future__ import annotations`.
- [ ] **Module Typedef Output (`WriteModuleTypedefs`)**:
  - Assemble all extracted types into `types.py`.
  - Include common base types (`ServiceOptions`, `Resp[T]`, `i6JsonEmbedded`).
  - Export `__all__` list with all generated model type names.
- [ ] **Unit Tests**:
  - `forgepyschema_unit_test.go`: Verify mapping of primitives, slices, nested structs, maps, pointers, and enum types.

### 1.3 API Method Generator (`forgepyapi`)
- [ ] **Service & Method Generation**:
  - Inspect `apispecz.Api` list from `apicollectionz.Apis()`.
  - Group endpoints either under a single `Apis` class or modular service classes (`UserApis`, `DatasetApis`, `IngestApis`, `PipeApis`, `TfApis`).
- [ ] **Endpoint Invocation Logic**:
  - Extract route labels (e.g. `servname` for service base URL routing).
  - Interpolate URL path parameters (e.g. `/api/v1/store/{store_id}` $\to$ dynamically populated from request).
  - Query parameter serialization (filtering `None` values, formatting arrays).
  - Header mapping (content-type, custom headers, authorization).
  - Body payload serialization (JSON encoding).
- [ ] **Method Signatures & Typing**:
  - Produce type-hinted methods:
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
  - Provide both synchronous and asynchronous invocation variants (e.g. `get_dataset(...)` and `aget_dataset(...)`, or separate `AsyncApis` client).
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
  - Generate clean import headers (`from typing import ...`, `from ..invoker import Invoker, Resp, ServiceOptions`).
  - Output files:
    - `sdkapis/__init__.py`: Package re-exports.
    - `sdkapis/types.py`: Generated data schemas.
    - `sdkapis/apis.py`: Generated low-level API endpoints.
    - `sdkapis/sdkfuncs.py`: High-level wrapper workflows.
- [ ] **Formatting Integration**:
  - Ensure generated code adheres to PEP 8 standards.
  - Add optional hook to run `ruff format` or `black` if installed in the environment.
- [ ] **Unit Tests**:
  - `forgepybundle_unit_test.go`: End-to-end bundling and generation validation test.

---

## Phase 2: Python SDK Runtime (`sdkpy`)

### 2.1 Package Setup & Tooling Configuration
- [ ] **Module & Directory Layout**:
  - Create `sdkpy/` Go module (`go.mod`, `main.go`, `cmdsdkpy/`).
  - Create Python source root under `sdkpy/src/i6sdk/`.
- [ ] **Modern Python Packaging (`pyproject.toml`)**:
  - Build backend: `hatchling` (or `flit_core`).
  - Python version target: `>= 3.10`.
  - Dependencies:
    - `httpx >= 0.25.0` (modern HTTP client supporting HTTP/1.1, HTTP/2, connection pooling, sync and async).
    - Optional: `pydantic >= 2.0` (if Pydantic models are enabled).
  - Optional dependency groups (`dev`, `test`):
    - `pytest`, `pytest-asyncio`, `respx` (HTTP mocking).
    - `ruff` (linter and formatter).
    - `mypy` (strict static type checker).
- [ ] **Type Marker**:
  - Add `src/i6sdk/py.typed` to signal PEP 561 compliance to type checkers.

### 2.2 Client Core Runtime (`sdkpy/src/i6sdk/`)
- [ ] **Configuration (`config.py`)**:
  - `I6SdkConfig` dataclass:
    - `base_url`: Default API gateway URL (e.g. `https://decsuite.sandbox.rs.infinity6.ai`).
    - `serv_urls`: Dict mapping service names (e.g. `tfinfer`, `userz`) to service-specific host URLs.
    - `app`: Application identifier.
    - `token`: Static bearer token / API key.
    - `timeout`: Default request timeout in seconds (default `30.0`).
    - `max_retries`: Automatic retry attempts for idempotent endpoints on network blips.
  - Auto-configuration from environment variables:
    - `I6_BASE_URL`, `I6_API_KEY`, `I6_APP`, `I6_TIMEOUT`.
- [ ] **Exceptions (`exceptions.py`)**:
  - `I6Error` (base SDK exception).
  - `I6ApiError` (server returned non-2xx; includes `status_code`, `error_code`, `details`, `request_id`).
  - `I6AuthError` (authentication failed or token expired).
  - `I6ConnectionError` (network failures, timeouts).
  - `I6ValidationError` (invalid client-side parameters).
- [ ] **HTTP Invoker (`invoker.py`)**:
  - `Invoker` class encapsulating `httpx.Client` (sync) and `httpx.AsyncClient` (async).
  - Request preparation:
    - Base URL and path joining (handling leading/trailing slashes).
    - Header injection (User-Agent, Authorization: Bearer, Content-Type: application/json).
    - Query parameter encoding.
  - Response handling:
    - Parse JSON response or raw bytes based on content type.
    - Wrap into `Resp[T]` container (status code, headers, data, raw response).
    - Raise typed `I6ApiError` on HTTP $\ge$ 400.
  - Connection lifecycle:
    - Context manager support (`with I6Client() as client:`, `async with AsyncI6Client() as client:`).
    - Proper closing of underlying connection pools.

### 2.3 Authentication & Token Management (`sdkpy/src/i6sdk/auth/` & `tokens/`)
- [ ] **Token Storage Provider (`tokens/storage.py`)**:
  - `TokenStorage` abstract base class: `get(key)`, `set(key, val)`, `delete(key)`.
  - `MemoryTokenStorage`: Default in-memory cache.
  - `FileTokenStorage`: Persistent storage in `~/.config/i6/tokens.json` with secure file permissions (`0600`).
- [ ] **Auth Client (`auth/auth.py`)**:
  - Device code flow (`userapiauthdevicecode` & `userapiauthdevicetoken`).
  - Refresh token rotation (`userapiauthrefresh`).
  - Token attenuation / scoped delegation (`userapidatasettokenaccessattenuate`).
  - Automatic token refresh interceptor inside `Invoker`.

### 2.4 Service Container (`services.py`) & Top-Level SDK API
- [ ] **Service Orchestrator (`services.py`)**:
  - `Services` container matching the TypeScript `Services` pattern:
    - Holds reference to shared `I6SdkConfig` and `Invoker`.
    - Lazy-initializes service clients (`auth`, `tokens`, `ingest`, `infer`, `dash`, `apis`).
- [ ] **Top-Level Package Interface (`src/i6sdk/__init__.py`)**:
  - Export primary user-facing entry points:
    ```python
    from i6sdk import I6Client, AsyncI6Client, I6SdkConfig
    from i6sdk.exceptions import I6Error, I6ApiError
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
  - Add `cmd_test_py()` to run `pytest` under `sdkpy/`.
  - Include `sdkpy` in `./build.sh codegen`.
- [ ] **`comp.sh` Integration**:
  - Support `./comp.sh sdkpy run generate`.
  - Support `./comp.sh sdkpy test`.

### 3.2 Testing Strategy
- [ ] **Go Unit Tests (`forgez/forgepy/`)**:
  - Ensure 100% test pass on all schema parsers, AST mappings, and generator outputs.
- [ ] **Python Syntax & Type Validation**:
  - `ruff check sdkpy/src/` (linting).
  - `ruff format --check sdkpy/src/` (formatting).
  - `mypy --strict sdkpy/src/` (strict type verification).
- [ ] **Python Unit & Integration Tests (`sdkpy/tests/`)**:
  - `test_config.py`: Environment variable loading, priority resolution.
  - `test_invoker.py`: URL building, header attachment, timeout handling, error status parsing with mock HTTP responses (`respx`).
  - `test_auth.py`: Token storage, refresh cycle, device code polling simulation.
  - `test_apis.py`: Smoke test calling generated API methods against mocked endpoints to ensure correct serialization and deserialization.
  - `test_async.py`: Ensure async clients function symmetrically with sync clients.

### 3.3 Documentation & Examples
- [ ] **`sdkpy/README.md`**:
  - Installation instructions (`pip install .`).
  - Quickstart examples (initialization, authentication, querying a dataset, running predictions).
  - Configuration reference.
- [ ] **Sample Scripts**:
  - Add `sdkpy/examples/quickstart.py` demonstrating basic workflow.
