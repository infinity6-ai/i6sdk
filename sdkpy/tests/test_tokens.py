"""Unit tests for token storage and token manager."""

import tempfile
import time
import unittest
from pathlib import Path

from i6sdk import I6Client, I6SdkConfig, create_client
from i6sdk.tokens import FileTokenStorage, MemoryTokenStorage, SqliteTokenStorage, UserToken


class TestTokens(unittest.TestCase):
    def test_memory_token_storage(self) -> None:
        storage = MemoryTokenStorage()
        self.assertIsNone(storage.get_item("key1"))
        storage.set_item("key1", "val1")
        self.assertEqual(storage.get_item("key1"), "val1")

        # Prefix list & delete
        storage.set_item("i6_ds:orders", "ord")
        storage.set_item("i6_ds:users", "usr")
        self.assertEqual(len(storage.list_prefix("i6_ds:")), 2)
        storage.delete_prefix("i6_ds:")
        self.assertEqual(len(storage.list_prefix("i6_ds:")), 0)

        storage.remove_item("key1")
        self.assertIsNone(storage.get_item("key1"))

    def test_sqlite_token_storage(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            file_path = Path(tmpdir) / "tokens.db"
            storage = SqliteTokenStorage(file_path)
            storage.set_item("user", "user-token-abc")
            self.assertEqual(storage.get_item("user"), "user-token-abc")

            # Check permissions are 0600
            mode = file_path.stat().st_mode & 0o777
            self.assertEqual(mode, 0o600)

            # Prefix list & delete
            storage.set_item("i6_ds:orders", "ord")
            storage.set_item("i6_ds:users", "usr")
            self.assertEqual(len(storage.list_prefix("i6_ds:")), 2)
            storage.delete_prefix("i6_ds:")
            self.assertEqual(len(storage.list_prefix("i6_ds:")), 0)
            self.assertEqual(storage.get_item("user"), "user-token-abc")

            # Re-read from another instance
            storage2 = SqliteTokenStorage(file_path)
            self.assertEqual(storage2.get_item("user"), "user-token-abc")
            storage.close()
            storage2.close()

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
            token_path = Path(tmpdir) / "test.tokens.db"
            client = create_client(storage=str(token_path))
            storage = client.tokens().storage
            self.assertIsInstance(storage, SqliteTokenStorage)
            assert isinstance(storage, SqliteTokenStorage)
            self.assertEqual(str(storage.path), str(token_path))
            storage.close()

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

    def test_dataset_refresh_token_management(self) -> None:
        client = I6Client(I6SdkConfig())
        tokens = client.tokens()

        self.assertIsNone(tokens.get_dataset_refresh_token("ds-robot"))
        tokens.set_dataset_refresh_token("ds-robot", "secret-robot-key-abc")
        self.assertEqual(tokens.get_dataset_refresh_token("ds-robot"), "secret-robot-key-abc")

        tokens.remove_dataset_refresh_token("ds-robot")
        self.assertIsNone(tokens.get_dataset_refresh_token("ds-robot"))

    def test_machine_robot_dataset_refresh_token_auto_renewal(self) -> None:
        from cryptography.hazmat.primitives import serialization
        from cryptography.hazmat.primitives.asymmetric import ed25519
        import base64
        import json

        priv = ed25519.Ed25519PrivateKey.generate()
        raw_priv = priv.private_bytes(
            serialization.Encoding.Raw, serialization.PrivateFormat.Raw, serialization.NoEncryption()
        )
        raw_pub = priv.public_key().public_bytes(
            serialization.Encoding.Raw, serialization.PublicFormat.Raw
        )
        key_b64 = base64.urlsafe_b64encode(raw_priv + raw_pub).decode().rstrip("=")

        calls = []

        def mock_fetcher(
            url: str, method: str, headers: dict, body: any, timeout: float
        ) -> tuple[int, dict, str]:
            calls.append({"url": url, "method": method, "headers": headers, "body": body})
            if "/api/userz/refresh-dataset" in url:
                # Verify that x-i6-signature header is present
                sig = headers.get("x-i6-signature") or headers.get("X-I6-Signature")
                self.assertIsNotNone(sig)
                self.assertEqual(len(sig), 162)  # 121 raw bytes packed base64url

                return 200, {"content-type": "application/json"}, json.dumps({
                    "access_token": "robot-access-tok-123",
                    "dataset": "dectest",
                    "expires_at": 9999999999999,
                    "grants": 15,
                })
            return 404, {}, ""

        client = I6Client(I6SdkConfig(
            serv_urls={"decsuite": "https://api.i6.test", "userz": "https://auth.i6.test"},
            fetcher=mock_fetcher,
        ))

        # Configure machine robot credentials (no user credentials needed!)
        client.tokens().set_dataset_refresh_token("dectest", key_b64)

        # First call: automatically mints unbounded dataset access token using Ed25519 signature
        tok = client.tokens().get_valid_dataset_access_token("dectest")
        self.assertEqual(tok, "robot-access-tok-123")
        self.assertEqual(len(calls), 1)
        req_body = json.loads(calls[0]["body"]) if isinstance(calls[0]["body"], str) else (json.loads(calls[0]["body"].decode("utf-8")) if calls[0]["body"] else {})
        self.assertTrue(req_body.get("unbound"))

        # Second call: uses cached token without network dispatch
        cached_tok = client.tokens().get_valid_dataset_access_token("dectest")
        self.assertEqual(cached_tok, "robot-access-tok-123")
        self.assertEqual(len(calls), 1)

    def test_granular_dataset_tokens_crud_and_clear(self) -> None:
        client = I6Client(I6SdkConfig())
        tokens = client.tokens()

        # Initially empty
        self.assertIsNone(tokens.get_dataset_token("ds1"))
        self.assertEqual(tokens.get_dataset_tokens_map(), {})

        # Set dataset token
        dt1 = {"dataset": "ds1", "access_token": "tok1", "expires_at": 9999999999999, "grants": 1}
        dt2 = {"dataset": "ds2", "access_token": "tok2", "expires_at": 9999999999999, "grants": 2}
        tokens.set_dataset_token("ds1", dt1)
        tokens.set_dataset_token("ds2", dt2)

        self.assertEqual(tokens.get_dataset_token("ds1"), dt1)
        self.assertEqual(tokens.get_dataset_token("ds2"), dt2)
        m = tokens.get_dataset_tokens_map()
        self.assertEqual(len(m), 2)
        self.assertEqual(m["ds1"]["access_token"], "tok1")
        self.assertEqual(m["ds2"]["access_token"], "tok2")

        # Mutual exclusivity: Setting user token clears robot keys and derived dataset tokens
        tokens.set_dataset_refresh_token("ds1", "rf1")
        self.assertEqual(tokens.get_dataset_refresh_token("ds1"), "rf1")
        tokens.set_user_token({"user_access_token": "u1", "expires_at": 9999999999999})
        self.assertIsNone(tokens.get_dataset_refresh_token("ds1"))
        self.assertIsNotNone(tokens.get_user_token())

        # Mutual exclusivity: Setting robot key clears user token and derived dataset tokens
        tokens.set_dataset_token("ds1", dt1)
        tokens.set_dataset_refresh_token("ds1", "rf1")
        self.assertIsNone(tokens.get_user_token())
        self.assertIsNone(tokens.get_dataset_token("ds1"))
        self.assertEqual(tokens.get_dataset_refresh_token("ds1"), "rf1")

        # Clear tokens
        tokens.clear_tokens()
        self.assertIsNone(tokens.get_user_token())
        self.assertIsNone(tokens.get_dataset_refresh_token("ds1"))
        self.assertEqual(tokens.get_dataset_tokens_map(), {})
        self.assertEqual(tokens.get_dataset_refresh_tokens_map(), {})

    def test_ingest_token_ed25519(self):
        client = I6Client(I6SdkConfig())
        tokens = client.tokens()
        from cryptography.hazmat.primitives.asymmetric import ed25519
        from cryptography.hazmat.primitives import serialization
        import base64

        priv = ed25519.Ed25519PrivateKey.generate()
        raw_bytes = priv.private_bytes(
            serialization.Encoding.Raw,
            serialization.PrivateFormat.Raw,
            serialization.NoEncryption(),
        )
        priv_b64 = base64.urlsafe_b64encode(raw_bytes).decode("utf-8").rstrip("=")

        res = tokens.ingest_token(priv_b64, dataset="my-ds")
        self.assertEqual(res["type"], "dataset_refresh")
        self.assertEqual(res["dataset"], "my-ds")
        self.assertEqual(tokens.get_dataset_refresh_token("my-ds"), priv_b64)

    def test_ingest_token_jwt_offline(self):
        client = I6Client(I6SdkConfig())
        tokens = client.tokens()
        import base64

        header = base64.urlsafe_b64encode(b'{"alg":"none"}').decode("utf-8").rstrip("=")
        payload = base64.urlsafe_b64encode(b'{"dataset":"fashion-ds","exp":1895000000}').decode("utf-8").rstrip("=")
        jwt_ds = f"{header}.{payload}.sig"

        res = tokens.ingest_token(jwt_ds)
        self.assertEqual(res["type"], "dataset_bearer")
        self.assertEqual(res["dataset"], "fashion-ds")
        self.assertIsNotNone(tokens.get_dataset_token("fashion-ds"))

        # User access JWT
        u_payload = base64.urlsafe_b64encode(b'{"user_id":"usr_777","exp":1895000000}').decode("utf-8").rstrip("=")
        jwt_u = f"{header}.{u_payload}.sig"
        res_u = tokens.ingest_token(jwt_u, offline=True)
        self.assertEqual(res_u["type"], "user_access")
        self.assertEqual(res_u["user_id"], "usr_777")
        self.assertIsNotNone(tokens.get_user_token())


if __name__ == "__main__":
    unittest.main()
