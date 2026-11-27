"""Unit tests for internal bitmapz."""

import unittest

from i6sdkpy.internal.bitmapz import contains_all, contains_any, join, set_bits, split


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

    def test_serialize_deserialize_vectors(self) -> None:
        from i6sdkpy.internal.bitmapz import deserialize, serialize

        cases = [
            (0, bytes([0x00])),
            (1, bytes([0x01])),
            (15, bytes([0x0F])),
            (16, bytes([0x10, 0x10])),
            (255, bytes([0x10, 0xFF])),
            (256, bytes([0x20, 0x00, 0x01])),
            (0x123456, bytes([0x30, 0x56, 0x34, 0x12])),
            (0xFFFFFFFFFFFFFFFF, bytes([0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF])),
        ]

        for val, want in cases:
            encoded = serialize(val)
            self.assertEqual(encoded, want, f"failed to serialize {val}")
            decoded = deserialize(encoded)
            self.assertEqual(decoded, val, f"failed to deserialize {val}")

    def test_grants_enum_and_exports(self) -> None:
        from i6sdkpy import Grants, bitmapz

        self.assertEqual(Grants.SYSTEM, 1)
        self.assertEqual(Grants.ADMIN, 2)
        self.assertEqual(Grants.WRITE, 4)
        self.assertEqual(Grants.READ, 8)
        self.assertEqual(Grants.ALL, 15)

        # Ensure TitleCase aliases are removed
        for alias in ("System", "Admin", "Write", "Read", "All"):
            self.assertFalse(hasattr(Grants, alias))

        self.assertTrue(bitmapz.contains_all(Grants.ALL, Grants.READ, Grants.WRITE))
        self.assertEqual(bitmapz.join(Grants.READ, Grants.WRITE), 12)

    def test_deserialize_validation_errors(self) -> None:
        from i6sdkpy.internal.bitmapz import deserialize, serialize

        with self.assertRaisesRegex(ValueError, "data is empty"):
            deserialize(b"")

        with self.assertRaisesRegex(ValueError, "cannot exceed 8 bytes"):
            deserialize(bytes([0x90, 0x01]))

        with self.assertRaisesRegex(ValueError, "expected 1 byte for 4-bit value"):
            deserialize(bytes([0x01, 0x02]))

        with self.assertRaisesRegex(ValueError, "lower 4 bits must be 0"):
            deserialize(bytes([0x11, 0x02]))

        with self.assertRaisesRegex(ValueError, "invalid data length"):
            deserialize(bytes([0x20, 0x01]))

        with self.assertRaisesRegex(ValueError, "invalid data length"):
            deserialize(bytes([0x10, 0x01, 0x02]))

    def test_serialize_validation_errors(self) -> None:
        from i6sdkpy.internal.bitmapz import serialize

        with self.assertRaisesRegex(ValueError, "must be non-negative"):
            serialize(-1)

        with self.assertRaisesRegex(ValueError, "cannot exceed 8 bytes"):
            serialize(1 << 65)

    def test_set_alias_removed(self) -> None:
        from i6sdkpy import bitmapz

        # Ensure built-in set is not shadowed and set_bits is used directly
        self.assertFalse(hasattr(bitmapz, "set"))
        self.assertEqual(bitmapz.set_bits(0, True, 4, 8), 12)
        self.assertEqual(bitmapz.set_bits(12, False, 4), 8)


if __name__ == "__main__":
    unittest.main()
