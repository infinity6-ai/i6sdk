"""Example: Uploading multiple large files (10 x 100MB) using the Infinity6 Ingestion API.

This script demonstrates streaming multi-file / multi-part data batches to the
Infinity6 ingestion service in parallel.

Target Table: prcbidding-margin-catalog (from lakezprcbiddingmain/prcbiddingmain_catalog.go)
Table Partitions: ingested_at, store_id
Files: 10 large NDJSON files (~100MB each) located in the repository's `gen/` directory.
"""

from __future__ import annotations

import argparse
import datetime
import json
import os
import pathlib
import sys
import time
from typing import List

# Ensure local i6sdk package is importable when running example directly
_EXAMPLES_DIR = pathlib.Path(__file__).resolve().parent
_REPO_ROOT = _EXAMPLES_DIR.parent.parent
_SRC_DIR = _EXAMPLES_DIR.parent / "src"
if str(_SRC_DIR) not in sys.path:
    sys.path.insert(0, str(_SRC_DIR))

# Default target directory for generated files
DEFAULT_GEN_DIR = _REPO_ROOT / "gen"
TARGET_FILE_COUNT = 10
TARGET_FILE_SIZE_BYTES = 100 * 1024 * 1024  # 100 MiB per file (~1 GB total across 10 files)


def generate_single_file(filepath: pathlib.Path, file_index: int, target_size: int) -> int:
    """Generates a single NDJSON file filled with valid catalog records up to target_size."""
    template = (
        '{"id":"SKU-%02d-%07d","store_id":"1",'
        '"desc":"Performance Catalog Item %02d-%07d",'
        '"supplier":"Acme Global Suppliers Ltd","status":"ACTIVE","exclusive_sale":"NO",'
        '"hierarchy":[{"name":"department","value":"Apparel"},'
        '{"name":"category","value":"Outerwear"},{"name":"subcategory","value":"Jackets"}]}\n'
    )
    batch_size = 10000
    written = 0
    item_idx = 0

    with open(filepath, "wb") as f:
        while written < target_size:
            chunk = "".join(
                template % (file_index, item_idx + i, file_index, item_idx + i)
                for i in range(batch_size)
            ).encode("utf-8")
            f.write(chunk)
            written += len(chunk)
            item_idx += batch_size

    return written


def ensure_large_files(
    gen_dir: pathlib.Path,
    count: int = TARGET_FILE_COUNT,
    target_size_bytes: int = TARGET_FILE_SIZE_BYTES,
) -> List[pathlib.Path]:
    """Ensures that the required large NDJSON files exist in gen_dir, generating any that are missing."""
    gen_dir.mkdir(parents=True, exist_ok=True)
    file_paths: List[pathlib.Path] = []

    print(f"Checking for {count} large files (target: {target_size_bytes / (1024 * 1024):.0f} MB each) in {gen_dir}...")
    for i in range(1, count + 1):
        filename = f"ingestion_large_{i:02d}.json"
        path = gen_dir / filename
        file_paths.append(path)

        if path.is_file() and path.stat().st_size >= target_size_bytes:
            current_mb = path.stat().st_size / (1024 * 1024)
            print(f"  [{i:02d}/{count:02d}] {filename} exists ({current_mb:.1f} MB) - OK")
        else:
            print(
                f"  [{i:02d}/{count:02d}] Generating {filename} (~{target_size_bytes / (1024 * 1024):.0f} MB)...",
                end="",
                flush=True,
            )
            t0 = time.time()
            size = generate_single_file(path, file_index=i, target_size=target_size_bytes)
            elapsed = time.time() - t0
            print(f" done ({size / (1024 * 1024):.1f} MB in {elapsed:.2f}s)")

    total_bytes = sum(p.stat().st_size for p in file_paths)
    print(f"All {count} files ready. Total size: {total_bytes / (1024 * 1024):,.1f} MB ({total_bytes:,} bytes).\n")
    return file_paths


