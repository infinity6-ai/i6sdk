# SDK Docs Site — Plan

We build an interactive HTML reference for the i6 SDK, like Swagger. Every method has its own static page with its docs and a code box: the VS Code editor (Monaco), with autocomplete for the methods and for the request JSON. The user signs in, changes the call and presses **Run**, and the call runs for real in the browser.

We build it once, on **sdkweb**, because sdkweb runs in the browser. There is no page per language.

**The pages are static.** All the docs are plain HTML written at build time, nothing is rendered in the browser. Google can index every page, and every method has its own URL to share. JavaScript only adds the interactive parts on top (the editor, **Run**, sign in, search).

**How it works, in one line:** forge already has every API spec, and a small script reads the handwritten modules from the SDK types; one Go command turns both into static HTML pages.

```
apicollectionz ─► forge (parse, in process) ──────────┐
                                                      ├─► go run . docs ─► target/docs/*.html  (static, one page per method)
sdkweb src ─► tsc ─► dts-bundle-generator ─► i6sdk-web.d.ts
                                    │                 │
                                    └─► docgen.mjs ─► sdk.json (handwritten modules)
                                    └─► Monaco autocomplete (loaded by the pages)
```

- **Scope:** sdkweb only.
- **Customer API does not change.** No public name, signature or behaviour changes. We add the docs build and edit JSDoc comments.
- **Hosting:** GitHub Pages, through the existing `i6pages` repo (`cdn.infinity6.ai/docs/`), the same way `tokens/` is released today. There is one copy, for the latest release; older versions are not kept.
- **Out of scope:** the Python and Go SDKs, guides and tutorials (they stay in the READMEs).

---

## What changed from the first plan

| Before | Now |
| --- | --- |
| A Markdown file, `sdkpy/SDK.md` | Static HTML pages for sdkweb, one per method |
| One extractor per language (`docgen.py`, …) | Forge for the APIs (in process), one small script for the handwritten modules |
| Read-only | Code box with autocomplete, sign in, and a **Run** button |
| HTTP method and path hidden (old decision 8) | Shown, like Swagger: forge has them (`POST /api/ds/dataset/{dataset}/portal`) |
| Python examples | JS examples, written the way customers call the SDK: `sdk.apis().dszCreateCustomerPortal(...)` |

---

## Decisions

