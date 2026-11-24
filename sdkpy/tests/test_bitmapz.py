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

    def test_serialize_deserialize_vectors(self) -> None:
        from i6sdk.internal.bitmapz import deserialize, serialize

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
        from i6sdk import Grants, bitmapz

        self.assertEqual(Grants.System, 1)
        self.assertEqual(Grants.Admin, 2)
        self.assertEqual(Grants.Write, 4)
        self.assertEqual(Grants.Read, 8)
        self.assertEqual(Grants.All, 15)

        self.assertEqual(Grants.SYSTEM, 1)
        self.assertEqual(Grants.ADMIN, 2)
        self.assertEqual(Grants.WRITE, 4)
        self.assertEqual(Grants.READ, 8)
        self.assertEqual(Grants.ALL, 15)

        self.assertTrue(bitmapz.contains_all(Grants.All, Grants.Read, Grants.Write))
        self.assertEqual(bitmapz.join(Grants.Read, Grants.Write), 12)



if __name__ == "__main__":
    unittest.main()
