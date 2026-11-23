"""Unit tests for internal bitmapz."""

import unittest

from i6sdk.internal.bitmapz import contains_all, contains_any, join, set_bits, split


class TestBitmapz(unittest.TestCase):
    def test_join(self) -> None:
        self.assertEqual(join(1, 2, 4), 7)
        self.assertEqual(join(4, 8), 12)

    def test_contains_all(self) -> None:
        self.assertTrue(contains_all(15, 1, 2, 4, 8))
        self.assertTrue(contains_all(7, 1, 2))
        self.assertFalse(contains_all(7, 8))

    def test_contains_any(self) -> None:
        self.assertTrue(contains_any(7, 4, 8))
        self.assertFalse(contains_any(7, 8, 16))

    def test_split(self) -> None:
        self.assertEqual(split(0), [])
        self.assertEqual(split(7), [1, 2, 4])
        self.assertEqual(split(12), [4, 8])

    def test_set_bits(self) -> None:
        self.assertEqual(set_bits(7, False, 2), 5)
        self.assertEqual(set_bits(5, True, 2), 7)


if __name__ == "__main__":
    unittest.main()