| # | Decision | Why |
| --- | --- | --- |
| 1 | Only sdkweb | It runs in the browser, so the pages can run real calls with it |
| 2 | **Static pages, built at build time.** Docs are never rendered in the browser | Google indexes them, and links can be shared. The pages work (read-only) with JavaScript off |
| 3 | One page per method, plus one page per section and an index | Every method gets its own URL, title and description, so it shows up in search results and can be shared: `cdn.infinity6.ai/docs/apis/dsz-create-customer-portal.html` |
| 4 | One Go command, `go run . docs` in sdkweb, writes every page with `html/template` | Go, like the rest of forge. `html/template` escapes everything, so no doc text can break the page |
| 5 | The API docs come from forge, in process: `docs` calls the same parse as `generate` | Same parse that writes `sdkapis.js`, so the docs can't drift from the code. No JSON file in between |
| 6 | The handwritten modules (auth, tokens, ingest, infer, dash, `I6Sdk`) are read from the bundled `.d.ts` by `sdkweb/docgen.mjs`, which writes `sdk.json` for `docs` | They are not in forge. The `.d.ts` is also what Monaco uses, so the docs and the autocomplete always agree |
| 7 | One real bundled `i6sdk-web.d.ts`, built with `dts-bundle-generator` (already a dev dependency) | The current one is the module files joined with `cat`; its relative imports don't resolve, so Monaco would show `any`. The new one also fixes the file customers download for their editors |
| 8 | TypeScript 5 next to TypeScript 7, as a dev dependency. TypeScript 7 keeps building the code (`tsc`); TypeScript 5 is used by `dts-bundle-generator` and `docgen.mjs` | `dts-bundle-generator` needs the TypeScript 5 API: with TypeScript 7 it crashes on load. TypeScript 7 has an API only under `typescript/unstable/…`, which can break in any 7.x release. Since TypeScript 5 is needed for the bundler anyway, `docgen.mjs` uses it too |
| 9 | The code box is a **JS** box per method, prefilled with the call, not a separate JSON body editor | Monaco completes the methods *and* the request fields from the same `.d.ts` |
| 10 | The example in the box is written into the static page at build time, from the request typedef, with an empty value per field | Google sees it, and it is readable without JavaScript. Monaco takes it over when **Try it** is opened |
| 11 | The SDK is set up in the page the way we tell customers to: `new I6Sdk({ servUrls, app, userExchangeCallbackUrl })`. A bar at the top of every page always shows that config and the sign-in state | Users see exactly what the calls use, and the example is the code they will write |
| 12 | Sign in is the SDK's own flow (`sdk.auth().requireUser()`), from the top bar | Nothing new to build; the page gets real user and dataset tokens |
| 13 | No confirm before calls that change data | The user signs in with their own account and presses **Run** on purpose |
| 14 | Missing docs fail the build | Docs can't get worse over time |
| 15 | Monaco is loaded from jsDelivr at a pinned version. Doc text is rendered at build time by a small Markdown renderer in `docs` (Step 3), not a library | No npm or bundler for the pages, like the rest of `src/site/`. The docs use only a few Markdown features, so about 80 lines cover them with no new dependency, and every word is in the static HTML |
| 16 | Generated pages are not committed: `docs` writes to `target/docs/`, and the release copies them to `i6pages` | They are built from code we already commit, and `i6pages` keeps the published copy |
| 17 | **Try it** has an `inspect(x)` function that prints everything about what it is given: a call shows the HTTP requests the SDK made and what came back, a method shows its docs | Swagger shows the raw request and response; `inspect` gives the same, and more, from the code box. It exists only in the docs page, so the SDK's customer API does not change |

---

## Which methods get a page, and which can run

The question was whether the handwritten modules belong in the page. They do: login, upload, inference and dashboards are the main things customers do with the SDK. They don't all get **Try it**, though.

| Module | Page | Try it | Why |
| --- | --- | --- | --- |
| API reference (`sdk.apis()`, about 63 methods) | yes | yes | Plain JSON in, JSON out. The core of the Swagger part |
| Functions (`sdk.funcs()`) | yes | yes | Same as the APIs |
| `I6Sdk`, config | yes | no | Shown in the top bar instead (decision 11) |
| Auth (`sdk.auth()`) | yes | no | The top bar already does sign in and sign out. Running `requireUser` from a box would redirect away from the page |
| Tokens (`sdk.tokens()`) | yes | no | The SDK manages tokens by itself; customers only read about them |
| Inference (`sdk.infer().run`, …) | yes | yes | JSON options, JSON result |
| Dashboards (`sdk.dash().query`) | yes | yes | JSON options, JSON result |
| Ingestion (`sdk.ingest().upload`) | yes | no | An upload writes real data into a dataset |
| Errors, `Resp`, `bitmapz`, `Grants` | yes | no | Types and helpers, nothing to call on the server |

---

## The page

