# Downloader & OPFS Cache Enhancement Plan

## Objective
Implement stale-while-revalidate caching and background preloading for assets stored in OPFS (`fs.js` and `downloader.js`). This enables instant cached reads while keeping assets up to date in the background, prevents duplicate in-flight network requests, and tracks file creation timestamps.

---

## 1. Plan

### Step 1: Track `checkedat.txt` in `fs.js` [COMPLETED]
- **Location:** Under `i6/fs/{name}/checkedat.txt`.
- **Content:** An ISO timestamp (e.g. `new Date().toISOString()`) representing when the file was last verified/checked against remote.
- **Write Triggers:**
  - Written on `fs.release()` when a new version becomes active.
  - Written via `fs.touch({ name, checkedAt? })` to update verification time.
- **Expose on `FileRef`:**
  - `fs.get({ name })` returns `checkedAt: string` (ISO timestamp) on the released `FileRef`.
  - Backwards-compatible: if `checkedat.txt` is missing, `checkedAt` is omitted, allowing consumers to fall back to `version`.

### Step 2: Track In-Flight Download Promises in `downloader.js` [COMPLETED]
- **Deduplication Map:** Added `this.downloads = new Map<string, Promise<FileRef>>()` to the `Downloader` instance.
- **Promise Re-use:** When `downloader.update()` runs:
  - If `this.downloads.has(opts.name)`, re-uses and returns the existing in-flight promise.
  - Otherwise, initiates `_fetchAndUpdate(opts)`, tracks it in `this.downloads`, and deletes it in `.finally()` once settled.
- **On ETag Match:** When the remote ETag matches `current.etag`, calls `await fs.touch({ name: opts.name })` to update `checkedAt` and returns `{ ...current, checkedAt }`.

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
