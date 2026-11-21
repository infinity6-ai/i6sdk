((QUnit) => {

    QUnit.module("bitmapz", function (hooks) {

        let bitmapz
        let Grants

        hooks.beforeEach(async function () {
            ({ bitmapz, Grants } = await import('../main/index.js'))
        })

        QUnit.test("join combines bits using bitwise OR", function (assert) {
            assert.strictEqual(bitmapz.join(), 0)
            assert.strictEqual(bitmapz.join(1), 1)
            assert.strictEqual(bitmapz.join(1, 2, 4), 7)
            assert.strictEqual(bitmapz.join(Grants.Read, Grants.Write), 12)
            assert.strictEqual(bitmapz.join(Grants.System, Grants.Admin, Grants.Write, Grants.Read), Grants.All)
        })

        QUnit.test("containsAll checks if all required bits are present", function (assert) {
            assert.true(bitmapz.containsAll(15, Grants.Read, Grants.Write))
            assert.true(bitmapz.containsAll(12, Grants.Read, Grants.Write))
            assert.false(bitmapz.containsAll(12, Grants.Read, Grants.Admin))
            assert.true(bitmapz.containsAll(0), "empty expected returns true")
        })

        QUnit.test("containsAny checks if at least one bit is present", function (assert) {
            assert.true(bitmapz.containsAny(12, Grants.Read, Grants.Admin))
            assert.true(bitmapz.containsAny(8, Grants.Read))
            assert.false(bitmapz.containsAny(12, Grants.System, Grants.Admin))
            assert.false(bitmapz.containsAny(0, Grants.Read), "0 current returns false")
            assert.false(bitmapz.containsAny(15), "empty expected returns false")
        })

        QUnit.test("split decomposes bitmask into ascending powers of two", function (assert) {
            assert.deepEqual(bitmapz.split(0), [])
            assert.deepEqual(bitmapz.split(1), [1])
            assert.deepEqual(bitmapz.split(5), [1, 4])
            assert.deepEqual(bitmapz.split(12), [4, 8])
            assert.deepEqual(bitmapz.split(15), [1, 2, 4, 8])
            assert.deepEqual(bitmapz.split(Grants.All), [1, 2, 4, 8])
        })

        QUnit.test("set adds or removes bits", function (assert) {
            assert.strictEqual(bitmapz.set(0, true, Grants.Read, Grants.Write), 12)
            assert.strictEqual(bitmapz.set(12, false, Grants.Write), 8)
            assert.strictEqual(bitmapz.set(8, true, Grants.Admin), 10)
        })

        QUnit.test("serialize and deserialize roundtrip test vectors matching Go bitmapz", function (assert) {
            const cases = [
                { value: 0, want: new Uint8Array([0x00]) },
                { value: 1, want: new Uint8Array([0x01]) },
                { value: 15, want: new Uint8Array([0x0F]) },
                { value: 16, want: new Uint8Array([0x10, 0x10]) },
                { value: 255, want: new Uint8Array([0x10, 0xFF]) },
                { value: 256, want: new Uint8Array([0x20, 0x00, 0x01]) },
                { value: 0x123456, want: new Uint8Array([0x30, 0x56, 0x34, 0x12]) },
            ]

            for (const c of cases) {
                const encoded = bitmapz.serialize(c.value)
                assert.deepEqual(Array.from(encoded), Array.from(c.want), `serialized ${c.value} matches expected bytes`)
                const decoded = bitmapz.deserialize(encoded)
                assert.strictEqual(decoded, c.value, `deserialized ${c.value} matches original`)
            }

            // 64-bit BigInt test
            const max64 = 0xFFFFFFFFFFFFFFFFn
            const encoded64 = bitmapz.serialize(max64)
            assert.deepEqual(Array.from(encoded64), [0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF])
            const decoded64 = bitmapz.deserialize(encoded64)
            assert.strictEqual(BigInt(decoded64), max64)
        })

        QUnit.test("deserialize validates corrupted or invalid headers", function (assert) {
            assert.throws(() => bitmapz.deserialize(new Uint8Array([])), /data is empty/)
            assert.throws(() => bitmapz.deserialize(new Uint8Array([0x90])), /invalid size/)
            assert.throws(() => bitmapz.deserialize(new Uint8Array([0x05, 0x01])), /invalid data length/)
            assert.throws(() => bitmapz.deserialize(new Uint8Array([0x15, 0x01])), /lower 4 bits must be 0/)
            assert.throws(() => bitmapz.deserialize(new Uint8Array([0x20, 0x01])), /invalid data length/)
        })
    })

})(QUnit)
