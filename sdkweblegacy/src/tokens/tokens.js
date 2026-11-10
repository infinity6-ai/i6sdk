/**
 * @fileoverview Token management client for i6 SDK.
 * Manages storage, retrieval, exchange, and silent renewal of authentication and authorization tokens.
 */

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * User token credentials returned by the POST /api/userz/auth/exchange
 * and POST /api/userz/auth/refresh endpoints, stored in localStorage under 'i6_user_user_token'.
 * @typedef {Object} UserToken
 * @property {string} user_access_token - Signed EdDSA user access token JWT.
 * @property {string} user_refresh_token - Opaque user refresh token used for session extension.
 * @property {string} user_id - Authenticated user unique identifier UUID.
 * @property {number} expires_at - Access token expiration timestamp in epoch milliseconds.
 * @property {number} [refresh_token_expires_at] - Refresh token expiration timestamp in epoch milliseconds.
 */

const USER_TOKEN_KEY = "i6_user_user_token";

/**
 * Client for managing authentication tokens, local storage persistence, token exchanges,
 * and automatic token renewal.
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

    /**
     * In-flight refresh promise to prevent duplicate refresh calls during concurrent operations.
     * @private
     * @type {Promise<UserToken|null>|null}
     */
    this._refreshPromise = null;
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
   * Checks whether the given user access token is expired or will expire within the given skew window.
   *
   * @param {UserToken|null} [token] - The token object to inspect. If omitted, inspects stored token.
   * @param {number} [skewSeconds=30] - Safety margin in seconds before actual expiration timestamp.
   * @returns {boolean} True if the access token is missing, expired, or expiring within skewSeconds.
   */
  isAccessTokenExpired(token = this.getUserToken(), skewSeconds = 30) {
    if (!token || !token.user_access_token || typeof token.expires_at !== "number") {
      return true;
    }
    return Date.now() >= (token.expires_at - skewSeconds * 1000);
  }

  /**
   * Checks whether the user refresh token is expired.
   *
   * @param {UserToken|null} [token] - The token object to inspect. If omitted, inspects stored token.
   * @returns {boolean} True if the refresh token is missing or its expiration timestamp has elapsed.
   */
  isRefreshTokenExpired(token = this.getUserToken()) {
    if (!token || !token.user_refresh_token) {
      return true;
    }
    if (typeof token.refresh_token_expires_at === "number") {
      return Date.now() >= token.refresh_token_expires_at;
    }
    return false;
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
   * //   "refresh_token_expires_at": 1825874065000,
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
   * Retrieves the raw stored user access token JWT string without checking expiration or refreshing.
   * Use {@link Tokens#getValidAccessToken} when fresh/valid credentials are required.
   *
   * @returns {string|undefined} The signed EdDSA user access token JWT,
   *   or undefined if no user is authenticated.
   */
  getAccessToken() {
    return this.getUserToken()?.user_access_token;
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

  /**
   * Refreshes the user access token using the stored user refresh token via `POST /api/userz/auth/refresh`.
   * Enforces token family rotation by storing the newly rotated refresh token and fresh access token.
   * Deduplicates concurrent in-flight refresh requests to prevent token reuse conflicts.
   *
   * @returns {Promise<UserToken|null>} The refreshed token credentials, or null if refresh failed / token revoked.
   */
  async refresh() {
    if (this._refreshPromise) {
      return this._refreshPromise;
    }

    this._refreshPromise = (async () => {
      try {
        const token = this.getUserToken();
        if (!token || !token.user_refresh_token || this.isRefreshTokenExpired(token)) {
          this.clearTokens();
          return null;
        }

        const userzUrl = this._getUserzUrl();
        const refreshUrl = `${userzUrl}/api/userz/auth/refresh`;

        const resp = await fetch(refreshUrl, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            refresh_token: token.user_refresh_token,
          }),
        });

        if (!resp.ok) {
          // If unauthorized, token is revoked or family reuse was detected; wipe stored session
          if (resp.status === 401 || resp.status === 400) {
            this.clearTokens();
          }
          return null;
        }

        const data = await resp.json();
        this.setUserToken(data);
        return data;
      } catch {
        return null;
      } finally {
        this._refreshPromise = null;
      }
    })();

    return this._refreshPromise;
  }

  /**
   * Retrieves a guaranteed-valid user access token JWT string.
   *
   * If the stored access token is still valid (not expired within `skewSeconds`), it is returned immediately.
   * If the access token is expired or close to expiring, an automatic silent refresh is performed
   * using the user refresh token.
   * If the user is unauthenticated or the refresh token has expired/been revoked, returns null.
   *
   * @param {number} [skewSeconds=30] - Safety margin in seconds before expiration.
   * @returns {Promise<string|null>} The valid user access token JWT, or null if unauthenticated or session expired.
   */
  async getValidAccessToken(skewSeconds = 30) {
    const token = this.getUserToken();
    if (!token) {
      return null;
    }

    if (!this.isAccessTokenExpired(token, skewSeconds)) {
      return token.user_access_token;
    }

    // Access token is expired or expiring soon; refresh session
    const refreshed = await this.refresh();
    return refreshed?.user_access_token || null;
  }
}

export { Tokens };
