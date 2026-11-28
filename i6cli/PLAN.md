# i6cli Architectural Plan & Roadmap

## 1. Executive Summary

`i6cli` provides **`i6`**, the official standalone static binary command-line interface for the Infinity6 platform.

### Key Objectives
- **Zero-Dependency Customer Tooling**: Customers can run `i6` on any machine, container, or CI/CD runner without installing Go, Python, or Node.js.
- **100% Feature Parity with SDKs**: Anything possible via `sdkgo`, `sdkpy`, or `sdkweb` can be performed from the terminal or shell scripts.
- **Built on `sdkgo`**: `i6cli` acts as a presentation and UX layer over `github.com/infinity6-ai/i6sdk/sdkgo`, eliminating duplicate protocol, auth, cache, or ingestion code.
- **True Static Binaries (`CGO_ENABLED=0`)**: Pure Go implementation using `modernc.org/sqlite` and Go standard library crypto/networking. Cross-compiles for Linux (`amd64`, `arm64`), macOS (`arm64`, `amd64`), and Windows without dynamic `libc`/`glibc` dependencies.
- **Script-Friendly & CI-Ready**: Flexible output formatting (`--format=table|tsv|json|quiet|raw`), strict stdout/stderr separation, predictable exit codes, stdin/stdout piping, and non-interactive execution modes.

---

## 2. System Architecture

