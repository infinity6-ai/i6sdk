/**
 * @fileoverview Error types thrown by the i6 Web Legacy SDK.
 */

/**
 * @typedef {"http"|"network"|"invalid_argument"|"no_upload_url"|"asset"|"model"|"query"} I6ErrorCode
 */

/**
 * Error thrown by the SDK when a call fails.
 *
 * `code` tells what kind of failure it was, so callers can branch on it instead of parsing
 * the message. `status` and `body` are set for `"http"` failures.
 */
class I6Error extends Error {
  /**
   * @param {string} message - Human readable description.
   * @param {Object} [details]
   * @param {I6ErrorCode} [details.code] - Failure kind. Defaults to "http" when a status is given.
   * @param {number} [details.status] - HTTP status code, when the server answered.
   * @param {*} [details.body] - Error body of the response (parsed JSON, or raw text).
   * @param {*} [details.cause] - The underlying error, if any.
   */
  constructor(message, { code, status, body, cause } = {}) {
    super(message, cause === undefined ? undefined : { cause });
    this.name = "I6Error";
    /** @type {I6ErrorCode|undefined} */
    this.code = code ?? (status === undefined ? undefined : "http");
    /** @type {number|undefined} */
    this.status = status;
    /** @type {*} */
    this.body = body;
  }

  /**
   * Builds an "http" error from a failed generated-client result ({status, error}).
   * @param {string} message
   * @param {{status: number, error?: *}} result
   * @returns {I6Error}
   */
  static fromResult(message, result) {
    return new I6Error(`${message} with status ${result.status}`, {
      code: "http",
      status: result.status,
      body: result.error,
    });
  }
}

export { I6Error };
