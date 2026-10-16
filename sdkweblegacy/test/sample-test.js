((QUnit) => {

    QUnit.module("i6sdk-web", function () {

        QUnit.test("hello", async function (assert) {
            // const { I6Sdk } = await import('../src/index.js')
            assert.true(true)
        });

        QUnit.test("setConfig overwrites the given keys at any time", async function (assert) {
            const { I6Sdk } = await import('../src/index.js')
            const sdk = new I6Sdk({ sessionId: 's1', userId: 'u1' })
            const baseUrl = sdk.config().baseUrl
            sdk.setConfig({ userId: 'u2' })
            assert.deepEqual(sdk.config(), { baseUrl, sessionId: 's1', userId: 'u2' }, 'other keys are kept')
            sdk.setConfig({ baseUrl: 'http://x' })
            assert.strictEqual(sdk.auth()._services.config().baseUrl, 'http://x', 'clients see the new config')
        });

    });

})(QUnit)
