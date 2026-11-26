# SDK Docs Generator — Plan

We generate a customer-facing `sdkpy/SDK.md` from the docstrings we already have. It is built by `./build.sh docs`, with no LLM and nothing written by hand. The design works for any language: adding one means writing one small extractor script.

**How it works, in one line:** a per-language script reads the code and prints JSON in one shared format; one Go command turns that JSON into Markdown.

```
sdkpy/docgen.py  ──JSON (stdin)──►  forgez docs  ──►  sdkpy/SDK.md
(Python ast, stdlib only)           (Go: check + render)
```

- **Scope of v1:** sdkpy only.
- **Customer experience does not change.** No public name, signature or behaviour changes. We only add `SDK.md` and edit docstrings and comments.
- **Missing docs fail the build from day 1.**
- **Out of scope:** an HTML site, guides and tutorials (they stay in `README.md`), and other SDKs (see "Adding a language").

---

## Decisions

| # | Decision | Why |
| --- | --- | --- |
| 1 | The extractor reads **all** code, generated (`sdkapis/`) included | One path for everything. No Go-side spec reader, per-language adapters or drift tests |
| 2 | One Go package, `forgez/forgedoc`, holds the format, the check and the renderer | One place to read, one place to test |
| 3 | The renderer is plain Go (`strings.Builder` + `fmt`), not `text/template` | About 200 lines anyone can follow; no template funcs |
| 4 | The extractor is piped into the renderer from a shell script | No `go:embed`, no temp files, no exec wrapper |
| 5 | The extractor lives in the SDK (`sdkpy/docgen.py`), outside `src/` | It ships with the code it reads, and isn't in the wheel |
| 6 | Request and response fields are written inline under each method, flattened with dotted names | Customers see what to send and what comes back without jumping around |
| 7 | `SDK.md` is committed next to the code; the release regenerates it and refuses a dirty tree | Reviewable in PRs, readable on GitHub, and needs no extra check code |
| 8 | HTTP method and path are not shown | The Python code doesn't have them, and customers don't see them today either |
| 9 | Methods that return an `internal/` class are hidden (`Services.fs()`, `downloader()`, `invoker()`) | They are plumbing, not customer API |
| 10 | Constructors (`__init__`) are documented only for classes in `__all__` | Customers get `Auth`, `Infer`, … from `client.auth()` and never build them |
| 11 | Wrapper fields (`params`, `payload`, `query`, `headers`, `input`, `partitions`, `output`) need no doc | Their type is a record that is flattened into dotted rows, so they never appear as a row |
| 12 | camelCase twins (`modelName` next to `model_name`) are merged into one row, "also `modelName`", like method aliases | Customers see one entry per thing, and the twin needs no separate doc |
| 13 | Overriding methods get their own docstrings; no doc inheritance | No extra rule in the extractor; only 8 short docstrings |
| 14 | `bitmapz` stays public and documented | It is in `__all__` and customers use it for grants |
| 15 | Sections have friendly titles ("Inference", "API reference") taken from the first line of each module docstring | Customers think in features, not import paths; we control the words in the code |
| 16 | Method headings show how customers call them: `client.infer().run`, not `Infer.run` | Copy-paste ready; matches the README's `client = i6sdk.create_client()` |
| 17 | Sections follow the customer journey: client, configuration, auth, tokens, ingestion, inference, dashboards, API reference, workflows, then errors and helpers | It falls out of the walk order, so no config file is needed |

---

## What we do first

1. **Agree on what is public** (Python):
    - Names in `sdkpy/src/i6sdk/__init__.py` `__all__`.
    - Classes named in the parameter or return annotations of a public method, followed transitively. `Services.apis() -> Apis` and `Services.infer() -> Infer` are how `Apis`, `Infer`, `Tokens`, `Auth`, `Ingest`, `Dash` and `SdkFuncs` get in, and `Infer.run(opts: InferQueryOptions)` is how `InferQueryOptions` gets in.
    - Never: names starting with `_`; dunders, except `__init__` on classes in `__all__`; classes defined under `internal/` unless `__all__` re-exports them (`bitmapz`); and methods whose return type is such an internal class.