```
┌───────────────────────────────────────────────────────────────────────────┐
│ i6 SDK docs   [search]                                                    │
│ ┌ SDK config ───────────────────────────────────────────┐  ┌ Account ───┐ │
│ │ const sdk = new I6Sdk({                               │  │ Signed in  │ │
│ │   app: "default",                                     │  │ ana@acme   │ │
│ │   servUrls: { decsuite: "https://decsuite.sandbox…",  │  │ [Sign out] │ │
│ │               tfinfer: "…", userz: "…" },             │  └────────────┘ │
│ │   userExchangeCallbackUrl: "…/docs/callback.html",    │                 │
│ │ })                                          [Edit]    │                 │
│ └───────────────────────────────────────────────────────┘                 │
├────────────┬──────────────────────────────────────────────────────────────┤
│ Auth       │ POST  /api/ds/dataset/{dataset}/portal       token: dataset  │
│ Tokens     │ sdk.apis().dszCreateCustomerPortal                           │
│ Ingestion  │                                                              │
│ Inference  │ Create Customer Portal                                       │
│ Dashboards │ Creates a Stripe billing portal session URL where ...        │
│ API ▾      │                                                              │
│  dsz       │ Request                                                      │
│  pipez     │ | Field           | Type   | Description             |       │
│  ...       │ | params.dataset  | string | Dataset of the customer |       │
│ Functions  │ | payload.locale  | string | Locale for the portal … |       │
│            │                                                              │
│            │ Response                                                     │
│            │ | payload.url     | string | URL of the Stripe …     |       │
│            │                                                              │
│            │ ┌ Try it ───────────────────────────────────────────────┐    │
│            │ │ return await sdk.apis().dszCreateCustomerPortal({      │    │
│            │ │   params: { dataset: "" },                            │    │
│            │ │   payload: { locale: "", return_url: "" },            │    │
│            │ │ })                                                    │    │
│            │ └──────────────────────────────────────── [Run] [Reset]─┘    │
│            │ ┌ Result ───────────────────────────────────────────────┐    │
│            │ │ 200  { "payload": { "url": "https://..." } }          │    │
│            │ └───────────────────────────────────────────────────────┘    │
└────────────┴──────────────────────────────────────────────────────────────┘
```

**Static parts (in the HTML file):** the nav, the method's HTTP method and path, the call path, the token type it needs (`labels.token_types`), name, summary and description text (rendered from Markdown at build time, decision 15), the request and response fields as flat tables with dotted names (nested records are followed once, so cycles stop), the example code, and the page `<title>` and `<meta name="description">` from the method's name and summary. Each field row also has an anchor, so a link can point to one field.

**Interactive parts (JavaScript on top):**

