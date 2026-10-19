/**
 * @fileoverview Error types thrown by the i6 Web Legacy SDK.
 */

/**
 * @typedef {"http"|"network"|"invalid_argument"|"no_upload_url"|"asset"|"model"|"query"} I6ErrorCode
 */

/**
 * Why an API call failed, derived from the HTTP status alone; the server message is never parsed.
 * @typedef {"unauthenticated"|"forbidden"|"not_found"|"bad_request"|"unavailable"|"network"} I6ApiErrorCode
 */

/**
 * Error of an API call that did not succeed, as returned in `I6Result.error`.
 * @typedef {Object} I6ApiError
 * @property {number} status - HTTP status; 0 when the server did not answer (network failure).
 * @property {I6ApiErrorCode} code - Failure kind, from the status: 401 "unauthenticated" (session expired),
 *   403 "forbidden" (no permission), 404 "not_found" (no data), other 4xx "bad_request", 5xx "unavailable".
 * @property {string} message - The server response text as is (e.g. the Stripe reason), or "HTTP <status>" if empty.
 * @property {*} [cause] - The raw response text, or the underlying error on a network failure.
 */

/**
 * What every generated API call resolves with; it does not throw on an HTTP error or a network failure.
 * Check `ok` to tell the two shapes apart. `body` is the parsed JSON of an API that has no declared payload.
 * @template R - The generated response type of the API (`{payload}`).
 * @typedef {({status: number, ok: true, body?: *} & R) | {status: number, ok: false, error: I6ApiError}} I6Result
 */

/**
 * Error thrown by the SDK when a call fails.
 *
 * `code` tells what kind of failure it was, so callers can branch on it instead of parsing
 * the message. `status` and `body` are set for `"http"` failures. `"asset"` failures also have `status`
 * when a file download was answered with an error (404: the file does not exist), and always `cause`.
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
      body: result.error?.message,
      cause: result.error,
    });
  }
}

/**
 * @param {number} status - HTTP status of a failed response.
 * @returns {I6ApiErrorCode}
 */
function codeFromStatus(status) {
  if (status === 401) return "unauthenticated"
  if (status === 403) return "forbidden"
  if (status === 404) return "not_found"
  if (status >= 500) return "unavailable"
  return "bad_request"
}

/**
 * Builds the error of a failed HTTP response, keeping the server text untouched.
 * @param {number} status - HTTP status.
 * @param {string} text - Response body as text.
 * @returns {I6ApiError}
 */
function apiErrorFromResponse(status, text) {
  return { status, code: codeFromStatus(status), message: text || `HTTP ${status}`, cause: text }
}

/**
 * Builds the error of a request that got no answer.
 * @param {*} cause - The error thrown by fetch.
 * @returns {I6ApiError}
 */
function apiErrorFromNetwork(cause) {
  return { status: 0, code: "network", message: cause?.message || "network error", cause }
}

export { I6Error, apiErrorFromResponse, apiErrorFromNetwork };
