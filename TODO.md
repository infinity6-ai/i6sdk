# Downloader & OPFS Cache Enhancement Plan

## Objective
Implement stale-while-revalidate caching and background preloading for assets stored in OPFS (`fs.js` and `downloader.js`). This enables instant cached reads while keeping assets up to date in the background, prevents duplicate in-flight network requests, and tracks file creation timestamps.

---

## 1. Plan

### Step 1: Track `checkedat.txt` in `fs.js`
- **Location:** Under `i6/fs/{name}/checkedat.txt`.
- **Content:** An ISO timestamp or epoch milliseconds (e.g. `Date.now().toString()`) representing when the file was last verified/checked against remote.
- **Write Triggers:**
  - Written when a new version is created/released.
  - Updated via a new method `fs.touch({ name })` whenever a `HEAD` check confirms the remote ETag matches the released version.
- **Expose on `FileRef`:**
  - `fs.get({ name })` reads `checkedat.txt` and returns `checkedAt: number` (epoch ms) on `FileRef`.
  - **Backward Compatibility:** If `checkedat.txt` does not exist, `checkedAt` falls back to `new Date(version).getTime()`.

### Step 2: Track In-Flight Download Promises in `downloader.js`
- **Deduplication Map:** Add `this.downloads = new Map<string, Promise<FileRef>>()` to the `Downloader` instance.
- **Promise Re-use:** When `downloader.update()` starts an update check or download:
  - If `this.downloads.has(opts.name)`, re-use and return the existing in-flight promise.
  - Otherwise, initiate the update promise, store it in `this.downloads`, and delete it in a `.finally()` callback once settled.
- **Background Detach Safety:** When a download/check runs in the background (stale-while-revalidate), attach a `.catch()` logger to the background promise so network drops do not trigger unhandled promise rejections.
- **On ETag Match:** When the remote ETag matches `current.etag`, call `await fs.touch({ name: opts.name })` to reset `checkedAt` so redundant `HEAD` checks are avoided for the next `minPreloadTime` seconds across page reloads.

### Step 3: Lifecycle Management in `downloader.update(opts)`
- **Parameters:** Accept `{ name, url, minPreloadTime }` (`minPreloadTime` in seconds).
- **Core Principle:** If a file is older than `minPreloadTime`, kick off the background update check first, then immediately return the existing cached file. Only block if there is no file at all.
- **Execution Flow:**
  1. Retrieve `current = await fs.get({ name: opts.name })`.
  2. Compute `age = current ? (Date.now() - current.checkedAt) / 1000 : Infinity`.
  3. **Step A — Kick off download/check if missing or stale:**
     - If `!current || (opts.minPreloadTime !== undefined && age >= opts.minPreloadTime)`:
       - Ensure an update promise is started and tracked in `this.downloads` (if not already running).
  4. **Step B — Return cached file if available:**
     - If `current` exists: return `current` immediately (non-blocking). Any check/download started in Step A continues running in the background.
  5. **Step C — Await download if no file existed:**
     - If `!current`: await and return the update promise tracked in `this.downloads`.

---

## 2. Questions & Clarifications

1. **The `tag` parameter:**
   - In `{ name, tag, minPreloadTime }`, is `tag` an optional known ETag/version (e.g. to bypass the initial `HEAD` request if already known by the caller)? Or is `url` sufficient?

2. **Default behavior when `minPreloadTime` is omitted:**
   - If `minPreloadTime` is omitted or `undefined`: should it check and block on update (current behavior), or should it never preload in the background?
