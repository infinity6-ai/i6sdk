"""Conditional HTTP downloader matching sdkweb's Downloader with versioned ETag caching."""

from __future__ import annotations

import datetime
import pathlib
import threading
import time
import urllib.error
import urllib.request
from typing import Dict, Optional

from .fs import FileRef, Fs


class Downloader:
    """Downloads remote files into the versioned file store matching sdkweb's Downloader."""

    def __init__(self, fs: Optional[Fs] = None) -> None:
        self.fs = fs or Fs()
        self._lock = threading.RLock()

    def update(
        self,
        name: str,
        url: str,
        min_preload_time: Optional[float] = None,
        timeout: float = 60.0,
    ) -> pathlib.Path:
        """Downloads or revalidates a file by name, returning the local blob.bin path."""
        current = self.fs.get(name)

        if min_preload_time is not None and current is not None:
            last_checked = 0.0
            if current.checked_at:
                try:
                    dt = datetime.datetime.fromisoformat(current.checked_at.replace("Z", "+00:00"))
                    last_checked = dt.timestamp()
                except Exception:
                    pass
            age = time.time() - last_checked
            if age < min_preload_time:
                return self.fs.resolve(name, current.version)

        return self._fetch_and_update(name, url, current, timeout=timeout)

    def _fetch_and_update(
        self,
        name: str,
        url: str,
        current: Optional[FileRef],
        timeout: float = 60.0,
    ) -> pathlib.Path:
        headers: Dict[str, str] = {
            "User-Agent": "i6sdkpy/0.0.1",
        }
        if current and current.etag:
            headers["If-None-Match"] = current.etag

        req = urllib.request.Request(url=url, headers=headers, method="GET")

        new_ref: Optional[FileRef] = None
        try:
            with urllib.request.urlopen(req, timeout=timeout) as resp:
                resp_etag = resp.headers.get("etag") or resp.headers.get("ETag") or ""
                new_ref = self.fs.create(name, etag=resp_etag)
                target_path = self.fs.resolve(name, new_ref.version)
                tmp_path = target_path.with_suffix(".tmp")

                with open(tmp_path, "wb") as f:
                    while chunk := resp.read(65536):
                        f.write(chunk)
                tmp_path.replace(target_path)

                released = self.fs.release(name, new_ref.version)
                return self.fs.resolve(name, released.version)

        except urllib.error.HTTPError as e:
            if e.code == 304 and current:
                self.fs.touch(name)
                return self.fs.resolve(name, current.version)
            if new_ref:
                self.fs.discard(name, new_ref.version)
            if current:
                return self.fs.resolve(name, current.version)
            raise RuntimeError(f"Download failed for {name} ({url}): HTTP {e.code}") from e
        except Exception as e:
            if new_ref:
                self.fs.discard(name, new_ref.version)
            if current:
                return self.fs.resolve(name, current.version)
            raise RuntimeError(f"Download failed for {name} ({url}): {e}") from e
