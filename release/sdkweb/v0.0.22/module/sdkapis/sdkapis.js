class Apis {
  /**
   * Constructs an Ingest client instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    this._services = services;
  }
  /**
   * Create Customer Portal
   * 
   * Creates a customer portal session for managing subscription and billing.
   * 
   * @param {apiDszCreateCustomerPortalReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszCreateCustomerPortalResp>>}
   */
  async dszCreateCustomerPortal(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/portal`,
      json: req.payload,
      opts
    });
  }
  /**
   * Create Ingestion Token
   * 
   * Creates a legacy ingestion token for dataset ingestion.
   * 
   * @param {apiDszCreateIngestionTokenReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszCreateIngestionTokenResp>>}
   */
  async dszCreateIngestionToken(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/ingestiontoken/create`,
      json: req.payload,
      opts
    });
  }
  /**
   * Create Recsys Token
   * 
   * Creates a legacy recsys token for real-time recommendation querying.
   * 
   * @param {apiDszCreateRecsysTokenReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszCreateRecsysTokenResp>>}
   */
  async dszCreateRecsysToken(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/recsystoken/create`,
      json: req.payload,
      opts
    });
  }
  /**
   * Create Store
   * 
   * Registers a new store in a dataset.
   * 
   * @param {apiDszCreateStoreReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszCreateStoreResp>>}
   */
  async dszCreateStore(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/store/create`,
      json: req.payload,
      opts
    });
  }
  /**
   * List Active Dataset Solutions
   * 
   * Lists all active solution domains enabled for a dataset.
   * 
   * @param {apiDszDatasetSolutionsReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszDatasetSolutionsResp>>}
   */
  async dszDatasetSolutions(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/solutions`,
      opts
    });
  }
  /**
   * Delete Dataset Token
   * 
   * Revokes a legacy dataset token.
   * 
   * @param {apiDszDeleteTokenReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszDeleteTokenResp>>}
   */
  async dszDeleteToken(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "DELETE",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/token/delete`,
      query: req.query,
      opts
    });
  }
  /**
   * Get Store Solution Domain
   * 
   * Retrieves the solution domain for a store within a dataset.
   * 
   * @param {apiDszDomainGetReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszDomainGetResp>>}
   */
  async dszDomainGet(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/store/solution-domain`,
      query: req.query,
      opts
    });
  }
  /**
   * Get Subscription
   * 
   * Gets the plan and status of the Stripe subscription for a dataset solution.
   * 
   * @param {apiDszGetSubscriptionReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszGetSubscriptionResp>>}
   */
  async dszGetSubscription(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/solution/${encodeURIComponent(String(req.params.solution))}/subscription`,
      opts
    });
  }
  /**
   * List Datasets
   * 
   * Lists all datasets available to the user.
   * 
   * @param {apiDszListDatasetsReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszListDatasetsResp>>}
   */
  async dszListDatasets(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "user",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/user-id/${encodeURIComponent(String(req.params.userId))}/datasets`,
      opts
    });
  }
  /**
   * List Ingestion Tokens
   * 
   * Lists all ingestion tokens provisioned for a dataset.
   * 
   * @param {apiDszListIngestionTokensReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszListIngestionTokensResp>>}
   */
  async dszListIngestionTokens(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/ingestiontoken/list`,
      opts
    });
  }
  /**
   * List Payments
   * 
   * Lists all Stripe payment records for a dataset.
   * 
   * @param {apiDszListPaymentsReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszListPaymentsResp>>}
   */
  async dszListPayments(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/payments`,
      opts
    });
  }
  /**
   * List Subscription Plans
   * 
   * Lists all available subscription pricing plans for a solution product.
   * 
   * @param {apiDszListPlansReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszListPlansResp>>}
   */
  async dszListPlans(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/solution/${encodeURIComponent(String(req.params.solution))}/plans`,
      opts
    });
  }
  /**
   * List Recsys Tokens
   * 
   * Lists all recsys tokens provisioned for a dataset.
   * 
   * @param {apiDszListRecsysTokensReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszListRecsysTokensResp>>}
   */
  async dszListRecsysTokens(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/recsystoken/list`,
      opts
    });
  }
  /**
   * List Stores
   * 
   * Lists all stores configured for a dataset.
   * 
   * @param {apiDszListStoresReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszListStoresResp>>}
   */
  async dszListStores(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/store/list`,
      opts
    });
  }
  /**
   * Get Model Signed URLs
   * 
   * Generates presigned URLs for accessing model files in the specified dataset.
   * 
   * @param {apiDszModelGetReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszModelGetResp>>}
   */
  async dszModelGet(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/model/dataset/${encodeURIComponent(String(req.params.dataset))}/signed-get/${encodeURIComponent(String(req.params.model_name))}`,
      json: req.payload,
      opts
    });
  }
  /**
   * Load Last Batch Summary Result
   * 
   * Load a table file from the datalake as JSON and return it.
   * 
   * @param {apiDszResultLastBatchSummaryReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszResultLastBatchSummaryResp>>}
   */
  async dszResultLastBatchSummary(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/results/last-batch-summary`,
      query: req.query,
      opts
    });
  }
  /**
   * Load Relevance Fashion Catalog Result
   * 
   * Load a table file from the datalake as JSON and return it.
   * 
   * @param {apiDszResultRelevanceFashionCatalogReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszResultRelevanceFashionCatalogResp>>}
   */
  async dszResultRelevanceFashionCatalog(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/results/relevance-fashion-catalog`,
      query: req.query,
      opts
    });
  }
  /**
   * Generate Signed URLs
   * 
   * Generates pre-signed storage URLs for dataset lake table paths.
   * 
   * @param {apiDszSignedUrlReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszSignedUrlResp>>}
   */
  async dszSignedUrl(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/signed-url`,
      json: req.payload,
      opts
    });
  }
  /**
   * Create Store Display
   * 
   * Creates or updates a store display configuration.
   * 
   * @param {apiDszStoreDisplayCreateReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszStoreDisplayCreateResp>>}
   */
  async dszStoreDisplayCreate(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/store/display/create`,
      json: req.payload,
      opts
    });
  }
  /**
   * List Store Displays
   * 
   * Lists displays configured for a store in a dataset.
   * 
   * @param {apiDszStoreDisplayListReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszStoreDisplayListResp>>}
   */
  async dszStoreDisplayList(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/store/display/list`,
      query: req.query,
      opts
    });
  }
  /**
   * Create Stripe Checkout
   * 
   * Creates a Stripe checkout session for a subscription and returns its ID and URL.
   * 
   * @param {apiDszStripeCheckoutReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszStripeCheckoutResp>>}
   */
  async dszStripeCheckout(req, opts) {
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/stripe/checkout`,
      json: req.payload,
      opts
    });
  }
  /**
   * Get Dataset Table File
   * 
   * Get table file by redirecting to its signed URL.
   * 
   * @param {apiDszTableGetReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszTableGetResp>>}
   */
  async dszTableGet(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/table-file/{path...}`,
      opts
    });
  }
  /**
   * List Dataset Tables
   * 
   * Lists files and tables in the specified dataset.
   * 
   * @param {apiDszTableListReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszTableListResp>>}
   */
  async dszTableList(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/table-list`,
      query: req.query,
      opts
    });
  }
  /**
   * Update Payment Auto-Renew
   * 
   * Updates the auto-renewal status for a dataset subscription solution.
   * 
   * @param {apiDszUpdatePaymentAutoRenewReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszUpdatePaymentAutoRenewResp>>}
   */
  async dszUpdatePaymentAutoRenew(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/solution/${encodeURIComponent(String(req.params.solution))}/payment/autorenew`,
      json: req.payload,
      opts
    });
  }
  /**
   * Update Subscription Plan
   * 
   * Updates the plan (price) of the Stripe subscription for a dataset product/solution.
   * 
   * @param {apiDszUpdateSubscriptionPlanReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiDszUpdateSubscriptionPlanResp>>}
   */
  async dszUpdateSubscriptionPlan(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ds/dataset/${encodeURIComponent(String(req.params.dataset))}/solution/${encodeURIComponent(String(req.params.solution))}/subscription/plan`,
      json: req.payload,
      opts
    });
  }
  /**
   * Ingest Relevance Fashion Catalog Event
   * 
   * Ingests catalog item data into the dataset catalog table.
   * 
   * @param {apiIngestzEventRelevanceFashionCatalogReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiIngestzEventRelevanceFashionCatalogResp>>}
   */
  async ingestzEventRelevanceFashionCatalog(req, opts) {
    opts = opts || {};
    opts.labels = {
      "sinkable": "true",
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ingest/dataset/${encodeURIComponent(String(req.params.dataset))}/event/relevance-fashion-catalog/text.txt`,
      json: req.payload,
      opts,
      sinkable: true
    });
  }
  /**
   * Ingest Relevance Fashion Event
   * 
   * Ingests telemetry event data into the dataset event table.
   * 
   * @param {apiIngestzEventRelevanceFashionEventReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiIngestzEventRelevanceFashionEventResp>>}
   */
  async ingestzEventRelevanceFashionEvent(req, opts) {
    opts = opts || {};
    opts.labels = {
      "sinkable": "true",
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ingest/dataset/${encodeURIComponent(String(req.params.dataset))}/event/relevance-fashion-event/text.txt`,
      json: req.payload,
      opts,
      sinkable: true
    });
  }
  /**
   * Get Ingestion Upload URLs
   * 
   * Generates presigned PUT URLs and an ingestion token for uploading files into a target dataset table and partitions.
   * 
   * @param {apiIngestzGetUrlReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiIngestzGetUrlResp>>}
   */
  async ingestzGetUrl(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ingest/dataset/${encodeURIComponent(String(req.params.dataset))}/get-url`,
      json: req.payload,
      opts
    });
  }
  /**
   * Get Last Successful Ingestion Data Summary
   * 
   * Returns the date of ingestion, rows count and distinct and null columns count.
   * 
   * @param {apiIngestzLastBatchSummaryReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiIngestzLastBatchSummaryResp>>}
   */
  async ingestzLastBatchSummary(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ingest/dataset/${encodeURIComponent(String(req.params.dataset))}/last-batch`,
      query: req.query,
      opts
    });
  }
  /**
   * Ingest Relevance Fashion Catalog Stream
   * 
   * Ingests a batch of catalog items into the dataset catalog table.
   * 
   * @param {apiIngestzStreamRelevanceFashionCatalogReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiIngestzStreamRelevanceFashionCatalogResp>>}
   */
  async ingestzStreamRelevanceFashionCatalog(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ingest/dataset/${encodeURIComponent(String(req.params.dataset))}/stream/relevance-fashion-catalog`,
      json: req.payload,
      opts
    });
  }
  /**
   * Ingest Relevance Fashion Event Stream
   * 
   * Ingests a batch of telemetry events into the dataset event table.
   * 
   * @param {apiIngestzStreamRelevanceFashionEventReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiIngestzStreamRelevanceFashionEventResp>>}
   */
  async ingestzStreamRelevanceFashionEvent(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/ingest/dataset/${encodeURIComponent(String(req.params.dataset))}/stream/relevance-fashion-event`,
      json: req.payload,
      opts
    });
  }
  /**
   * Get Pipeline
   * 
   * Gets a pipeline of a dataset by its id.
   * 
   * @param {apiPipezGetReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezGetResp>>}
   */
  async pipezGet(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/id/${encodeURIComponent(String(req.params.id))}`,
      opts
    });
  }
  /**
   * List Pipelines
   * 
   * Lists the pipelines of a dataset, optionally filtered, with cursor pagination.
   * 
   * @param {apiPipezListReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezListResp>>}
   */
  async pipezList(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/list`,
      query: req.query,
      opts
    });
  }
  /**
   * Activate Relevance Fashion Batch Ingestion Pipeline Scheduler
   * 
   * Activates a scheduler for the relevance fashion batch ingestion pipeline on a dataset.
   * 
   * @param {apiPipezSchedulerActivateRelevanceFashionRelevanceFashionBatchIngestReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezSchedulerActivateRelevanceFashionRelevanceFashionBatchIngestResp>>}
   */
  async pipezSchedulerActivateRelevanceFashionRelevanceFashionBatchIngest(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/scheduler/activate/relevance-fashion/relevance-fashion-batch-ingest`,
      json: req.payload,
      opts
    });
  }
  /**
   * Activate Relevance Fashion Main Pipeline Scheduler
   * 
   * Activates a scheduler for the relevance fashion main pipeline on a dataset.
   * 
   * @param {apiPipezSchedulerActivateRelevanceFashionRelevanceFashionMainReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezSchedulerActivateRelevanceFashionRelevanceFashionMainResp>>}
   */
  async pipezSchedulerActivateRelevanceFashionRelevanceFashionMain(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/scheduler/activate/relevance-fashion/relevance-fashion-main`,
      opts
    });
  }
  /**
   * Deactivate Relevance Fashion Batch Ingestion Pipeline Scheduler
   * 
   * Deactivates a scheduler for the relevance fashion batch ingestion pipeline on a dataset.
   * 
   * @param {apiPipezSchedulerDeactivateRelevanceFashionRelevanceFashionBatchIngestReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezSchedulerDeactivateRelevanceFashionRelevanceFashionBatchIngestResp>>}
   */
  async pipezSchedulerDeactivateRelevanceFashionRelevanceFashionBatchIngest(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/scheduler/deactivate/relevance-fashion/relevance-fashion-batch-ingest`,
      json: req.payload,
      opts
    });
  }
  /**
   * Deactivate Relevance Fashion Main Pipeline Scheduler
   * 
   * Deactivates a scheduler for the relevance fashion main pipeline on a dataset.
   * 
   * @param {apiPipezSchedulerDeactivateRelevanceFashionRelevanceFashionMainReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezSchedulerDeactivateRelevanceFashionRelevanceFashionMainResp>>}
   */
  async pipezSchedulerDeactivateRelevanceFashionRelevanceFashionMain(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/scheduler/deactivate/relevance-fashion/relevance-fashion-main`,
      opts
    });
  }
  /**
   * Start Batch Ingestion Pipeline
   * 
   * Starts a batch ingestion pipeline on a dataset with the given params and returns its id.
   * 
   * @param {apiPipezStartRelevanceFashionRelevanceFashionBatchIngestReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezStartRelevanceFashionRelevanceFashionBatchIngestResp>>}
   */
  async pipezStartRelevanceFashionRelevanceFashionBatchIngest(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/start/relevance_fashion/relevance-fashion-batch-ingest`,
      json: req.payload,
      opts
    });
  }
  /**
   * Start Main Pipeline
   * 
   * Starts the main relevance fashion pipeline on a dataset and returns its id.
   * 
   * @param {apiPipezStartRelevanceFashionRelevanceFashionMainReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiPipezStartRelevanceFashionRelevanceFashionMainResp>>}
   */
  async pipezStartRelevanceFashionRelevanceFashionMain(req, opts) {
    opts = opts || {};
    opts.labels = {
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/pipe/dataset/${encodeURIComponent(String(req.params.dataset))}/start/relevance_fashion/relevance-fashion-main`,
      opts
    });
  }
  /**
   * Sample Fraction
   * 
   * Sample API that performs a fraction operation
   * 
   * @param {apiSampleFractionReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiSampleFractionResp>>}
   */
  async sampleFraction(req, opts) {
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/gox/routez/sample/fraction/${encodeURIComponent(String(req.params.numerator))}/${encodeURIComponent(String(req.params.denominator))}`,
      query: req.query,
      headers: req.headers,
      json: req.payload,
      opts
    });
  }
  /**
   * Get Model Metadata
   * 
   * Gets the input and output tensors of a model version.
   * 
   * @param {apiTfzModelMetadataReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiTfzModelMetadataResp>>}
   */
  async tfzModelMetadata(req, opts) {
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["decsuite"] || ""}/api/metadata/${encodeURIComponent(String(req.params.model_name))}/${encodeURIComponent(String(req.params.version))}`,
      opts
    });
  }
  /**
   * Model Predict
   * 
   * Runs a prediction of a dataset model version on the inference service (tfinfer), sending the model inputs and the model partitions and returning the model outputs.
   * 
   * @param {apiTfzModelPredictReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiTfzModelPredictResp>>}
   */
  async tfzModelPredict(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "tfinfer",
      "token_types": "dataset",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["tfinfer"] || ""}/api/model/dataset/${encodeURIComponent(String(req.params.dataset))}/${encodeURIComponent(String(req.params.model_name))}/${encodeURIComponent(String(req.params.version))}/predict`,
      json: req.payload,
      opts
    });
  }
  /**
   * Device Authorization Code Initiation
   * 
   * Initiates an RFC 8628 device authorization flow for CLI and remote headless SSH sessions.
   * 
   * @param {apiUserAuthDeviceCodeReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthDeviceCodeResp>>}
   */
  async userAuthDeviceCode(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/device/code`,
      json: req.payload,
      opts
    });
  }
  /**
   * Device Token Polling & Exchange
   * 
   * Polled by CLI clients to check status and exchange an authorized device code for active user tokens.
   * 
   * @param {apiUserAuthDeviceTokenReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthDeviceTokenResp>>}
   */
  async userAuthDeviceToken(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/device/token`,
      json: req.payload,
      opts
    });
  }
  /**
   * OAuth 2.0 PKCE Code Exchange
   * 
   * Exchanges an authorization code and PKCE code_verifier for authenticated user tokens.
   * 
   * @param {apiUserAuthExchangeReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthExchangeResp>>}
   */
  async userAuthExchange(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/exchange`,
      json: req.payload,
      opts
    });
  }
  /**
   * Google One-Tap / ID Token Sign-In
   * 
   * Authenticates a user via verified Google ID token / GIS One-Tap credentials.
   * 
   * @param {apiUserAuthGoogleReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthGoogleResp>>}
   */
  async userAuthGoogle(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/google`,
      json: req.payload,
      opts
    });
  }
  /**
   * User Logout
   * 
   * Revokes an active human user refresh token and terminates session.
   * 
   * @param {apiUserAuthLogoutReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthLogoutResp>>}
   */
  async userAuthLogout(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/logout`,
      json: req.payload,
      opts
    });
  }
  /**
   * Current User Profile
   * 
   * Retrieves current authenticated user identity profile metadata.
   * 
   * @param {apiUserAuthMeReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthMeResp>>}
   */
  async userAuthMe(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/me`,
      headers: req.headers,
      opts
    });
  }
  /**
   * Password Reset Initiation
   * 
   * Initiates an out-of-band password reset email delivery.
   * 
   * @param {apiUserAuthPasswordResetReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthPasswordResetResp>>}
   */
  async userAuthPasswordReset(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/password/reset`,
      json: req.payload,
      opts
    });
  }
  /**
   * Email/Password Sign-In
   * 
   * Authenticates a user via email and password credentials.
   * 
   * @param {apiUserAuthPasswordSigninReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthPasswordSigninResp>>}
   */
  async userAuthPasswordSignin(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/signin/password`,
      json: req.payload,
      opts
    });
  }
  /**
   * Email/Password Registration
   * 
   * Registers a new user account with email and password credentials.
   * 
   * @param {apiUserAuthPasswordSignupReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthPasswordSignupResp>>}
   */
  async userAuthPasswordSignup(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/auth/signup/password`,
      json: req.payload,
      opts
    });
  }
  /**
   * Human User Session Refresh
   * 
   * Exchanges a human user refresh token for fresh access and rotated refresh tokens.
   * 
   * @param {apiUserAuthRefreshReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserAuthRefreshResp>>}
   */
  async userAuthRefresh(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/refresh-user`,
      json: req.payload,
      opts
    });
  }
  /**
   * Public Settings Discovery
   * 
   * Retrieves public operational policies and token expiration parameters.
   * 
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserConfSettingsResp>>}
   */
  async userConfSettings(opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/conf/settings.json`,
      opts
    });
  }
  /**
   * Attenuate Dataset Access Token
   * 
   * Derives an in-flight attenuated child access token with dropped privileges from an active access token.
   * 
   * @param {apiUserDatasetTokenAccessAttenuateReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserDatasetTokenAccessAttenuateResp>>}
   */
  async userDatasetTokenAccessAttenuate(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/attenuate`,
      headers: req.headers,
      json: req.payload,
      opts
    });
  }
  /**
   * Create Dataset Access Token
   * 
   * Mints a short-lived dataset access token for an authenticated user.
   * 
   * @param {apiUserDatasetTokenAccessCreateReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserDatasetTokenAccessCreateResp>>}
   */
  async userDatasetTokenAccessCreate(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/dataset-access-token/${encodeURIComponent(String(req.params.dataset))}`,
      headers: req.headers,
      json: req.payload,
      opts
    });
  }
  /**
   * Refresh Dataset Access Token
   * 
   * Refreshes and attenuates a dataset access token using an Ed25519 signature assertion.
   * 
   * @param {apiUserDatasetTokenRefreshReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserDatasetTokenRefreshResp>>}
   */
  async userDatasetTokenRefresh(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/refresh-dataset`,
      headers: req.headers,
      json: req.payload,
      opts
    });
  }
  /**
   * Provision Dataset Refresh Token
   * 
   * Provisions a long-lived machine dataset refresh token.
   * 
   * @param {apiUserDatasetTokenRefreshCreateReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserDatasetTokenRefreshCreateResp>>}
   */
  async userDatasetTokenRefreshCreate(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/dataset-refresh-token/${encodeURIComponent(String(req.params.dataset))}`,
      headers: req.headers,
      json: req.payload,
      opts
    });
  }
  /**
   * Revoke Dataset Refresh Token
   * 
   * Instantly revokes a machine dataset refresh token by JTI.
   * 
   * @param {apiUserDatasetTokenRefreshDeleteReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserDatasetTokenRefreshDeleteResp>>}
   */
  async userDatasetTokenRefreshDelete(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "DELETE",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/dataset-refresh-token/${encodeURIComponent(String(req.params.dataset))}/jti/${encodeURIComponent(String(req.params.jti))}`,
      headers: req.headers,
      opts
    });
  }
  /**
   * List Dataset Refresh Tokens
   * 
   * Lists all active machine refresh tokens provisioned for a dataset.
   * 
   * @param {apiUserDatasetTokenRefreshListReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserDatasetTokenRefreshListResp>>}
   */
  async userDatasetTokenRefreshList(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "GET",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/dataset-refresh-token/${encodeURIComponent(String(req.params.dataset))}`,
      headers: req.headers,
      opts
    });
  }
  /**
   * Universal Token & Signature Diagnostics
   * 
   * Auto-detects, validates, and diagnoses any credential supported by authz.Filter without state side-effects.
   * 
   * @param {apiUserTokenCheckReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiUserTokenCheckResp>>}
   */
  async userTokenCheck(req, opts) {
    opts = opts || {};
    opts.labels = {
      "servname": "userz",
      ...opts.labels || {}
    };
    return this._services.invoker().invokeResult({
      method: "POST",
      url: `${this._services.config().servUrls?.["userz"] || ""}/api/userz/token/check`,
      headers: req.headers,
      json: req.payload,
      opts
    });
  }
}
export { Apis };
class SdkFuncs {
  /**
   * Constructs an SdkFuncs instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    this._services = services;
  }
  /**
   * Get Attached Revenue
   * 
   * Gets the revenue attached to recommendations for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashAttachedRevenueReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashAttachedRevenueOutput>}
   */
  async productRelevanceFashionDashAttachedRevenueV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-attached-revenue",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * Get Coverage
   * 
   * Gets the recommendation coverage of the catalog for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashCoverageReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashCoverageOutput>}
   */
  async productRelevanceFashionDashCoverageV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-coverage",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * Get Display Metrics Status
   * 
   * Gets recommendation metrics, revenue, session counts, and uplift status by display for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashDisplayMetricsStatusReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashDisplayMetricsStatusOutput>}
   */
  async productRelevanceFashionDashDisplayMetricsStatusV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-display-metrics-status",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * Get Funnel
   * 
   * Gets the event counts of the recommendation funnel for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashFunnelReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashFunnelOutput>}
   */
  async productRelevanceFashionDashFunnelV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-funnel",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * Get Metrics by Category
   * 
   * Gets the metrics impacted by recommendations by showcase category for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryOutput>}
   */
  async productRelevanceFashionDashImpactedMetricsByShowcaseCategoryV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-impacted-metrics-by-showcase-category",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * Get Metrics by SKU
   * 
   * Gets the metrics impacted by recommendations by SKU for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashImpactedMetricsBySkuReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashImpactedMetricsBySkuOutput>}
   */
  async productRelevanceFashionDashImpactedMetricsBySkuV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-impacted-metrics-by-sku",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * List Showcase Categories
   * 
   * Lists the showcase categories of a store.
   * 
   * @param {sdkFuncProductRelevanceFashionDashListShowcaseCategoriesReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashListShowcaseCategoriesOutput>}
   */
  async productRelevanceFashionDashListShowcaseCategoriesV1(req, opts) {
    const dash = this._services.dash();
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-list-showcase-categories",
      partitions: { "version": "v1", ...req.partitions },
      serviceOptions: opts
    });
    return ret;
  }
  /**
   * Get Bought Together
   * 
   * Infers the SKUs frequently bought together with the given SKUs for a store.
   * 
   * @param {sdkFuncProductRelevanceFashionFbtReq} req
   * @param {(req: sdkFuncProductRelevanceFashionFbtReq, results: sdkFuncProductRelevanceFashionFbtOutput) => Promise<sdkFuncProductRelevanceFashionFbtSelected>} [selectResults] - picks which of the recommended skus are shown and may attach your own data to each in enriched; default: for each result, in order, the first sku not in req.input.skus (the cart) and not already picked, with the rating of that result; no duplicates, and results with no such sku are skipped
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionFbtSelected>}
   */
  async productRelevanceFashionFbtV1(req, selectResults, opts) {
    const infer = this._services.infer();
    const results = await infer.run({
      dataset: req.input.dataset,
      modelName: "relevance-fashion-fbt",
      params: req.input,
      partitions: { "version": "v1", ...req.partitions },
      sessionId: req.input.session_id,
      userId: req.input.user_id,
      serviceOptions: opts
    });
    return infer.selectResults({
      dataset: req.input.dataset,
      modelName: "relevance-fashion-fbt",
      modelVersion: "v1",
      storeId: req.partitions.store_id,
      req,
      results,
      recId: results.rec_id,
      sessionId: req.input.session_id,
      userId: req.input.user_id,
      selectResults,
      serviceOptions: opts
    });
  }
}
export { SdkFuncs };
