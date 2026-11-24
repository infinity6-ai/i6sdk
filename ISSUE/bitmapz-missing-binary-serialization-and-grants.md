# Issue: `sdkpy` Missing `bitmapz` Binary Serialization and `Grants` Enum

**Component:** `sdkpy` (`i6sdk.internal.bitmapz`, `i6sdk`)  
**Severity:** Medium (Developer ergonomics and protocol compatibility)  
**Status:** Resolved  

---

## Summary

In `sdkweb`, the SDK exports an explicit `Grants` constant object and comprehensive `bitmapz` utilities, including binary serialization (`serialize` and `deserialize`) matching the core Go implementation in `commonz/bitmapz`. In `sdkpy`, `bitmapz.py` previously omitted `serialize()` and `deserialize()`, and the SDK defined no `Grants` constants, forcing Python developers to deal with raw magic numbers (`1`, `2`, `4`, `8`, `15`) when handling permissions and dataset tokens.

---

## Resolution

1. Added `serialize()` and `deserialize()` to `sdkpy/src/i6sdk/internal/bitmapz.py` matching Go `commonz/bitmapz` and `sdkweb/src/main/internal/bitmapz/bitmapz.js`.
2. Created `Grants` enum (`IntEnum`) in `sdkpy/src/i6sdk/grants.py` with standard UPPERCASE values (`SYSTEM = 1`, `ADMIN = 2`, `WRITE = 4`, `READ = 8`, `ALL = 15`) and TitleCase aliases (`System`, `Admin`, `Write`, `Read`, `All`).
3. Exported `Grants` and `bitmapz` at the root of `i6sdk` (`from i6sdk import Grants, bitmapz`), as well as attributes on `I6Client.Grants` and `I6Client.bitmapz`.
4. Added comprehensive test coverage in `sdkpy/tests/test_bitmapz.py` validating standard serialization vectors and error conditions.
