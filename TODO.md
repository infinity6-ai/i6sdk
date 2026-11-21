# Integrating `brwjsz` to Execute `sdkweblegacy` QUnit Tests

Migrate the web SDK test runner from legacy `webchecker` to the `brwjsz` browser test framework in `gox`.

---

### Motivation & Goals
- **Native Go Test Integration:** Map all 12 QUnit test suites directly into Go sub-tests (`t.Run`) via `brwqunit.Test`.
- **Remove Legacy Dependencies:** Deprecate and remove `go.code.infinity6.ai/web/api/webchecker` and related legacy platform packages.
- **Fast Developer Workflow:** Enable dual execution via automated `go test -run '^TestWeb.*'` and the standalone `brwjsz qunit` CLI tool.
- **Full Module & Bare Specifier Support:** Ensure browser-side ES module imports (`sql.js`, `@litertjs/core`) resolve cleanly via import maps.

---

### Architecture & Serving Layout

```text
sdkweblegacy/ (served as Site root)
├── src/                      # ES modules (index.js, litert.js, etc.)
├── test/                     # QUnit test suites (*-test.js)
├── static/                   # Static assets, wasm files, html
└── sdkweb_web_test.go        # Go browser runner invoking brwqunit
```

---

### Implementation Checklist

- [x] **Phase 1: Extend `brwjsz` / `brwqunit` for Import Maps & Module Specifiers (`gox/brwjsz`)**
  - [x] **1A (Import Map Support in Runner)**:
    - Update `brwjsz/staticz/qunit/runner.js` to accept optional `importmap` definitions.
    - Inject `<script type="importmap">` into the test iframe before loading scripts.
  - [x] **1B (`brwqunit` Go API)**:
    - Add `ImportMap map[string]string` (or options struct) to `brwqunit.Run` / `brwqunit.Test`.
    - Support passing custom import maps from Go test suites.
  - [x] **1C (`cmdqunit` CLI)**:
    - Add `--importmap` flag (file path or inline JSON) to `brwjsz qunit` command.
  - [x] **1D (Verification in `gox`)**:
    - Add unit/web tests in `gox/brwjsz/brwqunit` verifying that bare specifiers imported inside test modules resolve properly.
    - Run `cmd_unit_test.sh` in `gox/brwjsz`.

- [x] **Phase 2: Integrate `brwjsz` in `sdkweblegacy` (`i6sdk/sdkweblegacy`)**
  - [x] **2A (Dependency Configuration)**:
    - Add `github.com/infinity6-ai/gox/brwjsz` to `sdkweblegacy/go.mod`.
  - [x] **2B (Web Test Implementation)**:
    - Implement `TestWebSdk(t *testing.T)` in `sdkweblegacy` using `brwjsz.New` and `brwqunit.Test`.
    - Configure `Site` to serve `sdkweblegacy` directory.
    - Configure import map for:
      - `sql.js` -> `https://cdn.jsdelivr.net/npm/sql.js@1.14.2/+esm`
      - `@litertjs/core` -> `https://cdn.jsdelivr.net/npm/@litertjs/core@2.5.3/+esm`
    - Execute discovery or run individual suites (`sample-test.js`, `b64-test.js`, `auth-tokens-test.js`, `dash-test.js`, etc.).
  - [x] **2C (Headless & Visual Support)**:
    - Respect `BRWJSZ_DISABLE_HEADLESS` configuration for debugging and visual test inspection.

- [x] **Phase 3: Tooling & Script Updates**
  - [x] **3A (`build.sh`)**:
    - Update `cmd_test_web` in `sdkweblegacy/build.sh` to run `go test -run '^TestWeb.*' ./...`.
    - Update `cmd_test_visual` to export `BRWJSZ_DISABLE_HEADLESS=true`.
  - [x] **3B (CLI Verification)**:
    - Verify `brwjsz qunit --site .` works directly from `sdkweblegacy`.

- [x] **Phase 4: Deprecate & Clean Up Legacy `webchecker`**
  - [x] **4A**: Remove `internal/appstack/appstack_web_test.go` and references to `webchecker.RunWeb`.
  - [x] **4B**: Clean up unused `webchecker.js` from `static/test/webchecker.js` or retain standalone `test.html` without server dependencies.
  - [x] **4C**: Run `go mod tidy` in `sdkweblegacy`.

- [ ] **Phase 5: Full Workspace Verification**
  - [ ] Run `go test -run '^(TestUnit.*|TestWeb.*)$' ./...` in `sdkweblegacy`.
  - [ ] Run `go vet ./...` in both `gox` and `i6sdk`.
