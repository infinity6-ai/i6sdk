# i6sdk - Official Infinity6 Python SDK

Lightweight, high-performance, and **zero-dependency** Python SDK for the Infinity6 platform.

## Features

- **Zero Runtime Dependencies**: Relies exclusively on the Python Standard Library (Python $\ge 3.10$).
- **Strictly Typed**: Includes PEP 561 `py.typed` marker and TypedDict schemas for full IDE autocompletion and type checking (`mypy`, `pyright`).
- **Complete Platform Coverage**: Authentication, token management, data ingestion, ML model inference, and analytics dashboards.
- **Pluggable Transport**: Standard `urllib.request` out of the box, with support for custom fetcher adapters (`httpx`, `aiohttp`, `requests`).

## Installation

```bash
pip install .
```

## Quick Start

```python
import i6sdk

# Initialize client using environment variables (I6_SDK_DECSUITE_URL)
client = i6sdk.create_client()

# Or pass configuration explicitly
config = i6sdk.I6SdkConfig(
    serv_urls={"decsuite": "https://decsuite.sandbox.rs.infinity6.ai"},
)
client = i6sdk.I6Client(config)

# Call API endpoints
resp = client.apis().dsz_list_datasets()
if resp.ok:
    data = resp.parse()
    print("Datasets:", data)
else:
    print("Error:", resp.error)
```

## Architecture

The Python SDK mirrors the architecture of the Infinity6 TypeScript/Web SDK:

- `client.apis()`: Auto-generated, type-safe API methods for every platform service.
- `client.sdk_funcs()`: High-level recommendation and dashboard workflows.
- `client.auth()`: PKCE OAuth, device code authentication, and token management.
- `client.tokens()`: Token cache and scoped attenuation.
- `client.infer()`: Offline SQLite query execution + server-side ML model predictions.
- `client.ingest()`: Batch event streaming and data ingestion.
- `client.dash()`: Analytics query execution.
