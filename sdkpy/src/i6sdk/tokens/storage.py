"""Token storage interface and implementations."""

from __future__ import annotations

import abc
import json
import os
import pathlib
import threading
from typing import Any, Dict, Optional, Union


class TokenStorage(abc.ABC):
    """Abstract interface for token persistence."""

    @abc.abstractmethod
    def get_item(self, key: str) -> Optional[str]:
        """Retrieve stored string value for key."""
        pass

    @abc.abstractmethod
    def set_item(self, key: str, value: str) -> None:
        """Store string value for key."""
        pass

    @abc.abstractmethod
    def remove_item(self, key: str) -> None:
        """Remove value for key."""
        pass


class MemoryTokenStorage(TokenStorage):
    """Thread-safe in-memory token storage."""

    def __init__(self) -> None:
        self._data: Dict[str, str] = {}
        self._lock = threading.Lock()

    def get_item(self, key: str) -> Optional[str]:
        with self._lock:
            return self._data.get(key)

    def set_item(self, key: str, value: str) -> None:
        with self._lock:
            self._data[key] = value

    def remove_item(self, key: str) -> None:
        with self._lock:
            self._data.pop(key, None)


class FileTokenStorage(TokenStorage):
    """File-backed token storage with secure file permissions (0600)."""

    def __init__(self, path: Optional[Union[pathlib.Path, str]] = None) -> None:
        if path is None:
            config_dir = pathlib.Path.home() / ".config" / "i6"
            config_dir.mkdir(parents=True, exist_ok=True)
            self.path = config_dir / "tokens.json"
        else:
            self.path = pathlib.Path(path)
            self.path.parent.mkdir(parents=True, exist_ok=True)
        self._lock = threading.Lock()

    def _read_data(self) -> Dict[str, str]:
        if not self.path.exists():
            return {}
        try:
            with open(self.path, "r", encoding="utf-8") as f:
                return json.load(f)
        except Exception:
            return {}

    def _write_data(self, data: Dict[str, str]) -> None:
        tmp_path = self.path.with_suffix(".tmp")
        flags = os.O_WRONLY | os.O_CREAT | os.O_TRUNC
        fd = os.open(tmp_path, flags, 0o600)
        with open(fd, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)
        tmp_path.replace(self.path)

    def get_item(self, key: str) -> Optional[str]:
        with self._lock:
            return self._read_data().get(key)

    def set_item(self, key: str, value: str) -> None:
        with self._lock:
            data = self._read_data()
            data[key] = value
            self._write_data(data)

    def remove_item(self, key: str) -> None:
        with self._lock:
            data = self._read_data()
            if key in data:
                del data[key]
                self._write_data(data)
