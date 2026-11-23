"""Unit tests for HTTP invoker."""

import json
import unittest
from typing import Any, Dict, Tuple

from i6sdk import I6Client, I6SdkConfig
from i6sdk.exceptions import I6ApiError, I6ConnectionError
from i6sdk.invoker import Invoker, Resp, code_from_status


class TestInvoker(unittest.TestCase):
    def test_code_from_status(self) -> None:
        self.assertEqual(code_from_status(401), "unauthenticated")
        self.assertEqual(code_from_status(403), "forbidden")
        self.assertEqual(code_from_status(404), "not_found")
        self.assertEqual(code_from_status(500), "unavailable")
        self.assertEqual(code_from_status(503), "unavailable")
        self.assertEqual(code_from_status(400), "bad_request")

    def test_resp_success(self) -> None:
        raw_json = json.dumps({"dataset": "my-ds", "status": "active"})
        resp = Resp(
            status=200,
            headers={"content-type": "application/json"},
            parser=lambda: json.loads(raw_json),
            text=raw_json,
        )
        self.assertTrue(resp.ok)
        self.assertIsNone(resp.error)
        data = resp.parse()
        self.assertEqual(data["dataset"], "my-ds")

    def test_resp_error(self) -> None:
        raw_error = json.dumps({"error": "Dataset not found"})
        resp = Resp(
            status=404,
            headers={"content-type": "application/json"},
            parser=lambda: (_ for _ in ()).throw(
                I6ApiError(message="Not Found", status=404, code="not_found", cause=raw_error)
            ),
            text=raw_error,
        )
        self.assertFalse(resp.ok)
        self.assertIsNotNone(resp.error)
        self.assertEqual(resp.error["code"], "not_found")
        self.assertEqual(resp.error["status"], 404)

        with self.assertRaises(I6ApiError) as ctx:
            resp.parse()
        self.assertEqual(ctx.exception.status, 404)
        self.assertEqual(ctx.exception.code, "not_found")

    def test_custom_fetcher_hook(self) -> None:
        recorded_calls = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> Tuple[int, Dict[str, str], str]:
            recorded_calls.append(
                {"url": url, "method": method, "headers": headers, "body": body}
            )
            return 200, {"content-type": "application/json"}, json.dumps({"count": 42})

        config = I6SdkConfig(
            base_url="https://api.i6.test",
            token="secret-token",
            fetcher=mock_fetcher,
        )
        client = I6Client(config)
        resp = client.invoker().invoke_result(
            method="POST",
            path="api/test",
            query={"limit": 10},
            json={"filter": "all"},
        )
        self.assertTrue(resp.ok)
        self.assertEqual(resp.parse(), {"count": 42})

        self.assertEqual(len(recorded_calls), 1)
        call = recorded_calls[0]
        self.assertEqual(call["url"], "https://api.i6.test/api/test?limit=10")
        self.assertEqual(call["method"], "POST")
        self.assertEqual(call["headers"]["Authorization"], "Bearer secret-token")
        self.assertEqual(call["headers"]["Content-Type"], "application/json")


if __name__ == "__main__":
    unittest.main()
