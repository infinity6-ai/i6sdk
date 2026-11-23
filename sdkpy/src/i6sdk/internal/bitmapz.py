"""Compact variable-length bitmask operations, serialization, and deserialization.

Direct Python port of gox/commonz/bitmapz and sdkweb/src/main/internal/bitmapz/bitmapz.js.
"""

from __future__ import annotations

from typing import List, Union


def join(*values: Union[int, float]) -> int:
    """Combines all given values into a single bitmask using bitwise OR."""
    result = 0
    for v in values:
        result |= int(v)
    return result


def contains_all(current: int, *expected: int) -> bool:
    """Reports whether current contains all bits in expected."""
    all_bits = join(*expected)
    return (current & all_bits) == all_bits


def contains_any(current: int, *expected: int) -> bool:
    """Reports whether current contains at least one of the bits in expected."""
    any_bits = join(*expected)
    return (current & any_bits) != 0


def split(current: int) -> List[int]:
    """Decomposes current into a list of individual set bits in ascending order."""
    if not current:
        return []
    result: List[int] = []
    c = current
    while c > 0:
        bit = c & -c
        result.append(bit)
        c &= ~bit
    return result


def set_bits(current: int, enable: bool, *values: int) -> int:
    """Adds or removes the given values from current based on enable."""
    bits = join(*values)
    if enable:
        return current | bits
    return current & ~bits
