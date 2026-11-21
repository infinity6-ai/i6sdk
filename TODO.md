# Proposed Migration to gox

Evaluation of moving `genjsz` modules to `../gox/schemaz`, `../gox/routez`, or `../gox/commonz`:

| Package | Move to `gox`? | Target in `gox` | Verdict & Rationale |
| :--- | :---: | :--- | :--- |
| **`jstypez`** | **Yes** | `gox/schemaz/jstypez` | Pure, domain-agnostic JavaScript/TypeScript AST. Zero business logic. |
| **`jsschema`** (Core) | **Yes** *(with caveats)* | `gox/schemaz/jsschemaz` | Translating `schemaz.Schema` into JSDoc/TS `@typedef` is a generic schema capability. |
| **`jsapi`** | **Partially** | `gox/routez/jsapiz` | Parsing `routez.ApiSpec` into a JS API AST makes sense in `routez`, but the generator is currently hardcoded to the `i6` invoker runtime. |
| **`jssdkfunc`** | **No** | *(Keep in `i6sdk`)* | Strictly coupled to `i6sdk`'s inference, dashboards, and `sdkfuncz`. |
| **`jswriter`** | **No** | *(Keep in `i6sdk`)* | Application-level orchestrator generating the legacy web SDK bundle. |
| **`camelstring`** | **Delete / Replace** | `gox/commonz/stringz/camelz` | Already implemented in `gox/commonz/stringz/camelz`; local copy in `i6sdk` is redundant. |

---

### Key Considerations Before Migrating

1. **`jsschema` Caveats**:
   - **BigQuery dependency**: `jsschema_known.go` currently imports `cloud.google.com/go/bigquery`. Provide a type-registration hook (e.g. `RegisterKnownType(...)`) instead of importing BigQuery into `gox/schemaz`.
   - **SDK-specific typedef headers**: Keep `WriteModuleTypedefs` (which imports `../services/services.js` and `../invoker/invoker.js`) in `i6sdk`; move only the generic `FormatTypedefs` to `gox/schemaz`.

2. **`jsapi` Invoker Runtime**:
   - `ParseApi(ctx, apiz.Api)` can live in `gox/routez/jsapiz`.
   - The generator producing `this._services.invoker().invokeResult(...)` is specific to `i6sdk` and should either remain in `i6sdk` or use a pluggable invoker template (e.g., standard `fetch`).

3. **Migration Steps**:
   - **Step 1**: Replace `i6sdk/camelstring` with `gox/commonz/stringz/camelz`.
   - **Step 2**: Move `jstypez` and generic schema reflection (`FormatTypedefs`) to `gox/schemaz/jsschemaz`.
   - **Step 3**: Move `jsapi` parser logic to `gox/routez/jsapiz`.
