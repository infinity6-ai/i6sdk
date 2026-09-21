/**
 * @fileoverview Authentication client.
 */

/**
 * @typedef {import("../index.js").I6SdkConfig} I6SdkConfig
 */

/**
 * Authentication client for handling user sessions and authentication requests.
 */
class Auth {
  constructor(services) {
    this._services = services;
  }

  /**
   * Redirects the browser to the login page with the current URL as backurl.
   * @private
   */
  #redirectToLogin() {
    const backurl = location.href
    const encoded = encodeURIComponent(backurl)
    const loginUrl = `${this._services.config().baseUrl}/api/decsuite/auth/login?backurl=${encoded}`

    location = loginUrl
  }

  /**
   * Fetches the current authenticated user profile.
   * @returns {Promise<Object|null>} The user profile data, or null if redirected to login.
   * @throws {Error} If the HTTP request fails.
   */
  async me() {
    const resp = await this._services.invoker().invoke({
      path: `api/auth/me`,
    });

    // https://decsuite.sandbox.rs.infinity6.ai/api/decsuite/_/auth/me


    // if (window.gambeta) {
    //     this.#redirectToLogin()
    //     return null
    // } 

    const user = await resp.json();

    if (!user) {
      this.#redirectToLogin();
      return null;
    }
    return user;
  }
}

export { Auth };

