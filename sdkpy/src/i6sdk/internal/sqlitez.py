"""SQLite database helper wrapping Python's built-in sqlite3."""

from __future__ import annotations

import json
import pathlib
import sqlite3
from typing import Any, Dict, List, Optional, Union


class SqliteConn:
    """An open SQLite connection with query and named-query helpers."""

    def __init__(self, db: sqlite3.Connection) -> None:
        self.db: Optional[sqlite3.Connection] = db
        self.db.row_factory = sqlite3.Row

    def __enter__(self) -> SqliteConn:
        return self

    def __exit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        self.close()

    def close(self) -> None:
        """Closes the underlying database connection."""
        if self.db is not None:
            self.db.close()
            self.db = None

    def query(self, query: str, params: Optional[List[Any]] = None) -> List[Dict[str, Any]]:
        """Executes a SQL query and returns rows as dictionaries."""
        if self.db is None:
            raise RuntimeError("sqlitez: connection is closed")
        cur = self.db.cursor()
        try:
            cur.execute(query, params or [])
            rows = cur.fetchall()
            return [dict(row) for row in rows]
        finally:
            cur.close()

    def named_query(
        self,
        name: str,
        params: Optional[List[Any]] = None,
        ignore_if_not_exist: bool = False,
    ) -> List[Dict[str, Any]]:
        """Executes a named query defined in the namedqueries table."""
        try:
            found = self.query(
                "SELECT query FROM namedqueries WHERE name = ?",
                [name],
            )
        except Exception:
            if ignore_if_not_exist:
                return []
            raise

        if not found:
            if ignore_if_not_exist:
                return []
            raise LookupError(f"sqlitez: no named query {name!r}")

        sql_text = str(found[0]["query"])
        return self.query(sql_text, params)

    def named_string_query(self, name: str, params: Optional[List[Any]] = None) -> Optional[str]:
        """Binds JSON-stringified params to positional placeholders and returns column 'ret'."""
        bound_params = [json.dumps(p) for p in (params or [])]
        rows = self.named_query(name=name, params=bound_params)
        if not rows:
            return None
        row = rows[0]
        if "ret" not in row:
            raise LookupError(f"sqlitez: named query {name!r} returned no 'ret' column")
        ret_val = row["ret"]
        if ret_val is None:
            return None
        return str(ret_val)

    def named_json_query(self, name: str, params: Optional[List[Any]] = None) -> Any:
        """Executes named_string_query and parses the result as JSON."""
        raw = self.named_string_query(name=name, params=params)
        if raw is None:
            return None
        return json.loads(raw)


class Sqlitez:
    """Factory for opening SQLite database connections."""

    @staticmethod
    def open(path: Union[pathlib.Path, str]) -> SqliteConn:
        """Opens a connection to the SQLite database at path."""
        p = pathlib.Path(path)
        if not p.is_file():
            raise FileNotFoundError(f"sqlitez: {path} is not a file")
        # Open in read-only mode for safety
        uri_path = f"file:{p.resolve().as_posix()}?mode=ro"
        conn = sqlite3.connect(uri_path, uri=True, check_same_thread=False)
        return SqliteConn(conn)


sqlitez = Sqlitez()
