# Proof of Concept: Explicit Envelope `ApiDataRefs` Generation & Manual Go SDK

## Executive Summary

This POC proposes a two-part architectural evolution for the Infinity6 SDK mono-repo:
1. **Automate `ApiDataRefs` via Explicit Envelopes**: Standardize all API definitions around explicit request/response sub-struct envelopes (`Req{Params, Query, Headers, Payload}` and `Resp{Headers, Payload}`) with standard `json:"..."` tags. An automated generator (`schemagen`) reads standard Go doc comments (`// ...`) and generates `ApiDataRefs()`.
2. **Build `sdkgo` Manually**: Rather than attempting full AST code generation for the Go client SDK (`forgego`), implement `sdkgo` manually with idiomatic Go ergonomics, thread-safe token management, streaming file ingestion, and local SQLite inference, consuming the generated API definitions.

---

## 1. Problem Statement & Motivation

### Current State
Every endpoint in `apiuserz/`, `apidatasetz/`, `apiingestz/`, `apipipez/`, and `apitfz/` must handwrite `ApiDataRefs()`:

```go
func (d *Api) ApiDataRefs() *apiz.DataRefs {
    return &apiz.DataRefs{
        PathParams: &schemaz.Schema{
            Object: func(read bool) map[string]*schemaz.Schema {
                if read && d.Req == nil { return nil }
                if d.Req == nil { d.Req = &DszListDatasetsReq{} }
                return map[string]*schemaz.Schema{
                    "userId": {
                        Desc:     &schemaz.Desc{Summary: "User identifier (id)."},
                        UrlValue: schemaz.ParseStr(&d.Req.UserId),
                    },
                }
            },
        },
        RespBody: &schemaz.Schema{
            Object: func(read bool) map[string]*schemaz.Schema {
                if read && d.Resp == nil { return nil }
                if d.Resp == nil { d.Resp = &DszListDatasetsResp{} }
                return map[string]*schemaz.Schema{
                    "datasets": {
                        Desc: &schemaz.Desc{Summary: "List of datasets available to the user."},
                        Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Dataset, *Dataset]{
                            Slice: func() *[]*Dataset { return &d.Resp.Datasets },
                            // ... nested mapping ...
                        }),
                    },
                }
            },
        },
    }
}
```

### Pain Points
- **Boilerplate overhead**: `ApiDataRefs()` accounts for ~70–80% of all code in API packages.
- **Documentation duplication**: Descriptions written inside `schemaz.Desc{Summary: "..."}` are invisible to Go developer tooling (`gopls`, hover documentation, godoc).
- **Fragile pointer binding**: Manual closure mechanics (`read bool`, pointer dereferencing, `schemaz.NewArray`) are error-prone and tedious to maintain.
- **Flattened Request Ambiguity**: Currently, many `Req` structs lump path parameters, headers, and body properties into a single flat struct, risking field name collisions and payload pollution during JSON marshaling.

---

## 2. Proposed Architecture

```
┌────────────────────────────────────────────────────────┐
│  Go Source File (e.g. dszapidatasetlist.go)            │
│  - Clean struct definitions                            │
│  - Go doc comments (// ...)                            │
│  - Standard Req{Params, Query, Headers, Payload}       │
│  - Standard Resp{Headers, Payload}                     │
│  - Pure `json:"..."` tags (no custom DSL)              │
└───────────────────────────┬────────────────────────────┘
                            │
              `schemagen` (go:generate)
              - go/parser + go/doc
              - 1:1 Envelope mapping
                            ▼
┌────────────────────────────────────────────────────────┐
│  Generated File (e.g. dszapidatasetlist.gen.go)        │
│  - func (d *Api) ApiDataRefs() *apiz.DataRefs          │
└───────────────────────────┬────────────────────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        ▼                                       ▼
┌───────────────────────────┐       ┌───────────────────────────┐
│ Existing Generators       │       │ New Manual sdkgo Client   │
│ - apicollectionz          │       │ - I6Client / Config       │
│ - forgejs  -> sdkweb      │       │ - Tokens (auto-refresh)   │
│ - forgepy  -> sdkpy       │       │ - Ingest / Infer / Dash   │
│ (Zero breaking changes!)  │       │ - client.Do(ctx, api)     │
└───────────────────────────┘       └───────────────────────────┘
```

