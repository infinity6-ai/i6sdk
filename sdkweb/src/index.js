import { Ingestz } from "./ingestz/ingestz";

class I6Sdk {
  // async config(config) { //subscription e loja
    // if (!config) {
    //   throw new Error('SDK configuration is missing.');
    // }
    // this._config = config;
  // }

  hello() {
    return 'Hello from i6sdk-web!';
  }

  Ingestz() {
    return new Ingestz()
  }
}

window.I6Sdk = I6Sdk;
export { I6Sdk };
