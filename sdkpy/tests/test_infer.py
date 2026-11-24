"""Unit tests for Infer client schema and telemetry events."""

import json
import unittest
from typing import Any, Dict, List, Tuple

from i6sdk import I6Client, I6SdkConfig
from i6sdk.exceptions import I6ValidationError


class TestInfer(unittest.TestCase):
    def test_sqlite_url_matches_server_schema(self) -> None:
        calls = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            calls.append({"url": url, "method": method, "headers": headers, "body": body})
            if "/api/model/dataset/dectest/signed-get/relevance-fashion-fbt" in url:
                # Backend returns urls list
                return 200, {"content-type": "application/json"}, json.dumps({
                    "urls": [
                        {"name": "sqlite.bin", "url": "https://storage.test/sqlite.bin"}
                    ]
                })
            return 404, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        sqlite_url = client.infer()._sqlite_url({
            "dataset": "dectest",
            "model_name": "relevance-fashion-fbt",
            "partitions": {"store_id": "1", "version": "v1"},
        })

        self.assertEqual(sqlite_url, "https://storage.test/sqlite.bin")
        self.assertEqual(len(calls), 1)
        # Verify request sent POST body with model_version and model_partitions
        sent_body = json.loads(calls[0]["body"].decode("utf-8") if isinstance(calls[0]["body"], (bytes, bytearray)) else calls[0]["body"])
        self.assertEqual(sent_body.get("model_version"), "v1")
        self.assertEqual(sent_body.get("model_partitions"), {"store_id": "1", "version": "v1"})

    def test_run_requires_session_id_and_user_id(self) -> None:
        client = I6Client(I6SdkConfig())
        opts = {
            "dataset": "dectest",
            "model_name": "relevance-fashion-fbt",
            "partitions": {"store_id": "1", "version": "v1"},
            "params": {},
        }
        with self.assertRaises(I6ValidationError):
            client.infer().run(opts)

    def test_select_results_emits_recshow_event(self) -> None:
        events = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            if "relevance-fashion-event" in url:
                payload = json.loads(body.decode("utf-8") if isinstance(body, (bytes, bytearray)) else body)
                events.append(payload)
                return 200, {"content-type": "application/json"}, "{}"
            return 404, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        ))

        req = {"input": {"dataset": "dectest", "skus": [{"sku_id": "cart-1"}]}}
        results = {
            "results": [
                {"skus": [{"sku_id": "cart-1"}, {"sku_id": "rec-1"}], "rating": 0.9},
                {"skus": [{"sku_id": "rec-2"}], "rating": 0.8},
            ]
        }

        selected = client.infer().select_results({
            "dataset": "dectest",
            "model_name": "relevance-fashion-fbt",
            "partitions": {"store_id": "1", "version": "v1"},
            "req": req,
            "results": results,
            "rec_id": "rec-123",
            "session_id": "sess-abc",
            "user_id": "user-xyz",
        })

        self.assertEqual(selected["results"], [
            {"sku_id": "rec-1", "rating": 0.9},
            {"sku_id": "rec-2", "rating": 0.8},
        ])
        self.assertEqual(len(events), 1)
        event = events[0]
        self.assertEqual(event["store_id"], "1")
        self.assertEqual(event["session_id"], "sess-abc")
        self.assertEqual(event["user_id"], "user-xyz")
        self.assertIn("recshow", event)
        recshow = event["recshow"]
        self.assertEqual(recshow["rec_id"], "rec-123")
        self.assertEqual(recshow["model_id"], "relevance-fashion-fbt")
        self.assertEqual(recshow["model_version"], "v1")
        shown_payload = json.loads(recshow["payload"])
        self.assertEqual(shown_payload["skus"], [{"sku_id": "rec-1"}, {"sku_id": "rec-2"}])


if __name__ == "__main__":
    unittest.main()
