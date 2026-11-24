# Issue: `sdkpy` Missing `bitmapz` Binary Serialization and `Grants` Enum

**Component:** `sdkpy` (`i6sdk.internal.bitmapz`, `i6sdk`)  
**Severity:** Medium (Developer ergonomics and protocol compatibility)  
**Status:** Open  

---

## Summary

In `sdkweb`, the SDK exports an explicit `Grants` constant object and comprehensive `bitmapz` utilities, including binary serialization (`serialize` and `deserialize`) matching the core Go implementation in `commonz/bitmapz`. In `sdkpy`, `bitmapz.py` omits `serialize()` and `deserialize()`, and the SDK defines no `Grants` constants, forcing Python developers to deal with raw magic numbers (`1`, `2`, `4`, `8`, `15`) when handling permissions and dataset tokens.

---

## Root Cause Analysis

### 1. Missing Serialization in `sdkpy/src/i6sdk/internal/bitmapz.py`
In `sdkweb/src/main/internal/bitmapz/bitmapz.js`:
* `serialize(v)`: Encodes a bitmask into a compact variable-length `Uint8Array` (1-byte optimization for `<= 15`, multi-byte with length prefix for larger values).
* `deserialize(bytes)`: Decodes the byte array back to a bitmask number/BigInt.

In `sdkpy/src/i6sdk/internal/bitmapz.py`:
* Only bitwise helpers (`join`, `contains_all`, `contains_any`, `split`, `set_bits`) are implemented.
* `serialize` and `deserialize` are completely absent.

### 2. Missing `Grants` Constants in `sdkpy`
In `sdkweb/src/main/index.js`:
```javascript
const Grants = Object.freeze({
  System: 1,
  Admin: 2,
  Write: 4,
  Read: 8,
  All: 15,
});
export { Grants };
```
In `sdkpy`:
* No `Grants` class or enum is defined.
* `bitmapz` is not exported from `i6sdk` package root.

---

## Proposed Remediation

1. Port `serialize()` and `deserialize()` to `sdkpy/src/i6sdk/internal/bitmapz.py` matching `sdkweb` and Go `commonz/bitmapz`.
2. Add an `IntEnum` or frozen class `Grants` in `sdkpy`:
```python
from enum import IntEnum

class Grants(IntEnum):
    SYSTEM = 1
    ADMIN = 2
    WRITE = 4
    READ = 8
    ALL = 15
```
3. Export `Grants` and `bitmapz` from `i6sdk/__init__.py`.