```
 ┌────────────────────────────────────────────────────────────────────────┐
 │                         i6 Command Line Tool                           │
 │   Global flags / Env: --profile / I6_CLI_CONFIG_PROFILE, JSON override │
 └───────────────────────────────────┬────────────────────────────────────┘
                                     │
 ┌───────────────────────────────────▼────────────────────────────────────┐
 │                     Cobra Subcommand Taxonomy                          │
 │  auth │ tokens │ dataset │ ingest │ dash │ infer │ pipe │ funcs │ apis │
 │  config │ completion                                                   │
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
- `i6 auth login`: Interactive browser login screen (OAuth 2.0 Device Code flow). Displays authorization code and URL, polls `/api/userz/auth/device/token`, and saves full user session in `tokens.db`.
- `i6 auth token [<token-string> | -] [--offline]`: Universal token ingestion command. Automatically detects the token type via local formatting and `/api/userz/token/check` (with `--offline`, skips network calls and ingests locally):
  - **Ed25519 Dataset Robot Key**: Ingests as machine robot root credential (`i6_rf:<dataset>`).
  - **User Refresh Token**: Refreshes via `/api/userz/auth/refresh` to establish a full user session (`i6_user_user_token`).
  - **Standalone User Access Token**: Saved as a non-refreshable user session in `i6_user_user_token`.
  - **Standalone Dataset Bearer Token**: Saved as a non-refreshable dataset access token in `i6_ds:<dataset>`.
- `i6 auth logout`: Revokes session token family on server via `UserAuthLogout` and wipes local credentials from active profile.
- `i6 auth status` / `i6 auth me`: Displays active authentication identity (User Session vs. Robot Service Account vs. Standalone Bearer vs. Unauthenticated), tenant, dataset, and token expiration.
- `i6 auth inspect [<token-string | ->]`: Performs deep cryptographic and claims diagnostics on a User JWT access token via `/api/userz/token/check`. If `<token-string>` is omitted, audits the active profile's stored user session access token.

### B. Local Profile Token Management (`i6 tokens`)

Manages the active credentials stored in the local profile's `tokens.db`.

*Root Credential Invariant*: A profile has either a **Human User Session** (`tokens user`) OR a **Machine Robot Key** (`tokens dataset-refresh`), never both. This eliminates ambiguity in identity, signing, and permissions.

*Security Invariant*: `list` and `get` commands **never** display sensitive `secret` strings, tokens, or private keys. Revealing secrets is strictly restricted to explicit `expose` and `export` commands.

#### Profile Tokens Health Audit (`i6 tokens check [--refreshed]`)
- `i6 tokens check [--refreshed]`: Performs an aggregated health audit of all tokens stored in the active profile's `tokens.db`:
  - Validates the active root credential (`tokens user` or `tokens dataset-refresh`).
  - Iterates through and validates all cached derived tokens (`tokens dataset-access`).
  - Calls `/api/userz/token/check` for each token and displays an aggregated status report.

#### Deep Profile Tokens Inspection (`i6 tokens inspect [--refreshed]`)
- `i6 tokens inspect [--refreshed]`: Performs a deep cryptographic and diagnostic inspection of every credential currently stored in the active profile's `tokens.db` against `/api/userz/token/check`:
  - Inspects active root credential (user JWT or robot key assertion).
  - Inspects all cached derived dataset access tokens.
  - Displays detailed diagnostics: token classification, mathematical validity, server summary message, expiration, remaining TTL, effective grants, and hints.

#### 1. Human User Root Credential (`i6 tokens user`)
Controls the authenticated human user session in `tokens.db`:
- `i6 tokens user get`: Displays user profile metadata (user_id, tenant, expiration; omits tokens).
- `i6 tokens user expose`: Prints raw active JWT access token to stdout (for curl / scripts).
- `i6 tokens user export`: Prints the active user refresh token to stdout (for backup or environment transfer).
- `i6 tokens user attenuation --dataset <dataset> [--grants <grants>] [--exp <seconds>] [--solution <solution>]`: Mints an in-flight attenuated child access token for `--dataset` with reduced permissions and shortened TTL.
- `i6 tokens user delete`: Revokes session on server and removes credentials from `tokens.db` (logout).

#### 2. Machine Robot Root Credential (`i6 tokens dataset-refresh`)
Controls the machine service account Ed25519 refresh key (mutually exclusive with `tokens user`):
- `i6 tokens dataset-refresh get`: Displays robot key metadata for the dataset (safely omits the secret key).
- `i6 tokens dataset-refresh export [--out <path|->]`: Prints the active Ed25519 dataset refresh token to stdout or a file.
- `i6 tokens dataset-refresh set <key|file|-> [--dataset <dataset>]`: Sets or replaces the robot key for this profile.
- `i6 tokens dataset-refresh attenuation [--grants <grants>] [--exp <seconds>] [--solution <solution>]`: Uses the stored robot key to mint an attenuated child access token for `--dataset` with reduced permissions.
- `i6 tokens dataset-refresh delete`: Removes the robot credential from `tokens.db`.
- **Sender-Constrained PoP Operations (`i6 tokens dataset-refresh bounded`)**:
  Uses the profile's stored Ed25519 robot key without requiring manual key flags:
  - `i6 tokens dataset-refresh bounded create [--grants=0x0F] [--exp=3600]`: Mints a sender-constrained Proof-of-Possession token cryptographically bound to the profile's robot public key.
  - `i6 tokens dataset-refresh bounded sign-request --method <METHOD> --path <PATH> [--body <str|->]`: Generates the 121-byte Ed25519 signature envelope (`X-I6-Signature` header) over the target request using the profile's stored robot key.

#### 3. Derived Access Tokens Cache (`i6 tokens dataset-access [--dataset <dataset>]`)
Manages short-lived bearer access tokens cached in `tokens.db` (a profile can hold multiple tokens across different datasets):
- `i6 tokens dataset-access list`: Lists all cached dataset access tokens across all datasets in `tokens.db`.
- `i6 tokens dataset-access get [--refreshed]`: Displays metadata for the target dataset (expiry, grants, solution) without revealing the raw token. With `--refreshed`, automatically refreshes using the root credential if missing or expired.
- `i6 tokens dataset-access expose [--refreshed]`: Prints raw Bearer JWT to stdout. With `--refreshed`, automatically refreshes if missing or expired before printing.
- `i6 tokens dataset-access attenuation [--grants <grants>] [--exp <seconds>] [--solution <solution>]`: Calls `/api/userz/dataset-token/access/attenuate` to derive an in-flight child access token with further reduced permissions from the active cached token.
- `i6 tokens dataset-access check [--refreshed]`: Cryptographically validates the token against `/api/userz/token/check`.
- `i6 tokens dataset-access delete`: Evicts the cached bearer token for `--dataset` from local `tokens.db`.

#### 4. Stateless Zero-Disk Operations (`i6 tokens raw`)
Pure in-memory cryptographic and token exchange operations that never read or write to `tokens.db`:
- `i6 tokens raw check <token-string | ->`: Validates and debugs an arbitrary token string or stdin pipe via `/api/userz/token/check` without saving or affecting `tokens.db`.
- `i6 tokens raw dataset-bearer --refresh-token <key|-> --dataset <dataset>`: Statelessly signs an Ed25519 `RefreshAssertion`, exchanges it on the server, and prints the raw short-lived bearer access token to stdout without touching disk.
- `i6 tokens raw attenuation [<token|->] [--grants <mask|names>] [--solution <sol>] [--exp <seconds>]`: Calls `/api/userz/dataset-token/access/attenuate` to derive an in-flight child access token with attenuated permissions and shortened TTL, emitting the child token directly to stdout.
- `i6 tokens raw bounded-create --token <token|-> --key <privkey|->`: Mints a sender-constrained Proof-of-Possession token from an access token and private key.
- `i6 tokens raw sign-request --key <privkey|-> --method <METHOD> --path <PATH> [--body <str|->]`: Generates the 121-byte Ed25519 signature envelope (`X-I6-Signature` header) using the `reqsign` protocol.

### C. Datasets, Stores, Tables & Server Tokens (`i6 dataset`)

#### Dataset Administration
- `i6 dataset list`: Lists accessible datasets.
- `i6 dataset get <dataset>`: Displays dataset metadata, solutions, and configuration.
- `i6 dataset table list <dataset>`: Lists tables in dataset.
- `i6 dataset table get <dataset> <table>`: Displays table schema and partition specifications.
- `i6 dataset store list <dataset>`: Lists connected stores.
- `i6 dataset store create <dataset> --name <name>`: Registers a new store.

#### Remote Server Machine Key Registry (`i6 dataset keys [--dataset <dataset>]`)
Server-side key management where a dataset can have multiple registered machine robot keys with distinct JTIs:
- `i6 dataset keys list`: Calls server API (`/api/userz/dataset-token-refresh/list`) to list all registered machine robot keys for the dataset (showing JTI, solution, grants, expiration; safely omits private keys).
- `i6 dataset keys create [--solution <sol>] [--grants <g>] [--annotation <text>] [--save]`: Calls server API (`/api/userz/dataset-token-refresh/create`) to register a new machine robot Ed25519 key on the server. Dumps the one-time private key artifact; with `--save`, automatically stores it into active profile's `tokens.db`.
- `i6 dataset keys reqsign --dataset <dataset> --method <METHOD> --path <PATH> [--body <str|->] [--key <privkey|->]`: Generates the 121-byte Ed25519 request signature envelope (`X-I6-Signature` header) using the stored robot key (or an arbitrary private key).
- `i6 dataset keys delete <jti>`: Calls server API (`/api/userz/dataset-token-refresh/delete`) to revoke a machine robot key on the server.

### D. Data Ingestion (`i6 ingest`)
- `i6 ingest upload <dataset> <table> <files...>`:
  - Supports single/multiple files or shell globbing.
  - Automatically acquires upload URLs using `c.Ingest().GetUrl()`.
  - Concurrently uploads files via `c.Ingest().Upload()`.
  - Renders live interactive terminal progress bars (multi-stream overall and per-file progress) with non-interactive TTY fallback to stderr.
- `i6 ingest get-url <dataset> <table> [--partitions key=val] [--amount N]`: Generates raw signed PUT upload URLs.
- `i6 ingest summary <dataset> <ingest-id>`: Checks batch lake validation status and summary.

### E. Analytics & Dashboards (`i6 dash`)
- `i6 dash query <query-name> [--dataset <dataset>] [--partitions key=val,...]`:
  - Downloads and caches analytical SQLite databases via `c.Dash()`.
  - Executes stored named queries and outputs formatted results (table, TSV, or JSON).

### F. Inference & Scoring (`i6 infer`)
- `i6 infer run <model-name> [--dataset <dataset>] [--input <json-or-file>]`:
  - Executes the full 3-step inference pipeline (candidate generation -> remote TF scoring -> hydration).
- `i6 infer metadata <model-name>`: Fetches model schema and signature metadata.

### G. Pipeline Orchestration (`i6 pipe`)
Manages data pipelines and recurring scheduled workloads via `apipipez`:
- `i6 pipe list [--dataset <dataset>]`: Lists all pipeline executions and their current statuses (`/api/pipez/list`).
- `i6 pipe get <pipeline-id>`: Retrieves execution status, step timings, and failure diagnostics (`/api/pipez/get`).
- `i6 pipe start <batch-ingest|main> [--dataset <dataset>]`: Triggers an immediate execution of a pipeline:
  - `batch-ingest`: Starts batch ingestion processing (`/api/pipez/start/batch-ingest`).
  - `main`: Starts main recommendation pipeline processing (`/api/pipez/start/main`).
- `i6 pipe scheduler activate <batch-ingest|main> [--dataset <dataset>]`: Activates recurring scheduled execution for the target pipeline.
- `i6 pipe scheduler deactivate <batch-ingest|main> [--dataset <dataset>]`: Deactivates recurring scheduled execution for the target pipeline.

### H. SdkFunc Workflows (`i6 funcs`)
Zero-codegen dynamic execution for all registered Infinity6 workflow functions (analytical queries and inference pipelines):
- `i6 funcs list`: Lists all available SdkFuncs (ID, execution type: `client_query` vs `client_infer`, and product domain).
- `i6 funcs describe <func-name>`: Displays input and output schemas, partition requirements, and functional description.
- `i6 funcs run <func-name> [--input <json|file|->] [--dataset <dataset>]`: Dynamically dispatches the workflow function:
  - If `client_query` $\rightarrow$ routes to `c.Dash().QueryRaw()` and formats the analytical result.
  - If `client_infer` $\rightarrow$ routes to `c.Infer().Run()` (3-step candidate generation, remote scoring, and hydration).

### I. Platform APIs (`i6 apis`)
Dynamic invocation of all 60+ registered Infinity6 HTTP APIs with automatic endpoint routing and bearer authentication:
- `i6 apis list [--service userz|decsuite|tfinfer]`: Lists all registered APIs (API ID, HTTP Method, Path, and target service).
- `i6 apis describe <api-id>`: Displays parameter, query, header, and payload schemas for the API.
- `i6 apis call <api-id> [--params key=val,...] [--query key=val,...] [--payload <json|file|->]`: Dynamically executes the API call through `c.ClientForPath()`, automatically injecting authentication headers and formatting the response.

### J. Configuration (`i6 config`)
- `i6 config list`: Lists all local configuration profiles discovered in `$I6_CLI_CONFIG_DIR`, showing default dataset, endpoint URLs, and token store status.
- `i6 config init [--box prod|staging|sandbox] [--dataset <dataset>] [-f/--force]`: Initializes the active profile directory and writes `config.json` with service endpoint presets and optional default dataset. Protected with `-f/--force` against accidental overwrites.
- `i6 config show`: Displays active profile, directory paths, default dataset, and resolved service endpoint URLs.
- `i6 config set <key> <value>`: Updates a configuration field (e.g. `box`, `default_dataset`, `userz_url`, `decsuite_url`, `tfinfer_url`) in the active profile's `config.json`.

### K. Shell Completion & Utilities (`i6 completion`)
Provided natively by Cobra's default autocompletion engine:
- `i6 completion [bash|zsh|fish|powershell]`: Generates shell completion scripts for the specified shell.
- Includes dynamic completion support for flags, profiles, and subcommands.

---

## 4. Configuration & State Layout

All CLI configuration, credentials, and local caches are organized per profile and managed through `github.com/infinity6-ai/gox/commonz/configz`:

- **Config Base Directory**: `PropConfigDir` (`I6_CLI_CONFIG_DIR`, defaults to `$HOME/.config/i6/`).
- **Active Profile Selection Hierarchy**:
  1. Persistent root flag: `--profile <name>` (or `-p <name>`)
  2. Configuration property: `PropConfigProfile` (`I6_CLI_CONFIG_PROFILE`)
  3. Default profile: `"main"`
- **Inline JSON Override**: `PropConfigJSON` (`I6_CLI_CONFIG_JSON`, optional inline JSON string for zero-disk CI/CD execution).
- **Dataset Override**: `PropDataset` (`I6_DATASET`, optional default dataset override).

```
$I6_CLI_CONFIG_DIR/ (defaults to $HOME/.config/i6/)
└── [profile]/ (defaults to "main", selected via --profile / I6_CLI_CONFIG_PROFILE)
    ├── config.json          # Profile configuration: default_dataset, custom service URLs
    ├── tokens.db            # Persistent SQLite token store via tokenstore.NewSqlite (WAL, 0600)
    └── cache/               # Profile-scoped SQLite blobs, queries, models, ETags (fszcache)
        ├── dash/
        └── infer/
