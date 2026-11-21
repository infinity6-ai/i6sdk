((QUnit) => {

    QUnit.module("b64", function () {

        QUnit.test("encodeUrlSafe and decodeUrlSafe", async function (assert) {
            const { encodeUrlSafe, decodeUrlSafe, decodeUrlSafeToBytes } = await import('../src/internal/b64/b64.js');
            
            const originalString = "Hello, world! 🌍 \x00 \xFF";
            
            const encoded = encodeUrlSafe(originalString);
            assert.strictEqual(typeof encoded, 'string', 'encoded is string');
            // Ensure URL safe characters (no + or / or =)
            assert.true(!encoded.includes('+'), 'no plus');
            assert.true(!encoded.includes('/'), 'no slash');
            assert.true(!encoded.includes('='), 'no equals');

            const decodedString = decodeUrlSafe(encoded);
            assert.strictEqual(decodedString, originalString, 'decodes to original string');
        });

    });

})(QUnit);
