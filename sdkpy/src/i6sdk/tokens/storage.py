"""Token storage interface and SQLite/Memory implementations."""

from __future__ import annotations

import abc
import os
import pathlib
import sqlite3
import threading
import time
from typing import Dict, Optional, Union


class TokenStorage(abc.ABC):
    """Abstract interface for token persistence."""

    @abc.abstractmethod
    def get_item(self, key: str) -> Optional[str]:
        """Retrieve stored string value for key."""
        pass

    @abc.abstractmethod
    def set_item(self, key: str, value: str) -> None:
        """Store string value for key."""
        pass

    @abc.abstractmethod
    def remove_item(self, key: str) -> None:
        """Remove value for key."""
        pass

    @abc.abstractmethod
    def delete_prefix(self, prefix: str) -> None:
        """Remove all keys starting with prefix."""
        pass

    @abc.abstractmethod
    def list_prefix(self, prefix: str) -> Dict[str, str]:
        """List all key-value pairs where key starts with prefix."""
        pass

    def close(self) -> None:
        """Close storage resources if applicable."""
        pass


class MemoryTokenStorage(TokenStorage):
    """Thread-safe in-memory token storage."""

    def __init__(self) -> None:
        self._data: Dict[str, str] = {}
        self._lock = threading.Lock()

    def get_item(self, key: str) -> Optional[str]:
        with self._lock:
            return self._data.get(key)

    def set_item(self, key: str, value: str) -> None:
        with self._lock:
            self._data[key] = value

    def remove_item(self, key: str) -> None:
        with self._lock:
            self._data.pop(key, None)

    def delete_prefix(self, prefix: str) -> None:
        with self._lock:
            keys_to_del = [k for k in self._data if k.startswith(prefix)]
            for k in keys_to_del:
                del self._data[k]

    def list_prefix(self, prefix: str) -> Dict[str, str]:
        with self._lock:
            return {k: v for k, v in self._data.items() if k.startswith(prefix)}

    def close(self) -> None:
        pass


class SqliteTokenStorage(TokenStorage):
    """SQLite-backed token storage with WAL mode and secure file permissions (0600)."""

    def __init__(self, path: Optional[Union[pathlib.Path, str]] = None) -> None:
        if path is None:
            config_dir = pathlib.Path.home() / ".config" / "i6"
            config_dir.mkdir(parents=True, exist_ok=True)
            self.path = config_dir / "tokens.db"
        else:
            self.path = pathlib.Path(path)
            if str(self.path) != ":memory:":
                self.path.parent.mkdir(parents=True, exist_ok=True)

        self._is_memory = str(self.path) == ":memory:"
        self._lock = threading.Lock()

        self._conn = sqlite3.connect(
            str(self.path),
            check_same_thread=False,
            timeout=5.0,
            isolation_level=None,  # Autocommit mode
        )
        self._init_db()

    def _init_db(self) -> None:
        with self._lock:
            cur = self._conn.cursor()
            try:
                cur.execute("PRAGMA journal_mode=WAL;")
                cur.execute("PRAGMA busy_timeout=5000;")
                cur.execute("PRAGMA synchronous=NORMAL;")
                cur.execute(
                    """
                    CREATE TABLE IF NOT EXISTS tokens (
                        key TEXT PRIMARY KEY,
                        value TEXT NOT NULL,
                        updated_at INTEGER NOT NULL
                    );
                    """
                )
            finally:
                cur.close()

        if not self._is_memory and self.path.exists():
            try:
                os.chmod(self.path, 0o600)
            except OSError:
                pass

    def get_item(self, key: str) -> Optional[str]:
        with self._lock:
            cur = self._conn.cursor()
            try:
                cur.execute("SELECT value FROM tokens WHERE key = ?", (key,))
                row = cur.fetchone()
                return row[0] if row else None
            finally:
                cur.close()

    def set_item(self, key: str, value: str) -> None:
        now_ms = int(time.time() * 1000)
        with self._lock:
            cur = self._conn.cursor()
            try:
                cur.execute(
                    """
                    INSERT INTO tokens (key, value, updated_at) VALUES (?, ?, ?)
                    ON CONFLICT(key) DO UPDATE SET value=excluded.value, updated_at=excluded.updated_at
                    """,
                    (key, value, now_ms),
                )
            finally:
                cur.close()

    def remove_item(self, key: str) -> None:
        with self._lock:
            cur = self._conn.cursor()
            try:
                cur.execute("DELETE FROM tokens WHERE key = ?", (key,))
            finally:
                cur.close()

    def delete_prefix(self, prefix: str) -> None:
        with self._lock:
            cur = self._conn.cursor()
            try:
                cur.execute(
                    "DELETE FROM tokens WHERE substr(key, 1, length(?)) = ?",
                    (prefix, prefix),
                )
            finally:
                cur.close()

    def list_prefix(self, prefix: str) -> Dict[str, str]:
        with self._lock:
            cur = self._conn.cursor()
            try:
                cur.execute(
                    "SELECT key, value FROM tokens WHERE substr(key, 1, length(?)) = ?",
                    (prefix, prefix),
                )
                return {row[0]: row[1] for row in cur.fetchall()}
            finally:
                cur.close()

    def close(self) -> None:
        with self._lock:
            self._conn.close()


# Backward compatibility alias
FileTokenStorage = SqliteTokenStorage
