"""Unit tests for config module."""

import os
import unittest
from unittest import mock

import i6sdk


class TestConfig(unittest.TestCase):
    def test_default_config(self) -> None:
        cfg = i6sdk.DEFAULT_CONFIG
        self.assertEqual(cfg.serv_urls.get("decsuite"), "https://decsuite.sandbox.rs.infinity6.ai")
        self.assertEqual(cfg.app, "default")
        self.assertEqual(cfg.timeout, 30.0)

    def test_load_from_env(self) -> None:
        env_vars = {
            "I6_SDK_DECSUITE_URL": "https://decsuite.test.i6.ai",
            "I6_SDK_TFINFER_URL": "https://tfinfer.test.i6.ai",
            "I6_SDK_USERZ_URL": "https://userz.test.i6.ai",
            "I6_SDK_CACHE_DIR": "/tmp/custom-cache",
        }
        with mock.patch.dict(os.environ, env_vars, clear=False):
            cfg = i6sdk.load_config_from_env()
            self.assertEqual(cfg.serv_urls.get("decsuite"), "https://decsuite.test.i6.ai")
            self.assertEqual(cfg.serv_urls.get("tfinfer"), "https://tfinfer.test.i6.ai")
            self.assertEqual(cfg.serv_urls.get("userz"), "https://userz.test.i6.ai")
            self.assertEqual(cfg.cache_dir, "/tmp/custom-cache")
            self.assertEqual(cfg.app, "default")
            self.assertEqual(cfg.timeout, 30.0)

    def test_create_client_factory(self) -> None:
        client = i6sdk.create_client(
            serv_urls={"decsuite": "https://custom.i6.ai"},
            app="my-app",
            timeout=45.5,
            cache_dir="/tmp/custom-cache",
        )
        self.assertEqual(client.config().app, "my-app")
        self.assertEqual(client.config().timeout, 45.5)
        self.assertEqual(client.config().cache_dir, "/tmp/custom-cache")
        self.assertEqual(str(client.fs().base_dir), "/tmp/custom-cache/sdkpy/fs")
        self.assertEqual(client.config().serv_urls.get("decsuite"), "https://custom.i6.ai")

    def test_set_config(self) -> None:
        client = i6sdk.create_client()
        client.set_config({"app": "updated-app"})
        self.assertEqual(client.config().app, "updated-app")
        self.assertFalse(hasattr(client, "setConfig"))


if __name__ == "__main__":
    unittest.main()
