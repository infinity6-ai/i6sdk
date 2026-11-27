# i6cli Architectural Plan & Roadmap

## 1. Executive Summary

`i6cli` provides **`i6`**, the official standalone static binary command-line interface for the Infinity6 platform.

### Key Objectives
- **Zero-Dependency Customer Tooling**: Customers can run `i6` on any machine, container, or CI/CD runner without installing Go, Python, or Node.js.
- **100% Feature Parity with SDKs**: Anything possible via `sdkgo`, `sdkpy`, or `sdkweb` can be performed from the terminal or shell scripts.
- **Built on `sdkgo`**: `i6cli` acts as a presentation and UX layer over `github.com/infinity6-ai/i6sdk/sdkgo`, eliminating duplicate protocol, auth, cache, or ingestion code.
- **True Static Binaries (`CGO_ENABLED=0`)**: Pure Go implementation using `modernc.org/sqlite` and Go standard library crypto/networking. Cross-compiles for Linux (`amd64`, `arm64`), macOS (`arm64`, `amd64`), and Windows without dynamic `libc`/`glibc` dependencies.
- **Script-Friendly & CI-Ready**: Flexible output formatting (`--format=table|tsv|json|quiet|raw`), predictable exit codes, stdin/stdout piping, and non-interactive execution modes.

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
- `i6 auth login`: Interactive browser login screen (OAuth 2.0 Device Code flow). Displays authorization code and URL, polls server, and saves full user session in `tokens.db`.
- `i6 auth token [<token-string> | -] [--offline]`: Universal token ingestion command. Automatically detects the token type via local formatting and `/api/userz/token/check` (with `--offline`, skips network calls and ingests locally):
  - **Ed25519 Dataset Robot Key**: Ingests as machine robot root credential (`i6_rf:<dataset>`).
  - **User Refresh Token**: Refreshes via `/api/userz/auth/refresh` to establish a full user session (`i6_user_user_token`).
  - **Standalone User Access Token**: Saved as a non-refreshable user session in `i6_user_user_token`.
  - **Standalone Dataset Bearer Token**: Saved as a non-refreshable dataset access token in `i6_ds:<dataset>`.
- `i6 auth logout`: Revokes session token family on server via `UserAuthLogout` and wipes local credentials.
- `i6 auth status` / `i6 auth me`: Displays active authentication identity (User Session vs. Robot Service Account vs. Standalone Bearer vs. Unauthenticated), tenant, dataset, and token expiration.

### B. Local Profile Token Management (`i6 tokens`)

Manages the active credentials stored in the local profile's `tokens.db`.

*Root Credential Invariant*: A profile has either a **Human User Session** (`tokens user`) OR a **Machine Robot Key** (`tokens dataset-refresh`), never both. This eliminates ambiguity in identity, signing, and permissions.

*Security Invariant*: `list` and `get` commands **never** display sensitive `secret` strings, tokens, or private keys. Revealing secrets is strictly restricted to explicit `expose` and `export` commands.

#### Profile Tokens Health Audit (`i6 tokens check [--refreshed]`)
- `i6 tokens check [--refreshed]`: Performs a comprehensive diagnostic audit of **all tokens stored in the active profile's `tokens.db`**:
  - Validates the active root credential (`tokens user` or `tokens dataset-refresh`).
  - Iterates through and validates all cached derived tokens (`tokens dataset-access`).
  - Calls `/api/userz/token/check` for each token and displays an aggregated diagnostic report (credential type, dataset/user ID, server validity status, remaining TTL, grants).
  - With `--refreshed`, automatically refreshes any expired or missing access tokens before reporting.

#### 1. Human User Root Credential (`i6 tokens user`)
Controls the authenticated human user session in `tokens.db`:
- `i6 tokens user get`: Displays user profile metadata (user_id, tenant, expiration; omits tokens).
- `i6 tokens user expose`: Prints raw active JWT access token to stdout (for curl / scripts).
- `i6 tokens user export`: Prints the active user refresh token to stdout (for backup or environment transfer).
- `i6 tokens user delete`: Revokes session on server and removes credentials from `tokens.db` (logout).

