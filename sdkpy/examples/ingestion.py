"""Example: Ingesting catalog data using an Ed25519 Dataset Refresh Token.

Target Table: prcbidding-margin-catalog (from lakezprcbiddingmain/prcbiddingmain_catalog.go)
Table Partitions: ingested_at, store_id
"""

from __future__ import annotations

import datetime
import io
import json
import os
import sys

def main() -> None:
    # 1. Read DATASET_REFRESH_TOKEN from environment
    raw_token_json = os.environ.get("DATASET_REFRESH_TOKEN")
    if not raw_token_json:
        print(
            "Error: DATASET_REFRESH_TOKEN environment variable is required.\n"
            "Please source secrets/envs.secret.sh or export DATASET_REFRESH_TOKEN.",
            file=sys.stderr,
        )
        sys.exit(1)
    try:
        token_data = json.loads(raw_token_json)
    except json.JSONDecodeError as err:
        print(f"Error parsing DATASET_REFRESH_TOKEN JSON: {err}", file=sys.stderr)
        sys.exit(1)

    dataset = token_data.get("dataset", "dectest")
    print(f"Initializing Infinity6 SDK for dataset: {dataset} (JTI: {token_data.get('jti')})")

    # 2. Initialize the SDK client pointing to the STAGING environment
    from i6sdk import STAGING_SERV_URLS, create_client

    print("Configuring SDK with STAGING service endpoints:")
    for svc, endpoint in sorted(STAGING_SERV_URLS.items()):
        print(f"  • {svc:10s}: {endpoint}")

    client = create_client(serv_urls=STAGING_SERV_URLS)

    # 3. Configure the machine robot's dataset refresh token
    # The SDK automatically uses this to sign Ed25519 RefreshAssertions and mint
    # unbounded dataset access tokens on demand without requiring a user session.
    client.tokens().set_dataset_refresh_token(dataset, token_data)

    # 4. Target Table & Partitions from lakezprcbiddingmargin/prcbiddingmargin_catalog.go:
    # Table: "prcbidding-margin-catalog"
    # Table definition partitions: ["ingested_at", "store_id"]
    # Note: 'ingested_at' is automatically assigned by the server from the generated IngestId.
    # The client specifies the operational partitions:
    table_name = "prcbidding-margin-catalog"
    store_id = "1"

    partitions = {
        "store_id": store_id,
    }

    now_iso = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT00:00:00Z")

    # 5. Read sample catalog data from JSON Lines (NDJSON) file
    import pathlib
    data_path = pathlib.Path(__file__).resolve().parent / "ingestion.json"
    if not data_path.is_file():
        print(f"Error: Data file not found: {data_path}", file=sys.stderr)
        sys.exit(1)

    with open(data_path, "rb") as f:
        ndjson_data = f.read()

    # 6. Wrap in a BinaryIO stream for uploading
    stream = io.BytesIO(ndjson_data)
    record_count = len([line for line in ndjson_data.split(b"\n") if line.strip()])

    print(f"Read {record_count} catalog records ({len(ndjson_data)} bytes) from {data_path.name}.")
    print(f"Target: Table='{table_name}', Partitions={partitions}")

    from i6sdk.ingest import UploadProgress

    # 7. Progress callback
    def on_progress(p: UploadProgress) -> None:
        percent = p.get("percent", 0.0)
        loaded = p.get("loaded", 0)
        total = p.get("total", 0)
        completed = p.get("completed_streams", 0)
        total_streams = p.get("total_streams", 0)
        print(
            f"\rUpload Progress: {percent:5.1f}% | "
            f"Loaded: {loaded:,} / {total:,} bytes | "
            f"Completed streams: {completed}/{total_streams}",
            end="",
            flush=True,
        )

    # 8. Upload the stream concurrently via signed PUT URL
    try:
        result = client.ingest().upload(
            dataset=dataset,
            table=table_name,
            partitions=partitions,
            streams=[stream],
            content_type="application/x-ndjson",
            on_progress=on_progress,
        )
        print(f"\nUpload complete!")
        print(f"Ingestion ID : {result['ingest_id']}")
        print(f"Stream Count : {result['count']}")
    except Exception as err:
        print(f"\nUpload failed: {err}", file=sys.stderr)
        raise


if __name__ == "__main__":
    main()
