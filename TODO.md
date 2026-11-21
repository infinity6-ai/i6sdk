# Multi-Language Architecture: Mold -> Forge

Architecture for generating multi-language SDK clients (JS, Python, Java) from language-agnostic models.

### Metaphor
- **`mold` (Language-Agnostic ASTs & Parsers)**: Captures the shape and schema of types, endpoints, and SDK functions independent of target language.
- **`forgez` (Language-Specific Emitters)**: Pours the mold into target code (e.g., JavaScript/TypeScript, Python, Java).

---

### Target Layout

```text
gox/
├── schemaz/
│   ├── typeastz/                 # Universal type AST (Tuple, Typedef, CustomType)
│   └── schemamold/               # Schema reflection engine (Go struct/schema -> typeastz)
│
└── routez/
    └── apispecz/                 # HTTP route specification AST (Api, Apis, Req, Resp, Parser)

i6sdk/
├── mold/
│   └── moldsdkfunc/              # SDK function specification AST & Parser (SdkFunc, SdkFuncs, Parser)
│
└── forgez/
    ├── forgejs/                  # JavaScript / TypeScript Emitters
    │   ├── forgejsschema/        # Emits JSDoc / TS @typedef blocks
    │   ├── forgejsapi/           # Emits class Apis
    │   ├── forgejssdkfunc/       # Emits class SdkFuncs
    │   └── forgejsbundle/        # Assembles sdkapis.js bundle
    │
    ├── forgepy/                  # Python Emitters (Future)
    │   └── ...
    │
    └── forgejava/                # Java Emitters (Future)
        └── ...
```

---

### Migration Checklist

- [x] **Phase 1: Language-Agnostic Molds in `gox`**
  - [x] **1A (`gox/schemaz/typeastz`)**:
    - Created `gox/schemaz/typeastz` with `Tuple`, `Tuples`, `Typedef`, `Typedefs`, `Typ*` constants, and language-agnostic `CustomType` / `CustomTypeSimple` interfaces.
    - Updated `jsschemaz`, `jsapiz`, and all 38 files in `i6sdk` to import `typeastz`.
    - Removed `gox/schemaz/jstypez`.
    - Validated unit tests in `gox/schemaz`, `gox/routez`, and `i6sdk`.
  - [x] **1B (`gox/schemaz/schemamold`)**:
    - Renamed `gox/schemaz/jsschemaz` -> `gox/schemaz/schemamold`.
    - Updated references to `schemamold` in `gox/routez/jsapiz` and `i6sdk`.
    - Validated all unit tests in `gox/schemaz`, `gox/routez`, and `i6sdk`.
  - [x] **1C (`gox/routez/apispecz`)**:
    - Renamed `gox/routez/jsapiz` -> `gox/routez/apispecz`.
    - Updated references to `apispecz` in `i6sdk`.
    - Validated unit tests in `gox/routez` and `i6sdk`.

- [ ] **Phase 2: SDK Function Mold in `i6sdk`**
  - [ ] **2A (`i6sdk/mold/moldsdkfunc`)**:
    - Extract `SdkFunc`, `SdkFuncs`, `Req`, `Resp`, `SdkFuncSpec` from `genjsz/jssdkfunc` into `mold/moldsdkfunc`.
    - Move `Parser`, `New()`, `ParseSdkFunc`, and `ParseSdkFuncs` into `moldsdkfunc`.
    - Add unit tests in `mold/moldsdkfunc`.

- [ ] **Phase 3: JavaScript Forge in `i6sdk/forgez/forgejs`**
  - [ ] **3A (`forgez/forgejs/forgejsschema`)**:
    - Move `genjsz/jsschema` -> `forgez/forgejs/forgejsschema`.
    - Update imports to `gox/schemaz/schemamold` and `gox/schemaz/typeastz`.
  - [ ] **3B (`forgez/forgejs/forgejsapi`)**:
    - Move `genjsz/jsapi` -> `forgez/forgejs/forgejsapi`.
    - Update imports to `gox/routez/apispecz`.
  - [ ] **3C (`forgez/forgejs/forgejssdkfunc`)**:
    - Move writer from `genjsz/jssdkfunc` -> `forgez/forgejs/forgejssdkfunc`.
    - Update imports to `mold/moldsdkfunc`.
  - [ ] **3D (`forgez/forgejs/forgejsbundle`)**:
    - Move `genjsz/jswriter` -> `forgez/forgejs/forgejsbundle`.
    - Update generator orchestration in `sdkgo/internal/cmd/cmdsdkjs`.
  - [ ] **3E**: Delete `genjsz/` directory.

- [ ] **Phase 4: Workspace Verification**
  - Update all imports in `sdkgo/...` (unit tests and models).
  - Run `go vet ./...` and `go test ./...` in both `gox` and `i6sdk`.
