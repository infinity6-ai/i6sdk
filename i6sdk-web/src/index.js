class I6Sdk {
  async config(config) {
    if (!config) {
      throw new Error('SDK configuration is missing.');
    }
    this._config = config;
  }

  hello() {
    return 'Hello from i6sdk-web!';
  }
}

// Expose the SDK to the global scope
window.I6Sdk = I6Sdk;
export { I6Sdk };
