"""Unit tests for PKCE and auth helpers."""

import unittest

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


if __name__ == "__main__":
    unittest.main()
