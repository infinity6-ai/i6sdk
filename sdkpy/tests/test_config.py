"""Unit tests for config module."""

import os
import unittest
from unittest import mock

import i6sdk


class TestConfig(unittest.TestCase):
    def test_default_config(self) -> None:
        cfg = i6sdk.DEFAULT_CONFIG
        self.assertEqual(cfg.base_url, "https://decsuite.sandbox.rs.infinity6.ai")
        self.assertEqual(cfg.app, "default")
        self.assertEqual(cfg.timeout, 30.0)
        self.assertIsNone(cfg.token)

    def test_load_from_env(self) -> None:
        env_vars = {
            "I6_BASE_URL": "https://api.test.i6.ai",
            "I6_API_KEY": "test-key-123",
            "I6_APP": "my-app",
            "I6_TIMEOUT": "45.5",
        }
        with mock.patch.dict(os.environ, env_vars, clear=False):
            cfg = i6sdk.load_config_from_env()
            self.assertEqual(cfg.base_url, "https://api.test.i6.ai")
            self.assertEqual(cfg.token, "test-key-123")
            self.assertEqual(cfg.app, "my-app")
            self.assertEqual(cfg.timeout, 45.5)

    def test_create_client_factory(self) -> None:
        client = i6sdk.create_client(api_key="custom-key", base_url="https://custom.i6.ai")
        self.assertEqual(client.config().token, "custom-key")
        self.assertEqual(client.config().base_url, "https://custom.i6.ai")


if __name__ == "__main__":
    unittest.main()