---

## 3. Explicit Envelope Architecture

### A. The Standard Request & Response Envelopes

Every API defines its HTTP contract using explicit sub-structs corresponding directly to HTTP transport slots.

**Key Rule: Unused slots are simply omitted from the struct.** If an endpoint has no body (e.g. `GET`), no query string, or no headers, omit the field completely.

```go
package dszapidatasetlist

import "github.com/infinity6-ai/i6sdk/sdkcorez/grants"

// Params contains path parameters embedded in the route URL.
type Params struct {
    // Unique user identifier.
    UserId string `json:"userId"` 
}

// Dataset represents an Infinity6 dataset container and permissions.
type Dataset struct {
    // Dataset unique identifier.
    Id string `json:"id"`

    // Display name of the dataset.
    Name string `json:"name"`

    // Permission bitmask for the current caller.
    Grants grants.Grant `json:"grants"`
}

// Payload represents the JSON HTTP response payload.
type Payload struct {
    // Datasets accessible to the requester.
    Datasets []*Dataset `json:"datasets"`
}

// Req encapsulates all incoming HTTP data for the endpoint.
// Unused slots (Query, Headers, Payload) are omitted.
type Req struct {
    Params *Params `json:"params,omitempty"`
}

// Resp encapsulates all outgoing HTTP data for the endpoint.
// Unused slots (Headers) are omitted.
type Resp struct {
    Payload *Payload `json:"payload"`
}

// List Datasets
// Lists all datasets available to the user.
type Api struct {
    apiz.UnimplementedApi `json:"-"`
    Req                   *Req  `json:"req"`
    Resp                  *Resp `json:"resp"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
    return apiz.ApiSpec{
        Id:     "dsz-list-datasets",
        Method: "GET",
        Path:   "/api/ds/user-id/{userId}/datasets",
        Desc: func() *schemaz.Desc {
            return &schemaz.Desc{
                Labels: map[string]string{
                    speclabeltokentype.Key: speclabeltokentype.User,
                },
            }
        },
    }
}
```

Both `ApiSpec()` and `ApiDataRefs()` are generated into `<package>.gen.go`:
- `ApiSpec()` calls `ApiRoute() apiz.ApiSpec` and **enriches only non-already informed values**:
  - `Name` and `Summary` are populated from the doc comments above `type Api struct` only if they were not explicitly defined in `ApiRoute()`.
  - All existing `Labels`, custom IDs, and route configurations in `ApiRoute()` are preserved intact.
  - Zero generation-specific types are required in `gox` (standard `apiz.ApiSpec` is used directly).
- `ApiDataRefs()` is generated from the `{Params, Query, Headers, Payload}` envelope sub-structs.
- `apiz.UnimplementedApi` satisfies the `apiz.Api` interface at compile time prior to running `go generate`, panicking with an actionable message if invoked before generation.

---

### B. Why Explicit Envelopes Are Superior

1. **Zero Custom Tag DSL**:
   - No need to invent, parse, or maintain custom tags like `apiz:"param,..."` or `apiz:"header,..."`.
   - Standard `json:"..."` tags handle naming and serialization natively.
2. **100% Multi-Language Parity (`Payload`)**:
   - `sdkweb` (`forgejs`) and `sdkpy` (`forgepy`) already generate types with `{ params, query, headers, payload }`:
     - Python: `req = {"params": {...}, "query": {...}, "payload": {...}}`
     - TypeScript: `req = { params: {...}, query: {...}, payload: {...} }`
   - Using `Payload` in Go brings all three SDKs into complete cross-language naming and structural symmetry. (Note: `schemagen` can also treat `Body` as an alias for `Payload` if preferred).
3. **Trivial Generator Implementation**:
   Mapping to `apiz.DataRefs` is 1:1 and deterministic:
   - `Req.Params`   $\rightarrow$ `DataRefs.PathParams`
   - `Req.Query`    $\rightarrow$ `DataRefs.QueryParams`
   - `Req.Headers`  $\rightarrow$ `DataRefs.ReqHeaders`
   - `Req.Payload`  $\rightarrow$ `DataRefs.ReqBody`
   - `Resp.Headers` $\rightarrow$ `DataRefs.RespHeaders`
   - `Resp.Payload` $\rightarrow$ `DataRefs.RespBody`
4. **Clean Serialization & No Payload Leaks**:
   - `json.Marshal(req.Payload)` marshals *strictly* the body payload. Path parameters and headers can never accidentally leak into the HTTP body.
5. **No Parameter Collisions**:
   - Eliminates naming ambiguities when a path parameter and a body field share the same identifier (e.g. `/datasets/{id}` with `{ "id": ... }` in the payload).

---

### C. `schemagen` Generator Engine

The generator runs via `go generate ./...` or `./build.sh codegen`:

1. **Parse Go AST & Doc Comments**:
   Uses standard library `go/parser` (with `parser.ParseComments`) and `go/doc` to extract doc comments preceding each type and field.
2. **Direct Envelope Mapping**:
   Inspects `Req` and `Resp` fields:
   - Evaluates populated envelope fields (`Params`, `Query`, `Headers`, `Payload`). Unpopulated or omitted slots automatically map to `nil` in `apiz.DataRefs`.
   - Recursively walks child struct fields to generate standard `schemaz.Schema` closures and pointer bindings.
   - Maps field doc comments directly into `schemaz.Desc{Summary: docText}`.
3. **Domain Types & `checked.Value[T]` Resolution**:
   Infinity6 APIs frequently use domain-validated types (e.g. `lakedsz.Dataset`, `solution.Solution`, `grants.Grant`):
   - **`checked.Value[T]` / `JSType()`**: Types implementing `JSType() string` or embedding `checked.Value[T]` are recognized automatically; `schemagen` targets their underlying primitive schema (`TypString`, `TypInt`, etc.) and binds their `MustParse` / `Get()` accessors.
   - **Primitive Aliases**: Types like `type Grant uint8` are mapped directly to their primitive base types with full enum comment extraction.
4. **Output**:
   Emits `<package>.gen.go` implementing `func (d *Api) ApiDataRefs() *apiz.DataRefs`.

---

## 4. Manual `sdkgo` Client Implementation

With `ApiDataRefs()` generated, `sdkgo` is implemented manually to provide an idiomatic, high-performance Go experience:

### A. Core Architecture

```
sdkgo/
├── client.go           # I6Client facade and constructor (New, LoadConfigFromEnv)
├── config.go           # Config, Server URLs, timeout, HTTP client
├── invoker/
│   ├── invoker.go      # Dispatches calls via gox/routez/apiclientz or net/http
│   └── errors.go       # I6ApiError, I6AuthError, I6ConnectionError
├── tokens/
│   ├── tokens.go       # Double-checked locking cache with silent token auto-renewal
│   └── storage.go      # In-memory and file token persistence
├── ingest/
│   └── ingest.go       # Multi-stream chunked file uploads with progress & cancellation
├── infer/
│   └── infer.go        # SQLite model query execution (modernc.org/sqlite or native) + TF prediction
└── dash/
    └── dash.go         # Dashboard SQLite execution runner
