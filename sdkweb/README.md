# sdkweb: distribution

Two ways to ship the web SDK without the npm registry. Each release is versioned.

| # | Strategy | Example URL |
|---|----------|-------------|
| 1 | GitHub Releases | `https://github.com/infinity6-ai/i6sdk/releases/download/sdkweb-dist-v0.0.4/i6sdk-web.min.js` |
| 2 | jsDelivr (from GitHub) | `https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.4/sdkweb/dist/bundle/i6sdk-web.min.js` |

## Versioning

- Version comes from `versionz`: `./build.sh version` gives something like `v0.0.4`.
- Each release is the git tag `sdkweb-dist-<version>`. Don't confuse it with `sdkweb/<version>`,
  which is the Go module tag.
- Never re-release a version. Bump it with `./build.sh version-inc`.

## Build (both strategies)

From the repo root:

```bash
V="$(./build.sh version)"
./sdkweb/build.sh bundle
```

---

## 1. GitHub Releases

Attaches the built files to a release. They are downloads; for a `<script src>` URL use strategy 2.

Create (needs the [GitHub CLI](https://cli.github.com/), `gh`):

```bash
gh release create "sdkweb-dist-$V" sdkweb/dist/bundle/* --title "sdkweb $V"
```

If the tag doesn't exist yet, `gh` creates it on the current commit of `main`. If you also use
strategy 2, do strategy 2 first, so this command reuses its tag.

Use:

```bash
curl -LO https://github.com/infinity6-ai/i6sdk/releases/download/sdkweb-dist-v0.0.4/i6sdk-web.min.js
```

---

## 2. jsDelivr (from GitHub)

jsDelivr serves any file committed at a tag of a public GitHub repo. `dist` is git-ignored, so the
release tag points to a commit that exists only to hold `dist`. Nothing is added to `main`.

Create:

```bash
git checkout --detach
git add -f sdkweb/dist
git commit -m "sdkweb dist $V"
git tag "sdkweb-dist-$V"
git push origin "sdkweb-dist-$V"
git checkout -
```

Use:

```html
<script src="https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.4/sdkweb/dist/bundle/i6sdk-web.min.js"></script>
<script>
  const sdk = new I6Sdk();
</script>
```

Types: `https://cdn.jsdelivr.net/gh/infinity6-ai/i6sdk@sdkweb-dist-v0.0.4/sdkweb/dist/bundle/i6sdk-web.d.ts`
