((QUnit) => {

    QUnit.module("gzip", function () {

        QUnit.test("gzip and gunzipString for raw strings", async function (assert) {
            const { gzip, gunzipString, gunzip, gzipStringToBase64Url, gunzipStringFromBase64Url } = await import('../src/internal/gzip/gzip.js');

            const text = "Hello World! This is an arbitrary text with UTF-8 characters: 🚀🌟🔥";

            const compressed = await gzip(text);
            assert.true(compressed instanceof Uint8Array, 'gzip returns Uint8Array');

            const decompressedRaw = await gunzip(compressed);
            assert.true(decompressedRaw instanceof Uint8Array, 'gunzip returns Uint8Array');

            const decompressed = await gunzipString(compressed);
            assert.strictEqual(decompressed, text, 'gunzipString decompresses back to original text');

            const b64 = await gzipStringToBase64Url(text);
            assert.strictEqual(typeof b64, 'string', 'gzipStringToBase64Url returns string');

            const decompressedFromB64 = await gunzipStringFromBase64Url(b64);
            assert.strictEqual(decompressedFromB64, text, 'gunzipStringFromBase64Url restores original string');
        });

        QUnit.test("gzipJson and gunzipJson for objects", async function (assert) {
            const { gzipJson, gunzipJson, gzipJsonToBase64Url, gunzipJsonFromBase64Url } = await import('../src/internal/gzip/gzip.js');

            const originalData = { hello: "world", a: 1, nested: { b: [2, 3, "foo"] } };

            const compressedBytes = await gzipJson(originalData);
            assert.true(compressedBytes instanceof Uint8Array, 'returns Uint8Array');

            const decompressedData = await gunzipJson(compressedBytes);
            assert.deepEqual(decompressedData, originalData, 'decompresses to original data');

            const b64 = await gzipJsonToBase64Url(originalData);
            assert.strictEqual(typeof b64, 'string', 'b64 string returned');

            const decompressedFromB64 = await gunzipJsonFromBase64Url(b64);
            assert.deepEqual(decompressedFromB64, originalData, 'decompresses from b64 to original data');
        });

    });

})(QUnit);
