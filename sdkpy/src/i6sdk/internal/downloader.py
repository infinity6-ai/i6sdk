"""Conditional HTTP downloader with ETag caching using urllib.request."""

from __future__ import annotations

import pathlib
import time
import urllib.error
import urllib.request
from typing import Any, Dict, Optional, Union

from .fs import Fs


class Downloader:
    """Downloads remote files into local cache, downloading only when ETag changes."""

    def __init__(self, fs: Optional[Fs] = None) -> None:
        self.fs = fs or Fs()

    def update(
        self,
        name: str,
        url: str,
        min_preload_time: Optional[float] = None,
        timeout: float = 60.0,
    ) -> pathlib.Path:
        """Downloads or revalidates a file by name and returns its local filesystem path."""
        target_path = self.fs.resolve(name)
        meta = self.fs.get_meta(name)

        if min_preload_time is not None and meta and target_path.exists():
            last_checked = float(meta.get("checked_at", 0))
            age = time.time() - last_checked
            if age < min_preload_time:
                return target_path

        headers: Dict[str, str] = {
            "User-Agent": "i6sdk-python/0.0.1",
        }
        if meta and (etag := meta.get("etag")):
            headers["If-None-Match"] = etag

        req = urllib.request.Request(url=url, headers=headers, method="GET")

        try:
            with urllib.request.urlopen(req, timeout=timeout) as resp:
                resp_etag = resp.headers.get("etag") or resp.headers.get("ETag") or ""
                # Stream write to tmp file
                tmp_path = target_path.with_suffix(".tmp")
                with open(tmp_path, "wb") as f:
                    while chunk := resp.read(65536):
                        f.write(chunk)
                tmp_path.replace(target_path)

                self.fs.set_meta(
                    name,
                    {
                        "url": url,
                        "etag": resp_etag,
                        "checked_at": time.time(),
                        "size": target_path.stat().st_size,
                    },
                )
                return target_path
        except urllib.error.HTTPError as e:
            if e.code == 304 and target_path.exists():
                # Not modified, update checked_at timestamp
                if meta:
                    meta["checked_at"] = time.time()
                    self.fs.set_meta(name, meta)
                return target_path
            raise RuntimeError(f"Download failed for {name} ({url}): HTTP {e.code}") from e
        except Exception as e:
            if target_path.exists():
                # Fallback to existing cached file on network error
                return target_path
            raise RuntimeError(f"Download failed for {name} ({url}): {e}") from e
