((QUnit) => {

    QUnit.module("gzip", function () {

        QUnit.test("gzipJson and gunzipJson", async function (assert) {
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
