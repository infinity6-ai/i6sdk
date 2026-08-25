class RaasBoughtTogether {

  async init(config) {
    if (!config || !config.subscription || !config.store) {
      throw new Error('RaasBoughtTogether configuration is missing.');
    }
    this.config = config
  }

  async predict(cartItems, opts) {
    opts = {
      maxResults: 3, 
      deterministic: false,
      ...opts,
    }

    if (cartItems.length === 5) {
      throw 'mock error'
    }

    const predictions = cartItems.map((item, idx) => ({
      products: [{id: item.id}],
      score: Math.random() + (opts.deterministic ? idx : 0)
    }));
    const sorted = predictions.sort((a, b) => b.score - a.score);
    return sorted.slice(0, opts.maxResults);
  }
}

export { RaasBoughtTogether }