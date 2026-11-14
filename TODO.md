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
- **Parameters:** Accept `{ name, url, tag?, minPreloadTime, minAge }` (times in seconds).
- **Execution Flow:**
  1. Retrieve `current = await fs.get({ name: opts.name })`.
  2. Compute `age = current ? (Date.now() - current.createdAt) / 1000 : Infinity`.
  3. **Trigger Preload (if stale or missing):**
     - If `!current || age >= minPreloadTime`: ensure a download promise is active in `this.downloads`.
  4. **Fast Return (if too new / fresh):**
     - If `current && age < minAge`: return `current` immediately (non-blocking, zero network delay).
  5. **Fallback:**
     - If `!current`: return and await the in-flight download promise.
     - If `current` exists but is expired (`age >= minAge`): await the download promise to provide fresh data.

---

## 2. Questions & Clarifications

1. **Meaning of `minAge` vs `minPreloadTime`:**
   - Does `minAge` represent the **maximum allowed freshness** (i.e. if `age < minAge`, return cached immediately without waiting; if `age >= minAge`, block until the new version is downloaded)?
   - Does `minPreloadTime` represent the **preload threshold** (i.e. start downloading in the background if `age >= minPreloadTime`, but if `age < minAge` we can still return the stale file immediately)?

2. **The `tag` parameter:**
   - In `{ name, tag, minPreloadTime, minAge }`, is `tag` an optional known ETag or version (to bypass the initial `HEAD` request if already known by the caller)?
   - Is `url` still required as the remote endpoint?

3. **Updating `createdat.txt` on ETag match (304 / unchanged):**
   - When a `HEAD` request returns the same ETag as the released version, should we update `createdat.txt` to `Date.now()` so the file is considered fresh again, preventing redundant `HEAD` checks on every subsequent call?

4. **Default values:**
   - When `minPreloadTime` and `minAge` are omitted, should `downloader.update()` default to current behavior (always check ETag via `HEAD` and block until verified/downloaded)?