```

### Profile-Scoped SQLite Path Resolution
In `sdkgo/tokenstore/sqlite.go`, `tokenstore.DefaultSqliteUrl()` defaults to `$HOME/.config/i6/tokens.db`. To guarantee profile isolation, `i6cli`'s `cliconfig` packages MUST explicitly format the target URL rather than passing `nil`:
```go
dbPath := filepath.Join(profileDir, "tokens.db")
targetUrl, err := urlz.Parse("file://" + filepath.ToSlash(dbPath))
store, err := tokenstore.NewSqlite(targetUrl)
```

### Dataset Resolution Order
When a command requires a target dataset, it resolves the dataset identifier using the following priority:
1. Command flag: `--dataset <name>`
2. Profile configuration: `default_dataset` in `config.json`
3. Environment override: `I6_DATASET`
4. If unresolved: fail fast with exit code `2` and an actionable error:
   `"error: dataset required. Provide --dataset <name> or configure default via 'i6 config init --dataset <name>'"`

### Core Invariant: Mutually Exclusive Profile Identity
Each CLI configuration profile represents either a **Human User** or a **Machine Service Account**, never both:
$$\text{Profile Root Credential} \in \{\text{User Session}, \text{Dataset Refresh Key}, \emptyset\}$$

- User session credentials (`i6_user_user_token`) and machine robot keys (`i6_rf:<dataset>`) cannot coexist within the same profile.
- If a profile has an active user session, running `i6 tokens dataset-refresh set` returns an error directing the user to use a separate profile (e.g. `i6 --profile worker ...`).
- If a profile has an active robot key, running `i6 auth login` returns an error directing the user to use a separate profile (e.g. `i6 --profile developer ...`).
- Derived short-lived dataset bearer tokens (`i6_ds:<dataset>`) can be cached for multiple datasets within any active profile.

### Configuration Hierarchy & Resolution Order
1. **Zero-Disk Inline JSON**: If `I6_CLI_CONFIG_JSON` is populated, parse in-memory immediately using `tokenstore.NewMemory()` (bypasses filesystem reads, ideal for CI/CD runners).
2. **Environment & Profile Defaults**:
   - Target Directory: `$I6_CLI_CONFIG_DIR/<profile>/`
   - Config file: `$I6_CLI_CONFIG_DIR/<profile>/config.json`
   - Token store: `$I6_CLI_CONFIG_DIR/<profile>/tokens.db` (via explicit `file://` URL)
   - Cache store: `$I6_CLI_CONFIG_DIR/<profile>/cache/`

