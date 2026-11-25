"""Unit tests for generated APIs and SdkFuncs methods."""

import json
import unittest
from typing import Any, Dict, Tuple

from i6sdk import I6Client, I6SdkConfig


class TestApis(unittest.TestCase):
    def test_generated_apis_call(self) -> None:
        last_request = {}

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            last_request["url"] = url
            last_request["method"] = method
            last_request["headers"] = headers
            last_request["body"] = body
            resp_body = {
                "payload": {
                    "datasets": [
                        {"id": "ds_fashion", "name": "Fashion Dataset"}
                    ]
                }
            }
            return 200, {"content-type": "application/json"}, json.dumps(resp_body)

        config = I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        )
        client = I6Client(config)
        client.tokens().set_user_token({"user_access_token": "bearer-token-xyz"})

        # Call generated method
        resp = client.apis().dsz_list_datasets(req={"params": {"userId": "user_42"}})
        self.assertTrue(resp.ok)
        data = resp.parse()
        payload = data.get("payload") or {}
        datasets = payload.get("datasets") or []
        self.assertEqual(len(datasets), 1)
        self.assertEqual(datasets[0].get("id"), "ds_fashion")

        self.assertEqual(last_request["method"], "GET")
        self.assertIn("/api/ds/user-id/user_42/datasets", last_request["url"])
        self.assertIn("bearer-token-xyz", last_request["headers"]["Authorization"])

    def test_generated_pipez_start_batch_ingest_envelope(self) -> None:
        last_request = {}

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            last_request["url"] = url
            last_request["method"] = method
            last_request["body"] = body
            resp_body = {"payload": {"id": "run_batch_123"}}
            return 200, {"content-type": "application/json"}, json.dumps(resp_body)

        config = I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        )
        client = I6Client(config)

        resp = client.apis().pipez_start_relevance_fashion_relevance_fashion_batch_ingest(
            req={
                "params": {"dataset": "my-fashion-dataset"},
                "payload": {"ingest_id": "ingest_456", "table": "catalog"},
            }
        )
        self.assertTrue(resp.ok)
        self.assertEqual(resp.parse(), {"payload": {"id": "run_batch_123"}})
        self.assertIn("/api/pipe/dataset/my-fashion-dataset/start/relevance_fashion/relevance-fashion-batch-ingest", last_request["url"])
        self.assertEqual(last_request["method"], "POST")
        raw_body = last_request.get("body")
        parsed_body = json.loads(raw_body.decode("utf-8") if isinstance(raw_body, bytes) else raw_body)
        self.assertEqual(parsed_body, {"ingest_id": "ingest_456", "table": "catalog"})

    def test_generated_pipez_start_main_envelope(self) -> None:
        last_request = {}

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            last_request["url"] = url
            last_request["method"] = method
            resp_body = {"payload": {"id": "run_main_789"}}
            return 200, {"content-type": "application/json"}, json.dumps(resp_body)

        config = I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        )
        client = I6Client(config)

        resp = client.apis().pipez_start_relevance_fashion_relevance_fashion_main(
            req={
                "params": {"dataset": "my-fashion-dataset"},
            }
        )
        self.assertTrue(resp.ok)
        self.assertEqual(resp.parse(), {"payload": {"id": "run_main_789"}})
        self.assertIn("/api/pipe/dataset/my-fashion-dataset/start/relevance_fashion/relevance-fashion-main", last_request["url"])
        self.assertEqual(last_request["method"], "POST")

    def test_generated_api_path_interpolation(self) -> None:
        last_request = {}

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            last_request["url"] = url
            last_request["method"] = method
            return 200, {"content-type": "application/json"}, json.dumps({"status": "ok"})

        config = I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test"},
            fetcher=mock_fetcher,
        )
        client = I6Client(config)
        client.tokens().set_user_token({"user_access_token": "token-123"})

        resp = client.apis().sample_fraction(
            req={
                "params": {"numerator": 10.0, "denominator": 2.0},
                "query": {"precision": 2},
            }
        )
        self.assertTrue(resp.ok)
        self.assertIn("/10.0/2.0", last_request["url"])
        self.assertIn("precision=2", last_request["url"])

        # Also verify camelCase alias
        self.assertTrue(hasattr(client.apis(), "sampleFraction"))

    def test_funcs_accessor(self) -> None:
        client = I6Client(I6SdkConfig())
        # client.funcs() must return an SdkFuncs instance
        funcs_instance = client.funcs()
        self.assertIsNotNone(funcs_instance)
        self.assertEqual(client.funcs(), funcs_instance)

        # sdk_funcs and sdkFuncs must be removed
        self.assertFalse(hasattr(client, "sdk_funcs"))
        self.assertFalse(hasattr(client, "sdkFuncs"))


if __name__ == "__main__":
    unittest.main()
