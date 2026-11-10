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
                    expires_at: 1900000000000,
                    refresh_token_expires_at: 2000000000000
                }), { status: 200 })
            }

            const exchanged = await tokens.exchange("auth-code-123", "verifier-456")
            assert.strictEqual(capturedUrl, "http://userz.test/api/userz/auth/exchange")
            assert.deepEqual(capturedBody, { code: "auth-code-123", code_verifier: "verifier-456" })
            assert.strictEqual(exchanged.user_access_token, "jwt-exchanged-abc")
            assert.strictEqual(tokens.getAccessToken(), "jwt-exchanged-abc")
            assert.strictEqual(tokens.getUserToken().refresh_token_expires_at, 2000000000000)
        })

        QUnit.test("expiration checks correctly identify expired vs fresh tokens", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            assert.true(tokens.isAccessTokenExpired(null), "null token is expired")
            assert.true(tokens.isRefreshTokenExpired(null), "null refresh token is expired")

            const now = Date.now()
            const freshToken = {
                user_access_token: "jwt-fresh",
                user_refresh_token: "rt-fresh",
                user_id: "u-1",
                expires_at: now + 3600 * 1000, // 1 hour in future
                refresh_token_expires_at: now + 86400 * 1000 // 1 day in future
            }
            assert.false(tokens.isAccessTokenExpired(freshToken), "future token is not expired")
            assert.false(tokens.isRefreshTokenExpired(freshToken), "future refresh token is not expired")

            // Within skew window (e.g. expiring in 10s with 30s skew)
            const almostExpiredToken = {
                ...freshToken,
                expires_at: now + 10 * 1000
            }
            assert.true(tokens.isAccessTokenExpired(almostExpiredToken, 30), "token within skew is considered expired")
            assert.false(tokens.isAccessTokenExpired(almostExpiredToken, 5), "token with smaller skew is not expired")

            const expiredToken = {
                ...freshToken,
                expires_at: now - 1000,
                refresh_token_expires_at: now - 500
            }
            assert.true(tokens.isAccessTokenExpired(expiredToken), "past token is expired")
            assert.true(tokens.isRefreshTokenExpired(expiredToken), "past refresh token is expired")
        })

        QUnit.test("getValidAccessToken returns existing token when fresh, refreshes when expired", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            const now = Date.now()
            tokens.setUserToken({
                user_access_token: "jwt-already-valid",
                user_refresh_token: "rt-valid",
                user_id: "u-1",
                expires_at: now + 3600 * 1000
            })

            let refreshCalls = 0
            globalThis.fetch = async () => {
                refreshCalls++
                return new Response("{}", { status: 200 })
            }

            const validToken = await tokens.getValidAccessToken()
            assert.strictEqual(validToken, "jwt-already-valid")
            assert.strictEqual(refreshCalls, 0, "no refresh fetch performed when token is valid")

            // Now simulate expired access token
            tokens.setUserToken({
                user_access_token: "jwt-old-expired",
                user_refresh_token: "rt-rot-1",
                user_id: "u-1",
                expires_at: now - 10000
            })

            let capturedRefreshUrl = null
            let capturedRefreshBody = null
            globalThis.fetch = async (url, opts) => {
                refreshCalls++
                capturedRefreshUrl = url
                capturedRefreshBody = JSON.parse(opts.body)
                return new Response(JSON.stringify({
                    user_access_token: "jwt-renewed-fresh",
                    user_refresh_token: "rt-rot-2",
                    user_id: "u-1",
                    expires_at: now + 3600 * 1000,
                    refresh_token_expires_at: now + 86400 * 1000
                }), { status: 200 })
            }

            const renewed = await tokens.getValidAccessToken()
            assert.strictEqual(capturedRefreshUrl, "http://userz.test/api/userz/auth/refresh")
            assert.deepEqual(capturedRefreshBody, { refresh_token: "rt-rot-1" })
            assert.strictEqual(renewed, "jwt-renewed-fresh")
            assert.strictEqual(tokens.getAccessToken(), "jwt-renewed-fresh")
            assert.strictEqual(tokens.getUserToken().user_refresh_token, "rt-rot-2", "refresh token rotated")
        })

        QUnit.test("concurrent getValidAccessToken deduplicates refresh in flight", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            tokens.setUserToken({
                user_access_token: "jwt-expired",
                user_refresh_token: "rt-initial",
                user_id: "u-1",
                expires_at: Date.now() - 5000
            })

            let refreshCount = 0
            globalThis.fetch = async () => {
                refreshCount++
                // Small delay to simulate network call
                await new Promise(r => setTimeout(r, 10))
                return new Response(JSON.stringify({
                    user_access_token: "jwt-deduped",
                    user_refresh_token: "rt-rotated",
                    user_id: "u-1",
                    expires_at: Date.now() + 3600 * 1000
                }), { status: 200 })
            }

            // Launch 3 concurrent calls
            const [t1, t2, t3] = await Promise.all([
                tokens.getValidAccessToken(),
                tokens.getValidAccessToken(),
                tokens.getValidAccessToken()
            ])

            assert.strictEqual(t1, "jwt-deduped")
            assert.strictEqual(t2, "jwt-deduped")
            assert.strictEqual(t3, "jwt-deduped")
            assert.strictEqual(refreshCount, 1, "exactly one refresh network call was made")
        })

        QUnit.test("refresh 401 response wipes stored session and returns null", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            tokens.setUserToken({
                user_access_token: "jwt-expired",
                user_refresh_token: "rt-revoked",
                user_id: "u-1",
                expires_at: Date.now() - 1000
            })

            globalThis.fetch = async () => {
                return new Response("token revoked", { status: 401 })
            }

            const token = await tokens.getValidAccessToken()
            assert.strictEqual(token, null, "null returned on revoked refresh")
            assert.strictEqual(tokens.getUserToken(), null, "tokens cleared from storage on 401")
        })

        QUnit.test("authv2.user automatically refreshes expired access token", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            // Expired access token but valid refresh token
            tokens.setUserToken({
                user_access_token: "jwt-stale",
                user_refresh_token: "rt-valid",
                user_id: "u-1",
                expires_at: Date.now() - 5000
            })

            const calls = []
            globalThis.fetch = async (url, opts) => {
                calls.push({ url, headers: opts?.headers })
                if (url.endsWith("/api/userz/auth/refresh")) {
                    return new Response(JSON.stringify({
                        user_access_token: "jwt-newly-refreshed",
                        user_refresh_token: "rt-new",
                        user_id: "u-1",
                        expires_at: Date.now() + 3600 * 1000
                    }), { status: 200 })
                }
                if (url.endsWith("/api/userz/auth/me")) {
                    return new Response(JSON.stringify({
                        user_id: "u-1",
                        email_redacted: "test@infinity6.ai"
                    }), { status: 200 })
                }
                return new Response("", { status: 404 })
            }

            const user = await sdk.authv2().user()
            assert.strictEqual(calls.length, 2, "made refresh call then auth/me call")
            assert.strictEqual(calls[0].url, "http://userz.test/api/userz/auth/refresh")
            assert.strictEqual(calls[1].url, "http://userz.test/api/userz/auth/me")
            assert.strictEqual(calls[1].headers['Authorization'], "Bearer jwt-newly-refreshed")
            assert.deepEqual(user, { user_id: "u-1", email_redacted: "test@infinity6.ai" })
        })

    })

})(QUnit)
