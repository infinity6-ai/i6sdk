# Proposed Migration to gox

Evaluation of remaining candidates for migration to `../gox`:

| Package | Move to `gox`? | Target in `gox` | Verdict & Rationale |
| :--- | :---: | :--- | :--- |
| **`jstypez`** | **Yes** | `gox/schemaz/jstypez` | Pure, domain-agnostic JavaScript/TypeScript AST. Zero business logic. |
| **`jsschema`** (Core) | **Yes** *(with caveats)* | `gox/schemaz/jsschemaz` | Translating `schemaz.Schema` into JSDoc/TS `@typedef` is a generic schema capability. |
| **`jsapi`** (Parser) | **Partially** | `gox/routez/jsapiz` | Parsing `routez.ApiSpec` into a JS API AST makes sense in `routez`. The generator (`this._services.invoker()`) remains in `i6sdk`. |

*(Note: `jssdkfunc` and `jswriter` stay in `i6sdk` as they are coupled to the i6 SDK runtime).*

---

### Migration Checklist

- [x] **Step 1: Replace `camelstring` with `gox/commonz/stringz/camelz`**
  - Removed local `camelstring` package.
  - Updated all call sites in `jsapi`, `jsschema`, and `jssdkfunc` to use `gox/camelz`.
- [ ] **Step 2: Move `jstypez` and generic schema reflection (`FormatTypedefs`) to `gox/schemaz`**
  - Create `gox/schemaz/jstypez` with core AST and `CustomType` interfaces.
  - Create `gox/schemaz/jsschemaz` with schema-to-AST translation and `FormatTypedefs`.
  - Provide a type registration hook for external types (e.g., BigQuery types) rather than importing cloud libraries into `gox/schemaz`.
  - Keep `WriteModuleTypedefs` (SDK-specific imports) in `i6sdk`.
- [ ] **Step 3: Move `jsapi` route parser to `gox/routez/jsapiz`**
  - Move `ParseApi(ctx, apiz.Api)` to `gox/routez/jsapiz`.
  - Retain `i6sdk`-specific client class generation (`Apis`) in `i6sdk`.
