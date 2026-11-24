"""Unit tests for token storage and token manager."""

import tempfile
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
            self.assertIsInstance(client.tokens().storage, FileTokenStorage)
            self.assertEqual(str(client.tokens().storage.path), str(token_path))


if __name__ == "__main__":
    unittest.main()
