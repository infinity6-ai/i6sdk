"""Unit tests for Infer client schema and telemetry events."""

import json
import unittest
from typing import Any, Dict, List, Tuple

from i6sdkpy import I6Client, I6SdkConfig
from i6sdkpy.exceptions import I6ValidationError


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

    def test_run_validation_options(self) -> None:
        client = I6Client(I6SdkConfig())
        base = {
            "dataset": "dectest",
            "model_name": "relevance-fashion-fbt",
            "partitions": {"store_id": "1", "version": "v1"},
            "session_id": "sess-1",
            "user_id": "u-1",
            "params": {},
        }
        with self.assertRaises(I6ValidationError):
            client.infer().run({})

        for key in ["dataset", "model_name", "session_id", "user_id"]:
            invalid = dict(base)
            invalid[key] = ""
            with self.assertRaises(I6ValidationError):
                client.infer().run(invalid)

        with self.assertRaises(I6ValidationError):
            client.infer().run({**base, "partitions": "not-a-dict"})  # type: ignore

        with self.assertRaises(I6ValidationError):
            client.infer().run({**base, "params": None})

        with self.assertRaises(I6ValidationError):
            client.infer().run({**base, "limit": 0})

    def test_select_results_validation_errors(self) -> None:
        client = I6Client(I6SdkConfig())
        self.assertFalse(hasattr(client.infer(), "selectResults"))
        valid = {
            "dataset": "dectest",
            "model_name": "m1",
            "partitions": {"store_id": "1"},
            "req": {},
            "results": {"results": []},
            "rec_id": "rec-1",
            "session_id": "sess-1",
            "user_id": "user-1",
        }

        with self.assertRaises(I6ValidationError):
            client.infer().select_results({})

        with self.assertRaises(I6ValidationError):
            client.infer().select_results({**valid, "rec_id": ""})

        with self.assertRaises(I6ValidationError):
            client.infer().select_results({**valid, "session_id": ""})

        with self.assertRaises(I6ValidationError):
            client.infer().select_results({**valid, "user_id": ""})

        with self.assertRaises(TypeError):
            client.infer().select_results({**valid, "select_results": "not-callable"})

        with self.assertRaises(TypeError):
            client.infer().select_results({**valid, "select_results": lambda r, res: []})

        with self.assertRaises(TypeError):
            client.infer().select_results({**valid, "select_results": lambda r, res: {"results": [{}]}})

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

    def test_run_dispatches_recraw_event(self) -> None:
        import sqlite3
        import tempfile
        from pathlib import Path
        events = []

        with tempfile.TemporaryDirectory() as tmpdir:
            db_path = Path(tmpdir) / "test_model.sqlite"
            conn = sqlite3.connect(str(db_path))
            cur = conn.cursor()
            cur.execute("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT)")
            cur.execute(
                "INSERT INTO namedqueries (name, query) VALUES ('model-relevance-fashion-fbt-input', 'SELECT json_object(\"inputs\", ?1) as ret')"
            )
            cur.execute(
                "INSERT INTO namedqueries (name, query) VALUES ('model-relevance-fashion-fbt-output', 'SELECT json_object(\"results\", json_array(json_object(\"sku_id\", \"sku-100\", \"rating\", 0.95))) as ret WHERE ?1 IS NOT NULL AND ?2 > 0')"
            )
            conn.commit()
            conn.close()

            def mock_fetcher(
                url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
            ) -> Tuple[int, Dict[str, str], str]:
                if "predict" in url:
                    return 200, {"content-type": "application/json"}, json.dumps({
                        "outputs": json.dumps({"raw": 123})
                    })
                elif "relevance-fashion-event" in url:
                    payload = json.loads(body.decode("utf-8") if isinstance(body, (bytes, bytearray)) else body)
                    events.append(payload)
                    return 200, {"content-type": "application/json"}, "{}"
                return 404, {}, ""

            client = I6Client(I6SdkConfig(
                serv_urls={"decsuite": "https://api.i6.test", "tfinfer": "https://api.i6.test"},
                fetcher=mock_fetcher,
            ))

            # Monkey-patch update to return our local SQLite db directly
            client.infer().update = lambda opts: db_path  # type: ignore

            results = client.infer().run({
                "dataset": "dectest",
                "model_name": "relevance-fashion-fbt",
                "partitions": {"store_id": "1", "version": "v1"},
                "session_id": "sess-test",
                "user_id": "user-test",
                "params": {"cart": []},
            })

            self.assertIn("rec_id", results)
            rec_id = results["rec_id"]
            self.assertEqual(len(events), 1)
            event = events[0]
            self.assertEqual(event["store_id"], "1")
            self.assertEqual(event["session_id"], "sess-test")
            self.assertEqual(event["user_id"], "user-test")
            self.assertIn("recraw", event)
            recraw = event["recraw"]
            self.assertEqual(recraw["rec_id"], rec_id)
            self.assertEqual(recraw["model_id"], "relevance-fashion-fbt")
            self.assertEqual(recraw["model_version"], "v1")
            recraw_payload = json.loads(recraw["payload"])
            self.assertEqual(recraw_payload["rec_id"], rec_id)


if __name__ == "__main__":
    unittest.main()
