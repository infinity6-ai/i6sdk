"""Ingest client for uploading streams, files, and batches to the ingestion service."""

from __future__ import annotations

import concurrent.futures
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
    List,
    Optional,
    Sequence,
    Tuple,
    TypedDict,
    Union,
)

if TYPE_CHECKING:
    from ..invoker import ServiceOptions
    from ..services import Services

from ..exceptions import I6ApiError, I6ValidationError


class UploadProgress(TypedDict, total=False):
    """Progress status of an in-flight upload."""

    loaded: int
    total: int
    percent: float
    completed_streams: int
    total_streams: int
    stream_index: Optional[int]


class UploadResult(TypedDict):
    """Result of a successful upload."""

    ingest_id: str
    count: int


class SignedUrlResult(TypedDict):
    """Result of a signed upload URL request."""

    ingest_id: str
    url: str
    urls: List[str]


StreamInput = Union[
    BinaryIO,
    str,
    bytes,
    bytearray,
    Tuple[Union[BinaryIO, str, bytes, bytearray], int],
    Dict[str, Any],
]


class _NormalizedStream:
    """Normalized stream descriptor with size, content type, and cleanup tracking."""

    def __init__(
        self,
        stream_obj: BinaryIO,
        size: int,
        content_type: str,
        should_close: bool,
    ) -> None:
        self.stream_obj = stream_obj
        self.size = size
        self.content_type = content_type
        self.should_close = should_close


class ProgressStreamReader:
    """Wraps a binary stream to track upload progress and enable cancellation."""

    def __init__(
        self,
        stream: BinaryIO,
        total_size: int,
        on_chunk: Optional[Callable[[int], None]] = None,
        cancel_event: Optional[Union[threading.Event, Any]] = None,
    ) -> None:
        self._stream = stream
        self._total_size = total_size
        self._on_chunk = on_chunk
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
            if self._on_chunk is not None:
                try:
                    self._on_chunk(self._bytes_read)
                except Exception:
                    pass
        return chunk

    def __len__(self) -> int:
        return self._total_size


def _normalize_stream(item: Any, default_content_type: str) -> _NormalizedStream:
    """Normalizes any supported stream input into a _NormalizedStream."""
    explicit_size: Optional[int] = None
    content_type = default_content_type
    source: Any = item

    if isinstance(item, tuple) and len(item) == 2:
        source, explicit_size = item
    elif isinstance(item, dict):
        source = item.get("stream")
        explicit_size = item.get("total_size") or item.get("totalSize") or item.get("size")
        content_type = item.get("content_type") or item.get("contentType") or default_content_type

    should_close = False
    if isinstance(source, (bytes, bytearray)):
        raw_bytes = bytes(source)
        calc_size = len(raw_bytes)
        stream_obj: BinaryIO = io.BytesIO(raw_bytes)
    elif isinstance(source, str):
        if os.path.isfile(source):
            calc_size = os.path.getsize(source)
            stream_obj = open(source, "rb")
            should_close = True
        else:
            raw_bytes = source.encode("utf-8")
            calc_size = len(raw_bytes)
            stream_obj = io.BytesIO(raw_bytes)
    elif hasattr(source, "read"):
        stream_obj = source
        if explicit_size is not None:
            calc_size = int(explicit_size)
        else:
            try:
                pos = stream_obj.tell()
                stream_obj.seek(0, io.SEEK_END)
                calc_size = stream_obj.tell() - pos
                stream_obj.seek(pos, io.SEEK_SET)
            except Exception:
                calc_size = 0
    else:
        raise I6ValidationError(
            f"Invalid stream item {type(source)!r}: expected BinaryIO, file path string, or bytes"
        )

    if explicit_size is not None:
        calc_size = int(explicit_size)

    return _NormalizedStream(
        stream_obj=stream_obj,
        size=calc_size,
        content_type=content_type,
        should_close=should_close,
    )


