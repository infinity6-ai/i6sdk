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
        self.assertEqual(user.get("user_id"), "u-123")

    def test_me_when_unauthorized_401(self) -> None:
        def mock_fetcher(
            url: str, method: str, headers: Dict[str, str], body: Any, timeout: float
        ) -> tuple[int, Dict[str, str], str]:
            return 401, {"content-type": "application/json"}, json.dumps({"error": "unauthorized"})

        client = i6sdk.I6Client(i6sdk.I6SdkConfig(fetcher=mock_fetcher))
        client.tokens().set_user_token({"user_access_token": "expired-token"})
        self.assertIsNone(client.auth().me())


if __name__ == "__main__":
    unittest.main()
