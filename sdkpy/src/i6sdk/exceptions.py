"""Exceptions for i6 Python SDK."""

from __future__ import annotations

from typing import Any, Optional


class I6Error(Exception):
    """Base exception for all Infinity6 SDK errors."""

    def __init__(self, message: str = "") -> None:
        super().__init__(message)
        self.message = message

    def __str__(self) -> str:
        return self.message or self.__class__.__name__


class I6ApiError(I6Error):
    """Exception raised when an API returns an HTTP error status (>= 400)."""

    def __init__(
        self,
        message: str = "",
        status: int = 0,
        code: str = "bad_request",
        cause: Any = None,
    ) -> None:
        super().__init__(message)
        self.status = status
        self.code = code
        self.cause = cause

    def __str__(self) -> str:
        return f"{self.message} (status={self.status}, code={self.code})"


class I6AuthError(I6Error):
    """Exception raised when authentication fails or credentials are invalid."""

    pass


class I6ConnectionError(I6Error):
    """Exception raised when network transport fails or times out."""

    pass


class I6ValidationError(I6Error):
    """Exception raised when arguments or parameters fail client-side validation."""

    pass
