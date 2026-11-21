# Proposed Migration to gox

Evaluation and progress of migrating modules to `../gox`:

| Package | Status | Target in `gox` | Details & Rationale |
| :--- | :---: | :--- | :--- |
| **`camelstring`** | **Done** | `gox/commonz/stringz/camelz` | Deleted local package; replaced with `gox/camelz`. |
| **`jstypez`** | **Done** | `gox/schemaz/jstypez` | Moved pure JS/TS type AST (`Tuple`, `Typedef`, etc.) and unit tests to `gox`. Local package deleted. |
| **`jsschema`** (Core) | **Done** | `gox/schemaz/jsschemaz` | Moved schema reflection engine (`Parser`, `AddKnownType`, `FormatTypedefs`) to `gox`. Stripped cloud dependencies. `i6sdk` now uses a clean `NewSdk()` factory. |
| **`jsapi`** (Parser) | **Pending (Step 3)** | `gox/routez/jsapiz` | Move `routez.ApiSpec` parser and AST to `gox/routez`. Retain client class generator in `i6sdk`. |

*(Note: `jssdkfunc` and `jswriter` stay in `i6sdk` as they are coupled to the i6 SDK runtime).*

---

### Migration Strategy: Thin Facade in `i6sdk`

To avoid churning 53+ test files across `sdkgo`, `i6sdk` retains thin facade packages (`genjsz/jsschema` and `genjsz/jsapi`) that wrap the constructors from `gox`:
- `jsschema.NewSdk()` in `i6sdk` returns `*jsschemaz.Parser` pre-configured with BigQuery and `fsz.FileStat` types.
- `jsapi.NewSdk()` in `i6sdk` will return `*jsapiz.Parser` configured with `jsschema.NewSdk()`.
- Existing tests and callers continue using `jsschema.NewSdk().ParseSchema(...)` and `jsapi.NewSdk().ParseApi(...)` without any import changes.

---

### Migration Checklist

- [x] **Step 1: Replace `camelstring` with `gox/commonz/stringz/camelz`**
  - Removed local `camelstring` package.
  - Updated all call sites in `jsapi`, `jsschema`, and `jssdkfunc` to use `gox/camelz`.

- [x] **Step 2: Move `jstypez` and generic schema reflection (`jsschemaz`) to `gox/schemaz`**
  - [x] **2A (`gox/schemaz/jstypez`)**:
    - Created `gox/schemaz/jstypez` with `Tuple`, `Tuples`, `Typedef`, `Typedefs`, `Typ*` constants, and `CustomType` / `CustomTypeSimple` interfaces.
    - Added `jstypez_typedefs_unit_test.go`.
    - Verified all unit tests pass in `gox/schemaz`.
  - [x] **2B (`gox/schemaz/jsschemaz`)**:
    - Created `gox/schemaz/jsschemaz` with `Parser`, `New()`, `AddKnownType()`, and instance methods (`ParseSchema`, `ParseAny`, `ParseSchemaObject`, etc.).
    - Added `FormatTypedefs` (JSDoc typedef block formatter).
    - Added `DescName`, `DescSummary`, and `Resolve`.
    - Stripped `logz` to ensure zero indirect dependencies in `schemaz/go.mod`.
    - Verified all unit tests pass in `gox/schemaz`.
  - [x] **2C (`i6sdk/genjsz/jsschema`)**:
    - Simplified `genjsz/jsschema/jsschema.go` into a clean factory providing strictly `NewSdk() *jsschemaz.Parser` with BigQuery and `FileStat` types registered.
    - Removed redundant pass-through aliases (`Parser`, `New`, `DescName`, `DescSummary`, `FormatTypedefs`) in favor of direct `jsschemaz` calls.
    - Retained `WriteModuleTypedefs` in `i6sdk` (SDK-specific imports for `Services` and `ServiceOptions`).
    - Removed `i6sdk/genjsz/jstypez` in favor of `gox/schemaz/jstypez`.
    - Validated all tests in `i6sdk`.

- [ ] **Step 3: Move `jsapi` route parser to `gox/routez/jsapiz`**
  - **3A (`gox/routez/jsapiz`)**:
    - Move `Api`, `Apis`, `Req`, `Resp`, `ApiSpec` data structures.
    - Move `Parser`, `New()`, `ParseApi`, `ParseApis`, and `ExtractTypedefs`.
    - Validate unit tests in `gox/routez`.
  - **3B (`i6sdk/genjsz/jsapi`)**:
    - Retain `Write` / `WriteV2` (generates the `i6`-specific `this._services.invoker()` client class).
    - Provide `NewSdk()` returning `*jsapiz.Parser` configured with `jsschema.NewSdk()`.
    - Validate all tests in `i6sdk`.

- [ ] **Step 4: Final Verification**
  - Run `go vet ./...` and `go test ./...` in both `gox` and `i6sdk`.
