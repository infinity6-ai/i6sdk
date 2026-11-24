# Issue: Method Naming Divergence for `SdkFuncs` Client Accessor (`funcs` vs `sdk_funcs`)

**Component:** `sdkpy` (`i6sdk.services`)  
**Severity:** Low-Medium (API inconsistency across SDKs)  
**Status:** Open  

---

## Summary

In `sdkweb`, the generated high-level SDK workflow functions are accessed via `sdk.funcs()`. In `sdkpy`, the accessor is named `client.sdk_funcs()` (with an alias `client.sdkFuncs()`), but **`client.funcs()` does not exist**. This divergence causes cross-language API inconsistency and confusion when reading documentation or code samples.

---

## Root Cause Analysis

In `sdkweb/src/main/index.js` and `services.js`:
```javascript
// sdkweb/src/main/index.js
funcs() {
  return this._service.funcs()
}

// sdkweb/src/main/services/services.js
funcs() {
  const ret = new SdkFuncs(this)
  if (ret.prepare) ret.prepare()
  return ret 
}
```

In `sdkpy/src/i6sdk/services.py`:
```python
def sdk_funcs(self) -> SdkFuncs:
    """Returns the generated high-level SDK workflow functions client."""
    if self._sdk_funcs is None:
        self._sdk_funcs = SdkFuncs(self)
    return self._sdk_funcs

sdkFuncs = sdk_funcs
# Notice: 'funcs' is not defined as an alias!
```

---

## Proposed Remediation

Add `funcs` alias on `Services` in `sdkpy/src/i6sdk/services.py`:

```python
    funcs = sdk_funcs
```
