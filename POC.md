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

type Api struct {
    Req  *Req  `json:"req"`
    Resp *Resp `json:"resp"`
}
```

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

### Phase 1: Prototype `schemagen`
- Build prototype CLI in `cmd/schemagen` or `schemaz/cmdschemagen`.
- Test on a single endpoint (e.g. `apidatasetz/dszapidatasetlist`).
- Verify that generated `ApiDataRefs()` matches the expected schema and passes unit tests.

### Phase 2: Verify Multi-Language Codegen
- Run `./comp.sh forgez test` and `./build.sh codegen`.
- Confirm that `sdkweb` (`sdkapis.js`) and `sdkpy` (`apis.py`, `types.py`) produce identical outputs with preserved JSDoc/docstrings.

### Phase 3: Rollout Across API Modules
- Refactor API packages to the `{Params, Query, Headers, Payload}` envelope convention (omitting unused slots).
- Add `//go:generate` directives and remove handwritten `ApiDataRefs()`.
- Update `comp.sh` and `build.sh` codegen steps.

### Phase 4: Build `sdkgo` Runtime
- Implement `I6Client`, `Invoker`, and `Tokens` (thread-safe double-checked renewal).
- Implement `Ingest`, `Infer`, and `Dash`.
- Add test coverage (`TestUnit*`) following mono-repo testing standards.
