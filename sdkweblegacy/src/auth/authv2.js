/**
 * @fileoverview Authentication client V2 using i6userz PKCE redirect flow.
 *
 * Flow Architecture:
 * 1. Login Initiation: {@link AuthV2#login} generates an RFC 7636 PKCE code verifier and S256 code
 *    challenge, appends the target application URL (`callback`) and `code_verifier` as query
 *    parameters to `config.userExchangeCallbackUrl`, and redirects the browser to `GET /api/userz/auth`.
 * 2. Identity Provider Authentication: The user signs in via Google or Email/Password on the i6userz UI.
 * 3. Authorization Code Redirect: i6userz redirects the browser to the intermediate exchange URL
 *    (`config.userExchangeCallbackUrl`) with a short-lived authorization `code` appended.
 * 4. Token Exchange: The exchange callback page invokes {@link AuthV2#handleCallback}, which extracts
 *    the `code` and `code_verifier`, calls `POST /api/userz/auth/exchange` on the userz service to exchange
 *    them for active tokens (access token JWT, refresh token, user ID, expiration), persists them in
 *    `localStorage` under `'i6_user_user_token'`, and redirects the browser to the original application URL.
 * 5. Profile Retrieval: {@link AuthV2#user} fetches user identity metadata via `GET /api/userz/auth/me`
 *    using Bearer token authentication.
 */

import { decodeUrlSafe } from "../internal/b64/b64.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * User profile metadata returned by the GET /api/userz/auth/me endpoint.
 * @typedef {Object} AuthPayload
 * @property {string} user_id - Authenticated user unique identifier UUID.
 * @property {string} [email_redacted] - Redacted email address for display (e.g. "u***@infinity6.ai").
 * @property {string[]} [providers] - List of linked identity providers (e.g. ["google", "password"]).
 * @property {number} [created_at] - Account creation timestamp in epoch milliseconds.
 * @property {string} [sub] - User ID (legacy alias for user_id).
 * @property {string} [display] - User display name (legacy alias).
 * @property {string} [email] - User email address (legacy alias).
 */

/**
 * User token credentials returned by the POST /api/userz/auth/exchange endpoint
 * and stored in localStorage under the key 'i6_user_user_token'.
 * @typedef {Object} UserToken
 * @property {string} user_access_token - Signed EdDSA user access token JWT.
 * @property {string} user_refresh_token - Opaque user refresh token used for session extension.
 * @property {string} user_id - Authenticated user unique identifier UUID.
 * @property {number} expires_at - Access token expiration timestamp in epoch milliseconds.
 */

/**
 * Authentication client V2 for handling browser redirects, PKCE, and token exchange with i6userz.
 * Configuration (callbackUrl, userExchangeCallbackUrl, app) is read directly from `i6sdk.config()`.
 */
