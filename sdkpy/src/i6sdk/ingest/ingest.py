"""Ingest client for uploading streams, files, and batches to the ingestion service."""

from __future__ import annotations

import io
import os
import threading
import urllib.request
from typing import (
    TYPE_CHECKING,
    Any,
    BinaryIO,
    Callable,
    Dict,
    Optional,
    TypedDict,
    Union,
)

if TYPE_CHECKING:
    from ..invoker import ServiceOptions
    from ..services import Services

from ..exceptions import I6ApiError, I6ValidationError


class UploadProgress(TypedDict):
    """Progress status of an in-flight upload."""

    loaded: int
    total: int
    percent: float


class UploadResult(TypedDict):
    """Result of a successful upload."""

    ingest_id: str


class SignedUrlResult(TypedDict):
    """Result of a signed upload URL request."""

    ingest_id: str
    url: str


class ProgressStreamReader:
    """Wraps a binary stream to track upload progress and enable cancellation."""

    def __init__(
        self,
        stream: BinaryIO,
        total_size: int,
        on_progress: Optional[Callable[[UploadProgress], None]] = None,
        cancel_event: Optional[Union[threading.Event, Any]] = None,
    ) -> None:
        self._stream = stream
        self._total_size = total_size
        self._on_progress = on_progress
        self._cancel_event = cancel_event
        self._bytes_read = 0

    def read(self, size: int = -1) -> bytes:
        if self._cancel_event is not None and getattr(self._cancel_event, "is_set", lambda: False)():
            raise I6ValidationError("Upload aborted by caller")

        if size == -1:
            chunks = []
            while True:
                chunk = self._read_chunk(65536)
                if not chunk:
                    break
                chunks.append(chunk)
            return b"".join(chunks)

        return self._read_chunk(size)

    def _read_chunk(self, size: int) -> bytes:
        if self._cancel_event is not None and getattr(self._cancel_event, "is_set", lambda: False)():
            raise I6ValidationError("Upload aborted by caller")

        chunk = self._stream.read(size)
        if chunk:
            self._bytes_read += len(chunk)
            if self._on_progress is not None:
                total = self._total_size if self._total_size > 0 else self._bytes_read
                pct = (self._bytes_read / total) * 100.0 if total > 0 else 100.0
                try:
                    self._on_progress({
                        "loaded": self._bytes_read,
                        "total": total,
                        "percent": min(100.0, pct),
                    })
                except Exception:
                    pass
        return chunk

    def __len__(self) -> int:
        return self._total_size