#### 2. Machine Robot Root Credential (`i6 tokens dataset-refresh`)
Controls the machine service account Ed25519 refresh key (mutually exclusive with `tokens user`):
- `i6 tokens dataset-refresh get`: Displays robot key metadata for the dataset (safely omits the secret key).
- `i6 tokens dataset-refresh export [--out <path|->]`: Prints the active Ed25519 dataset refresh token to stdout or a file.
- `i6 tokens dataset-refresh set <key|file|-> [--dataset <dataset>]`: Sets or replaces the robot key for this profile.
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
- `i6 tokens dataset-access check [--refreshed]`: Cryptographically validates the token against `/api/userz/token/check`.
- `i6 tokens dataset-access delete`: Evicts the cached bearer token for `--dataset` from local `tokens.db`.

#### 4. Stateless Zero-Disk Operations (`i6 tokens raw`)
Pure in-memory cryptographic and token exchange operations that never read or write to `tokens.db`:
- `i6 tokens raw check <token-string | ->`: Validates and debugs an arbitrary token string or stdin pipe via `/api/userz/token/check` without saving or affecting `tokens.db`.
- `i6 tokens raw dataset-bearer --refresh-token <key|-> --dataset <dataset>`: Statelessly signs an Ed25519 `RefreshAssertion`, exchanges it on the server, and prints the raw short-lived bearer access token to stdout without touching disk.
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

#### Remote Server Token Registry (`i6 dataset tokens [--dataset <dataset>]`)
Server-side token management where a dataset can have multiple registered tokens with distinct JTIs:
- `i6 dataset tokens list [--type all|ingestion|recsys|refresh]`: Calls server API to list all registered tokens for the dataset (showing JTI, labels, expiration; safely omits secrets).
- `i6 dataset tokens create [--type ingestion|recsys|refresh] [--annotation <text>] [--save]`: Calls server API to mint a new token on the server. Dumps one-time secret key; with `--save`, automatically stores it into active profile's `tokens.db`.
- `i6 dataset tokens delete <token-id-or-jti>`: Calls server API to revoke a specific registered token by its ID or JTI.

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

### Core Invariant: Mutually Exclusive Profile Identity
Each CLI configuration profile represents either a **Human User** or a **Machine Service Account**, never both:
$$\text{Profile Root Credential} \in \{\text{User Session}, \text{Dataset Refresh Key}, \emptyset\}$$

- User session credentials (`i6_user_user_token`) and machine robot keys (`i6_rf:<dataset>`) cannot coexist within the same profile.
- If a profile has an active user session, running `i6 tokens dataset-refresh set` returns an error directing the user to use a separate profile (e.g. `I6_CLI_CONFIG_PROFILE=worker`).
- If a profile has an active robot key, running `i6 auth login` returns an error directing the user to use a separate profile (e.g. `I6_CLI_CONFIG_PROFILE=developer`).
- Derived short-lived dataset bearer tokens (`i6_ds:<dataset>`) can be cached for multiple datasets within any active profile.