class AuthV2 {
  /**
   * Constructs an AuthV2 client instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    /**
     * @private
     * @type {Services}
     */
    this._services = services;
  }

  /**
   * Retrieves the base URL for the i6userz authentication service from SDK configuration.
   * Resolution precedence:
   * 1. `config.servUrls["userz"]`
   * 2. `config.userzUrl`
   * 3. `config.baseUrl`
   * 4. Empty string fallback
   * @private
   * @returns {string} The base URL for userz.
   */
  _getUserzUrl() {
    const config = this._services.config();
    return config.servUrls?.["userz"] || config.userzUrl || config.baseUrl || "";
  }

  /**
   * Generates a cryptographically random standard PKCE code verifier (RFC 7636 §4.1).
   * Uses unreserved URL characters (`[A-Z]`, `[a-z]`, `[0-9]`, `-`, `.`, `_`, `~`).
   * @private
   * @param {number} [length=64] - The length of the code verifier (between 43 and 128 characters).
   * @returns {string} The generated code verifier string.
   */
  _generateCodeVerifier(length = 64) {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~";
    const array = new Uint8Array(length);
    if (typeof crypto !== "undefined" && crypto.getRandomValues) {
      crypto.getRandomValues(array);
    } else {
      for (let i = 0; i < length; i++) {
        array[i] = Math.floor(Math.random() * 256);
      }
    }
    let verifier = "";
    for (let i = 0; i < array.length; i++) {
      verifier += charset[array[i] % charset.length];
    }
    return verifier;
  }

  /**
   * Generates an RFC 7636 S256 PKCE code challenge from a code verifier.
   * Computes SHA-256 over the ASCII code verifier and base64url-encodes the digest without padding.
   * @private
   * @param {string} verifier - The plaintext code verifier.
   * @returns {Promise<string>} The base64url-encoded code challenge.
   * @throws {Error} If Web Crypto API (`crypto.subtle.digest`) is unavailable.
   */
  async _generateCodeChallenge(verifier) {
    if (typeof crypto !== "undefined" && crypto.subtle && crypto.subtle.digest) {
      const encoder = new TextEncoder();
      const data = encoder.encode(verifier);
      const digest = await crypto.subtle.digest("SHA-256", data);
      return this._base64UrlEncode(new Uint8Array(digest));
    }
    throw new Error("Web Crypto API required for PKCE code challenge");
  }

  /**
   * Encodes a byte array as an unpadded base64url string (RFC 4648 §5).
   * Replaces '+' with '-', '/' with '_', and strips trailing '=' padding.
   * @private
   * @param {Uint8Array|ArrayBuffer} bytes - The bytes to encode.
   * @returns {string} The base64url-encoded string.
   * @throws {Error} If no base64 encoder (`btoa` or Node.js `Buffer`) is available.
   */
  _base64UrlEncode(bytes) {
    const bytesArray = new Uint8Array(bytes);
    let binary = "";
    for (let i = 0; i < bytesArray.byteLength; i++) {
      binary += String.fromCharCode(bytesArray[i]);
    }
    let base64;
    if (typeof btoa !== "undefined") {
      base64 = btoa(binary);
    } else if (typeof Buffer !== "undefined") {
      base64 = Buffer.from(binary, "binary").toString("base64");
    } else {
      throw new Error("No base64 encoder found");
    }
    return base64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
  }

  /**
   * Redirects the browser to the i6userz login page (`GET /api/userz/auth`) to initiate the PKCE login flow.
   *
   * Validates that `config.userExchangeCallbackUrl` is configured, clears any existing token from
   * `localStorage`, generates a random PKCE `code_verifier` and S256 `code_challenge`, and redirects
   * the browser with `app`, `redirect`, and `code_challenge` query parameters.
   *
   * @param {string} [callbackUrl] - The final destination URL to return to after login and token exchange complete.
   *   Defaults to `window.location.href`.
   * @returns {Promise<void>}
   * @throws {string} If `config.userExchangeCallbackUrl` is not defined in the SDK configuration.
   */
  async login(callbackUrl) {
    const config = this._services.config();
    if (!config.userExchangeCallbackUrl) {
      throw 'config.userExchangeCallbackUrl is required'
    }
    localStorage.removeItem('i6_user_user_token')

    const app = config.app || "default";
    const codeVerifier = this._generateCodeVerifier();
    const codeChallenge = await this._generateCodeChallenge(codeVerifier);

    // Build the exchange callback URL with 'callback' and 'code_verifier' query parameters.
    // The i6userz auth service will append '&code=...' upon successful authentication.
    let cbUrl = config.userExchangeCallbackUrl;
    console.log('xxx', cbUrl)
    if (cbUrl.includes('?')) {
      cbUrl += '&'
    } else {
      cbUrl += '?'      
    }
    cbUrl += 'callback=' + encodeURIComponent(callbackUrl || window.location.href) + '&'
    cbUrl += 'code_verifier=' + encodeURIComponent(codeVerifier)

    const userzUrl = this._getUserzUrl();

    const params = new URLSearchParams({
      app: app,
      redirect: cbUrl,
      code_challenge: codeChallenge,
    });

    const loginUrl = `${userzUrl}/api/userz/auth?${params.toString()}`;
    window.location.href = loginUrl;
  }

  /**
   * Processes the authorization code present in the URL query string, exchanging it for tokens via
   * `POST /api/userz/auth/exchange`, securely storing the {@link UserToken} object in `localStorage`
   * under `'i6_user_user_token'`, and redirecting to the post-exchange `callback` URL.
   *
   * Note: This method must be called explicitly by the exchange callback page
   * (configured via `config.userExchangeCallbackUrl`). It is not called automatically by `user()`.
   *
   * @returns {Promise<boolean>} True if handling tokens completed (code was exchanged and redirect initiated),
   *   false if not an auth callback (no authorization `code` in the URL query parameters).
   * @throws {Error} If the token exchange request fails or the server returns an error status.
   */
  async handleCallback() {
    if (typeof window === "undefined") return false;

    const params = new URLSearchParams(window.location.search);
    const code = params.get("code");
    if (!code) {
      return false;
    }

    const codeVerifier = params.get('code_verifier')

    const userzUrl = this._getUserzUrl();
    const exchangeUrl = `${userzUrl}/api/userz/auth/exchange`;

    const resp = await fetch(exchangeUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        code: code,
        code_verifier: codeVerifier,
      }),
    });
    if (!resp.ok) {
      const errText = await resp.text();
      throw new Error(`Token exchange failed (${resp.status}): ${errText}`);
    }
    const data = await resp.json();
    console.log('aaa', data)
    localStorage.setItem('i6_user_user_token', JSON.stringify(data))
    const callback = params.get('callback')
    location = callback
    return true;
  }

  /**
   * Retrieves the stored user token credentials from localStorage.
   *
   * @param {boolean} [refreshed] - Optional flag indicating whether to return
   *   a freshly refreshed token.
   * @returns {UserToken|null} The stored user token object containing access and
   *   refresh tokens, or null if no valid session exists in storage.
   *
   * @example
   * // Stored object schema:
   * // {
   * //   "expires_at": 1794338065000,
   * //   "user_access_token": "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVC...",
   * //   "user_id": "u-172c6f15-155",
   * //   "user_refresh_token": "020c8d2e-d802-4de9-9077-9bc5def86c55..."
   * // }
   */
  getUserToken(refreshed) {
    const s = localStorage.getItem('i6_user_user_token')
    if (!s) {
      return null
    }
    let tk = JSON.parse(s)
    return tk
  }

  /**
   * Retrieves the active user access token JWT string for authorizing API requests.
   *
   * @returns {string|undefined} The signed EdDSA user access token JWT,
   *   or undefined if no user is authenticated.
   */
  getAccessToken() {
    return this.getUserToken(true)?.user_access_token
  }

  /**
   * Clears all stored authentication tokens and credentials from localStorage.
   * Removes the `'i6_user_user_token'` entry.
   * @returns {void}
   */
  clearTokens() {
    localStorage.removeItem("i6_user_user_token");
  }

  /**
   * Fetches the current authenticated user profile metadata (`GET /api/userz/auth/me`).
   * Sends the user access token in the `Authorization: Bearer <user_access_token>` header.
   * Does not redirect; use {@link AuthV2#requireUser} to automatically force a login redirect.
   *
   * Note: Expects tokens to already be present in storage. Call {@link AuthV2#handleCallback}
   * on the exchange page before attempting to fetch the user profile.
   *
   * @returns {Promise<AuthPayload|null>} The user profile data, or null if unauthenticated.
   * @throws {Error} If the network request fails or the userz service returns an error status.
   */
  async user() {
    const token = this.getAccessToken();
    if (!token) {
      return null;
    }
    const userzUrl = this._getUserzUrl()
    const resp = await fetch(`${userzUrl}/api/userz/auth/me`, {
      headers: {
        'Authorization': 'Bearer ' + token
      },
    });
    if (!resp.ok) {
      const errText = await resp.text();
      throw new Error(`Error getting logged user (${resp.status}): ${errText}`);
    }
    return resp.json()
  }

  /**
   * Fetches the current user profile, automatically redirecting the browser to the login page
   * via {@link AuthV2#login} when the user is unauthenticated.
   *
   * @returns {Promise<AuthPayload|null>} The user profile data, or null when redirecting to login.
   */
  async requireUser() {
    const user = await this.user();
    if (!user) {
      await this.login();
      return null;
    }
    return user;
  }

  /**
   * Logs out the current user by clearing stored tokens and navigating to `backurl`.
   *
   * @param {string} [backurl] - The URL to return to after logging out.
   *   Defaults to `config.callbackUrl` or the current window URL (`window.location.href`).
   * @returns {Promise<void>}
   */
  async logout(backurl = (typeof window !== "undefined" ? (this._services.config().callbackUrl || window.location.href) : "")) {
    this.clearTokens();
    window.location.href = backurl;
  }
}

export { AuthV2 };
