# Issue: `sdkgo` Lacks Client Runtime Implementation and Code Generation

**Component:** `sdkgo`, `forgez`  
**Severity:** Critical (Go SDK cannot be used to invoke services)  
**Status:** Open  

---

## Summary

While `sdkweb` and `sdkpy` provide fully functioning client runtimes with generated endpoint methods, authentication, token management, file ingestion, and local SQLite inference, **`sdkgo` is currently only a structural stub**. It contains no HTTP invocation engine, no authentication or token management client, and no code generation target in `forgez`.

---

## Root Cause Analysis

### 1. `sdkgo/sdkgo/client.go` Contains Only Static Spec Reflection
The existing Go `Client` struct:
```go
func (c *Client) Version() string { return version.Version() }
func (c *Client) Config() Config { return c.config }
func (c *Client) Apis() []apiz.Api { return apicollectionz.Apis() }
func (c *Client) Api(id string) apiz.Api { return apicollectionz.Api(id) }
func (c *Client) SdkFuncs() []sdkfuncz.Func { return apicollectionz.SdkFuncs() }
func (c *Client) SdkFunc(id string) sdkfuncz.Func { return apicollectionz.SdkFunc(id) }
```
* `Apis()` and `Api()` return server route descriptors (`apiz.Api`), not client calling methods.
* `SdkFuncs()` returns server `sdkfuncz.Func` interfaces.
* There is no `Invoker`, no HTTP request dispatching, no serialization of request bodies, and no response parsing.

### 2. Missing `forgego` Generator in `forgez`
The `forgez` module contains code generators for JavaScript (`forgejs`) and Python (`forgepy`):
* `forgez/forgejs/` -> outputs `sdkweb/src/main/sdkapis/sdkapis.js`
* `forgez/forgepy/` -> outputs `sdkpy/src/i6sdk/sdkapis/`
* **`forgego` does not exist.** There is no generator producing typed Go structs, HTTP request builders, and method signatures for calling the Infinity6 APIs.

### 3. Missing Client Subsystems
`sdkgo` has zero implementations for:
* **`auth`**: No device code flow, PKCE, or direct credentials.
* **`tokens`**: No user or dataset token cache, rotation, or auto-renewal.
* **`ingest`**: No multipart/signed URL file upload client.
* **`infer`**: No model bundle download, local SQLite named query execution, or server-side prediction bridge.
* **`dash`**: No dashboard query runner.

---

## Proposed Remediation

1. Implement `forgego` inside `forgez` to generate strongly-typed Go client methods and request/response structs from `apicollectionz.Apis()` and `apicollectionz.SdkFuncs()`.
2. Implement a Go `Invoker` in `sdkgo` leveraging standard `net/http.Client`.
3. Port `Auth`, `Tokens`, `Ingest`, `Infer`, and `Dash` clients to Go.