### Output Formatting (`--format [table|tsv|json|quiet|raw]`)
Commands producing data support the `--format` flag (with per-command contextual defaults):

- **`table`** (Default for `list`, `get`, `status`, `dash query`): Aligned ASCII tables for human terminal readability.
- **`tsv`**: Tab-separated values with headers, ideal for standard Unix pipelines (`cut`, `awk`, `xargs`).
- **`json`**: Formatted structured JSON payloads, ideal for `jq` and programmatic ingestion.
- **`quiet`**: Minimal output; prints only the primary ID or count (e.g. `dtr_abc123`), suppressing all tables and headers.
- **`raw`** (Default for `expose`, `export`, and binary outputs): Direct, unquoted string or byte output (e.g. raw JWT or Ed25519 seed string ready for `export TOKEN=$(...)`).

### Unix Composability & Stream Separation
To ensure clean piping and script integration:
- **Strict Stream Separation**:
  - `os.Stdout`: Reserved **exclusively** for command data payloads (tables, TSV, JSON, raw tokens).
  - `os.Stderr`: Used for all diagnostic, interactive, and UI output (progress bars, spinners, automatic token refresh logs, informational messages).
  - Enables clean composability: `export TOKEN=$(i6 tokens dataset-access expose)` without corrupted tokens.
