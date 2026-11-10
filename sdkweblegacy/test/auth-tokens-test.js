((QUnit) => {

    QUnit.module("i6sdk-web auth & tokens", function (hooks) {

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
            assert.deepEqual(sdk.auth().getUserToken(), sampleToken, "auth delegates getUserToken")
            assert.strictEqual(sdk.auth().getAccessToken(), "jwt-test-123", "auth delegates getAccessToken")

            // Clear tokens via auth or tokens
            sdk.auth().clearTokens()
            assert.strictEqual(tokens.getUserToken(), null, "tokens cleared")
            assert.strictEqual(sdk.auth().getAccessToken(), undefined, "access token cleared on auth")
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

        QUnit.test("auth.user automatically refreshes expired access token", async function (assert) {
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

            const user = await sdk.auth().user()
            assert.strictEqual(calls.length, 2, "made refresh call then auth/me call")
            assert.strictEqual(calls[0].url, "http://userz.test/api/userz/auth/refresh")
            assert.strictEqual(calls[1].url, "http://userz.test/api/userz/auth/me")
            assert.strictEqual(calls[1].headers['Authorization'], "Bearer jwt-newly-refreshed")
            assert.deepEqual(user, { user_id: "u-1", email_redacted: "test@infinity6.ai" })
        })

        QUnit.test("dataset tokens cache per dataset and clear on clearTokens and exchange", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            tokens.setDatasetToken("dataset-1", {
                access_token: "ds-tok-1",
                dataset: "dataset-1",
                grants: 14,
                expires_at: Date.now() + 3600 * 1000
            })
            tokens.setDatasetToken("dataset-2", {
                access_token: "ds-tok-2",
                dataset: "dataset-2",
                grants: 8,
                expires_at: Date.now() + 3600 * 1000
            })

            assert.strictEqual(tokens.getDatasetToken("dataset-1").access_token, "ds-tok-1")
            assert.strictEqual(tokens.getDatasetToken("dataset-2").access_token, "ds-tok-2")
            assert.strictEqual(tokens.getDatasetToken("nonexistent"), null)

            // clearTokens clears both user and dataset tokens
            tokens.setUserToken({ user_access_token: "u-tok", user_id: "u-1" })
            tokens.clearTokens()
            assert.strictEqual(tokens.getUserToken(), null)
            assert.strictEqual(tokens.getDatasetToken("dataset-1"), null)
            assert.strictEqual(tokens.getDatasetToken("dataset-2"), null)

            // exchange clears previous dataset tokens
            tokens.setDatasetToken("old-ds", { access_token: "old", dataset: "old-ds", grants: 8, expires_at: Date.now() + 100000 })
            globalThis.fetch = async () => new Response(JSON.stringify({
                user_access_token: "new-u",
                user_refresh_token: "new-rt",
                user_id: "u-new",
                expires_at: Date.now() + 100000
            }), { status: 200 })

            await tokens.exchange("code", "verifier")
            assert.strictEqual(tokens.getDatasetToken("old-ds"), null, "dataset tokens cleared before new exchange token saved")
        })

        QUnit.test("getValidDatasetAccessToken returns cached when fresh, requests when expired", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            // User must be authenticated
            tokens.setUserToken({
                user_access_token: "user-jwt-123",
                user_id: "u-1",
                expires_at: Date.now() + 3600 * 1000
            })

            let dsCalls = 0
            let capturedReq = null
            globalThis.fetch = async (url, opts) => {
                dsCalls++
                capturedReq = { url, headers: opts.headers, body: JSON.parse(opts.body) }
                return new Response(JSON.stringify({
                    access_token: "binary-dataset-token-abc",
                    dataset: "fashion-store",
                    solution: "relevance_fashion",
                    grants: 14,
                    expires_at: Date.now() + 14400 * 1000
                }), { status: 200 })
            }

            // 1. Initial fetch
            const tok1 = await tokens.getValidDatasetAccessToken("fashion-store")
            assert.strictEqual(tok1, "binary-dataset-token-abc")
            assert.strictEqual(dsCalls, 1)
            assert.strictEqual(capturedReq.url, "http://userz.test/api/userz/dataset/access-token")
            assert.strictEqual(capturedReq.headers['Authorization'], "Bearer user-jwt-123")
            assert.deepEqual(capturedReq.body, { dataset: "fashion-store", grants: 15 })

            // 2. Second call returns cached token without network request
            const tok2 = await tokens.getValidDatasetAccessToken("fashion-store")
            assert.strictEqual(tok2, "binary-dataset-token-abc")
            assert.strictEqual(dsCalls, 1, "no second fetch made while cached token is valid")

            // 3. Expire the token in cache and verify it re-fetches
            tokens.setDatasetToken("fashion-store", {
                access_token: "expired-ds-tok",
                dataset: "fashion-store",
                grants: 14,
                expires_at: Date.now() - 5000
            })

            const tok3 = await tokens.getValidDatasetAccessToken("fashion-store")
            assert.strictEqual(tok3, "binary-dataset-token-abc")
            assert.strictEqual(dsCalls, 2, "re-fetched when token expired")
        })

        QUnit.test("concurrent getValidDatasetAccessToken deduplicates request in flight", async function (assert) {
            const sdk = new I6Sdk({ servUrls: { userz: 'http://userz.test' } })
            const tokens = sdk.tokens()

            tokens.setUserToken({
                user_access_token: "user-jwt-123",
                user_id: "u-1",
                expires_at: Date.now() + 3600 * 1000
            })

            let callCount = 0
            globalThis.fetch = async () => {
                callCount++
                await new Promise(r => setTimeout(r, 10))
                return new Response(JSON.stringify({
                    access_token: "deduped-ds-token",
                    dataset: "products",
                    grants: 10,
                    expires_at: Date.now() + 3600 * 1000
                }), { status: 200 })
            }

            const [r1, r2, r3] = await Promise.all([
                tokens.getValidDatasetAccessToken("products"),
                tokens.getValidDatasetAccessToken("products"),
                tokens.getValidDatasetAccessToken("products")
            ])

            assert.strictEqual(r1, "deduped-ds-token")
            assert.strictEqual(r2, "deduped-ds-token")
            assert.strictEqual(r3, "deduped-ds-token")
            assert.strictEqual(callCount, 1, "only 1 fetch made for concurrent dataset token requests")
        })

        QUnit.test("tokens uses custom storage object when provided in config", async function (assert) {
            const memoryMap = new Map()
            const customStorage = {
                getItem: (key) => memoryMap.has(key) ? memoryMap.get(key) : null,
                setItem: (key, val) => memoryMap.set(key, val),
                removeItem: (key) => memoryMap.delete(key),
            }

            const sdk = new I6Sdk({ storage: customStorage })
            const tokens = sdk.tokens()

            const sampleToken = {
                user_access_token: "jwt-custom-mem",
                user_id: "u-custom-1",
                expires_at: 1800000000000
            }

            tokens.setUserToken(sampleToken)
            assert.deepEqual(tokens.getUserToken(), sampleToken, "token saved to custom storage")
            assert.ok(memoryMap.has("i6_user_user_token"), "token stored inside custom storage map")
            assert.strictEqual(localStorage.getItem("i6_user_user_token"), null, "did not touch global localStorage")

            tokens.setDatasetToken("test-ds", {
                access_token: "ds-tok-custom",
                dataset: "test-ds",
                grants: 15,
                expires_at: 1800000000000
            })
            assert.ok(memoryMap.has("i6_dataset_access_tokens"), "dataset token saved to custom storage")
            assert.strictEqual(localStorage.getItem("i6_dataset_access_tokens"), null, "global localStorage untouched")

            tokens.clearTokens()
            assert.strictEqual(tokens.getUserToken(), null, "tokens cleared from custom storage")
            assert.strictEqual(memoryMap.size, 0, "custom storage map is empty after clear")
        })

        QUnit.test("tokens uses custom storage factory function when provided in config", async function (assert) {
            const memoryMap = new Map()
            const customStorage = {
                getItem: (key) => memoryMap.has(key) ? memoryMap.get(key) : null,
                setItem: (key, val) => memoryMap.set(key, val),
                removeItem: (key) => memoryMap.delete(key),
            }

            const sdk = new I6Sdk({ storage: () => customStorage })
            const tokens = sdk.tokens()

            const sampleToken = {
                user_access_token: "jwt-custom-func",
                user_id: "u-custom-2",
                expires_at: 1800000000000
            }

            tokens.setUserToken(sampleToken)
            assert.deepEqual(tokens.getUserToken(), sampleToken, "token saved via storage function")
            assert.ok(memoryMap.has("i6_user_user_token"), "token stored in map returned by function")

            tokens.clearTokens()
            assert.strictEqual(tokens.getUserToken(), null, "tokens cleared")
        })

        QUnit.test("auth stores code_verifier in sessionStorage and omits it from redirect URL", async function (assert) {
            const localMap = new Map()
            const sessionMap = new Map()
            const localAdapter = {
                getItem: (k) => localMap.get(k) || null,
                setItem: (k, v) => localMap.set(k, v),
                removeItem: (k) => localMap.delete(k),
            }
            const sessionAdapter = {
                getItem: (k) => sessionMap.get(k) || null,
                setItem: (k, v) => sessionMap.set(k, v),
                removeItem: (k) => sessionMap.delete(k),
            }

            const sdk = new I6Sdk({
                userExchangeCallbackUrl: "https://my.app/auth/callback",
                servUrls: { userz: "https://userz.test" },
                localStorage: localAdapter,
                sessionStorage: sessionAdapter,
            })

            const auth = sdk.auth()
            let navigatedUrl = null
            auth._navigate = (url) => { navigatedUrl = url }

            await auth.login("https://my.app/dashboard")

            // 1. code_verifier stored in sessionStorage
            const codeVerifier = sessionMap.get("i6_auth_code_verifier")
            assert.ok(codeVerifier, "code_verifier was stored in sessionStorage")
            assert.strictEqual(localMap.has("i6_auth_code_verifier"), false, "not in localStorage")

            // 2. Redirect URL has code_challenge and callback, but NOT code_verifier
            assert.ok(navigatedUrl.includes("https://userz.test/api/userz/auth?"), "navigated to userz login")
            const urlObj = new URL(navigatedUrl)
            const redirectParam = urlObj.searchParams.get("redirect")
            assert.ok(redirectParam.includes("callback="), "redirect param contains callback")
            assert.false(redirectParam.includes("code_verifier="), "code_verifier is NOT in callback URL")
            assert.false(urlObj.searchParams.has("code_verifier"), "code_verifier is not in query params")
            assert.ok(urlObj.searchParams.has("code_challenge"), "code_challenge is present")
        })

        QUnit.test("auth.handleCallback exchanges code using code_verifier from sessionStorage and deletes it", async function (assert) {
            const sessionMap = new Map()
            sessionMap.set("i6_auth_code_verifier", "verifier-stored-in-session-999")

            let exchangedCode = null
            let exchangedVerifier = null
            let redirectedTo = null

            const sdk = new I6Sdk({
                sessionStorage: {
                    getItem: (k) => sessionMap.get(k) || null,
                    setItem: (k, v) => sessionMap.set(k, v),
                    removeItem: (k) => sessionMap.delete(k),
                }
            })

            // Mock tokens.exchange
            sdk.tokens().exchange = async (c, v) => {
                exchangedCode = c
                exchangedVerifier = v
                return {}
            }

            const auth = sdk.auth()
            auth._redirect = (url) => { redirectedTo = url }

            // Temporarily mock window.location.search
            const origSearch = window.location.search
            try {
                window.history.replaceState({}, '', '?code=auth-code-777&callback=' + encodeURIComponent('https://my.app/home'))

                const handled = await auth.handleCallback()
                assert.true(handled, "handleCallback completed successfully")
                assert.strictEqual(exchangedCode, "auth-code-777", "exchanged received code")
                assert.strictEqual(exchangedVerifier, "verifier-stored-in-session-999", "exchanged code_verifier retrieved from sessionStorage")
                assert.strictEqual(sessionMap.has("i6_auth_code_verifier"), false, "code_verifier was deleted from sessionStorage after exchange")
                assert.strictEqual(redirectedTo, "https://my.app/home", "redirected to post-exchange callback")
            } finally {
                window.history.replaceState({}, '', origSearch || '/')
            }
        })

        QUnit.test("auth generates and stores CSRF state on login, and passes it to userz auth URL", async function (assert) {
            const sessionMap = new Map()
            const sdk = new I6Sdk({
                userExchangeCallbackUrl: "https://my.app/auth/callback",
                servUrls: { userz: "https://userz.test" },
                sessionStorage: {
                    getItem: (k) => sessionMap.get(k) || null,
                    setItem: (k, v) => sessionMap.set(k, v),
                    removeItem: (k) => sessionMap.delete(k),
                }
            })

            const auth = sdk.auth()
            let navigatedUrl = null
            auth._navigate = (url) => { navigatedUrl = url }

            await auth.login("https://my.app/dashboard")

            const stateInSession = sessionMap.get("i6_auth_state")
            assert.ok(stateInSession, "state was stored in sessionStorage")

            const urlObj = new URL(navigatedUrl)
            assert.strictEqual(urlObj.searchParams.get("state"), stateInSession, "state in login URL matches stored state")
        })

        QUnit.test("auth.handleCallback verifies matching CSRF state and cleans it up", async function (assert) {
            const sessionMap = new Map()
            sessionMap.set("i6_auth_code_verifier", "ver-123")
            sessionMap.set("i6_auth_state", "valid-csrf-state-abc")

            const sdk = new I6Sdk({
                sessionStorage: {
                    getItem: (k) => sessionMap.get(k) || null,
                    setItem: (k, v) => sessionMap.set(k, v),
                    removeItem: (k) => sessionMap.delete(k),
                }
            })
            sdk.tokens().exchange = async () => ({})

            const auth = sdk.auth()
            auth._redirect = () => {}

            const origSearch = window.location.search
            try {
                window.history.replaceState({}, '', '?code=c123&state=valid-csrf-state-abc')
                const handled = await auth.handleCallback()
                assert.true(handled, "callback succeeds with valid state")
                assert.strictEqual(sessionMap.has("i6_auth_state"), false, "state removed from session after check")
            } finally {
                window.history.replaceState({}, '', origSearch || '/')
            }
        })

        QUnit.test("auth.handleCallback rejects state mismatch with CSRF error", async function (assert) {
            const sessionMap = new Map()
            sessionMap.set("i6_auth_code_verifier", "ver-123")
            sessionMap.set("i6_auth_state", "legitimate-state")

            const sdk = new I6Sdk({
                sessionStorage: {
                    getItem: (k) => sessionMap.get(k) || null,
                    setItem: (k, v) => sessionMap.set(k, v),
                    removeItem: (k) => sessionMap.delete(k),
                }
            })
            sdk.tokens().exchange = async () => ({})

            const auth = sdk.auth()
            auth._redirect = () => {}

            const origSearch = window.location.search
            try {
                window.history.replaceState({}, '', '?code=c123&state=attacker-state')
                let errorCaught = null
                try {
                    await auth.handleCallback()
                } catch (err) {
                    errorCaught = err
                }
                assert.ok(errorCaught, "threw error on state mismatch")
                assert.ok(errorCaught.message.includes("State mismatch / CSRF attempt"), "error mentions CSRF attempt")
                assert.strictEqual(sessionMap.has("i6_auth_state"), false, "state was removed from session storage")
            } finally {
                window.history.replaceState({}, '', origSearch || '/')
            }
        })

    })

})(QUnit)