class Ingest:
    """Client for uploading files and data batches to the i6 ingestion service."""

    def __init__(self, services: Services) -> None:
        self._services = services

    def get_url(
        self,
        dataset: str,
        table: str,
        partitions: Dict[str, str],
        amount: int = 1,
        service_options: Optional[Union[ServiceOptions, Dict[str, Any]]] = None,
    ) -> SignedUrlResult:
        """Requests signed upload URLs from the ingest service."""
        resp = self._services.apis().ingestz_get_url(
            req={
                "params": {
                    "dataset": dataset,
                },
                "payload": {
                    "table": table,
                    "partitions": partitions,
                    "amount": amount,
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
            "urls": uploads,
        }

    getUrl = get_url

    def upload(
        self,
        dataset: str,
        table: str,
        partitions: Dict[str, str],
        streams: Optional[Union[StreamInput, Sequence[StreamInput]]] = None,
        stream: Optional[StreamInput] = None,
        file: Optional[StreamInput] = None,
        total_size: Optional[int] = None,
        concurrency: int = 4,
        content_type: str = "application/octet-stream",
        on_progress: Optional[Callable[[UploadProgress], None]] = None,
        onProgress: Optional[Callable[[UploadProgress], None]] = None,
        cancel_event: Optional[threading.Event] = None,
        signal: Optional[Union[threading.Event, Any]] = None,
        timeout: Optional[float] = None,
        service_options: Optional[Union[ServiceOptions, Dict[str, Any]]] = None,
    ) -> UploadResult:
        """Uploads one or more streams concurrently to the ingest service using signed PUT URLs."""
        if not dataset:
            raise I6ValidationError("dataset is required")
        if not table:
            raise I6ValidationError("table is required")

        # Resolve inputs
        raw_items: List[Any] = []
        if streams is not None:
            if isinstance(streams, (list, tuple)):
                raw_items = list(streams)
            else:
                raw_items = [streams]
        elif stream is not None:
            raw_items = [stream]
        elif file is not None:
            raw_items = [file]
        else:
            raise I6ValidationError("streams or stream is required")

        if not raw_items or raw_items == [None]:
            raise I6ValidationError("streams cannot be empty")

        # Apply total_size override if single stream
        if len(raw_items) == 1 and total_size is not None:
            single = raw_items[0]
            if isinstance(single, tuple) and len(single) == 2:
                raw_items[0] = (single[0], total_size)
            elif isinstance(single, dict):
                single_dict = dict(single)
                single_dict["total_size"] = total_size
                raw_items[0] = single_dict
            else:
                raw_items[0] = (single, total_size)

        progress_cb = on_progress or onProgress
        abort_trigger = cancel_event if cancel_event is not None else signal

        if abort_trigger is not None and getattr(abort_trigger, "is_set", lambda: False)():
            raise I6ValidationError("Upload aborted by caller")

        normalized: List[_NormalizedStream] = [
            _normalize_stream(it, content_type) for it in raw_items
        ]
        count = len(normalized)
        total_bytes_all = sum(s.size for s in normalized)

        try:
            url_info = self.get_url(
                dataset=dataset,
                table=table,
                partitions=partitions,
                amount=count,
                service_options=service_options,
            )
            ingest_id = url_info["ingest_id"]
            upload_urls = url_info["urls"]

            if len(upload_urls) < count:
                raise I6ApiError(
                    f"Server returned {len(upload_urls)} upload URLs for {count} streams"
                )

            if abort_trigger is not None and getattr(abort_trigger, "is_set", lambda: False)():
                raise I6ValidationError("Upload aborted by caller")

            lock = threading.Lock()
            bytes_transferred_per_stream = [0] * count
            completed_streams_count = 0

            def make_stream_progress_cb(index: int) -> Callable[[int], None]:
                def cb(loaded_chunk_bytes: int) -> None:
                    with lock:
                        bytes_transferred_per_stream[index] = loaded_chunk_bytes
                        total_loaded = sum(bytes_transferred_per_stream)
                        pct = (
                            (total_loaded / total_bytes_all * 100.0)
                            if total_bytes_all > 0
                            else 100.0
                        )
                        if progress_cb is not None:
                            try:
                                progress_cb({
                                    "loaded": total_loaded,
                                    "total": total_bytes_all if total_bytes_all > 0 else total_loaded,
                                    "percent": min(100.0, pct),
                                    "completed_streams": completed_streams_count,
                                    "total_streams": count,
                                    "stream_index": index,
                                })
                            except Exception:
                                pass

                return cb

            request_timeout = timeout or self._services.config().timeout or 120.0
            config = self._services.config()

            def upload_single_stream(index: int, item: _NormalizedStream, target_url: str) -> None:
                nonlocal completed_streams_count
                on_chunk_progress = make_stream_progress_cb(index)

                reader = ProgressStreamReader(
                    stream=item.stream_obj,
                    total_size=item.size,
                    on_chunk=on_chunk_progress,
                    cancel_event=abort_trigger,
                )

                headers = {
                    "Content-Type": item.content_type,
                    "Content-Length": str(item.size),
                }

                if config.fetcher is not None:
                    body_bytes = reader.read()
                    status, _, text = config.fetcher(
                        target_url,
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
                        url=target_url,
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

                with lock:
                    completed_streams_count += 1
                on_chunk_progress(item.size)

            workers = max(1, min(concurrency, count))
            with concurrent.futures.ThreadPoolExecutor(max_workers=workers) as executor:
                futures = [
                    executor.submit(upload_single_stream, i, normalized[i], upload_urls[i])
                    for i in range(count)
                ]
                for f in concurrent.futures.as_completed(futures):
                    f.result()

            return {
                "ingest_id": ingest_id,
                "count": count,
            }
        finally:
            for it in normalized:
                if it.should_close:
                    try:
                        it.stream_obj.close()
                    except Exception:
                        pass
