/**
 * @fileoverview Pipe client for starting pipelines.
 */

import { pipezStart } from "../sdkapis/sdkapis.js";
import { I6Error } from "../errors.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Parameters for starting a pipe.
 * @typedef {Object} StartParams
 * @property {string} dataset - Dataset the pipe runs on.
 * @property {string} pipename - Name of the pipe to start.
 * @property {{name: string, value: string}[]} [params] - Pipe parameters.
 * @property {*} solutionDomain - Solution/domain the pipe runs in.
 */

/**
 * Result of a successful start.
 * @typedef {Object} StartResult
 * @property {string} id - Id of the started pipe run.
 */

/**
 * Client for starting pipes.
 */
class Pipe {
  /**
   * Constructs a Pipe client instance.
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
   * Starts a pipe (pipez-start API).
   * @param {StartParams} params - Start parameters.
   * @returns {Promise<StartResult>} The id of the started run.
   * @throws {I6Error} If the request fails.
   */
  async start({ dataset, pipename, params = [], solutionDomain }) {
    const resp = await pipezStart(this._services, {
      params: { dataset },
      payload: { pipename, params, solution_domain: solutionDomain },
    });
    if (!resp.ok) {
      throw I6Error.fromResult("Start pipe failed", resp);
    }
    return { id: resp.payload.id };
  }
}

export { Pipe };