- **SDK config bar**, on every page. It shows the `new I6Sdk({...})` call the page uses, with the SDK's defaults (read from `sdk.config()`) and the docs' own `userExchangeCallbackUrl`. **Edit** opens it in Monaco, so the user can point to another environment the same way a customer would. The edited config is kept in `localStorage`, so it stays while moving between pages, and **Reset** goes back to the defaults.
- **Account**, on every page: **Sign in** (the SDK's PKCE flow, coming back to the same page through `callback.html`), the signed-in user, and **Sign out**. **Run** without signing in shows "Sign in to run".
- **Try it:** Monaco takes over the static example when the user opens it, so pages stay light. **Run** runs the box as the body of an `async` function with `sdk` in scope. What it returns is shown as JSON; a `Resp` shows its status and payload. `console.log` in the box goes to the result panel too (reusing `src/site/console.js`). **Reset** puts the example back.
- **`inspect(x)`**, in every **Try it** box (see "The `inspect` function" below).
- **Search:** filters the nav links already in the page. No extra file.

**Also written by `docs`:** `sitemap.xml` (every page, for Google), `callback.html` (the PKCE return page), and `404.html`.

---

## The `inspect` function

`inspect` is a function every **Try it** box has, next to `sdk`. It prints everything about what it is given into the result panel, and returns it unchanged, so it can wrap any call:

```js
const resp = await inspect(sdk.apis().dszCreateCustomerPortal({
  params: { dataset: "acme" },
  payload: { locale: "en-US", return_url: "https://acme.com" },
}))
```

**What it prints, by what it is given:**

| Given | It prints |
| --- | --- |
| A call (the promise of an SDK method) | Every HTTP request the SDK made for that call, in order: method, full URL, headers, body, then the status, headers, body and time of each response. This includes the requests the SDK makes by itself, such as getting a dataset token. Then the value the call returned, or the error it threw |
| A `Resp` | Its status, headers and payload |
| An SDK method (`inspect(sdk.apis().dszCreateCustomerPortal)`) | Its docs, as on its page: summary, HTTP method and path, token type, the request and response fields, and a link to its page. The docs come from `methods.json`, written by `docs` |
| An error (`I6ApiError`, …) | Its class, message, status and details |
| Anything else | The value in full, as JSON, with nested objects open |

- **Tokens are hidden by default.** `Authorization` headers and token fields are printed as `Bearer eyJh…(hidden)`, so a screenshot or a copied result doesn't leak them. A **Show tokens** toggle on the result panel prints them in full; it is off on every page load.
- **How it sees the requests:** the page wraps the browser's `window.fetch` with a recorder before it creates `sdk`. Wrapping `fetch`, and not using the SDK's `fetcher` config, is needed because `auth.js` and `tokens.js` call `fetch` directly (token exchange, refresh, `me`), so a `fetcher` would miss them. The SDK and the user's config are not touched. `inspect` of a call shows the requests made since the **Run** started or since the previous `inspect` returned, so in code that awaits one call after another, each call shows its own requests.
- **How it knows which method it got:** it compares the function with the methods of the SDK objects the page has (`sdk`, `sdk.apis()`, `sdk.infer()`, …), not by its name. So `inspect(sdk.infer().run)` finds `Infer.run`, even though other classes may also have a `run`. A function that isn't an SDK method is printed like any other value.
- **Types:** the page adds `declare function inspect<T>(x: T): T` to Monaco, so it autocompletes and the type of the wrapped call is kept.
- **Explained in the docs:** every **Try it** box has a one-line hint under it: "Wrap a call in `inspect(...)` to see the requests it makes."

---

## Step 1 — The bundled `.d.ts`

- `sdkweb/package.json`: add TypeScript 5 for the tools that need its API (decision 8): `"typescript5": "npm:typescript@5"` for `docgen.mjs`, and `"overrides": { "dts-bundle-generator": { "typescript": "^5" } }` so `dts-bundle-generator` gets TypeScript 5 while the project keeps TypeScript 7.
- `sdkweb/build.sh bundle`: replace the `find … -exec cat` line with `dts-bundle-generator`, which writes one real `dist/bundle/i6sdk-web.d.ts` from the `tsc` output (`dist/types/index.d.ts`). Only what is reachable from the exports is in it, so unexported helper modules drop out by themselves.

---

## Step 2 — `sdkweb/docgen.mjs`: the handwritten modules

A Node script, about 150 lines. It reads `dist/bundle/i6sdk-web.d.ts` with the TypeScript compiler API and prints `target/sdk.json`. `tsc` has already turned the JSDoc into regular declarations, so the script only walks them: classes, methods, params, return types and typedefs, with their docs.

- **Dependency:** TypeScript 5, through the `typescript5` alias from Step 1 (`import ts from "typescript5"`). See decision 8.
- **What it skips:** `Apis` and `SdkFuncs`, because `docs` gets those from forge, and any member starting with `_` or `#private`.
- **Call paths:** `I6Sdk` is `sdk`. A class reached through a method with no params gets its parent's call path plus `.method()`, for example `sdk.infer()` for `Infer`.
- **Missing docs:** a method, param, typedef or field with no doc is printed as `name: missing doc`, and the script exits 1.
- **Test:** a small fixture `.d.ts` and the expected JSON, run by `go test` from `cmdsdkweb` (it runs the script with `node`).

---

## Step 3 — `go run . docs`: the static pages

A new command in `sdkweb/cmdsdkweb`, next to `generate`:

```bash
go run . docs --sdk target/sdk.json --out target/docs --base-url https://cdn.infinity6.ai/docs
```

1. Parses the APIs and SdkFuncs the same way `generate` does, with the same few calls (`apicollectionz.Apis()`, `apispecz.New(...).ParseApis`, …). `generate` is left untouched.
2. Reads `sdk.json`.
3. Checks the docs: an API, SdkFunc or leaf field with an empty summary is reported (`api id: field: missing doc`). If any are found, it prints them all and exits 1. Today every API field has a summary, so this only guards the future; the fix goes in the `apicollectionz` spec.
4. Writes the pages with `html/template` from templates in `sdkweb/cmdsdkweb/docstpl/` (embedded with `go:embed`): `index.html`, one page per section, one page per method, `sitemap.xml`, `methods.json`, `callback.html`, `404.html`. It also copies `src/site/docs/` (the JS and CSS), `dist/bundle/i6sdk-web.d.ts` and `dist/bundle/i6sdk-web.min.js` next to them.
5. Renders the doc text with a small Markdown renderer (`docsmd.go`, about 80 lines), at build time. It covers exactly what the docs use today:

    | Markdown | Uses today | HTML |
    | --- | --- | --- |
    | Paragraphs (blank lines) | everywhere | `<p>` |
    | Inline code `` `x` `` | 161 | `<code>` |
    | Lists (`- item`, `1. item`) | 35 | `<ul>`, `<ol>` |
    | Code blocks (three backticks) | 7 | `<pre><code>` |
    | Bold `**x**` | 1 | `<strong>` |
    | JSDoc links `{@link Auth#handleCallback}` | 9 | a link to that method's page |

    The text is escaped first and the tags are added after, so doc text can't inject HTML. Anything outside the table shows as plain text; if a doc needs more, we add it to the renderer.
6. The output is deterministic: no dates, and the order comes only from the code. The same code always gives the same files.

**Tests:** a golden test on one API page (`dsz-create-customer-portal.html`) and one handwritten page (`infer-run.html`), a table test for the check, one case per rule, and a table test for the Markdown renderer, one case per row of its table (plus one for escaping). `go test ./... -update` rewrites the golden files.

---

## Step 4 — The interactive parts

`sdkweb/src/site/docs/` (handwritten, plain JS, no build step):

| File | What it does |
| --- | --- |
| `docs.css` | Styles (same look as `console.html`) |
| `bar.js` | SDK config bar and account box: creates `sdk` from the saved or default config, **Edit**, **Reset**, **Sign in**, **Sign out** |
| `tryit.js` | Loads Monaco and `i6sdk-web.d.ts` on the first **Try it**, adds `declare const sdk: I6Sdk`, runs the code, shows the result |
| `inspect.js` | `inspect(x)` and the request recorder that wraps `window.fetch` |
| `search.js` | Filters the nav links in the page |

**Which SDK runs:** the `i6sdk-web.min.js` that `docs` copies next to the pages, so the pages always run the SDK they were built from, with no version to match and nothing to wait for on jsDelivr. Locally it is the same: rerun `build.sh docs` to see an SDK change.

---

## Step 5 — Build wiring and release

- `sdkweb/build.sh docs` (new): `bundle`, then `node docgen.mjs`, then `go run . docs`.
- `sdkweb/build.sh release`: runs `docs`, then one more `rsync` line, like `tokens/` today: `target/docs/` goes to `../../i6pages/docs/`.
- Local preview: the dev server gets a `/docs/*` route that serves `target/docs/` (`http://127.0.0.1:8043/docs/index.html`).
- A browser test in `sdkweb_web_test.go` (`TestWeb…`, like the existing ones): opens a few pages, checks the static content is there, and checks that every generated example has no type errors in Monaco. This catches a broken example or a broken `.d.ts` before release.
- One line in `sdkweb/README.md` linking to the docs.

---

## Step 6 — Doc fixes

Once Step 2 and Step 3 exist, the build prints the exact list of missing docs. We fix all of them in JSDoc comments only: no names, signatures or behaviour.

---

## Risks

- **Sign in and CORS.** **Run** only works once the docs origin is allowed by i6userz and the APIs. That is decided outside this plan; the docs and the autocomplete work without it.
- **Running code changes real data.** The page uses the user's own tokens, so it can only do what the user can already do.
- **Two TypeScript versions.** TypeScript 5 will get old. When the TypeScript 7 API is stable, `docgen.mjs` can move to it; only that script changes.
- **Monaco from a CDN.** It is pinned to one version, so it can't change under us. Reading the docs works offline from the cache, but **Try it** needs the network anyway.
