# Issue: Method Naming Divergence for `SdkFuncs` Client Accessor (`funcs` vs `sdk_funcs`)

**Component:** `sdkpy` (`i6sdk.services`)  
**Severity:** Low-Medium (API inconsistency across SDKs)  
**Status:** Resolved  

---

## Summary

In `sdkweb`, the generated high-level SDK workflow functions are accessed via `sdk.funcs()`. In `sdkpy`, the accessor was previously named `client.sdk_funcs()` (with an alias `client.sdkFuncs()`), but **`client.funcs()` did not exist**. This divergence caused cross-language API inconsistency and confusion when reading documentation or code samples.

---

## Resolution

`sdk_funcs` and `sdkFuncs` were removed in favor of `client.funcs()` in `sdkpy/src/i6sdk/services.py`:

```python
    def funcs(self) -> SdkFuncs:
        """Returns the generated high-level SDK workflow functions client."""
        if self._funcs is None:
            self._funcs = SdkFuncs(self)
        return self._funcs
```

All documentation, examples, and tests have been updated accordingly.
