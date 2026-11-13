((QUnit) => {

    QUnit.module("tokens-diagnostic", function () {

        QUnit.test("describeToken detects and parses user_jwt", async function (assert) {
            const { describeToken, base64UrlEncode } = await import('../static/tokens/tokens.js');

            const header = base64UrlEncode(new TextEncoder().encode(JSON.stringify({ alg: "EdDSA", typ: "JWT" })));
            const exp = Math.floor(Date.now() / 1000) + 3600;
            const payload = base64UrlEncode(new TextEncoder().encode(JSON.stringify({
                user_id: "usr_unit_test_123",
                email: "dev@company.com",
                aud: "default",
                exp: exp
            })));
            const sig = base64UrlEncode(new Uint8Array(64).fill(0x55));
            const jwt = `${header}.${payload}.${sig}`;

            const desc = describeToken(jwt);
            assert.strictEqual(desc.kind, "user_jwt", "token kind is user_jwt");
            assert.true(desc.checkable, "token is checkable on server");
            assert.strictEqual(desc.userId, "usr_unit_test_123", "user_id matches");
            assert.strictEqual(desc.payload.email, "dev@company.com", "email matches");
            assert.strictEqual(desc.isExpired, false, "token is not expired");
        });

        QUnit.test("buildCanonicalString sorts queries, strips i6_sig, and normalizes paths", async function (assert) {
            const { buildCanonicalString, sha256Hex } = await import('../static/tokens/tokens.js');

            const canon = await buildCanonicalString(
                "post",
                "/api/ds/dataset/dectest/solutions/",
                "b=2&i6_sig=strip_this&a=1",
                new TextEncoder().encode('{"item": 123}'),
                "tok_test_abc",
                1892934000,
                "00000000000000000000000000000000",
                false
            );

            const lines = canon.split("\n");
            assert.strictEqual(lines.length, 7, "canonical string has exactly 7 lines");
            assert.strictEqual(lines[0], "POST", "normalized method");
            assert.strictEqual(lines[1], "/api/ds/dataset/dectest/solutions", "normalized path without trailing slash");
            assert.strictEqual(lines[2], "a=1&b=2", "canonical query sorted and i6_sig removed");
            assert.strictEqual(lines[3], await sha256Hex('{"item": 123}'), "body sha256 matches");
            assert.strictEqual(lines[4], await sha256Hex('tok_test_abc'), "token sha256 matches");
            assert.strictEqual(lines[5], "1892934000", "exp matches");
            assert.strictEqual(lines[6], "00000000000000000000000000000000", "nonce matches");
        });

        QUnit.test("signRequest produces valid Type 2 envelope verifiable with public key", async function (assert) {
            const {
                signRequest,
                describeToken,
                importEd25519PrivateKey,
                getEd25519PublicKeyB64,
                base64UrlEncode,
                base64UrlDecodeToBytes
            } = await import('../static/tokens/tokens.js');

            const seed = new Uint8Array(32);
            crypto.getRandomValues(seed);
            const privKeyB64 = base64UrlEncode(seed);
            const privKey = await importEd25519PrivateKey(privKeyB64);
            const jti = await getEd25519PublicKeyB64(privKey);

            const exp = Math.floor(Date.now() / 1000) + 300;
            const res = await signRequest(
                privKeyB64,
                "POST",
                "/api/test",
                "param=value",
                new TextEncoder().encode("body payload"),
                "bound_access_token_123",
                exp,
                new Uint8Array(16),
                false
            );

            assert.ok(res.signatureHeader, "signature header generated");
            const parts = res.signatureHeader.split(".");
            assert.strictEqual(parts[0].length, 119, "envelope base64url length is exactly 119 characters");
            assert.strictEqual(parts[1], "bound_access_token_123", "bound access token preserved");

            // Local Describe verification
            const desc = describeToken(res.signatureHeader);
            assert.strictEqual(desc.kind, "request_signature", "detected request_signature");
            assert.strictEqual(desc.envelope.tokenType, 2, "token type is 2 (TypeSigEnvelope)");
            assert.strictEqual(desc.envelope.unsignedBody, false, "unsigned body is false");
            assert.strictEqual(desc.envelope.isTimeWindowed, true, "is time windowed replay mode");

            // Mathematical verification with Web Crypto
            const pubKeyBytes = base64UrlDecodeToBytes(jti);
            const pubCryptoKey = await crypto.subtle.importKey("raw", pubKeyBytes, { name: "Ed25519" }, false, ["verify"]);
            const sigBytes = base64UrlDecodeToBytes(res.envelopeB64).subarray(25, 89);
            const ok = await crypto.subtle.verify("Ed25519", pubCryptoKey, sigBytes, new TextEncoder().encode(res.canonicalString));
            assert.true(ok, "Ed25519 signature mathematically verified with public key");
        });

        QUnit.test("describeToken handles refresh assertions and dataset access tokens", async function (assert) {
            const { describeToken, base64UrlEncode } = await import('../static/tokens/tokens.js');

            // Refresh Assertion (121 bytes, Type 3)
            const assertionBuf = new Uint8Array(121);
            assertionBuf[0] = 3;
            const dv = new DataView(assertionBuf.buffer);
            dv.setBigUint64(1, 1992934000n, false);
            const assertionDesc = describeToken(base64UrlEncode(assertionBuf));
            assert.strictEqual(assertionDesc.kind, "refresh_assertion", "detected refresh_assertion");
            assert.strictEqual(assertionDesc.assertion.tokenType, 3, "assertion token type is 3");

            // Dataset Access Token (Type 1)
            const dsBuf = new Uint8Array([1, 10, 20, 30, 40]);
            const dsDesc = describeToken(base64UrlEncode(dsBuf));
            assert.strictEqual(dsDesc.kind, "dataset_access_token", "detected dataset_access_token");
            assert.strictEqual(dsDesc.byteLength, 5, "byte length matches");
        });

    });

})(QUnit);
