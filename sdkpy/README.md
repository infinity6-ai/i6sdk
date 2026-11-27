# i6sdk - Official Infinity6 Python SDK

Lightweight, high-performance Python SDK for the Infinity6 platform.

## Features

- **Minimal Runtime Dependencies**: Built for Python $\ge 3.10$ with minimal external dependencies (`cryptography` for Ed25519 request signing).
- **Strictly Typed**: Includes PEP 561 `py.typed` marker and TypedDict schemas for full IDE autocompletion and type checking (`mypy`, `pyright`).
- **Complete Platform Coverage**: Authentication, token management, data ingestion, ML model inference, and analytics dashboards.
- **Pluggable Transport**: Standard `urllib.request` out of the box, with support for custom fetcher adapters (`httpx`, `aiohttp`, `requests`).

## Installation

### Pre-built Wheels (GitHub Releases)

Pre-built wheels are published automatically on every release and edge build:

```bash
# Specific version tag (e.g. v0.0.8)
pip install "https://github.com/infinity6-ai/i6sdk/releases/download/v0.0.8/i6sdk-0.0.8-py3-none-any.whl"

# Latest official release
pip install --upgrade "https://github.com/infinity6-ai/i6sdk/releases/latest/download/i6sdk-latest-py3-none-any.whl"

# Edge development build (continuous from 'edge' branch)
pip install --upgrade "https://github.com/infinity6-ai/i6sdk/releases/download/edge/i6sdk-edge-py3-none-any.whl"
```

To pin in `requirements.txt`:
```text
i6sdk @ https://github.com/infinity6-ai/i6sdk/releases/latest/download/i6sdk-latest-py3-none-any.whl
```

### From Git Repository

```bash
# Specific version tag or branch
pip install "git+https://github.com/infinity6-ai/i6sdk.git@v0.0.8#subdirectory=sdkpy"
pip install --upgrade "git+https://github.com/infinity6-ai/i6sdk.git@edge#subdirectory=sdkpy"
```

### Local Development

Inside the cloned mono-repo:
```bash
pip install -e sdkpy
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
- `client.funcs()`: High-level recommendation and dashboard workflows.
- `client.auth()`: PKCE OAuth, device code authentication, and token management.
- `client.tokens()`: Token cache and scoped attenuation.
- `client.infer()`: Offline SQLite query execution + server-side ML model predictions.
- `client.ingest()`: Batch event streaming and data ingestion.
- `client.dash()`: Analytics query execution.
- `bitmapz`: Compact variable-length bitmask operations, serialization, and deserialization.
- `Grants`: Permission bitmasks (`System`, `Admin`, `Write`, `Read`, `All`).
