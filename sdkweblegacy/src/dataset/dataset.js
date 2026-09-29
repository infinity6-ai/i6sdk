/**
 * @fileoverview Dataset client for selecting the solution domain of a dataset.
 */

import { apis } from "../sdkapis/sdkapis.js";
import { I6Error } from "../errors.js";

/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * Parameters for selecting a solution domain.
 * @typedef {Object} SelectDomainParams
 * @property {string} dataset - Dataset to configure.
 * @property {*} domain - Domain to select.
 * @property {*} solution - Solution the domain belongs to.
 * @property {string} store - Store the selection applies to.
 */

/**
 * Client for dataset operations.
 */
class Dataset {
  /**
   * Constructs a Dataset client instance.
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
   * Selects the solution domain of a dataset (dsz-domain-select API).
   * @param {SelectDomainParams} params - Selection parameters.
   * @returns {Promise<void>} Resolves when the domain is selected.
   * @throws {I6Error} If the request fails.
   */
  async selectDomain({ dataset, domain, solution, store }) {
    const resp = await apis.apiDszDomainSelect(this._services, {
      params: { dataset },
      payload: { domain, solution, store },
    });
    if (!resp.ok) {
      throw I6Error.fromResult("Select domain failed", resp);
    }
  }
}

export { Dataset };
