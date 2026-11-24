"""Ingest client for uploading files and batches to the ingestion service."""

from __future__ import annotations

import pathlib
import urllib.request
from typing import TYPE_CHECKING, Any, Dict, Optional, TypedDict, Union

if TYPE_CHECKING:
    from ..invoker import ServiceOptions
    from ..services import Services

from ..exceptions import I6ApiError, I6ValidationError


class UploadResult(TypedDict):
    """Result of a successful upload."""

    ingest_id: str


class SignedUrlResult(TypedDict):
    """Result of a signed upload URL request."""

    ingest_id: str
    url: str


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
                    "table": table,
                },
                "query": partitions,
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
        return {
            "ingest_id": data.get("ingest_id", ""),
            "url": data.get("url", ""),
        }

    def upload(
        self,
        dataset: str,
        table: str,
        partitions: Dict[str, str],
        file: Union[pathlib.Path, str, bytes],
        content_type: str = "application/octet-stream",
        service_options: Optional[Union[ServiceOptions, Dict[str, Any]]] = None,
    ) -> UploadResult:
        """Uploads a file or bytes buffer to the ingest service using a signed URL."""
        if not dataset:
            raise I6ValidationError("dataset is required")
        if not table:
            raise I6ValidationError("table is required")

        url_info = self.get_url(
            dataset=dataset,
            table=table,
            partitions=partitions,
            service_options=service_options,
        )
        signed_url = url_info["url"]
        ingest_id = url_info["ingest_id"]

        # Read file payload
        if isinstance(file, (str, pathlib.Path)):
            p = pathlib.Path(file)
            if not p.is_file():
                raise FileNotFoundError(f"File not found: {file}")
            payload_bytes = p.read_bytes()
        elif isinstance(file, (bytes, bytearray)):
            payload_bytes = bytes(file)
        else:
            raise I6ValidationError("file must be a Path, file path string, or bytes")

        # Execute PUT to signed URL
        req = urllib.request.Request(
            url=signed_url,
            data=payload_bytes,
            headers={
                "Content-Type": content_type,
                "Content-Length": str(len(payload_bytes)),
            },
            method="PUT",
        )
        try:
            with urllib.request.urlopen(req, timeout=120.0) as resp:
                if not (200 <= resp.status < 300):
                    raise I6ApiError(
                        message=f"Upload PUT failed with HTTP {resp.status}",
                        status=resp.status,
                    )
        except Exception as e:
            raise I6ApiError(message=f"Upload failed: {e}") from e

        return {"ingest_id": ingest_id}
