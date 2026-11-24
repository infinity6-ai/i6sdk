"""Unit tests for PKCE and auth helpers."""

import json
import unittest
from typing import Any, Dict

import i6sdk
from i6sdk.auth import generate_code_challenge, generate_code_verifier


class TestAuth(unittest.TestCase):
    def test_pkce_generation(self) -> None:
        verifier = generate_code_verifier()
        self.assertGreaterEqual(len(verifier), 43)

        challenge = generate_code_challenge(verifier)
        self.assertTrue(len(challenge) > 0)
        # S256 challenge is base64url encoded without padding
        self.assertNotIn("=", challenge)
        self.assertNotIn("+", challenge)
        self.assertNotIn("/", challenge)

    def test_me_when_not_logged_in(self) -> None:
        client = i6sdk.create_client()
        self.assertIsNone(client.auth().me())

    def test_me_when_logged_in(self) -> None:
        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> tuple[int, Dict[str, str], str]:
            return 200, {"content-type": "application/json"}, json.dumps({"user_id": "u-123", "email": "test@infinity6.ai"})

        client = i6sdk.I6Client(i6sdk.I6SdkConfig(fetcher=mock_fetcher))
        client.tokens().set_user_token({"user_access_token": "valid-token"})
        user = client.auth().me()
        self.assertIsNotNone(user)
        assert user is not None
        self.assertEqual(user.get("user_id"), "u-123")

    def test_device_code_flow(self) -> None:
        calls = []

        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> tuple[int, Dict[str, str], str]:
            calls.append(url)
            if "/api/userz/auth/device/code" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "device_code": "dev-123",
                    "user_code": "CODE-XYZ",
                    "verification_uri": "https://auth.i6.test/device",
                    "verification_uri_complete": "https://auth.i6.test/device?code=CODE-XYZ",
                    "interval_sec": 1,
                    "expires_in_sec": 60,
                })
            elif "/api/userz/auth/device/token" in url:
                return 200, {"content-type": "application/json"}, json.dumps({
                    "status": "authorized",
                    "user_access_token": "dev-token-abc",
                    "user_refresh_token": "dev-refresh-xyz",
                    "user_id": "u-dev-1",
                    "expires_at": 9999999999999,
                })
            return 404, {}, ""

        client = i6sdk.I6Client(i6sdk.I6SdkConfig(fetcher=mock_fetcher))
        flow = client.auth().device_code_start()
        self.assertEqual(flow["device_code"], "dev-123")
        self.assertEqual(flow["user_code"], "CODE-XYZ")

        token = client.auth().device_code_poll(device_code="dev-123", interval=0.01, timeout=2.0)
        self.assertEqual(token.get("user_access_token"), "dev-token-abc")
        self.assertEqual(client.tokens().get_token(), "dev-token-abc")


if __name__ == "__main__":
    unittest.main()
