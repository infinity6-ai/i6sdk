(function (QUnit) {

    QUnit.config.autostart = false

    QUnit.on('runEnd', async runEnd => {
        await fetch('../api/i6web/webchecker/results', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json; charset=utf-8'
            },
            body: JSON.stringify(runEnd)
        });
    });

})(QUnit)