```

### B. Calling Ergonomics

`sdkgo` can provide both generic invocation and domain service facades:

#### 1. Generic Invocation (`client.Do`)
```go
api := &dszapidatasetlist.Api{
    Req: &dszapidatasetlist.Req{
        Params: &dszapidatasetlist.Params{UserId: "usr_123"},
    },
}

status, err := client.Do(ctx, api)
if err != nil {
    log.Fatalf("API error: %v", err)
}

fmt.Printf("Found %d datasets\n", len(api.Resp.Payload.Datasets))
```

#### 2. Service Facades (`client.Datasets()`, `client.Auth()`, etc.)
```go
resp, err := client.Datasets().List(ctx, &dszapidatasetlist.Req{
    Params: &dszapidatasetlist.Params{UserId: "usr_123"},
})
```

#### 3. Universal Rule for Generic `[T]` APIs & Backend Bridge Pattern (Option A / Path 2B)
To ensure 100% parity across TypeScript, Python, and Go SDKs, and full compatibility with `schemagen` and `apisgen`, **no generic `type Api[T ...]` structs remain in the public SDK surface**. Every polymorphic API family is resolved using:
1. **Concrete Envelope Packages (`i6sdk`)**: Each concrete instantiation (per table, pipe, or result) receives its own dedicated package with explicit envelopes, godoc, and `<MethodName>Id` stub struct.
2. **Backend Bridge Interface (`../*`)**: Sibling backend service runners remain generic `Service[A BridgeApi]`, adapting concrete APIs via a lightweight bridge interface with `newT[A]()` dynamic pointer allocation.

```go
// Example: In backend runner (i6pipez)
type PipeStarterApi interface {
    apiz.Api
    Dataset() lakedsz.Dataset
    Pipe() pipez.Pipe
    SetRespId(id string)
}

type PipezStartService[A PipeStarterApi] struct {
    api A
}
```

This pattern applies universally across all 5 polymorphic API families:
- `pipezstart.Api[T]` $\rightarrow$ `pipezstartbatchingest` & `pipezstartmain` (Bridge: `PipeStarterApi`) `[DONE]`
- `ingestzapievent.Api[T]` $\rightarrow$ `ingestzapieventcatalog` & `ingestzapieventevent` (Bridge: `EventIngestApi`) `[DONE]`
- `ingestzapistream.Api[T]` $\rightarrow$ `ingestzapistreamcatalog` & `ingestzapistreamevent` (Bridge: `StreamIngestApi`) `[DONE]`
- `scheduleractivate.Api[T]` $\rightarrow$ `scheduleractivatebatchingest` & `scheduleractivatemain` (Bridge: `SchedulerActivatorApi`) `[STEP 3]`
- `schedulerdeactivate.Api[T]` $\rightarrow$ `schedulerdeactivatebatchingest` & `schedulerdeactivatemain` (Bridge: `SchedulerDeactivatorApi`) `[STEP 3]`
- `dszapiloadresult.Api[T]` $\rightarrow$ `dszapiloadresultcatalog` & `dszapiloadresultbatchsummary` (Bridge: `LoadResultApi`) `[STEP 5]`

#### 4. Clean Checkouts & SdkApis Method Promotion (`<MethodName>Id`)
1. **Compile-time safety on fresh clones**: Every API package defines a concrete stub struct `type <MethodName>Id struct{}` with an actionable panic.
2. **Method promotion**: `SdkApis` in `apicollectionz/apis.go` embeds these stub structs, statically promoting every API method to the client without requiring `.gen.go` files on disk.
3. **Automated Codegen (`mold apisgen`)**: Emits `apis.gen.go` with real delegating methods that cleanly shadow the embedded stubs.
4. **Single Source of Truth**: The Go method name on `<MethodName>Id` defines the API ID via `camelz.P(name).QL()`.
5. **Zero-Underscore Rule**: API IDs must **never contain underscores (`_`)**. Kebab-case with hyphens (`-`) is strictly enforced across all IDs.

---

## 5. Benefits & Downstream Compatibility

1. **Zero Generator Disruption**:
   `forgejs` and `forgepy` read `apicollectionz.Apis()`, which inspects `ApiDataRefs()`. Since the generated `.gen.go` files implement `ApiDataRefs()`, TypeScript and Python SDK generation remains completely uninterrupted.
2. **Superior Go DX**:
   Go developers writing or using APIs get first-class IDE completion and godoc hovers directly from source comments.
3. **Maintainable Go SDK**:
   Writing `sdkgo` manually avoids the complexity of building a full Go AST emitter (`forgego`), while still keeping `sdkgo` lean, concurrency-safe, and idiomatic.

---

## 6. Implementation Plan & Milestones

### Phase 1: Prototype `schemagen` [DONE]
- Build prototype CLI in `mold/schemagen`.
- Test on a single endpoint (`apidatasetz/dszapidatasetlist`).
- Verify that generated `ApiDataRefs()` matches the expected schema and passes unit tests.

### Phase 2: Verify Multi-Language Codegen [DONE]
- Run `./comp.sh forgez test` and `./build.sh codegen`.
- Confirm that `sdkweb` (`sdkapis.js`) and `sdkpy` (`apis.py`, `types.py`) produce identical outputs with preserved JSDoc/docstrings.

### Phase 2.5: SdkApis & Stub Architecture [DONE]
- Refactor `apicollectionz` to a concrete `SdkApis` struct with embedded stub promotion.
- Build `mold apisgen` to automate `apis.gen.go` generation.
- Enforce zero underscores across all API IDs.

### Phase 3: Rollout Across API Modules & Backend Service Projects [IN PROGRESS]

Each `api*z` module in `i6sdk` pairs with a corresponding backend service repository in the workspace that implements its handlers:

| Step | SDK API Module (`i6sdk`) | Backend Service Project (`../*`) | Scope & Key Endpoints | Status |
|:---:|:---|:---|:---|:---|
| **1** | `apitfz/` | `i6tfz/` | 2 endpoints (`tfzapimodelmetadata`, `tfzapimodelpredict`) | **DONE** (100% green in both repos) |
| **2** | `apiingestz/` | `i6ingestz/` | 6 concrete endpoints (split `event[T]` & `stream[T]`, `batchsummary`, `geturl`) | **DONE** (100% green in both repos) |
| **3** | `apipipez/` | `i6pipez/` | 8 concrete endpoints (`pipezget`, `pipezlist`, split `pipezstart[T]`, `scheduleractivate[T]` & `schedulerdeactivate[T]`) | **DONE** (100% green in both repos) |
| **4** | `apiuserz/` | `i6userz/` + `i6decsuite/` | 18 endpoints (`userapiauth*`, `userapidatasettoken*`, `userapitokencheck`) | **DONE** (100% green across all 3 repos) |
| **5** | `apidatasetz/` | `i6datasetz/` | 24 concrete endpoints (split `dszapiloadresult[T]`, remaining 22 dataset endpoints) | **NEXT** (Phase 3 Step 5) |

#### Synchronized Refactoring Workflow (Per Domain Module)
1. **Refactor API Contract (`i6sdk/api*z`)**:
   - For generic families, split into concrete packages following Option A.
   - Convert requests and responses to explicit `{Params, Query, Headers, Payload}` envelopes.
   - Replace handwritten `ApiSpec()` with `ApiRoute() apiz.ApiSpec` (omitting `Id`).
   - Define concrete `<MethodName>Id struct{}` with stub method.
   - Embed `apiz.UnimplementedApi` with `json:"-"`.
   - Add rich Go doc comments on structs and fields for automatic schema documentation.
   - Add `//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .`.
   - Run `go generate ./...` to produce `<endpoint>.gen.go`.
2. **Register in `apicollectionz`**:
   - Embed `<MethodName>Id` in `SdkApis`.
   - Run `go generate ./...` in `apicollectionz` to execute `apisgen` and emit `apis.gen.go`.
3. **Update Backend Service Handler (`i6*`)**:
   - For generic families, implement Path 2B bridge interface runner.
   - Update handler to read inputs from `d.api.Req.Params`, `d.api.Req.Payload`, etc.
   - Update handler to write response to `d.api.Resp = &<endpoint>.Resp{ Payload: ... }`.
   - Update service unit tests to instantiate envelopes.
4. **Validate Pair**:
   - Run `./comp.sh <module> test` and `./comp.sh apicollectionz test` in `i6sdk`.
   - Run `go test -v ./...` in the backend service project.
   - Run `./build.sh codegen && ./build.sh test` to confirm multi-language parity.

### Phase 4: Build `sdkgo` Runtime
- Implement `I6Client`, `Invoker`, and `Tokens` (thread-safe double-checked renewal).
- Implement `Ingest`, `Infer`, and `Dash`.
- Add test coverage (`TestUnit*`) following mono-repo testing standards.