### Configuration Hierarchy & Resolution Order
1. **Zero-Disk Inline JSON**: If `I6_CLI_CONFIG_JSON` is populated, parse in-memory immediately using `tokenstore.NewMemory()` (bypasses filesystem reads, ideal for CI/CD runners).
2. **Environment & Profile Defaults**:
   - Target Directory: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/`
   - Config file: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/config.json`
   - Token store: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/tokens.db` (via `tokenstore.NewSqlite`)
   - Cache store: `$I6_CLI_CONFIG_DIR/$I6_CLI_CONFIG_PROFILE/cache/`

Specific subcommands take positional arguments or command-scoped flags, preventing flag pollution on the root command.

### Output Formatting (`--format [table|tsv|json|quiet|raw]`)
Commands producing data support the `--format` flag (with per-command contextual defaults):

- **`table`** (Default for `list`, `get`, `status`, `dash query`): Aligned ASCII tables for human terminal readability.
- **`tsv`**: Tab-separated values with headers, ideal for standard Unix pipelines (`cut`, `awk`, `xargs`).
- **`json`**: Formatted structured JSON payloads, ideal for `jq` and programmatic ingestion.
- **`quiet`**: Minimal output; prints only the primary ID or count (e.g. `dtr_abc123`), suppressing all tables and headers.
- **`raw`** (Default for `expose`, `export`, and binary outputs): Direct, unquoted string or byte output (e.g. raw JWT or Ed25519 seed string ready for `export TOKEN=$(...)`).

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
    ├── outfmt/                   # Table, TSV, JSON, quiet, and raw output formatters
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
  - [ ] Implement `internal/outfmt` (table, tsv, json, quiet, and raw formatters).

- [ ] **Phase 2: Authentication & Local Profile Tokens (`auth`, `tokens`, `config`)**
  - [ ] Implement `i6 config init [--box sandbox|staging|prod] [-f/--force]` and `i6 config show`.
  - [ ] Implement `i6 auth login` (Device code flow with polling; enforces mutual exclusivity).
  - [ ] Implement `i6 auth token [<token-string> | -] [--offline]` (universal token ingestion with auto-detection).
  - [ ] Implement `i6 auth logout` and `i6 auth status`.
  - [ ] Implement `i6 tokens check [--refreshed]` (audits all credentials stored in active profile's tokens.db via /api/userz/token/check).
  - [ ] Implement `i6 tokens user [get|expose|export|delete]` (human user session root credential; expose prints access token, export prints refresh token).
  - [ ] Implement `i6 tokens dataset-refresh [get|export|set|delete]` (machine robot root credential; export prints robot refresh key):
    - `bounded [create|sign-request]` (uses stored robot key to mint PoP token or sign HTTP requests)
  - [ ] Implement `i6 tokens dataset-access [--dataset <dataset>]` (multiple derived cached tokens):
    - `list` (lists all cached dataset tokens in tokens.db)
    - `get [--refreshed]` (metadata only; refreshes if needed)
    - `expose [--refreshed]` (raw Bearer JWT; refreshes if needed)
    - `check [--refreshed]` (validates against server; refreshes if needed)
    - `delete` (evicts cached bearer token)
  - [ ] Implement `i6 tokens raw` (stateless in-memory operations):
    - `check <token-string | ->` (statelessly debugs arbitrary token)
    - `dataset-bearer --refresh-token <key|-> --dataset <dataset>`
    - `bounded-create --token <token|-> --key <privkey|->`
    - `sign-request --key <privkey|-> --method <METHOD> --path <PATH> [--body <str|->]`

- [ ] **Phase 3: Ingestion & Dataset Management (`ingest`, `dataset`)**
  - [ ] Implement `i6 ingest upload` with multi-stream parallel worker execution and live terminal progress bars.
  - [ ] Implement `i6 ingest get-url` and `summary`.
  - [ ] Implement `i6 dataset list`, `get`, `table`, and `store`.
  - [ ] Implement `i6 dataset tokens [--dataset <dataset>]` for remote server token management:
    - `list [--type all|ingestion|recsys|refresh]` (lists all registered tokens by JTI on the server)
    - `create [--type ingestion|recsys|refresh] [--annotation <text>] [--save]` (mints token on the server)
    - `delete <token-id-or-jti>` (revokes token on the server by ID or JTI)

- [ ] **Phase 4: Analytics, Inference & Typed Workflows (`dash`, `infer`, `pipe`, `run`)**
  - [ ] Implement `i6 dash query` with table output.
  - [ ] Implement `i6 infer run` and `metadata`.
  - [ ] Implement `i6 pipe list`, `get`, `start`.
  - [ ] Implement `i6 run` for typed `SdkFuncs` dispatch.

- [ ] **Phase 5: Release Packaging & Automation**
  - [ ] Add static binary cross-compilation matrix to `.github/workflows/release.yml`.
  - [ ] Add `install.sh` shell script for instant installation.
  - [ ] Write `i6cli/README.md` with complete command reference and quickstart examples.