2. **Agree on the docstring format.** It's the one the generated code already uses:
    - Module docstring: the first line is the section title, short and customer-facing ("Inference"), and the rest is the section intro.
    - The text before the first `:field:` line is the description, kept as Markdown.
    - `:param name: text` documents a parameter, and `:return: text` documents the return value.
    - Record fields (`TypedDict`, `@dataclass`) are documented with the `#` comment line(s) directly above the field. In generated types, the `[Go: …] - ` prefix is stripped, so customers never see Go type names.
3. **Fix the docs listed in "Doc fixes" below.** Once Steps 1 and 2 exist, `./build.sh docs` prints the exact list (`file:line: name: missing doc`). We fix all of it before Step 3 turns the check on in the release.

---

## Doc fixes

This list comes from a quick scan of the current code on 2026-11-25 using the rules above. The check in Step 1 gives the exact list. Docstrings and comments are the only changes: no names, signatures or behaviour.

**Module docstrings**: the first line becomes the section title. The current one-line docstring moves to the second paragraph as the intro. Proposed titles, in the order they appear in `SDK.md`:

| File | Title |
| --- | --- |
| `__init__.py` | Getting a client |
| `services.py` | Client |
| `config.py` | Configuration |
| `auth/auth.py` | Authentication |
| `tokens/tokens.py` | Tokens |
| `ingest/ingest.py` | Data ingestion |
| `infer/infer.py` | Inference |
| `dash/dash.py` | Dashboards |
| `sdkapis/apis.py` (generated: `forgepybundle.go`, line 91) | API reference |
| `sdkapis/sdkfuncs.py` (generated: `forgepybundle.go`, line 103) | Workflows |
| `sdkapis/types.py` (generated: `forgepyschema.go`, line 279) | API types (only types not already inlined under a method) |
| `exceptions.py` | Errors |
| `invoker.py` | Responses |
| `tokens/storage.py` | Token storage |
| `grants.py` | Permissions |
| `internal/bitmapz.py` | Bitmasks |

**Generated code (good shape).** 766 fields; every leaf field in the API specs already has a description. Only two hardcoded types lack docs:

| Where | What to fix |
| --- | --- |
| `forgez/forgepy/forgepyschema/forgepyschema.go` (the hardcoded `FileStat` block) | `#` comments on 7 fields: `url`, `content_type`, `md5`, `size`, `etag`, `created_at`, `updated_at` |
| `ServiceOptions` | Nothing in the generator: the duplicate in `sdkapis/types.py` is skipped in favour of the one in `invoker.py` (below) |

**Handwritten methods**: missing `:param` lines.

| File | Methods |
| --- | --- |
| `__init__.py` | `create_client` (5 params) |
| `services.py` | `Services.__init__` (no docstring), `set_config` |
| `auth/auth.py` | `get_valid_access_token`, `sign_in`, `sign_up`, `device_code_start`, `device_code_poll` |
| `tokens/tokens.py` | `set_user_token`, `get/set/remove_dataset_refresh_token`, `is_access_token_expired`, `is_refresh_token_expired`, `get_valid_access_token`, `get/set/remove_dataset_token`, `is_dataset_token_expired`, `refresh_dataset_access_token` (5 params), `get_valid_dataset_access_token`; the `storage` property (no docstring) |
| `tokens/storage.py` | `TokenStorage.get/set/remove_item`; `MemoryTokenStorage` and `FileTokenStorage`: `__init__` and `get/set/remove_item` (no docstrings) |
| `infer/infer.py` | `update`, `run`, `select_results` (`opts`) |
| `ingest/ingest.py` | `get_url` (5 params), `upload` (9 params; `onProgress` is merged into `on_progress`, and `signal` is a different name, so it gets its own doc) |
| `dash/dash.py` | `query` (`opts`) |
| `exceptions.py` | `I6Error.__init__`, `I6ApiError.__init__` (no docstrings) |
| `invoker.py` | `Resp.__init__` (no docstring) |
| `internal/bitmapz.py` | `contains_all`, `contains_any`, `split`, `set_bits`, `serialize`, `deserialize` |

**Handwritten records**: missing `#` comments on fields.

| File | Type (fields) |
| --- | --- |
| `config.py` | `I6SdkConfig` (7) |
| `invoker.py` | `ServiceOptions` (4) |
| `infer/infer.py` | `InferQueryOptions` (9; the 5 camelCase twins are merged and need no doc) |
| `dash/dash.py` | `DashQueryOptions` (6) |
| `ingest/ingest.py` | `UploadProgress` (6), `UploadResult` (2), `SignedUrlResult` (3) |
| `tokens/tokens.py` | `UserToken` (5), `DatasetToken` (4) |

