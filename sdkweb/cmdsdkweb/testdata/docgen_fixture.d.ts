/**
 * Main SDK class providing access to i6 services.
 */
export class I6Sdk {
    /**
     * Returns an Infer service instance.
     * @returns {Infer} An instance of the Infer client.
     */
    infer(): Infer;
    /**
     * Returns the generated Apis client.
     * @returns {Apis} An instance of the Apis client.
     */
    apis(): Apis;
    /**
     * Returns a copy of the configuration.
     * @returns {SampleConfig} The current configuration.
     */
    config(): SampleConfig;
    /**
     * Internal helper that should be skipped.
     */
    private _internalHelper(): void;
}

/**
 * Service for running inference models.
 */
declare class Infer {
    /**
     * Runs the model in the browser and returns its recommendation result.
     * @param {SampleOptions} opts - Inference options.
     * @returns {Promise<SampleResult>} The results of the model.
     */
    run(opts: SampleOptions): Promise<SampleResult>;
    /**
     * Internal method that should be skipped.
     */
    private _skipMe(): void;
    #private;
}

/**
 * Generated API client.
 */
declare class Apis {
    /**
     * A generated API method.
     */
    dszList(): void;
}

/**
 * Configuration options.
 */
export type SampleConfig = {
    /**
     * - Application identifier.
     */
    app: string;
};

/**
 * Options for running inference.
 */
export type SampleOptions = {
    /**
     * - Model identifier name.
     */
    modelName: string;
    /**
     * - Optional query limit.
     */
    limit?: number;
};

/**
 * Result of inference run.
 */
export type SampleResult = {
    /**
     * - Recommendation identifier.
     */
    recId: string;
};
