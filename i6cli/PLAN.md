# i6cli Architectural Plan & Roadmap

## 1. Executive Summary

`i6cli` provides **`i6`**, the official standalone static binary command-line interface for the Infinity6 platform.

### Key Objectives
- **Zero-Dependency Customer Tooling**: Customers can run `i6` on any machine, container, or CI/CD runner without installing Go, Python, or Node.js.
- **100% Feature Parity with SDKs**: Anything possible via `sdkgo`, `sdkpy`, or `sdkweb` can be performed from the terminal or shell scripts.
- **Built on `sdkgo`**: `i6cli` acts as a presentation and UX layer over `github.com/infinity6-ai/i6sdk/sdkgo`, eliminating duplicate protocol, auth, cache, or ingestion code.
- **True Static Binaries (`CGO_ENABLED=0`)**: Pure Go implementation using `modernc.org/sqlite` and Go standard library crypto/networking. Cross-compiles for Linux (`amd64`, `arm64`), macOS (`arm64`, `amd64`), and Windows without dynamic `libc`/`glibc` dependencies.
- **Script-Friendly & CI-Ready**: Configurable output formats (`--output=json|yaml|table|csv|quiet`), predictable exit codes, stdin/stdout piping, and non-interactive execution modes.

---

## 2. System Architecture

```
 ┌────────────────────────────────────────────────────────────────────────┐
 │                         i6 Command Line Tool                           │
 │   Environment: I6_CLI_CONFIG_DIR, I6_CLI_CONFIG_PROFILE, JSON override │
 └───────────────────────────────────┬────────────────────────────────────┘
                                     │
 ┌───────────────────────────────────▼────────────────────────────────────┐
 │                     Cobra Subcommand Taxonomy                          │
 │  auth │ tokens │ dataset │ ingest │ infer │ dash │ pipe │ run │ config  │
 └───────────────────────────────────┬────────────────────────────────────┘
                                     │ calls
 ┌───────────────────────────────────▼────────────────────────────────────┐
 │                     sdkgo Engine (I6Client)                            │
 ├────────────────┬─────────────────┬──────────────────┬──────────────────┤
 │ c.Apis()       │ c.Tokens()      │ c.Ingest()       │ c.Infer()        │
 │ c.Funcs()      │ c.Dash()        │ tokenstore.Sqlite│ fszcache.Fs      │
 └────────────────┴─────────────────┴──────────────────┴──────────────────┘
```

---

## 3. Command Taxonomy

Commands follow modern CLI standards (`i6 <domain> <action> [flags]`):

### A. Authentication & Session (`i6 auth`)
- `i6 auth login`: Interactive login supporting OAuth 2.0 Device Code flow (prints code + URL, polls server automatically) or prompt for API token.
- `i6 auth logout`: Revokes session token family on server via `UserAuthLogout` and wipes local credentials.
- `i6 auth status` / `i6 auth me`: Displays authenticated user profile, email, account age, and active token expiration.
- `i6 auth token`: Prints raw active JWT access token to stdout (ideal for piping to `curl` or downstream scripts).

### B. Token Management (`i6 tokens`)

#### General Token Operations
- `i6 tokens refresh`: Manually forces a silent token refresh rotation for active credentials.
- `i6 tokens dataset-access create <dataset> [--grants=0x0F]`: Mints or derives a scoped dataset access token.

