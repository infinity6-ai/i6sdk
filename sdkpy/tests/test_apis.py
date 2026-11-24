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
        resp = client.apis().dsz_list_datasets()
        self.assertTrue(resp.ok)
        data = resp.parse()
        payload = data.get("payload") or {}
        datasets = payload.get("datasets") or []
        self.assertEqual(len(datasets), 1)
        self.assertEqual(datasets[0].get("id"), "ds_fashion")

        self.assertEqual(last_request["method"], "GET")
        self.assertIn("bearer-token-xyz", last_request["headers"]["Authorization"])

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


if __name__ == "__main__":
    unittest.main()
