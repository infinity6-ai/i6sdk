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


# Alias matching sdkweb
set = set_bits


def serialize(v: int) -> bytes:
    """Encodes the bitmask v into a compact variable-length byte array.

    If v fits in 4 bits (0..15), it is serialized in 1 byte:
    the upper 4 bits are 0000, and the lower 4 bits store the value.
    Otherwise, the upper 4 bits specify the number of following payload bytes (1..8),
    the lower 4 bits are 0000, and the payload bytes follow in little-endian order.
    """
    if v < 0:
        raise ValueError("bitmapz: value must be non-negative")
    if v <= 0x0F:
        return bytes([v])

    bit_len = v.bit_length()
    size = (bit_len + 7) // 8
    if size > 8:
        raise ValueError(f"bitmapz: invalid size {size}: bitmask size cannot exceed 8 bytes")

    buf = bytearray(1 + size)
    buf[0] = size << 4
    for i in range(size):
        buf[1 + i] = (v >> (8 * i)) & 0xFF
    return bytes(buf)


def deserialize(data: Union[bytes, bytearray]) -> int:
    """Decodes a bitmask from data according to the variable-length encoding."""
    if not data or len(data) == 0:
        raise ValueError("bitmapz: data is empty")

    header = data[0]
    size = header >> 4

    if size > 8:
        raise ValueError(f"bitmapz: invalid size {size}: bitmask size cannot exceed 8 bytes")

    if size == 0:
        if len(data) != 1:
            raise ValueError(
                f"bitmapz: invalid data length: expected 1 byte for 4-bit value, got {len(data)}"
            )
        return header & 0x0F

    if (header & 0x0F) != 0:
        raise ValueError(
            f"bitmapz: invalid header 0x{header:02x}: lower 4 bits must be 0 when size > 0"
        )

    expected_len = 1 + size
    if len(data) != expected_len:
        raise ValueError(
            f"bitmapz: invalid data length: expected {expected_len} bytes for size {size}, got {len(data)}"
        )

    val = 0
    for i in range(size):
        val |= data[1 + i] << (8 * i)
    return val


__all__ = [
    "join",
    "contains_all",
    "contains_any",
    "split",
    "set_bits",
    "set",
    "serialize",
    "deserialize",
]