#### Machine / Robot Dataset Refresh Tokens (`i6 tokens dataset-refresh`)
Dedicated management for long-lived machine/robot Ed25519 dataset refresh tokens. All subcommands accept the persistent flag `--dataset <dataset>` (falling back to the default dataset configured in the active profile's `config.json` if omitted):
- `i6 tokens dataset-refresh create [--dataset <dataset>]`: Creates a new dataset refresh token on the server for the target dataset and stores it in local `tokens.db`.
- `i6 tokens dataset-refresh list [--dataset <dataset>]`: Lists stored dataset refresh tokens (optionally filtered by `--dataset`). Does **not** display the sensitive `secret` column to keep terminal output safe.
- `i6 tokens dataset-refresh delete <jti> [--dataset <dataset>]`: Revokes a dataset refresh token on the server and deletes it from local `tokens.db`.
- `i6 tokens dataset-refresh expose <jti> --out <path|-> [--dataset <dataset>]`: Dumps the full JSON of the key identified by `[jti]` (including the sensitive `secret`). Requires `--out` (can be `-` for stdout or a file path). Takes exactly one `jti` argument.
- `i6 tokens dataset-refresh export [<jti...>] [--dataset <dataset>] [--all]`: Exports dataset refresh token definitions. Accepts one or more positional `jti` arguments, or the `--all` flag to export all stored refresh tokens (providing zero arguments without `--all` is an error).
- `i6 tokens dataset-refresh import [<file...>] [--dataset <dataset>] [--all]`: Imports dataset refresh tokens into the local `tokens.db`. Accepts one or more file paths, stdin (`-`), or extracts all keys from an export payload (providing zero arguments is an error).

### C. Datasets, Stores & Tables (`i6 dataset`)
- `i6 dataset list`: Lists accessible datasets.
- `i6 dataset get <dataset>`: Displays dataset metadata, solutions, and configuration.
- `i6 dataset table list <dataset>`: Lists tables in dataset.
- `i6 dataset table get <dataset> <table>`: Displays table schema and partition specifications.
- `i6 dataset store list <dataset>`: Lists connected stores.
- `i6 dataset store create <dataset> --name <name>`: Registers a new store.

### D. Data Ingestion (`i6 ingest`)
- `i6 ingest upload <dataset> <table> <files...>`:
  - Supports single/multiple files or shell globbing.
  - Automatically acquires upload URLs using `c.Ingest().GetUrl()`.
  - Concurrently uploads files via `c.Ingest().Upload()`.
  - Renders live interactive terminal progress bars (multi-stream overall and per-file progress).
- `i6 ingest get-url <dataset> <table> [--partitions key=val] [--amount N]`: Generates raw signed PUT upload URLs.
- `i6 ingest summary <dataset> <ingest-id>`: Checks batch lake validation status and summary.

### E. Analytics & Dashboards (`i6 dash`)
- `i6 dash query <query-name> --dataset <dataset> [--partitions key=val,...]`:
  - Downloads and caches analytical SQLite databases via `c.Dash()`.
  - Executes stored named queries and outputs formatted results (table or JSON).

### F. Inference & Scoring (`i6 infer`)
- `i6 infer run <model-name> --dataset <dataset> [--input <json-or-file>]`:
  - Executes the full 3-step inference pipeline (candidate generation -> remote TF scoring -> hydration).
- `i6 infer metadata <model-name>`: Fetches model schema and signature metadata.

### G. Typed Product Workflows (`i6 run`)
- `i6 run <workflow-name> --input <json-or-file>`:
  - Invokes typed `SdkFuncs` (e.g. `relevance-fashion-fbt`, `dash-funnel`) directly via `c.Funcs()`.

### H. Configuration (`i6 config`)
- `i6 config init [--box sandbox|staging|prod] [-f/--force]`: Initializes the active profile directory and writes `config.json` configured for the selected box (defaults to `prod`). Protected with `-f/--force` against accidental overwrites.
- `i6 config show`: Displays active profile, directory paths, box environment, and resolved endpoint URLs.

---

## 4. Configuration & State Layout

All CLI configuration, credentials, and local caches are organized per profile:

- **Config Base Directory**: `I6_CLI_CONFIG_DIR` (defaults to `$HOME/.config/i6/`).
- **Active Profile**: `I6_CLI_CONFIG_PROFILE` (defaults to `"main"`).
- **Inline JSON Override**: `I6_CLI_CONFIG_JSON` (optional inline JSON string for zero-disk CI/CD execution).
- *Clean Root Command*: There are no global CLI flags on the root `i6` command. All configuration routing is strictly environment-driven or relies on sensible defaults.

```
$I6_CLI_CONFIG_DIR/ (defaults to $HOME/.config/i6/)
└── [profile]/ (defaults to "main", selected via I6_CLI_CONFIG_PROFILE)
    ├── config.json          # Profile configuration: box (sandbox/staging/prod), default dataset
    ├── tokens.db            # Persistent SQLite token store via tokenstore.NewSqlite (WAL, 0600)
    └── cache/               # Profile-scoped SQLite blobs, queries, models, ETags (fszcache)
        ├── dash/
        └── infer/
```

### Configuration Hierarchy & Resolution Order
1. **Zero-Disk Inline JSON**: If `I6_CLI_CONFIG_JSON` is populated, parse in-memory immediately using `tokenstore.NewMemory()` (bypasses filesystem reads, ideal for CI/CD runners).
2. **Environment & Profile Defaults**:
   - Target Directory: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/`
   - Config file: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/config.json`
   - Token store: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/tokens.db` (via `tokenstore.NewSqlite`)
   - Cache store: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/cache/`

Specific subcommands take positional arguments or command-scoped flags (e.g. `i6 dataset list --output=json`), preventing flag pollution on the root command.

---

## 5. Mono-Repo File Structure

```
i6cli/
├── go.mod
├── go.sum
├── main.go                       # Minimal entrypoint calling cmdcli.Execute()
├── PLAN.md                       # This plan
├── README.md                     # CLI installation and usage guide
├── cmdcli/
│   ├── cmdcli.go                 # Root command and version subcommand
│   ├── cmdcli_unit_test.go
│   ├── auth.go                   # i6 auth subcommands
│   ├── tokens.go                 # i6 tokens subcommands
│   ├── dataset.go                # i6 dataset subcommands
│   ├── ingest.go                 # i6 ingest subcommands
│   ├── infer.go                  # i6 infer subcommands
│   ├── dash.go                   # i6 dash subcommands
│   ├── pipe.go                   # i6 pipe subcommands
│   ├── run.go                    # i6 run workflow subcommand
│   └── config.go                 # i6 config subcommands
└── internal/
    ├── cliconfig/                # Config file parsing and profile management
    │   ├── cliconfig.go
    │   └── cliconfig_unit_test.go
    ├── outfmt/                   # Table, JSON, YAML, and CSV output formatters
    │   ├── outfmt.go
    │   └── outfmt_unit_test.go
    └── ui/                       # Terminal spinners, progress bars, and prompts
        ├── progress.go
        └── spinner.go
```

---

## 6. Static Binary Cross-Compilation (`CGO_ENABLED=0`)

All builds produce pure Go static binaries with zero external dynamic library dependencies:

```bash
# Linux (amd64 / x86_64)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
  -ldflags="-s -w -X 'github.com/infinity6-ai/i6sdk/versionz/version.GitCommit=${SHA}'" \
  -o dist/i6-linux-amd64 .

# Linux (arm64 / aarch64)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath \
  -ldflags="-s -w -X 'github.com/infinity6-ai/i6sdk/versionz/version.GitCommit=${SHA}'" \
  -o dist/i6-linux-arm64 .

# macOS (Apple Silicon arm64)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath \
  -ldflags="-s -w -X 'github.com/infinity6-ai/i6sdk/versionz/version.GitCommit=${SHA}'" \
  -o dist/i6-darwin-arm64 .

# macOS (Intel amd64)
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath \
  -ldflags="-s -w -X 'github.com/infinity6-ai/i6sdk/versionz/version.GitCommit=${SHA}'" \
  -o dist/i6-darwin-amd64 .

# Windows (amd64)
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath \
  -ldflags="-s -w -X 'github.com/infinity6-ai/i6sdk/versionz/version.GitCommit=${SHA}'" \
  -o dist/i6-windows-amd64.exe .
```

### GitHub Actions Release Pipeline
In `.github/workflows/release.yml`:
- Cross-compile static binaries across the 5 target platforms.
- Generate `checksums.txt` (SHA256).
- Compress into `.tar.gz` (POSIX) and `.zip` (Windows).
- Attach to GitHub Releases on both tagged releases (`v*`) and continuous `edge` builds.
- Provide a universal install script:
  ```bash
  curl -fsSL https://raw.githubusercontent.com/infinity6-ai/i6sdk/main/install.sh | sh
  ```

---

## 7. Phased Implementation Roadmap

- [ ] **Phase 1: Foundation & Scaffolding**
  - [x] Scaffold `i6cli/` directory, `go.mod`, `main.go`, and `cmdcli/cmdcli.go`.
  - [x] Integrate with `build.sh` (`cmd_comps_list`), `comp.sh`, and root `go.work`.
  - [x] Verify test execution via `./comp.sh i6cli test` and `./build.sh test`.
  - [ ] Implement `internal/cliconfig` (profile, base directory, and config.json resolution).
  - [ ] Implement `internal/outfmt` (table/JSON formatting for subcommands).

- [ ] **Phase 2: Authentication & Token Engine (`auth`, `tokens`, `config`)**
  - [ ] Implement `i6 config init [--box sandbox|staging|prod] [-f/--force]` and `i6 config show`.
  - [ ] Implement `i6 auth login` (Device code flow with polling + manual token input).
  - [ ] Implement `i6 auth logout`, `status`, and `token`.
  - [ ] Implement `i6 tokens refresh` and `i6 tokens dataset-access create`.
  - [ ] Implement `i6 tokens dataset-refresh [--dataset <dataset>]` persistent flag with subcommands:
    - `create` (mints on server, stores in local tokens.db)
    - `list` (safely omits secrets, optionally filtered by --dataset)
    - `delete <jti>` (revokes on server and removes locally)
    - `expose <jti> --out <path|->` (outputs full JSON with secret)
    - `export [<jti...>] [--all]` (requires at least one jti or --all)
    - `import [<file...>] [--all]` (supports files, stdin `-`, and payloads)

- [ ] **Phase 3: Ingestion & Dataset Management (`ingest`, `dataset`)**
  - [ ] Implement `i6 ingest upload` with multi-stream parallel worker execution and live terminal progress bars.
  - [ ] Implement `i6 ingest get-url` and `summary`.
  - [ ] Implement `i6 dataset list`, `get`, `table`, and `store`.

- [ ] **Phase 4: Analytics, Inference & Typed Workflows (`dash`, `infer`, `pipe`, `run`)**
  - [ ] Implement `i6 dash query` with table output.
  - [ ] Implement `i6 infer run` and `metadata`.
  - [ ] Implement `i6 pipe list`, `get`, `start`.
  - [ ] Implement `i6 run` for typed `SdkFuncs` dispatch.

- [ ] **Phase 5: Release Packaging & Automation**
  - [ ] Add static binary cross-compilation matrix to `.github/workflows/release.yml`.
  - [ ] Add `install.sh` shell script for instant installation.
  - [ ] Write `i6cli/README.md` with complete command reference and quickstart examples.
