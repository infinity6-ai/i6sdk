"""Permissions catalog and bitmask constants matching Go's grants package."""

from __future__ import annotations

from enum import IntEnum


class Grants(IntEnum):
    """Compact permission bitmask over datasets and solutions."""

    SYSTEM = 1
    ADMIN = 2
    WRITE = 4
    READ = 8
    ALL = 15

    # Case-insensitive / TitleCase aliases for parity with Go and sdkweb
    System = 1
    Admin = 2
    Write = 4
    Read = 8
    All = 15


__all__ = ["Grants"]
