((QUnit) => {

    QUnit.module("i6sdk-web", function () {

        QUnit.test("hello", async function (assert) {
            // const { I6Sdk } = await import('../main/index.js')
            assert.true(true)
        });

        QUnit.test("setConfig overwrites the given keys at any time", async function (assert) {
            const { I6Sdk } = await import('../main/index.js')
            const sdk = new I6Sdk({ servUrls: { decsuite: 'http://a' } })
            const servUrls = sdk.config().servUrls
            assert.strictEqual(sdk.config().servUrls?.["decsuite"], 'http://a')
            assert.ok(servUrls, 'defaults are merged')
            sdk.setConfig({ servUrls: { ...sdk.config().servUrls, decsuite: 'http://x' } })
            assert.strictEqual(sdk.config().servUrls?.["decsuite"], 'http://x', 'given key is overwritten')
            assert.strictEqual(sdk.auth()._services.config().servUrls?.["decsuite"], 'http://x', 'clients see the new config')
        });

        QUnit.test("warmupDataset validates its options", async function (assert) {
            const { I6Sdk } = await import('../main/index.js')
            const sdk = new I6Sdk()
            await assert.rejects(sdk.warmupDataset(), /opts.dataset/)
            await assert.rejects(sdk.warmupDataset({ dataset: 'd1' }), /opts.models/)
            await assert.rejects(sdk.warmupDataset({ dataset: 'd1', models: 'not-array' }), /opts.models/)
        });

    });

})(QUnit)