- **Diagnostic Logging via `gox/commonz/logz`**:
  - All internal logging across command packages uses `github.com/infinity6-ai/gox/commonz/logz` (`type tlogger logz.Type; var logger = logz.Create(tlogger(true))`).
  - Automatically writes exclusively to `os.Stderr`.
  - Supports automatic TTY color detection (colors suppressed when `os.Stderr` is redirected or in CI pipelines via `term.IsTerminal` / `I6_LOGZ_NO_COLOR=true`).
  - Supports structured logging (`I6_LOGZ_STRUCTURED=true`) for enterprise log collectors.
- **TTY Detection & Graceful Fallback**:
  - When `os.Stderr` is not a terminal (e.g. inside Docker or CI runners `!isatty(os.Stderr.Fd())`), interactive progress bars degrade gracefully into throttled periodic log lines.

### Standardized Predictable Exit Codes
Scripts and CI/CD pipelines can rely on deterministic exit codes:
- `0`: **Success**
- `1`: **General runtime error** (server 500, network failure, unhandled error)
- `2`: **Usage error** (invalid flags, missing arguments, invalid syntax)
- `3`: **Authentication error** (unauthenticated, missing token, expired refresh session)
- `4`: **Authorization error** (insufficient grants, forbidden operation, attenuated permission mismatch)
- `5`: **Not found** (dataset, table, store, or model does not exist)

