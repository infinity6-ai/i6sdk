"""Local filesystem cache manager for model bundles and SQLite databases."""

from __future__ import annotations

import json
import os
import pathlib
import threading
from typing import Any, Dict, Optional, Union


class Fs:
    """Manages the local cache directory for downloading and versioning assets."""

    def __init__(self, base_dir: Optional[Union[pathlib.Path, str]] = None) -> None:
        if base_dir is None:
            self.base_dir = pathlib.Path.home() / ".cache" / "i6"
        else:
            self.base_dir = pathlib.Path(base_dir)
        self.base_dir.mkdir(parents=True, exist_ok=True)
        self._meta_file = self.base_dir / "meta.json"
        self._lock = threading.Lock()

    def _read_meta(self) -> Dict[str, Dict[str, Any]]:
        if not self._meta_file.exists():
            return {}
        try:
            with open(self._meta_file, "r", encoding="utf-8") as f:
                return json.load(f)
        except Exception:
            return {}

    def _write_meta(self, meta: Dict[str, Dict[str, Any]]) -> None:
        tmp = self._meta_file.with_suffix(".tmp")
        with open(tmp, "w", encoding="utf-8") as f:
            json.dump(meta, f, indent=2)
        tmp.replace(self._meta_file)

    def resolve(self, name: str) -> pathlib.Path:
        """Returns the absolute file path for a named asset."""
        safe_name = name.replace("/", "_").replace("\\", "_")
        target = self.base_dir / safe_name
        return target

    def get_meta(self, name: str) -> Optional[Dict[str, Any]]:
        """Returns cached metadata for the asset if it exists on disk."""
        with self._lock:
            meta = self._read_meta().get(name)
            if meta and self.resolve(name).exists():
                return meta
            return None

    def set_meta(self, name: str, meta: Dict[str, Any]) -> None:
        """Saves metadata for the named asset."""
        with self._lock:
            all_meta = self._read_meta()
            all_meta[name] = meta
            self._write_meta(all_meta)

    def exists(self, name: str) -> bool:
        """Checks if the named asset exists in the cache."""
        return self.resolve(name).is_file()

    def delete(self, name: str) -> None:
        """Deletes the asset and its metadata from the cache."""
        with self._lock:
            all_meta = self._read_meta()
            if name in all_meta:
                del all_meta[name]
                self._write_meta(all_meta)
            target = self.resolve(name)
            if target.exists():
                try:
                    target.unlink()
                except OSError:
                    pass