def main() -> None:
    parser = argparse.ArgumentParser(description="Upload 10 large NDJSON files to Infinity6 Ingestion API")
    parser.add_argument(
        "--generate-only",
        action="store_true",
        help="Only generate the 10 large files in gen/ without uploading",
    )
    parser.add_argument(
        "--gen-dir",
        type=pathlib.Path,
        default=DEFAULT_GEN_DIR,
        help=f"Directory to store generated files (default: {DEFAULT_GEN_DIR})",
    )
    args = parser.parse_args()

    # 1. Ensure all 10 large files exist in gen/
    files = ensure_large_files(args.gen_dir, count=TARGET_FILE_COUNT, target_size_bytes=TARGET_FILE_SIZE_BYTES)
    if args.generate_only:
        print("Generation complete (--generate-only requested). Exiting.")
        return

    # 2. Read DATASET_REFRESH_TOKEN from environment
    raw_token_json = os.environ.get("DATASET_REFRESH_TOKEN")
    if not raw_token_json:
        print(
            "Error: DATASET_REFRESH_TOKEN environment variable is required to perform the upload.\n"
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

    # 3. Initialize the SDK client pointing to the STAGING environment
    from i6sdk import STAGING_SERV_URLS, create_client
    from i6sdk.ingest import UploadProgress

    print("Configuring SDK with STAGING service endpoints:")
    for svc, endpoint in sorted(STAGING_SERV_URLS.items()):
        print(f"  • {svc:10s}: {endpoint}")

    client = create_client(serv_urls=STAGING_SERV_URLS)

    # 4. Configure the machine robot's dataset refresh token
    # The SDK automatically uses this to sign Ed25519 RefreshAssertions and mint
    # unbounded dataset access tokens on demand without requiring a user session.
    client.tokens().set_dataset_refresh_token(dataset, token_data)

    # 5. Target Table & Partitions from lakezprcbiddingmargin/prcbiddingmargin_catalog.go:
    # Table: "prcbidding-margin-catalog"
    # Table definition partitions: ["ingested_at", "store_id"]
    table_name = "prcbidding-margin-catalog"
    store_id = "1"
    partitions = {
        "store_id": store_id,
    }

    total_bytes = sum(f.stat().st_size for f in files)
    print(f"Target: Table='{table_name}', Partitions={partitions}")
    print(f"Preparing to upload {len(files)} concurrent streams ({total_bytes / (1024 * 1024):,.1f} MB total)...")

    # 6. Progress reporting callback
    start_time = time.time()

    def on_progress(p: UploadProgress) -> None:
        percent = p.get("percent", 0.0)
        loaded = p.get("loaded", 0)
        total = p.get("total", 0)
        completed = p.get("completed_streams", 0)
        total_streams = p.get("total_streams", 0)

        elapsed = time.time() - start_time
        speed_mb = (loaded / (1024 * 1024)) / elapsed if elapsed > 0 else 0.0
        loaded_mb = loaded / (1024 * 1024)
        total_mb = total / (1024 * 1024)

        print(
            f"\rUpload Progress: {percent:5.1f}% | "
            f"{loaded_mb:,.1f} / {total_mb:,.1f} MB ({speed_mb:.1f} MB/s) | "
            f"Completed streams: {completed}/{total_streams}",
            end="",
            flush=True,
        )

    # 7. Open all 10 files and upload in parallel via signed PUT URLs
    file_handles = [open(f, "rb") for f in files]
    try:
        result = client.ingest().upload(
            dataset=dataset,
            table=table_name,
            partitions=partitions,
            streams=file_handles,
            content_type="application/x-ndjson",
            on_progress=on_progress,
            timeout=600.0,
        )
        total_elapsed = time.time() - start_time
        avg_speed = (total_bytes / (1024 * 1024)) / total_elapsed if total_elapsed > 0 else 0.0

        print(f"\n\nUpload complete successfully!")
        print(f"Ingestion ID  : {result['ingest_id']}")
        print(f"Stream Count  : {result['count']}")
        print(f"Total Bytes   : {total_bytes:,} ({total_bytes / (1024 * 1024):,.1f} MB)")
        print(f"Duration      : {total_elapsed:.2f}s ({avg_speed:.2f} MB/s)")
    except Exception as err:
        print(f"\nUpload failed: {err}", file=sys.stderr)
        raise
    finally:
        for fh in file_handles:
            try:
                fh.close()
            except Exception:
                pass


if __name__ == "__main__":
    main()