---

## 5. Mono-Repo File Structure

The public API boundary is strictly encapsulated: **`cmdcli` is the only public package** exported by `i6cli`. All command and subcommand implementations live within dedicated, modular packages under `internal/cmd/`:

```
i6cli/
├── go.mod
├── go.sum
├── main.go                       # Minimal entrypoint calling cmdcli.Execute()
├── PLAN.md                       # This architectural plan
├── README.md                     # CLI installation and usage guide
├── cmdcli/                       # ONLY public package (Prepare, Execute)
│   ├── cmdcli.go
│   └── cmdcli_unit_test.go
└── internal/
    ├── cmd/                      # Command taxonomy packages (internal)
    │   ├── cmdauth/              # i6 auth (cmdlogin, cmdtoken, cmdlogout, cmdstatus, cmdinspect)
    │   │   ├── cmdlogin/
    │   │   ├── cmdtoken/
    │   │   ├── cmdlogout/
    │   │   ├── cmdstatus/
    │   │   └── cmdinspect/
    │   ├── cmdconfig/            # i6 config (cmdinit, cmdlist, cmdshow, cmdset)
    │   │   ├── cmdinit/
    │   │   ├── cmdlist/
    │   │   ├── cmdshow/
    │   │   └── cmdset/
    │   ├── cmdtokens/            # i6 tokens (cmdcheck, cmdinspect, cmduser, cmddatasetrefresh, cmddatasetaccess, cmdraw)
    │   │   ├── cmdcheck/
    │   │   ├── cmdinspect/
    │   │   ├── cmduser/
    │   │   ├── cmddatasetrefresh/
    │   │   ├── cmddatasetaccess/
    │   │   └── cmdraw/
    │   ├── cmddataset/           # i6 dataset (cmdlist, cmdget, cmdtable, cmdstore, cmdtokens)
    │   ├── cmdingest/            # i6 ingest (cmdupload, cmdgeturl, cmdsummary)
    │   ├── cmdinfer/             # i6 infer (cmdrun, cmdmetadata)
    │   ├── cmddash/              # i6 dash (cmdquery)
    │   ├── cmdpipe/              # i6 pipe (cmdlist, cmdget, cmdstart, cmdscheduler)
    │   ├── cmdfuncs/             # i6 funcs (cmdlist, cmddescribe, cmdrun)
    │   └── cmdapis/              # i6 apis (cmdlist, cmddescribe, cmdcall)
    ├── cliconfig/                # Config file parsing, profile resolution, and SQLite path resolution
    │   ├── cliconfig.go
    │   └── cliconfig_unit_test.go
    ├── outfmt/                   # Table, TSV, JSON, quiet, and raw output formatters
    │   ├── outfmt.go
    │   └── outfmt_unit_test.go
    ├── exitcode/                 # Standardized CLI exit codes (0-5)
    │   ├── exitcode.go
    │   └── exitcode_unit_test.go
    └── ui/                       # Terminal spinners, progress bars (TTY-aware stderr)
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

- [x] **Phase 1: Foundation & Scaffolding**
  - [x] Scaffold `i6cli/` directory, `go.mod`, `main.go`, and `cmdcli/cmdcli.go`.
  - [x] Integrate with `build.sh` (`cmd_comps_list`), `comp.sh`, and root `go.work`.
  - [x] Verify test execution via `./comp.sh i6cli test` and `./build.sh test`.
  - [x] Implement `internal/cliconfig` (profile resolution, profile-isolated SQLite URL, dataset fallback, and config.json).
  - [x] Implement `internal/exitcode` (standardized exit codes 0-5).
  - [x] Implement `internal/outfmt` (table, tsv, json, quiet, and raw formatters with strict stdout routing).
  - [x] Implement persistent `--profile` / `-p` root flag and enable Cobra native shell autocompletion (`bash`, `zsh`, `fish`, `powershell`).

- [x] **Phase 2: Authentication & Local Profile Tokens (`auth`, `tokens`, `config`)**
  - [x] Implement `i6 config init [--box prod|staging|sandbox] [--dataset <dataset>] [-f/--force]`, `i6 config list`, `i6 config show`, and `i6 config set <k> <v>`.
  - [x] Implement `i6 auth login` (Device code flow with polling; enforces mutual exclusivity).
  - [x] Implement `i6 auth token [<token-string> | -] [--offline]` (universal token ingestion with auto-detection).
  - [x] Implement `i6 auth inspect [<token-string> | -] [--refreshed]` (deep diagnostics of arbitrary user JWT or active session).
  - [x] Implement `i6 auth logout` and `i6 auth status`.
  - [x] Implement `i6 tokens check [--refreshed]` (audits all credentials stored in active profile's tokens.db via /api/userz/token/check).
  - [x] Implement `i6 tokens inspect [--refreshed]` (deep diagnostic inspection of all profile credentials).
  - [x] Implement `i6 tokens user [get|expose|export|attenuation|delete]` (human user session root credential; attenuation mints attenuated dataset token).
  - [x] Implement `i6 tokens dataset-refresh [get|export|set|attenuation|delete]` (machine robot root credential; attenuation mints attenuated dataset token):
    - `bounded [sign-request]` (uses stored robot key to mint PoP token or sign HTTP requests)
  - [x] Implement `i6 tokens dataset-access [--dataset <dataset>]` (multiple derived cached tokens):
    - `list` (lists all cached dataset tokens in tokens.db)
    - `get [--refreshed]` (metadata only; refreshes if needed)
    - `expose [--refreshed]` (raw Bearer JWT; refreshes if needed)
    - `attenuation [--grants <g>] [--exp <s>] [--solution <sol>]` (derives further attenuated token)
    - `check [--refreshed]` (validates against server; refreshes if needed)
    - `delete` (evicts cached bearer token)
  - [x] Implement `i6 tokens raw` (stateless in-memory operations):
    - `check <token-string | ->` (statelessly debugs arbitrary token)
    - `dataset-bearer --refresh-token <key|-> --dataset <dataset>`
    - `attenuation [<token|->] [--grants <mask|names>] [--solution <sol>] [--exp <seconds>]`
    - `sign-request --key <privkey|-> --token <token|-> --method <METHOD> --path <PATH> [--body <str|->]`

- [ ] **Phase 3: Data Lake, Ingestion & Pipelines (`dataset`, `ingest`, `pipe`)**
  - [ ] Implement `i6 dataset list`, `get`, `table`, and `store`.
  - [x] Implement `i6 dataset keys [--dataset <dataset>]` for remote server machine robot key management:
    - [x] `list` (lists all registered machine robot keys by JTI on the server)
    - [x] `create [--solution <sol>] [--grants <g>] [--annotation <text>] [--save]` (registers machine key on server)
    - [x] `reqsign --dataset <dataset> --method <METHOD> --path <PATH>` (signs request using machine robot key)
    - [x] `delete <jti>` (revokes machine robot key on the server by JTI)
  - [ ] Implement `i6 ingest upload` with multi-stream parallel worker execution and live terminal progress bars (TTY-aware fallback).
  - [ ] Implement `i6 ingest get-url` and `summary`.
  - [ ] Implement `i6 pipe list`, `get`, `start <batch-ingest|main>`, and `scheduler <activate|deactivate> <batch-ingest|main>`.

- [ ] **Phase 4: Analytics & ML Inference (`dash`, `infer`)**
  - [ ] Implement `i6 dash query` with analytical SQLite bundle caching (`fszcache`) and formatted output (`table`, `tsv`, `json`).
  - [ ] Implement `i6 infer run` (3-step candidate generation -> remote TF scoring -> hydration) and `metadata`.

- [ ] **Phase 5: Dynamic Reflection & Platform APIs (`funcs`, `apis`)**
  - [ ] Implement `i6 funcs [list|describe|run]` for dynamic zero-codegen SdkFunc execution (`client_query` -> `c.Dash()`, `client_infer` -> `c.Infer()`).
  - [ ] Implement `i6 apis [list|describe|call]` for dynamic inspection and execution of all 60+ platform HTTP APIs via `schemaz.Desc` and `c.ClientForPath()`.

- [ ] **Phase 6: Release Packaging & Automation**
  - [ ] Add static binary cross-compilation matrix to `.github/workflows/release.yml`.
  - [ ] Add `install.sh` shell script for instant installation.
  - [ ] Write `i6cli/README.md` with complete command reference and quickstart examples.
