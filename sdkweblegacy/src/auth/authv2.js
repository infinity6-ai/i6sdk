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
  async login() {
    if (typeof window === "undefined") return;

    const config = this._services.config();
    const callbackUrl = config.callbackUrl || window.location.href;
    const exchangeCallbackUrl = config.exchangeCallbackUrl || null;
    const app = config.app || "default";

    const codeVerifier = this._generateCodeVerifier();
    const codeChallenge = await this._generateCodeChallenge(codeVerifier);

    sessionStorage.setItem("i6_auth_code_verifier", codeVerifier);
    sessionStorage.setItem("i6_auth_callback_url", callbackUrl);
    if (exchangeCallbackUrl) {
      sessionStorage.setItem("i6_auth_exchange_callback_url", exchangeCallbackUrl);
    }

    const userzUrl = this._getUserzUrl();
    const targetRedirect = exchangeCallbackUrl || callbackUrl;

    const params = new URLSearchParams({
      app: app,
      redirect: targetRedirect,
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

    const codeVerifier = sessionStorage.getItem("i6_auth_code_verifier");
    if (!codeVerifier) {
      return false;
    }

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

    if (data.user_refresh_token) {
      localStorage.setItem("i6_user_refresh_token", data.user_refresh_token);
    }
    if (data.user_access_token) {
      sessionStorage.setItem("i6_user_access_token", data.user_access_token);
    }
    if (data.user_id) {
      sessionStorage.setItem("i6_user_id", data.user_id);
    }
    if (data.expires_at) {
      sessionStorage.setItem("i6_user_expires_at", String(data.expires_at));
    }

    // Clean up code_verifier
    sessionStorage.removeItem("i6_auth_code_verifier");

    const exchangeCallbackUrl = sessionStorage.getItem("i6_auth_exchange_callback_url");
    const callbackUrl = sessionStorage.getItem("i6_auth_callback_url") || this._services.config().callbackUrl;

    sessionStorage.removeItem("i6_auth_exchange_callback_url");
    sessionStorage.removeItem("i6_auth_callback_url");

    if (exchangeCallbackUrl && callbackUrl && window.location.href !== callbackUrl) {
      window.location.replace(callbackUrl);
      await new Promise((r) => setTimeout(r, 5000));
      return true;
    } else {
      if (window.history && window.history.replaceState) {
        const cleanUrl = new URL(window.location.href);
        cleanUrl.searchParams.delete("code");
        window.history.replaceState({}, document.title, cleanUrl.toString());
      }
    }

    return true;
  }

  /**
   * Returns the current user access token stored in sessionStorage, or null if absent.
   * @returns {string|null}
   */
  getAccessToken() {
    return typeof sessionStorage !== "undefined" ? sessionStorage.getItem("i6_user_access_token") : null;
  }

  /**
   * Returns the current user refresh token stored in localStorage, or null if absent.
   * @returns {string|null}
   */
  getRefreshToken() {
    return typeof localStorage !== "undefined" ? localStorage.getItem("i6_user_refresh_token") : null;
  }

  /**
   * Clears all stored authentication tokens and credentials.
   * @returns {void}
   */
  clearTokens() {
    if (typeof localStorage !== "undefined") {
      localStorage.removeItem("i6_user_refresh_token");
    }
    if (typeof sessionStorage !== "undefined") {
      sessionStorage.removeItem("i6_user_access_token");
      sessionStorage.removeItem("i6_user_id");
      sessionStorage.removeItem("i6_user_expires_at");
    }
  }

  /**
   * Fetches the current authenticated user profile.
   * Automatically handles any code exchange if returning from login redirect.
   * Does not redirect; use {@link AuthV2#requireUser} to force a login.
   * @returns {Promise<AuthPayload|null>} The user profile data, or null if unauthenticated.
   */
  async user() {
    if (typeof window !== "undefined") {
      const params = new URLSearchParams(window.location.search);
      if (params.has("code")) {
        await this.handleCallback();
      }
    }

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