Plus a `#` comment above each public constant in `config.py` (`DEFAULT_CONFIG`, `STAGING_SERV_URLS`, `SANDBOX_SERV_URLS`). Records that are only reachable through hidden methods (`FileRef`, `InvokeOptions`, `ApiErrorDetail`) are not documented, so they need no fixes.

**Total:** 16 module titles, about 45 methods and 60 fields. Most fixes are one line each.

---

## Step 1 — `forgez/forgedoc`: format, check, render

### The format

These are the only types. The extractor prints exactly this JSON.

```go
package forgedoc

// Package is the whole doc of one SDK.
type Package struct {
	Lang    string   `json:"lang"` // code fence language: "python"
	Name    string   `json:"name"` // "i6sdk"
	Modules []Module `json:"modules"`
}

type Module struct {
	Path  string  `json:"path"`  // "i6sdk.infer"
	Title string  `json:"title"` // "Inference": first line of the module docstring
	Doc   string  `json:"doc,omitempty"` // the rest of the module docstring
	Funcs []Func  `json:"funcs,omitempty"`
	Types []Type  `json:"types,omitempty"`
	Vars  []Field `json:"vars,omitempty"` // public module constants
}

type Func struct {
	Name      string   `json:"name"`      // "dsz_list_datasets"
	Signature string   `json:"signature"` // in the SDK's language, without self
	Doc       string   `json:"doc"`       // Markdown
	Params    []Field  `json:"params,omitempty"`
	Returns   *Field   `json:"returns,omitempty"` // Name is empty
	Aliases   []string `json:"aliases,omitempty"` // ["dszListDatasets"]
}

// Type is a class, a record, an enum or an alias. The renderer looks at what is
// set, not at a kind: Fields and no Methods = record; Alias set = alias.
type Type struct {
	Name    string   `json:"name"`
	Doc     string   `json:"doc"`
	Call    string   `json:"call,omitempty"`  // how customers reach it: "client.infer()"
	Alias   string   `json:"alias,omitempty"` // I6Client = Services
	Bases   []string `json:"bases,omitempty"` // I6ApiError -> ["I6Error"]
	Fields  []Field  `json:"fields,omitempty"`
	Methods []Func   `json:"methods,omitempty"`
}

// Field is a record field, an enum member, a constant or a parameter.
type Field struct {
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`    // opaque, in the SDK's language
	Value   string   `json:"value,omitempty"`   // default, enum value or constant value
	Doc     string   `json:"doc,omitempty"`
	Aliases []string `json:"aliases,omitempty"` // ["modelName"]: same field, other spelling
}
```

There is no schema version. The extractor and the renderer live in the same repo and change in the same PR, and `Decode` uses `DisallowUnknownFields`, so a mismatch fails at once.

### The functions

```go
func Decode(r io.Reader) (*Package, error)    // json.Decoder + DisallowUnknownFields
func Check(p *Package) []string               // every missing doc; empty = ok
func Render(p *Package, title string) []byte  // Markdown, deterministic
```

**`Check`** returns *all* problems at once, so one run gives the whole to-do list. Each of these is reported:

- a `Module` with an empty `Title`;
- a `Func`, `Type` or method with an empty `Doc`;
- a `Param` with an empty `Doc`;
- a record `Field` or a `Var` with an empty `Doc`, except a field whose type is itself a record (it is flattened into dotted rows and never shown as a row);
- two types with the same `Name` (which would break links).

`Returns.Doc` is optional: "Returns the Infer client." already says it in the summary.

**`Render`** writes, in this order:

1. The title, then the `<!-- Generated by ./build.sh docs. Do not edit. -->` line, then a contents list linking every module by its `Title`.
2. One `## <Title>` section per module, in extractor order: the module doc, then constants, functions and types. A module with nothing left to show (for example, every type was inlined) is skipped.
3. For each func: a `###` heading, the signature in a code fence, the doc, a parameter table, the return type, and "Also: `dszListDatasets`" when it has aliases. The heading is `Call.name` when the type has a `Call` (`client.apis().dsz_list_datasets`), and `Type.name` otherwise (`Resp.parse`).
4. **Inline request and response.** For each param and for the return value, the renderer finds the record types whose names appear in the type string, for example `Optional[apiDszCreateCustomerPortalReq]` or `Resp[apiDszCreateCustomerPortalResp]`. It writes their fields as one flat table, following nested records with dotted names. Each type is followed once per func, so cycles stop.
5. Records that are inlined somewhere are not listed again. Every other type gets its own `###` with its doc, bases, and a fields table (enum values show in the `Value` column).
6. Every type name in a type string links to that type's anchor.
7. A field or param with `Aliases` shows them in its name cell: `model_name`<br>also `modelName`.

