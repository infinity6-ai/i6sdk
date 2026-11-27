# i6 Web SDK

The JavaScript SDK for the Infinity6 platform, for the browser. It is served by
[jsDelivr](https://www.jsdelivr.com/) straight from this GitHub repo: no npm, no install, no build step.

Reference docs, with a **Try it** box for every method: <https://cdn.infinity6.ai/docs/>.

## Quick start

```html
<script src="https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/bundle/i6sdk-web.min.js"></script>
<script>
  const sdk = new I6Sdk();
  // sdk.auth(), sdk.tokens(), sdk.ingest(), sdk.infer(), sdk.dash(), sdk.apis(), sdk.funcs()
</script>
```

Prefer ES modules? Use:

```html
<script type="module">
  import { I6Sdk } from "https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/module-min/index.js";
  const sdk = new I6Sdk();
</script>
```

---

## How the URLs work

```
https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@<version>/sdkweb/dist/<path>
                         │  │            │      │         │
                         │  │            │      │         └─ any file inside dist (see "What's inside dist")
                         │  │            │      └─ the version: a release tag, or the "latest" branch
                         │  │            └─ repo
                         │  └─ GitHub org
                         └─ jsDelivr serving files from GitHub
```

Everything inside `dist/` is available, including every subfolder.

---

## Versions (tags)

Every release is a git tag named:

```
sdkweb-dist-v<major>.<minor>.<patch>        e.g. sdkweb-dist-v0.0.8
```

- **A tag never changes.** A file under a tag is the same forever and jsDelivr caches it
  permanently. Use a tag in production.
- **Always write the full tag.** Short forms like `@1` or `@v0.0` don't work with these tags.
- Don't use the tags named `sdkweb/v…` or `v…`. They are for the Go modules and contain no `dist`
  folder.

See which versions exist:

- GitHub: https://github.com/infinity6-ai/i6sdk/tags (look for `sdkweb-dist-…`)
- jsDelivr API (JSON): https://data.jsdelivr.com/v1/packages/gh/infinity6-ai/i6sdk

Upgrading means changing the tag in your URLs:

```diff
- https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.7/sdkweb/dist/bundle/i6sdk-web.min.js
+ https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/bundle/i6sdk-web.min.js
```

### Latest

The branch **`sdkweb-dist`** always points to the newest release:

```html
<script src="https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist/sdkweb/dist/bundle/i6sdk-web.min.js"></script>
```

- Good for demos, prototypes and trying new versions.
- **Not for production.** It changes with every release, including breaking ones.
- jsDelivr caches branches for a while (hours), so a new release may take some time to show up
  under `@sdkweb-dist`.

---

## What's inside `dist`

```
sdkweb/dist/
├── bundle/        one file, the whole SDK. Easiest option.
├── module/        ES modules, one file per source file, readable
├── module-min/    same as module/, minified
└── types/         TypeScript declarations, one per source file
```

### `bundle/`: the whole SDK in one file

| File | What it is |
|------|------------|
| `i6sdk-web.js` | Full SDK, readable. Loaded with a classic `<script>`, it sets `window.I6Sdk`. |
| `i6sdk-web.min.js` | Same thing minified. **Use this one in production.** |
| `i6sdk-web.js.map`, `i6sdk-web.min.js.map` | Source maps. Browser dev tools load them by themselves, so you debug the original code. |
| `i6sdk-web.d.ts` | All type declarations in a single file, for editors, docs and LLMs. |

Besides `I6Sdk`, the bundle provides `I6Sdk.Resp`, `I6Sdk.bitmapz` and `I6Sdk.Grants`.

### `module/` and `module-min/`: ES modules

The same SDK split into files, with their imports between them. Load `index.js` and the browser
fetches the rest by itself:

```js
import { I6Sdk, Resp, bitmapz, Grants } from "https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/module-min/index.js";
```

| Path | Contents |
|------|----------|
| `index.js` | Entry point: `I6Sdk`, `Resp`, `bitmapz`, `Grants` |
| `auth/auth.js` | Login (PKCE), callback handling, user tokens |
| `tokens/tokens.js` | Storage, exchange and silent renewal of user and dataset tokens |
| `ingest/ingest.js` | Uploading files to the ingestion service |
| `infer/infer.js` | Running an i6 model in the browser |
| `dash/dash.js` | Running a dash's query in the browser and returning its result |
| `services/services.js` | Wiring of all services |
| `sdkapis/sdkapis.js` | Generated clients for every platform API |
| `invoker/invoker.js` | HTTP calls and the `Resp` type |
| `internal/…` | Internal helpers (b64, gzip, fs, sqlite, downloader, …). Their API may change at any time, so don't import them directly. |

Use `module-min/` in production. `module/` is the same code, easier to read.

These are many small files. If you want a single request, use `bundle/`.

### `types/`: TypeScript declarations

One `.d.ts` per module, same layout as `module/` (for example `types/auth/auth.d.ts`).
For a single file with everything, use `bundle/i6sdk-web.d.ts`.

To get autocomplete in your editor, download the declarations into your project:

```bash
curl -L -o i6sdk-web.d.ts https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/bundle/i6sdk-web.d.ts
```

### Browse the files

Open a folder URL ending in `/` to get a listing:

- https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/
- https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.8/sdkweb/dist/module/

---

## Good to know

- Some features load extra files at runtime. For example, the sqlite engine (sql.js wasm) used by
  `infer` is loaded from jsDelivr when it is first needed.
- By default the SDK points to the sandbox environment. Pass a config to `new I6Sdk({...})` to change
  it (see `I6SdkConfig` in `i6sdk-web.d.ts`).