class Ingest:
    """Client for uploading files and data batches to the i6 ingestion service."""

    def __init__(self, services: Services) -> None:
        self._services = services

    def get_url(
        self,
        dataset: str,
        table: str,
        partitions: Dict[str, str],
        service_options: Optional[Union[ServiceOptions, Dict[str, Any]]] = None,
    ) -> SignedUrlResult:
        """Requests a signed upload URL from the ingest service."""
        resp = self._services.apis().ingestz_get_url(
            req={
                "params": {
                    "dataset": dataset,
                },
                "payload": {
                    "table": table,
                    "partitions": partitions,
                    "amount": 1,
                },
            },
            opts=service_options,
        )
        if not resp.ok:
            raise I6ApiError(
                message=f"Failed to get upload URL: {resp.text}",
                status=resp.status,
                code=(resp.error or {}).get("code", "bad_request"),
            )
        data = resp.parse() or {}
        uploads = data.get("uploads") or []
        url = uploads[0] if uploads else ""
        return {
            "ingest_id": data.get("ingest_id", ""),
            "url": url,
        }

    getUrl = get_url

    def upload(
        self,
        dataset: str,
        table: str,
        partitions: Dict[str, str],
        stream: Optional[Union[BinaryIO, str, bytes]] = None,
        file: Optional[Union[BinaryIO, str, bytes]] = None,
        total_size: Optional[int] = None,
        content_type: str = "application/octet-stream",
        on_progress: Optional[Callable[[UploadProgress], None]] = None,
        onProgress: Optional[Callable[[UploadProgress], None]] = None,
        cancel_event: Optional[threading.Event] = None,
        signal: Optional[Union[threading.Event, Any]] = None,
        timeout: Optional[float] = None,
        service_options: Optional[Union[ServiceOptions, Dict[str, Any]]] = None,
    ) -> UploadResult:
        """Uploads a stream, file path, or bytes buffer to the ingest service using a signed PUT URL."""
        if not dataset:
            raise I6ValidationError("dataset is required")
        if not table:
            raise I6ValidationError("table is required")

        data_source = stream if stream is not None else file
        if data_source is None:
            raise I6ValidationError("stream is required")

        progress_cb = on_progress or onProgress
        abort_trigger = cancel_event if cancel_event is not None else signal

        if abort_trigger is not None and getattr(abort_trigger, "is_set", lambda: False)():
            raise I6ValidationError("Upload aborted by caller")

        # Resolve stream object and byte size
        should_close = False
        if isinstance(data_source, (bytes, bytearray)):
            raw_bytes = bytes(data_source)
            calc_size = len(raw_bytes)
            stream_obj: BinaryIO = io.BytesIO(raw_bytes)
        elif isinstance(data_source, str):
            if os.path.isfile(data_source):
                calc_size = os.path.getsize(data_source)
                stream_obj = open(data_source, "rb")
                should_close = True
            else:
                raw_bytes = data_source.encode("utf-8")
                calc_size = len(raw_bytes)
                stream_obj = io.BytesIO(raw_bytes)
        elif hasattr(data_source, "read"):
            stream_obj = data_source
            if total_size is not None:
                calc_size = total_size
            else:
                try:
                    pos = stream_obj.tell()
                    stream_obj.seek(0, io.SEEK_END)
                    calc_size = stream_obj.tell() - pos
                    stream_obj.seek(pos, io.SEEK_SET)
                except Exception:
                    calc_size = 0
        else:
            raise I6ValidationError("stream must be a BinaryIO stream, file path string, or bytes")

        if total_size is not None:
            calc_size = total_size

        try:
            url_info = self.get_url(
                dataset=dataset,
                table=table,
                partitions=partitions,
                service_options=service_options,
            )
            signed_url = url_info["url"]
            ingest_id = url_info["ingest_id"]

            if abort_trigger is not None and getattr(abort_trigger, "is_set", lambda: False)():
                raise I6ValidationError("Upload aborted by caller")

            reader = ProgressStreamReader(
                stream=stream_obj,
                total_size=calc_size,
                on_progress=progress_cb,
                cancel_event=abort_trigger,
            )

            headers = {
                "Content-Type": content_type,
                "Content-Length": str(calc_size),
            }

            request_timeout = timeout or self._services.config().timeout or 120.0
            config = self._services.config()

            if config.fetcher is not None:
                body_bytes = reader.read()
                status, _, text = config.fetcher(
                    signed_url,
                    "PUT",
                    headers,
                    body_bytes,
                    request_timeout,
                )
                if not (200 <= status < 300):
                    raise I6ApiError(
                        message=f"Upload PUT failed with HTTP {status}: {text}",
                        status=status,
                    )
            else:
                req = urllib.request.Request(
                    url=signed_url,
                    data=reader,
                    headers=headers,
                    method="PUT",
                )
                try:
                    with urllib.request.urlopen(req, timeout=request_timeout) as resp:
                        if not (200 <= resp.status < 300):
                            raise I6ApiError(
                                message=f"Upload PUT failed with HTTP {resp.status}",
                                status=resp.status,
                            )
                except I6ValidationError:
                    raise
                except Exception as e:
                    if abort_trigger is not None and getattr(abort_trigger, "is_set", lambda: False)():
                        raise I6ValidationError("Upload aborted by caller") from e
                    raise I6ApiError(message=f"Upload failed: {e}") from e

            return {"ingest_id": ingest_id}
        finally:
            if should_close:
                try:
                    stream_obj.close()
                except Exception:
                    pass
