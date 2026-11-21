(function($) {

    class FFNav {

         click(element) {
            element = $(element)
            if (!element.length) {
                throw 'no element to click';
            }
            element = element[0];
            var document = element.ownerDocument;
            if (document.dispatchEvent) { // W3C
                var oEvent = document.createEvent("MouseEvents");
                oEvent.initMouseEvent("click", true, true, window, 1, 1, 1, 1, 1, false, false, false, false, 0, element);
                element.dispatchEvent(oEvent);
            } else if (document.fireEvent) { // IE
                element.click();
            }
        }

        fixture() {
            let ret = $('#qunit-fixture iframe')
            if (ret.length === 0) {
                $('#qunit-fixture').html('<iframe src="about:blank"></iframe>')
                ret = $('#qunit-fixture iframe')
            }
            return ret
        }

        win() {
            return this.fixture()[0].contentWindow;
        }

        doc() {
            return this.win().document;
        }

        $(...args) {
            return this.win().jQuery(...args)
        }

        async get(selector, timeout = 5000, interval = 1) {
            const start = Date.now();

            return new Promise((resolve, reject) => {
                const check = () => {
                    const el = this.$(selector);
                    if (el.length) {
                        resolve(el);
                    } else if (Date.now() - start >= timeout) {
                        reject(new Error(`Element "${selector}" not found within ${timeout}ms`));
                    } else {
                        setTimeout(check, interval);
                    }
                };
                check();
            });
        }

        async open(url) {
            let fixture = this.fixture()

            return new Promise((resolve, reject) => {
                fixture.off('load.ffnav')
                fixture.on('load.ffnav', () => {
                    try {
                        let win = fixture[0].contentWindow
                        let $doc = $(win.document)
                        $doc.ready(() => {
                            if (!win.jQuery) {
                                reject('jquery is required')
                                return
                            }
                            resolve(win)
                        })
                    } catch (err) {
                        reject(err)
                    }
                })

                fixture.attr('src', url)
            })
        }

        close() {
            $('#qunit-fixture').html('')
        }

    }

    window.FFNav = FFNav

})(jQuery)
