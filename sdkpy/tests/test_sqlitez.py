"""Unit tests for internal sqlitez."""

import sqlite3
import tempfile
import unittest
from pathlib import Path

from i6sdk.internal.sqlitez import Sqlitez


class TestSqlitez(unittest.TestCase):
    def test_sqlitez_queries(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            db_path = Path(tmpdir) / "test.db"
            # Populate sqlite
            conn_init = sqlite3.connect(str(db_path))
            cur = conn_init.cursor()
            cur.execute("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
            cur.execute("INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob')")
            cur.execute("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT)")
            cur.execute(
                "INSERT INTO namedqueries (name, query) VALUES ('get_user', 'SELECT name FROM users WHERE id = ?1')"
            )
            cur.execute(
                "INSERT INTO namedqueries (name, query) VALUES ('model-test-output', 'SELECT json_object(\"result\", \"ok\", \"sku\", json(?1)) as ret')"
            )
            conn_init.commit()
            conn_init.close()

            with Sqlitez.open(db_path) as conn:
                rows = conn.query("SELECT * FROM users ORDER BY id")
                self.assertEqual(len(rows), 2)
                self.assertEqual(rows[0]["name"], "Alice")

                # Test named_json_query
                json_res = conn.named_json_query("model-test-output", params=["SKU123"])
                self.assertEqual(json_res, {"result": "ok", "sku": "SKU123"})


if __name__ == "__main__":
    unittest.main()
