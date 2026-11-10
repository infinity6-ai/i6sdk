((QUnit) => {

    QUnit.module("i6sdk-web authv2 & tokens", function (hooks) {

        let I6Sdk
        let realFetch
        let mockStorage

        hooks.beforeEach(async function () {
            ({ I6Sdk } = await import('../src/index.js'))
            realFetch = globalThis.fetch

            // Mock localStorage if needed
            const store = {}
            mockStorage = {
                getItem: (k) => store[k] || null,
                setItem: (k, v) => { store[k] = String(v) },
                removeItem: (k) => { delete store[k] },
                clear: () => { for (const k in store) delete store[k] }
            }
            if (typeof globalThis.localStorage === "undefined") {
                globalThis.localStorage = mockStorage
            } else {
                localStorage.clear()
            }
        })

        hooks.afterEach(function () {
            globalThis.fetch = realFetch
            if (typeof localStorage !== "undefined") {
                localStorage.clear()
            }
        })

        QUnit.test("tokens() is exposed on sdk and manages token storage", async function (assert) {
            const sdk = new I6Sdk({ baseUrl: 'http://test.api' })
            const tokens = sdk.tokens()
            assert.ok(tokens, "tokens service instance exists")

            assert.strictEqual(tokens.getUserToken(), null, "initially null")
            assert.strictEqual(tokens.getAccessToken(), undefined, "access token undefined")

            const sampleToken = {
                user_access_token: "jwt-test-123",
                user_refresh_token: "refresh-test-456",
                user_id: "u-uuid-789",
                expires_at: 1800000000000
            }

            tokens.setUserToken(sampleToken)
            assert.deepEqual(tokens.getUserToken(), sampleToken, "user token matches")
            assert.strictEqual(tokens.getAccessToken(), "jwt-test-123", "access token matches")

            // AuthV2 delegates to tokens
            assert.deepEqual(sdk.authv2().getUserToken(), sampleToken, "authv2 delegates getUserToken")
            assert.strictEqual(sdk.authv2().getAccessToken(), "jwt-test-123", "authv2 delegates getAccessToken")

            // Clear tokens via authv2 or tokens
            sdk.authv2().clearTokens()
            assert.strictEqual(tokens.getUserToken(), null, "tokens cleared")
            assert.strictEqual(sdk.authv2().getAccessToken(), undefined, "access token cleared on authv2")
        })

        QUnit.test("tokens.exchange sends code and verifier and saves token", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            let capturedUrl = null
            let capturedBody = null

            globalThis.fetch = async (url, opts) => {
                capturedUrl = url
                capturedBody = JSON.parse(opts.body)
                return new Response(JSON.stringify({
                    user_access_token: "jwt-exchanged-abc",
                    user_refresh_token: "refresh-exchanged-def",
                    user_id: "u-exchanged-1",
                    expires_at: 1900000000000
                }), { status: 200 })
            }

            const exchanged = await tokens.exchange("auth-code-123", "verifier-456")
            assert.strictEqual(capturedUrl, "http://userz.test/api/userz/auth/exchange")
            assert.deepEqual(capturedBody, { code: "auth-code-123", code_verifier: "verifier-456" })
            assert.strictEqual(exchanged.user_access_token, "jwt-exchanged-abc")
            assert.strictEqual(tokens.getAccessToken(), "jwt-exchanged-abc")
        })

        QUnit.test("authv2.user sends bearer token from tokens service", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            // When no token is stored, user() returns null without making a request
            let fetchCalled = false
            globalThis.fetch = async () => {
                fetchCalled = true
                return new Response("{}", { status: 200 })
            }
            const userNoToken = await sdk.authv2().user()
            assert.strictEqual(userNoToken, null)
            assert.strictEqual(fetchCalled, false)

            // Store a token
            tokens.setUserToken({
                user_access_token: "my-jwt-token",
                user_refresh_token: "my-refresh",
                user_id: "u-1",
                expires_at: 1900000000000
            })

            let authHeader = null
            let reqUrl = null
            globalThis.fetch = async (url, opts) => {
                reqUrl = url
                authHeader = opts.headers['Authorization']
                return new Response(JSON.stringify({
                    user_id: "u-1",
                    email_redacted: "u***@infinity6.ai"
                }), { status: 200 })
            }

            const user = await sdk.authv2().user()
            assert.strictEqual(reqUrl, "http://userz.test/api/userz/auth/me")
            assert.strictEqual(authHeader, "Bearer my-jwt-token")
            assert.deepEqual(user, { user_id: "u-1", email_redacted: "u***@infinity6.ai" })
        })

    })

})(QUnit)
