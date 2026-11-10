/**
 * @fileoverview Authentication client V2 using i6userz PKCE redirect flow.
 */

import { decodeUrlSafe } from "../internal/b64/b64.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Configuration options for the i6 SDK Auth payload.
 * @typedef {Object} AuthPayload
 * @property {string} [sub] - User ID.
 * @property {string} [user_id] - User ID.
 * @property {string} [display] - Display name.
 * @property {string} [email] - User email.
 */

/**
 * Authentication client V2 for handling browser redirects, PKCE, and token exchange with i6userz.
 * Configuration (callbackUrl, exchangeCallbackUrl, app) is read directly from `i6sdk.config()`.
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
   * Retrieves the base URL for the userz service from config.
   * @private
   * @returns {string} The base URL for userz.
   */
  _getUserzUrl() {
    const config = this._services.config();
    return config.servUrls?.["userz"] || config.userzUrl || config.baseUrl || "";
  }

  /**
   * Generates a random standard PKCE code verifier (64 chars).
   * @private
   * @param {number} [length=64] - The length of the code verifier.
   * @returns {string} The generated code verifier.
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
   * Generates a S256 PKCE code challenge from a code verifier.
   * @private
   * @param {string} verifier - The code verifier.
   * @returns {Promise<string>} The base64url-encoded code challenge.
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
   * Encodes a byte array as a base64url string without padding.
   * @private
   * @param {Uint8Array|ArrayBuffer} bytes - The bytes to encode.
   * @returns {string} The base64url-encoded string.
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
   * Redirects the browser to the i6userz login page to initiate the PKCE login flow.
   * Reads callbackUrl, exchangeCallbackUrl, and app directly from i6sdk.config().
   * @returns {Promise<void>}
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

    let cbUrl = config.userExchangeCallbackUrl  //+ '?' + callbackUrl || window.location.href;
    console.log('xxx', cbUrl)
    if (cbUrl.includes('?')) {
      cbUrl += '&'
    } else {
      cbUrl += '?'      
    }
    cbUrl += 'callback=' + encodeURIComponent(callbackUrl || window.location.href) + '&'
    cbUrl += 'code_verifier=' + encodeURIComponent(codeVerifier)

    // sessionStorage.setItem("i6_auth_code_verifier", codeVerifier);

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
   * Processes the authorization code present in the URL query string, exchanging it for tokens,
   * securely storing them (refresh token in localStorage, auth token in sessionStorage),
   * and redirecting to callbackUrl if an exchangeCallbackUrl was configured.
   * Automatically called by `user()` or `requireUser()`.
   * @returns {Promise<boolean>} True if handling tokens completed, false if not an auth callback.
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
    // sessionStorage.removeItem("i6_auth_code_verifier");
    return true;
  }

  getUserToken(refreshed) {
    const s = locatStorage.getItem('i6_user_user_token')
    if (!s) {
      return null
    }
    let tk = JSON.parse(s)
    // {"expires_at":1794338065000,"user_access_token":"eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCIsImtpZCI6Ijg3Qzc2RjQzMDI5NzMwMUUifQ.eyJqdGkiOiJmMTY3ZDFkMS1iMGFjLTQ1YjEtYjZkZi00NjIxYTlhOTJlZWUiLCJzdWIiOiJ1LTE3MmM2ZjE1LTE1NSIsImV4cCI6MTc5NDMzODA2NX0.T9p2Xdl9qnA3Dl1hFmGncN1J2cyRVd5BihTryaW1vsI8ykP2Ozh9wyaqHPQkGP3oNsiicKT1XoHXGfLAlpBpBw","user_id":"u-172c6f15-155","user_refresh_token":"020c8d2e-d802-4de9-9077-9bc5def86c55.1c18f90fbfd1299142cb7218a8435af29109e161c885b4189a5adfbffd5487e4"}
    return tk
  }

  getAccessToken() {
    return this.getUserToken(true)?.user_access_token
  }

  /**
   * Clears all stored authentication tokens and credentials.
   * @returns {void}
   */
  clearTokens() {
      localStorage.removeItem("i6_user_user_token");
  }

  /**
   * Fetches the current authenticated user profile.
   * Automatically handles any code exchange if returning from login redirect.
   * Does not redirect; use {@link AuthV2#requireUser} to force a login.
   * @returns {Promise<AuthPayload|null>} The user profile data, or null if unauthenticated.
   */
  async user() {
    const token = this.getAccessToken();
    if (!token) {
      return null;
    }

    const expiresAtStr = typeof sessionStorage !== "undefined" ? sessionStorage.getItem("i6_user_expires_at") : null;
    if (expiresAtStr) {
      const expiresAt = parseInt(expiresAtStr, 10);
      if (!isNaN(expiresAt) && Date.now() > expiresAt) {
        this.clearTokens();
        return null;
      }
    }

    const userId = typeof sessionStorage !== "undefined" ? sessionStorage.getItem("i6_user_id") : null;

    let payload = {};
    try {
      const parts = token.split(".");
      if (parts.length >= 2) {
        payload = JSON.parse(decodeUrlSafe(parts[1]));
      }
    } catch (e) {
      // Intentionally ignore JWT parse errors
    }

    return {
      sub: userId || payload.sub,
      user_id: userId || payload.sub,
      display: payload.display || payload.email || userId || payload.sub,
      ...payload,
    };
  }

  /**
   * Fetches the current user, redirecting the browser to the login page when there is none.
   * @returns {Promise<AuthPayload|null>} The user profile data, or null when redirecting.
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
   * @param {string} [backurl] - The URL to return to after logging out.
   *   Defaults to config.callbackUrl or current window URL.
   * @returns {Promise<void>}
   */
  async logout(backurl = (typeof window !== "undefined" ? (this._services.config().callbackUrl || window.location.href) : "")) {
    this.clearTokens();
    if (typeof window !== "undefined" && backurl) {
      window.location.href = backurl;
    }
  }
}

export { AuthV2 };
