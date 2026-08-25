class RaasBoughtTogether {

  async init(config) {
    if (!config || !config.subscription || !config.store) {
      throw new Error('RaasBoughtTogether configuration is missing.');
    }
    this.config = config
  }

  async predict(cartItems, maxResults = 3) {
    if (cartItems.length === 5) {
      throw 'mock error'
    }

    const predictions = cartItems.map(item => ({
      products: [{id: item.id}],
      score: Math.random()
    }));
    const sorted = predictions.sort((a, b) => b.score - a.score);
    return sorted.slice(0, maxResults);
  }
}

export { RaasBoughtTogether }