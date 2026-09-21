/**
 * @fileoverview Authentication client.
 */

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Authenticated user profile data.
 * @typedef {Record<string, any>} UserProfile
 */

/**
 * Authentication client for handling user sessions and authentication requests.
 */
class Auth {
  /**
   * Constructs an Auth client instance.
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
   * Redirects the browser to the login page with the current URL as backurl.
   * @private
   * @returns {void}
   */
  #redirectToLogin() {
    const backurl = location.href
    const encoded = encodeURIComponent(backurl)
    const loginUrl = `${this._services.config().baseUrl}/api/decsuite/auth/login?backurl=${encoded}`

    location = loginUrl
  }

  /**
   * Fetches the current authenticated user profile.
   * @returns {Promise<UserProfile|null>} The user profile data, or null if redirected to login.
   * @throws {Error} If the HTTP request fails.
   */
  async me() {
    const resp = await this._services.invoker().invoke({
      path: `api/auth/me`,
    });
    const user = await resp.json();
    if (!user) {
      this.#redirectToLogin();
      return null;
    }
    return user;
  }
}

export { Auth };

