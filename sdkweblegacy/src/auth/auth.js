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
    console.log("Backurl", backurl)

    const encoded = encodeURIComponent(backurl)
    console.log("Encoded", encoded)

    const loginUrl = `${backurl}api/decsuite/auth/login?backurl=${backurl}`
    console.log("LoginUrl", loginUrl)
    location = loginUrl
  }

  async me() {
    console.log("ME", window.gambeta)
    if (window.gambeta) {
        this.#redirectToLogin()
        return null
    } 

    const resp = await fetch(`${this._config.baseUrl}/api/decsuite/me`, {
      method: "GET",
    });

    if (!resp.ok) {
      throw new Error(`Get upload URL failed with status ${resp.status}`);
    }

    // if (window.gambeta) {
    //     this.#redirectToLogin()
    //     return null
    // } 

    const user = await resp.json();
    return user
  }
}

export { Auth };
