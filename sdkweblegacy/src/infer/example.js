const result = await new I6Sdk().infer().run({
    modelName: "relevance-fashion-fbt",
    modelVersion: "v1",
    storeId: "s1",
    dataset: "dectest",
    params: { "sku_ids": ["9132438814883","9132177490083","unknown-id"] },
    sqliteUrl: "i6tb://results.relevance-fashion-fbt-sqlite/store_id=s1/db.sqlite",
    tfliteUrl: "http://localhost:8000/src/infer/gen/model.tflite"
})