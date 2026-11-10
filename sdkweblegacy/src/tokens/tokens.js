/**
 * @fileoverview Token management client for i6 SDK.
 * Manages storage, retrieval, and exchange of authentication and authorization tokens.
 */

/**
 * @typedef {import("../services/services.js").Services} Services
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

const USER_TOKEN_KEY = "i6_user_user_token";

/**
 * Client for managing authentication tokens, local storage persistence, and token exchanges.
 */
class Tokens {
  /**
   * Constructs a Tokens client instance.
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
    if (typeof localStorage === "undefined") return null;
    const s = localStorage.getItem(USER_TOKEN_KEY);
    if (!s) {
      return null;
    }
    try {
      return JSON.parse(s);
    } catch {
      return null;
    }
  }

  /**
   * Stores user token credentials in localStorage under `'i6_user_user_token'`.
   *
   * @param {UserToken} token - The user token credentials to store.
   * @returns {void}
   */
  setUserToken(token) {
    if (typeof localStorage === "undefined") return;
    localStorage.setItem(USER_TOKEN_KEY, JSON.stringify(token));
  }

  /**
   * Retrieves the active user access token JWT string for authorizing API requests.
   *
   * @returns {string|undefined} The signed EdDSA user access token JWT,
   *   or undefined if no user is authenticated.
   */
  getAccessToken() {
    return this.getUserToken(true)?.user_access_token;
  }

  /**
   * Clears all stored authentication tokens and credentials from localStorage.
   * Removes the `'i6_user_user_token'` entry.
   * @returns {void}
   */
  clearTokens() {
    if (typeof localStorage === "undefined") return;
    localStorage.removeItem(USER_TOKEN_KEY);
  }

  /**
   * Exchanges an authorization code and PKCE code verifier for user token credentials
   * via `POST /api/userz/auth/exchange`, and persists them in localStorage.
   *
   * @param {string} code - The authorization code returned by i6userz auth.
   * @param {string} codeVerifier - The PKCE code verifier matching the login challenge.
   * @returns {Promise<UserToken>} The exchanged user token credentials.
   * @throws {Error} If the token exchange request fails or the server returns an error status.
   */
  async exchange(code, codeVerifier) {
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
    this.setUserToken(data);
    return data;
  }
}

export { Tokens };
