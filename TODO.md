# Proposed Migration to gox

Evaluation of remaining candidates for migration to `../gox`:

| Package | Move to `gox`? | Target in `gox` | Verdict & Rationale |
| :--- | :---: | :--- | :--- |
| **`jstypez`** | **Yes** | `gox/schemaz/jstypez` | Pure, domain-agnostic JavaScript/TypeScript AST. Zero business logic. |
| **`jsschema`** (Core) | **Yes** | `gox/schemaz/jsschemaz` | Translating `schemaz.Schema` into JSDoc/TS `@typedef` is a generic schema capability. Uses `AddKnownType` for decoupling external types. |
| **`jsapi`** (Parser) | **Partially** | `gox/routez/jsapiz` | Parsing `routez.ApiSpec` into a JS API AST makes sense in `routez`. The generator (`this._services.invoker()`) remains in `i6sdk`. |

*(Note: `jssdkfunc` and `jswriter` stay in `i6sdk` as they are coupled to the i6 SDK runtime).*

---

### Migration Strategy: Thin Facade in `i6sdk`

To avoid churning 53+ test files across `sdkgo`, `i6sdk` will retain thin facade packages (`genjsz/jsschema` and `genjsz/jsapi`) that alias the types and wrap the constructors from `gox`:
- `jsschema.NewSdk()` in `i6sdk` returns `jsschemaz.New()` pre-configured with BigQuery and `fsz.FileStat` types.
- `jsapi.NewSdk()` in `i6sdk` returns `jsapiz.New(jsschema.NewSdk())`.
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
    - Maintained cloud-agnostic design (zero BigQuery dependencies).
    - Verified all unit tests pass in `gox/schemaz`.
  - [x] **2C (`i6sdk/genjsz/jsschema`)**:
    - Updated `genjsz/jsschema` to alias `jsschemaz.Parser` and provided `NewSdk()` registering BigQuery and `FileStat`.
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
    - Provide `NewSdk()` aliasing `jsapiz.New(jsschema.NewSdk())`.
    - Validate all tests in `i6sdk`.

- [ ] **Step 4: Final Verification**
  - Run `go vet ./...` and `go test ./...` in both `gox` and `i6sdk`.
