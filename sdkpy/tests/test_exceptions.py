"""Unit tests for exceptions hierarchy."""

import unittest

from i6sdk.exceptions import (
    I6ApiError,
    I6AuthError,
    I6ConnectionError,
    I6Error,
    I6ValidationError,
)


class TestExceptions(unittest.TestCase):
    def test_i6_error_hierarchy(self) -> None:
        self.assertTrue(issubclass(I6ApiError, I6Error))
        self.assertTrue(issubclass(I6AuthError, I6Error))
        self.assertTrue(issubclass(I6ConnectionError, I6Error))
        self.assertTrue(issubclass(I6ValidationError, I6Error))

    def test_i6_api_error_attributes(self) -> None:
        err = I6ApiError(message="Forbidden", status=403, code="forbidden", cause="invalid_grant")
        self.assertEqual(err.status, 403)
        self.assertEqual(err.code, "forbidden")
        self.assertEqual(err.cause, "invalid_grant")
        self.assertIn("403", str(err))
        self.assertIn("forbidden", str(err))


if __name__ == "__main__":
    unittest.main()
