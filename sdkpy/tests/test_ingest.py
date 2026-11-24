"""Unit tests for Ingest client streaming, schema, and cancellation."""

import io
import json
import threading
import unittest
from typing import Any, Dict, List, Tuple

from i6sdk import I6Client, I6SdkConfig
from i6sdk.exceptions import I6ValidationError


class TestIngest(unittest.TestCase):
    def test_get_url_matches_backend_schema(self) -> None:
        calls = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            calls.append({"url": url, "method": method, "headers": headers, "body": body})
            if "/api/ingest/dataset/dectest/get-url" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "ingest_id": "ing-123",
                    "uploads": ["https://storage.test/upload/1"],
                })
            return 404, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        res = client.ingest().get_url(
            dataset="dectest",
            table="inbox_users",
            partitions={"date": "2026-11-24"},
        )

        self.assertEqual(res["ingest_id"], "ing-123")
        self.assertEqual(res["url"], "https://storage.test/upload/1")
        self.assertEqual(len(calls), 1)

        body_dict = json.loads(calls[0]["body"].decode("utf-8") if isinstance(calls[0]["body"], (bytes, bytearray)) else calls[0]["body"])
        self.assertEqual(body_dict.get("table"), "inbox_users")
        self.assertEqual(body_dict.get("partitions"), {"date": "2026-11-24"})
        self.assertEqual(body_dict.get("amount"), 1)

    def test_upload_binary_io_stream(self) -> None:
        calls = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            calls.append({"url": url, "method": method, "headers": headers, "body": body})
            if "/api/ingest/dataset/dectest/get-url" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "ingest_id": "ing-stream",
                    "uploads": ["https://storage.test/upload/stream"],
                })
            elif "storage.test/upload/stream" in url:
                return 200, {}, ""
            return 404, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        stream = io.BytesIO(b"binary data payload for streaming")
        res = client.ingest().upload(
            dataset="dectest",
            table="inbox_users",
            partitions={"date": "2026-11-24"},
            stream=stream,
        )

        self.assertEqual(res["ingest_id"], "ing-stream")

    def test_upload_progress_callback(self) -> None:
        progress_events = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            if "get-url" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "ingest_id": "ing-prog",
                    "uploads": ["https://storage.test/upload/prog"],
                })
            return 200, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        payload = b"A" * 1000
        res = client.ingest().upload(
            dataset="dectest",
            table="inbox_users",
            partitions={"date": "2026-11-24"},
            stream=payload,
            on_progress=progress_events.append,
        )

        self.assertEqual(res["ingest_id"], "ing-prog")
        self.assertTrue(len(progress_events) >= 1)
        last = progress_events[-1]
        self.assertEqual(last["loaded"], 1000)
        self.assertEqual(last["total"], 1000)
        self.assertEqual(last["percent"], 100.0)

    def test_upload_cancellation_pre_abort(self) -> None:
        cancel = threading.Event()
        cancel.set()

        client = I6Client(I6SdkConfig())
        with self.assertRaises(I6ValidationError):
            client.ingest().upload(
                dataset="dectest",
                table="inbox_users",
                partitions={"date": "2026-11-24"},
                stream=b"data",
                cancel_event=cancel,
            )

    def test_upload_cancellation_mid_stream(self) -> None:
        cancel = threading.Event()

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            if "get-url" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "ingest_id": "ing-cancel",
                    "uploads": ["https://storage.test/upload/cancel"],
                })
            return 200, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        def on_prog(p: Any) -> None:
            # Set abort when first chunk is processed
            cancel.set()

        with self.assertRaises(I6ValidationError):
            client.ingest().upload(
                dataset="dectest",
                table="inbox_users",
                partitions={"date": "2026-11-24"},
                stream=b"X" * 100000,
                on_progress=on_prog,
                cancel_event=cancel,
            )

    def test_upload_file_path_string(self) -> None:
        import tempfile
        from pathlib import Path

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            if "get-url" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "ingest_id": "ing-file",
                    "uploads": ["https://storage.test/upload/file"],
                })
            return 200, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        with tempfile.TemporaryDirectory() as tmpdir:
            file_path = str(Path(tmpdir) / "test_data.csv")
            with open(file_path, "w") as f:
                f.write("id,name\n1,alice\n")

            res = client.ingest().upload(
                dataset="dectest",
                table="inbox_users",
                partitions={"date": "2026-11-24"},
                stream=file_path,
            )
            self.assertEqual(res["ingest_id"], "ing-file")

    def test_upload_validation_errors(self) -> None:
        client = I6Client(I6SdkConfig())

        with self.assertRaises(I6ValidationError):
            client.ingest().upload(dataset="", table="t", partitions={}, stream=b"data")

        with self.assertRaises(I6ValidationError):
            client.ingest().upload(dataset="d", table="", partitions={}, stream=b"data")

        with self.assertRaises(I6ValidationError):
            client.ingest().upload(dataset="d", table="t", partitions={}, stream=None)


if __name__ == "__main__":
    unittest.main()
