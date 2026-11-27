# i6 Web SDK

The JavaScript/TypeScript SDK for the Infinity6 platform in the browser. Zero install, zero build step—served directly via CDN from GitHub Pages.

## CDN Endpoints

All built bundles, ES modules, and TypeScript types are hosted at:

```
https://sdk.lp.infinity6.ai/release/sdkweb/<channel>/<path>
```
*(Fallback host: `https://infinity6-ai.github.io/i6sdk/release/sdkweb/<channel>/<path>`)*

### Channels

| Channel | Description | Recommended For |
|---|---|---|
| `vX.Y.Z` (e.g. `v0.0.8`) | Immutable, pinned version release | Production |
| `latest` | Automatically updated with each tagged release | Demos & latest stable features |
| `edge` | Continuously updated from the `edge` branch | Development & preview builds |

---

## Quick Start

### 1. Single Bundle (Classic `<script>`)

Loads the entire SDK in a single minified bundle and exposes `window.I6Sdk`:

```html
<script src="https://sdk.lp.infinity6.ai/release/sdkweb/latest/bundle/i6sdkweb.min.js"></script>
<script>
  const sdk = new I6Sdk();
  // sdk.apis(), sdk.funcs(), sdk.auth(), sdk.tokens(), sdk.ingest(), sdk.infer(), sdk.dash()
</script>
```

### 2. ES Modules

Browser-native ES module imports:

```html
<script type="module">
  import { I6Sdk, Resp, bitmapz, Grants } from "https://sdk.lp.infinity6.ai/release/sdkweb/latest/module-min/index.js";
  const sdk = new I6Sdk();
</script>
```

---

## Distribution Contents (`dist/`)

Each channel hosts the complete build directory:

| Path | Format | Description |
|---|---|---|
| `bundle/i6sdkweb.min.js` | UMD / Global | Minified all-in-one browser bundle (recommended for `<script>`) |
| `bundle/i6sdkweb.js` | UMD / Global | Unminified bundle with source maps for debugging |
| `bundle/i6sdkweb.d.ts` | TypeScript | All type definitions stitched into a single file for IDEs |
| `module-min/` | ES Module | Minified individual ES modules, dynamically imported by the browser |
| `module/` | ES Module | Unminified individual ES modules for reading and inspection |
| `types/` | TypeScript | Granular TypeScript declaration files per module |

### Downloading TypeScript Definitions

To get IDE autocompletion in your project without an npm package:

```bash
curl -L -o i6sdkweb.d.ts https://sdk.lp.infinity6.ai/release/sdkweb/latest/bundle/i6sdkweb.d.ts
```

---

## Local Development & Testing

```bash
# Install dependencies
npm install

# Build bundles and types into dist/
./build.sh bundle

# Run browser tests
./build.sh test_visual
```

## Configuration

By default, the SDK connects to the Infinity6 sandbox environment. Pass an explicit config object to `new I6Sdk({...})` to change endpoints (see `I6SdkConfig` in `i6sdkweb.d.ts`).
