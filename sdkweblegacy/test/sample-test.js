((QUnit) => {

    QUnit.module("i6sdk-web", function () {

        QUnit.test("hello", async function (assert) {
            // const { I6Sdk } = await import('../src/index.js')
            assert.true(true)
        });

        QUnit.test("setConfig overwrites the given keys at any time", async function (assert) {
            const { I6Sdk } = await import('../src/index.js')
            const sdk = new I6Sdk({ baseUrl: 'http://a' })
            const servUrls = sdk.config().servUrls
            assert.strictEqual(sdk.config().baseUrl, 'http://a')
            assert.ok(servUrls, 'defaults are merged')
            sdk.setConfig({ baseUrl: 'http://x' })
            assert.strictEqual(sdk.config().baseUrl, 'http://x', 'given key is overwritten')
            assert.deepEqual(sdk.config().servUrls, servUrls, 'other keys are kept')
            assert.strictEqual(sdk.auth()._services.config().baseUrl, 'http://x', 'clients see the new config')
        });

    });

})(QUnit)
