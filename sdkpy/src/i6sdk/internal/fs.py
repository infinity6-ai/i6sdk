"""Versioned filesystem store matching sdkweb's OPFS layout."""

from __future__ import annotations

import datetime
import pathlib
import shutil
import threading
from dataclasses import dataclass
from typing import Optional, Union


@dataclass
class FileRef:
    """Reference to one version of a stored file, matching sdkweb's FileRef."""

    name: str
    version: str
    etag: str
    checked_at: Optional[str] = None


class Fs:
    """Manages the local cache directory for downloading and versioning assets matching sdkweb's OPFS layout."""

    def __init__(self, base_dir: Optional[Union[pathlib.Path, str]] = None) -> None:
        if base_dir is None:
            root = pathlib.Path.home() / ".cache" / "i6"
        else:
            root = pathlib.Path(base_dir)
        self.base_dir = root / "sdkpy" / "fs"
        self.base_dir.mkdir(parents=True, exist_ok=True)
        self._lock = threading.Lock()

    def _file_dir(self, name: str) -> pathlib.Path:
        safe_name = name.replace("/", "_").replace("\\", "_")
        return self.base_dir / safe_name

    def resolve(self, name: str, version: Optional[str] = None) -> pathlib.Path:
        """Returns the absolute file path for a version's blob.bin (defaults to released version)."""
        file_dir = self._file_dir(name)
        ver = version
        if not ver:
            version_file = file_dir / "version.txt"
            if version_file.exists():
                ver = version_file.read_text(encoding="utf-8").strip()
            else:
                ver = ""
        return file_dir / f"v{ver}" / "blob.bin"

    def get(self, name: str, version: Optional[str] = None) -> Optional[FileRef]:
        """Resolves a version of a file, or returns None if not found."""
        with self._lock:
            file_dir = self._file_dir(name)
            if not file_dir.is_dir():
                return None

            version_file = file_dir / "version.txt"
            released_version = version_file.read_text(encoding="utf-8").strip() if version_file.exists() else None
            ver = version or released_version
            if not ver:
                return None

            ver_dir = file_dir / f"v{ver}"
            blob_file = ver_dir / "blob.bin"
            if not ver_dir.is_dir() or not blob_file.is_file():
                return None

            etag_file = ver_dir / "etag.txt"
            etag = etag_file.read_text(encoding="utf-8").strip() if etag_file.exists() else ""

            checked_at = None
            if not version or version == released_version:
                checkedat_file = file_dir / "checkedat.txt"
                if checkedat_file.exists():
                    checked_at = checkedat_file.read_text(encoding="utf-8").strip()

            return FileRef(name=name, version=ver, etag=etag, checked_at=checked_at)

    def create(self, name: str, etag: str) -> FileRef:
        """Creates an unreleased version directory and writes etag.txt."""
        with self._lock:
            file_dir = self._file_dir(name)
            file_dir.mkdir(parents=True, exist_ok=True)

            # Check if latest existing version directory already matches this etag
            v_dirs = [d.name for d in file_dir.iterdir() if d.is_dir() and d.name.startswith("v")]
            latest = max(v_dirs) if v_dirs else None
            if latest:
                latest_etag_file = file_dir / latest / "etag.txt"
                if latest_etag_file.exists() and latest_etag_file.read_text(encoding="utf-8").strip() == etag:
                    return FileRef(name=name, version=latest[1:], etag=etag)

            version = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"
            ver_dir = file_dir / f"v{version}"
            ver_dir.mkdir(parents=True, exist_ok=True)
            (ver_dir / "etag.txt").write_text(etag, encoding="utf-8")
            return FileRef(name=name, version=version, etag=etag)

    def release(self, name: str, version: str) -> FileRef:
        """Makes a version the current released one, updates checkedat.txt, and cleans older versions."""
        with self._lock:
            ref = self.get(name, version=version)
            if not ref:
                raise LookupError(f"fs: version {version} of {name} does not exist")

            file_dir = self._file_dir(name)
            (file_dir / "version.txt").write_text(version, encoding="utf-8")
            checked_at = self.touch(name)
            self.clean_file(name)
            ref.checked_at = checked_at
            return ref

    def touch(self, name: str, checked_at: Optional[str] = None) -> str:
        """Updates checkedat.txt to now (or custom ISO timestamp)."""
        file_dir = self._file_dir(name)
        file_dir.mkdir(parents=True, exist_ok=True)
        ts = checked_at or datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"
        (file_dir / "checkedat.txt").write_text(ts, encoding="utf-8")
        return ts

    def discard(self, name: str, version: str) -> None:
        """Removes an unreleased version directory."""
        with self._lock:
            file_dir = self._file_dir(name)
            version_file = file_dir / "version.txt"
            if version_file.exists() and version_file.read_text(encoding="utf-8").strip() == version:
                raise RuntimeError(f"fs: version {version} of {name} is the released one")
            ver_dir = file_dir / f"v{version}"
            if ver_dir.exists():
                shutil.rmtree(ver_dir, ignore_errors=True)

    def clean_file(self, name: str) -> None:
        """Removes versions older than the released version."""
        file_dir = self._file_dir(name)
        if not file_dir.exists():
            return
        version_file = file_dir / "version.txt"
        if not version_file.exists():
            return
        released = "v" + version_file.read_text(encoding="utf-8").strip()
        for item in file_dir.iterdir():
            if item.is_dir() and item.name.startswith("v") and item.name < released:
                shutil.rmtree(item, ignore_errors=True)

    def clean(self) -> None:
        """Cleans older versions across all files."""
        for item in self.base_dir.iterdir():
            if item.is_dir():
                self.clean_file(item.name)

    def exists(self, name: str) -> bool:
        """Checks if a released version of the asset exists."""
        return self.get(name) is not None

    def delete(self, name: str) -> None:
        """Deletes the entire asset folder."""
        with self._lock:
            file_dir = self._file_dir(name)
            if file_dir.exists():
                shutil.rmtree(file_dir, ignore_errors=True)
