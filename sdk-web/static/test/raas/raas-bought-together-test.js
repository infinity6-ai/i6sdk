((QUnit) => {

    
    QUnit.module("raas-bought-together", function () {

        
        QUnit.test("mock", async function (assert) {
            const { I6Sdk } = await import('../../sdkweb/sdkweb.js')
            const i6sdk = new I6Sdk()
            await i6sdk.init({
                'subscription': 'sub1',
                'store': 'store1'
            })
            const rbt = await i6sdk.RaasBoughtTogether()
            const result = await rbt.predict([
                {id: 'a'}, 
                {id: 'b'}, 
                {id: 'c'}
            ], {maxResults: 2, deterministic: true})
            assert.equal(result[0].products[0].id, 'c')
            assert.true(result[0].score >= 0, `wrong score ${result[0].score}`)
            assert.equal(result[1].products[0].id, 'b')
            assert.true(result[1].score >= 0, `wrong score ${result[0].score}`)
            assert.equal(result.length, 2)
            
            
        });
        
    });

})(QUnit)