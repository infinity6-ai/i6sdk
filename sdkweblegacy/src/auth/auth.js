/**
 * @fileoverview Authentication client.
 */

/**
 * @typedef {Object} AuthConfig
 * @property {string} baseUrl - Base URL of the server.
 */

/**
 * Authentication client for handling user sessions and authentication requests.
 */
class Auth {
  /**
   * @param {AuthConfig} config - Configuration object.
   * @param {string} config.baseUrl - Base URL of the server.
   */
  constructor(config) {
    /**
     * @private
     * @type {AuthConfig}
     */
    this._config = config;
  }

  /**
   * Redirects the browser to the login page with the current URL as backurl.
   * @private
   */
  #redirectToLogin() {
    const backurl = location.href
    console.log("Backurl", backurl)

    const encoded = encodeURIComponent(backurl)
    console.log("Encoded", encoded)

    const loginUrl = `${this._config.baseUrl}/api/decsuite/auth/login?backurl=${encoded}`
    console.log("LoginUrl", loginUrl)
    location = loginUrl
  }

  /**
   * Fetches the current authenticated user profile.
   * @returns {Promise<Object|null>} The user profile data, or null if redirected to login.
   * @throws {Error} If the HTTP request fails.
   */
  async me() {
    const resp = await fetch(`${this._config.baseUrl}/api/decsuite/me`, {
      method: "GET",
    });

    if (!resp.ok) {
      throw new Error(`Get upload URL failed with status ${resp.status}`);
    }

    if (window.gambeta) {
        this.#redirectToLogin()
        return null
    } 

    const user = await resp.json();
    return user
  }
}

export { Auth };

