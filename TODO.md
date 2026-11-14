# Downloader & OPFS Cache Enhancement Plan

## Objective
Implement stale-while-revalidate caching and background preloading for assets stored in OPFS (`fs.js` and `downloader.js`). This enables instant cached reads while keeping assets up to date in the background, prevents duplicate in-flight network requests, and tracks file creation timestamps.

---

## 1. Plan

### Step 1: Track `createdat.txt` in `fs.js`
- **Write timestamp on creation:** When `fs.create({ name, etag })` creates `v{version}/`, write `createdat.txt` containing the creation timestamp (e.g., epoch milliseconds `Date.now().toString()`).
- **Expose on `FileRef`:** Update `fs.get()` to read `createdat.txt` and return `createdAt: number` (epoch ms) as part of `FileRef`.
- **Backward Compatibility:** If `createdat.txt` does not exist in an existing version folder, `fs.get()` falls back to `new Date(version).getTime()` since version directory names are ISO timestamps.
- **Cache Timestamp Touch:** Add a method (e.g. `fs.touch({ name, version })`) to update `createdat.txt` when a remote check confirms the asset ETag has not changed.

### Step 2: Track In-Flight Download Promises in `downloader.js`
- **Deduplication Map:** Add `this.downloads = new Map<string, Promise<FileRef>>()` to the `Downloader` instance.
- **Promise Re-use:** When `downloader.update()` decides a download is necessary:
  - If `this.downloads.has(opts.name)`, re-use and return the existing in-flight promise.
  - Otherwise, initiate the fetch/stream download promise, store it in `this.downloads`, and delete it in a `.finally()` callback once settled.
- **Background Detach Safety:** When a download runs in the background (stale-while-revalidate), attach a `.catch()` logger to the background promise so network drops do not trigger unhandled promise rejections.

### Step 3: Lifecycle Management in `downloader.update(opts)`
- **Parameters:** Accept `{ name, url, tag?, minPreloadTime }` (`minPreloadTime` in seconds).
- **Core Principle:** If a file is too old, kick off the background download first, then immediately return the existing cached file. Only block if there is no file at all.
- **Execution Flow:**
  1. Retrieve `current = await fs.get({ name: opts.name })`.
  2. Compute `age = current ? (Date.now() - current.createdAt) / 1000 : Infinity`.
  3. **Step A — Kick off download if missing or stale:**
     - If `!current || (opts.minPreloadTime !== undefined && age >= opts.minPreloadTime)`:
       - Ensure a download promise is started and tracked in `this.downloads` (if not already running).
  4. **Step B — Return cached file if available:**
     - If `current` exists: return `current` immediately (non-blocking). Any download started in Step A continues running in the background.
  5. **Step C — Await download if no file existed:**
     - If `!current`: await and return the download promise tracked in `this.downloads`.

---

## 2. Questions & Clarifications

1. **The `tag` parameter:**
   - In `{ name, tag, minPreloadTime }`, is `tag` an optional known ETag or version (to bypass the initial `HEAD` request if already known by the caller)?
   - Is `url` still required as the remote endpoint?

2. **Updating `createdat.txt` on ETag match (304 / unchanged):**
   - When a `HEAD` request returns the same ETag as the released version, should we update `createdat.txt` to `Date.now()` so the file is considered fresh again, preventing redundant `HEAD` checks on every subsequent call?

3. **Default behavior when `minPreloadTime` is omitted:**
   - If `minPreloadTime` is omitted or `undefined`: should it check and block on update (current behavior), or should it never preload in the background?