Output for one generated method:

~~~markdown
### client.apis().dsz_create_customer_portal

```python
def dsz_create_customer_portal(req: Optional[apiDszCreateCustomerPortalReq] = None, opts: Optional[ServiceOptions] = None) -> Resp[apiDszCreateCustomerPortalResp]
```

Create Customer Portal

Creates a Stripe billing portal session URL where the user can manage cards, auto-renew, and plans.

| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| req | Optional[apiDszCreateCustomerPortalReq] | None | apiDszCreateCustomerPortalReq request payload. |
| opts | Optional[[ServiceOptions](#serviceoptions)] | None | Optional per-request execution overrides. |

**Request** (`req`)

| Field | Type | Description |
| --- | --- | --- |
| params.dataset | str | Dataset of the customer |
| payload.locale | str | Locale for the portal interface, defaults to en-US if empty (e.g. pt-BR, en-US) |
| payload.return_url | str | URL where the user is redirected when leaving the portal |

**Returns** [Resp](#resp)[apiDszCreateCustomerPortalResp]

| Field | Type | Description |
| --- | --- | --- |
| payload.url | str | URL of the Stripe customer portal session |

Also: `dszCreateCustomerPortal`
~~~

**Deterministic output.** No dates and no versions, and order comes only from the JSON. The same code always gives the same bytes, so `git status` shows a change only when the docs changed.

### The command

`forgez docs` is added next to the existing `forgez py` in `forgez/cmdforgez`:

```bash
forgez docs --title "i6sdk — Python SDK reference" --out SDK.md   # reads the JSON on stdin
```

The command runs `Decode`, then `Check`. If `Check` finds problems, it prints every one and exits 1. Otherwise it runs `Render` and writes `--out` only if the bytes changed.

### Tests

- `testdata/pkg.json` → `testdata/pkg.md`: a golden test that covers every rule above, including nested inline records and a cycle. `go test ./... -update` rewrites the golden file.
- A table test for `Check`: one case per rule.

---

## Step 2 — `sdkpy/docgen.py`: the Python extractor

One script, Python 3.10+, standard library only, about 200 lines. It reads source files with `ast` and never imports the SDK, so there are no side effects.

```bash
python3 -I docgen.py src/i6sdk i6sdk   # prints the JSON on stdout
```

**The walk**

1. Read `__all__` from `src/i6sdk/__init__.py`. Follow each `from .x import y` to the module that defines `y`.
2. Walk **depth-first**, starting with the `__all__` names in order. From each public class, follow the parameter and return annotations of its public methods, in source order, to classes in the package. Skip `internal/` unless the class is in `__all__`. A method whose return type is a skipped internal class is hidden (`Services.fs()`, `downloader()`, `invoker()`). `__init__` is kept only for classes in `__all__`.
3. If two modules define the same name (for example, `ServiceOptions` exists both in `invoker.py` and in `sdkapis/types.py`), keep the one reached from `__all__`.
4. Output modules in the order they are first reached. `__all__` starts with `I6Client`, so the client comes first and its methods (`config`, `auth`, `tokens`, `ingest`, `infer`, `dash`, `apis`, `funcs`) set the order of the next sections. The rest of `__all__` (errors, `Resp`, storage, `Grants`, `bitmapz`) follows. This gives the order in the "Module docstrings" table.
5. **Call paths.** The client class (the target of `I6Client`) gets `Call = "client"`, the variable name the README uses. A class reached through a method with no params gets its parent's `Call` + `.method()`, for example `client.infer()` for `Infer`. Every other class has no `Call`.

**How Python maps to the format**

| Python | Format |
| --- | --- |
| `def` / method | `Func`. `Signature` = `ast.unparse` of the args and the return annotation, without `self`. `Params` keep their annotation as `Type` and their default as `Value` |
| `alias = method` in a class body (`dszListDatasets = dsz_list_datasets`) | added to the method's `Aliases` |
| a field or param whose name is the camelCase of another one in the same record or signature (`modelName` → `model_name`, `onProgress` → `on_progress`) | added to that field's `Aliases`, not emitted on its own. camelCase = split on `_`, then capitalise every part after the first |
| class with methods | `Type` with `Methods`; `__init__` params document the constructor |
| `TypedDict`, `@dataclass` | `Type` with `Fields`; a dataclass default goes in `Value` |
| `IntEnum` | `Type` with `Fields` (`Value` = the number); a member that repeats an earlier value is skipped (`Grants.System`) |
| `Name = Other` in `__all__` | `Type` with `Alias` (`I6Client = Services`) |
| class bases | `Bases` |
| constant in `__all__` | `Module.Vars`; `Value` = `ast.unparse`, cut at 200 characters |

**Docstrings:** in a module docstring, the first line goes in `Module.Title` and the rest in `Module.Doc`. Elsewhere, the text before the first `:field:` line goes in `Doc`; `:param x:` lines go in `Params[x].Doc`, and `:return:` goes in `Returns.Doc`. A param's type always comes from the annotation, never from the docstring.

**Field comments:** the `#` lines directly above a field or constant, read from the source by line number. A leading `[Go: …] - ` is removed.

**Errors:** a name in `__all__` that can't be resolved exits 1 with `file:line: message`.

**Tests:** `sdkpy/tests/test_docgen.py` runs the script on a small fixture package in `sdkpy/tests/docgen_fixture/` and compares the output with an expected JSON file. It runs with the existing unit tests (`cmd_unit_test.sh`), so nothing new is needed in CI.

---

## Step 3 — Build wiring

We follow the existing `cmd_unit_test.sh` pattern: a component with a `cmd_docs.sh` takes part in docs, and the others are skipped.

**`sdkpy/cmd_docs.sh`** (new):

```bash
#!/bin/bash -e
set -o pipefail
go run . generate
python3 -I docgen.py src/i6sdk i6sdk \
  | go run github.com/infinity6-ai/i6sdk/forgez docs --title "i6sdk — Python SDK reference" --out SDK.md
```

**`comp.sh`** (new function):

```bash
function cmd_docs() {
  [ ! -f cmd_docs.sh ] || ./cmd_docs.sh "$@"
}
```

**`build.sh`** (new command, plus one line in `cmd_release`):

```bash
function cmd_docs() {
  cmd_comps_run docs "$@"
}

function cmd_release() {
  ...
  cmd_docs                        # new: fails on missing docs; a stale SDK.md makes the tree dirty
  [ -z "$(git status -s "$@")" ]  # existing check, unchanged
  ...
}
```

**Also:** one line in `sdkpy/README.md` linking to `SDK.md`, and `*/SDK.md linguist-generated=true` in `.gitattributes` so GitHub folds the diff in PRs.

---

## Adding a language

Adding a language means writing one extractor script and one `cmd_docs.sh`. `forgedoc`, `comp.sh` and `build.sh` don't change.

| SDK | Extractor | Notes |
| --- | --- | --- |
| sdkweb (JS) | `sdkweb/docgen.mjs` using `acorn` + `comment-parser` (dev dependencies) | TypeScript 7 has no compiler API we can call from a script. JSDoc `@typedef`/`@property` become records |
| sdkgo | `go/doc` (stdlib), in-process | Waits until sdkgo has SDK code |

---

## Effort and risks

About **4.25 working days** to a complete `sdkpy/SDK.md`. These are estimates for one engineer.

| Step | Estimate |
| --- | --- |
| 1 — `forgedoc` + `forgez docs` + tests | 1.5 d |
| 2 — `docgen.py` + tests | 1.5 d |
| Doc fixes (16 module titles, about 45 methods, 60 fields) | 1 d |
| 3 — wiring | 0.25 d |

- **Docs are only as good as the docstrings.** The check fails on any gap from day 1, so docs can't get worse over time.
- **New API fields without descriptions** will fail the build. The fix goes in the `apicollectionz` spec, and every future SDK benefits.
- **Large file.** `types.py` alone is 2,525 lines. Inlining keeps related fields next to their method; if the file still gets too long, splitting it is a renderer-only change.
