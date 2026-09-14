((QUnit) => {

    QUnit.module("i6sdk-web", function () {

        QUnit.test("hello", async function (assert) {
            const { I6Sdk } = await import('../src/index.js')
            const i6sdk = new I6Sdk()
            await i6sdk.config({
                subscription: 'sub1',
                store: 'store1'
            })
            assert.equal(i6sdk.hello(), 'Hello from i6sdk-web!')
        });

        QUnit.test("config is required", async function (assert) {
            const { I6Sdk } = await import('../src/index.js')
            const i6sdk = new I6Sdk()
            let threw = false
            try {
                await i6sdk.config()
            } catch (e) {
                threw = true
            }
            assert.true(threw, 'config() should throw when no config is provided')
        });

    });

})(QUnit)
