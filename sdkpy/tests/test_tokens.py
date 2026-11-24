"""Unit tests for token storage and token manager."""

import tempfile
import time
import unittest
from pathlib import Path

from i6sdk import I6Client, I6SdkConfig, create_client
from i6sdk.tokens import FileTokenStorage, MemoryTokenStorage, UserToken


class TestTokens(unittest.TestCase):
    def test_memory_token_storage(self) -> None:
        storage = MemoryTokenStorage()
        self.assertIsNone(storage.get_item("key1"))
        storage.set_item("key1", "val1")
        self.assertEqual(storage.get_item("key1"), "val1")
        storage.remove_item("key1")
        self.assertIsNone(storage.get_item("key1"))

    def test_file_token_storage(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            file_path = Path(tmpdir) / "tokens.json"
            storage = FileTokenStorage(file_path)
            storage.set_item("user", "user-token-abc")
            self.assertEqual(storage.get_item("user"), "user-token-abc")

            # Check permissions are 0600
            mode = file_path.stat().st_mode & 0o777
            self.assertEqual(mode, 0o600)

            # Re-read from another instance
            storage2 = FileTokenStorage(file_path)
            self.assertEqual(storage2.get_item("user"), "user-token-abc")

    def test_tokens_manager(self) -> None:
        client = I6Client(I6SdkConfig())
        self.assertIsNone(client.tokens().get_token())

        user_token: UserToken = {
            "user_access_token": "dynamic-user-token",
            "user_refresh_token": "refresh-123",
            "user_id": "u-1",
            "expires_at": 9999999999999,
        }
        client.tokens().set_user_token(user_token)
        self.assertEqual(client.tokens().get_token(), "dynamic-user-token")

        client.tokens().remove_user_token()
        self.assertIsNone(client.tokens().get_token())

    def test_tokens_storage_string_config(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            token_path = Path(tmpdir) / "test.tokens.json"
            client = create_client(storage=str(token_path))
            storage = client.tokens().storage
            self.assertIsInstance(storage, FileTokenStorage)
            assert isinstance(storage, FileTokenStorage)
            self.assertEqual(str(storage.path), str(token_path))

    def test_silent_renewal_when_access_token_expired(self) -> None:
        calls = []

        def mock_fetcher(
            url: str, method: str, headers: dict, body: any, timeout: float
        ) -> tuple[int, dict, str]:
            calls.append({"url": url, "method": method, "headers": headers})
            if "/api/userz/refresh-user" in url:
                import json
                return 200, {"content-type": "application/json"}, json.dumps({
                    "user_access_token": "new-refreshed-token",
                    "user_refresh_token": "rotated-refresh-token",
                    "user_id": "u-1",
                    "expires_at": 9999999999999,
                })
            return 404, {}, ""

        client = I6Client(I6SdkConfig(fetcher=mock_fetcher))
        # Set expired access token but valid refresh token
        client.tokens().set_user_token({
            "user_access_token": "old-expired-token",
            "user_refresh_token": "valid-refresh-token",
            "user_id": "u-1",
            "expires_at": 1000,  # Far in the past
        })

        # get_valid_access_token should detect expiration and refresh silently
        token = client.tokens().get_valid_access_token()
        self.assertEqual(token, "new-refreshed-token")
        self.assertEqual(len(calls), 1)

    def test_is_access_token_expired(self) -> None:
        client = I6Client(I6SdkConfig())
        tokens = client.tokens()

        self.assertTrue(tokens.is_access_token_expired(None))
        self.assertTrue(tokens.is_access_token_expired({}))

        # Future token (not expired)
        fresh_token = {
            "user_access_token": "fresh",
            "expires_at": int(time.time() * 1000) + 60000,
        }
        self.assertFalse(tokens.is_access_token_expired(fresh_token))

        # Expiring within skew
        almost_expired = {
            "user_access_token": "almost",
            "expires_at": int(time.time() * 1000) + 10000,  # 10s left
        }
        self.assertTrue(tokens.is_access_token_expired(almost_expired, skew_seconds=30.0))
        self.assertFalse(tokens.is_access_token_expired(almost_expired, skew_seconds=5.0))

        # Past token
        expired = {
            "user_access_token": "expired",
            "expires_at": int(time.time() * 1000) - 10000,
        }
        self.assertTrue(tokens.is_access_token_expired(expired))

    def test_is_refresh_token_expired(self) -> None:
        client = I6Client(I6SdkConfig())
        tokens = client.tokens()

        self.assertTrue(tokens.is_refresh_token_expired(None))
        self.assertTrue(tokens.is_refresh_token_expired({}))

        valid_rf = {
            "user_access_token": "acc",
            "user_refresh_token": "rf",
            "refresh_token_expires_at": int(time.time() * 1000) + 100000,
        }
        self.assertFalse(tokens.is_refresh_token_expired(valid_rf))

        expired_rf = {
            "user_access_token": "acc",
            "user_refresh_token": "rf",
            "refresh_token_expires_at": int(time.time() * 1000) - 1000,
        }
        self.assertTrue(tokens.is_refresh_token_expired(expired_rf))

    def test_get_valid_dataset_access_token_silent_renewal(self) -> None:
        calls = []

        def mock_fetcher(
            url: str, method: str, headers: dict, body: any, timeout: float
        ) -> tuple[int, dict, str]:
            calls.append({"url": url, "method": method, "headers": headers})
            if "/api/userz/refresh-user" in url:
                import json
                return 200, {"content-type": "application/json"}, json.dumps({
                    "user_access_token": "new-user-token",
                    "user_refresh_token": "rotated-refresh-token",
                    "user_id": "u-1",
                    "expires_at": 9999999999999,
                })
            elif "/api/userz/dataset-access-token/test-ds" in url:
                import json
                # Verify that authorization header used the renewed user access token
                self.assertEqual(headers.get("Authorization"), "Bearer new-user-token")
                return 200, {"content-type": "application/json"}, json.dumps({
                    "dataset": "test-ds",
                    "access_token": "ds-token-xyz",
                    "expires_at": 9999999999999,
                    "grants": 15,
                })
            return 404, {}, ""

        client = I6Client(I6SdkConfig(fetcher=mock_fetcher))
        client.tokens().set_user_token({
            "user_access_token": "old-expired-token",
            "user_refresh_token": "valid-refresh-token",
            "user_id": "u-1",
            "expires_at": 1000,
        })

        ds_token = client.tokens().get_valid_dataset_access_token("test-ds")
        self.assertEqual(ds_token, "ds-token-xyz")
        self.assertEqual(len(calls), 2)

    def test_concurrent_silent_renewal_deduplication(self) -> None:
        import concurrent.futures
        import threading
        calls = []
        lock = threading.Lock()

        def mock_fetcher(
            url: str, method: str, headers: dict, body: any, timeout: float
        ) -> tuple[int, dict, str]:
            if "/api/userz/refresh-user" in url:
                import json
                import time
                time.sleep(0.05)  # simulate network delay
                with lock:
                    calls.append(url)
                return 200, {"content-type": "application/json"}, json.dumps({
                    "user_access_token": "deduped-token",
                    "user_refresh_token": "rotated-rf",
                    "user_id": "u-1",
                    "expires_at": 9999999999999,
                })
            return 404, {}, ""

        client = I6Client(I6SdkConfig(fetcher=mock_fetcher))
        client.tokens().set_user_token({
            "user_access_token": "expired-token",
            "user_refresh_token": "valid-rf",
            "user_id": "u-1",
            "expires_at": 1000,
        })

        with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
            futures = [executor.submit(client.tokens().get_valid_access_token) for _ in range(5)]
            results = [f.result() for f in futures]

        for res in results:
            self.assertEqual(res, "deduped-token")
        # Exactly one network refresh should have occurred
        self.assertEqual(len(calls), 1)



if __name__ == "__main__":
    unittest.main()
