class Auth {
  /**
   * @param {Object} config
   * @param {string} config.baseUrl - Default base URL of the ingestz server.
   */
  constructor(config) {
    this._config = config;
  }

  #redirectToLogin() {
    const backurl = location.href
    const loginUrl = `${backurl}/api/decsuite/auth/login?backurl=${encodeURIComponent(backurl)}`
    location = loginUrl
  }

  async me() {
    const resp = await fetch(`${baseUrl}/api/decsuite/me`, {
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
