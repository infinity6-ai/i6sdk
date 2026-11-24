# i6sdk - Official Infinity6 Python SDK

Lightweight, high-performance, and **zero-dependency** Python SDK for the Infinity6 platform.

## Features

- **Zero Runtime Dependencies**: Relies exclusively on the Python Standard Library (Python $\ge 3.10$).
- **Strictly Typed**: Includes PEP 561 `py.typed` marker and TypedDict schemas for full IDE autocompletion and type checking (`mypy`, `pyright`).
- **Complete Platform Coverage**: Authentication, token management, data ingestion, ML model inference, and analytics dashboards.
- **Pluggable Transport**: Standard `urllib.request` out of the box, with support for custom fetcher adapters (`httpx`, `aiohttp`, `requests`).

## Installation

### From GitHub Releases (Pre-built Wheels)

Release builds automatically publish pre-built `.whl` and `.tar.gz` packages to [GitHub Releases](https://github.com/infinity6-ai/i6sdk/releases).

#### Specific Version (`vX.Y.Z`)

To install a specific release version (e.g. `v0.1.0`), specify the release tag and wheel filename (note: the wheel filename omits the `v` prefix):

```bash
pip install "https://github.com/infinity6-ai/i6sdk/releases/download/v0.1.0/i6sdk-0.1.0-py3-none-any.whl"
```

Or install the latest official release:

```bash
pip install --upgrade "https://github.com/infinity6-ai/i6sdk/releases/latest/download/i6sdk-latest-py3-none-any.whl"
```

#### Edge Development Build (`edge`)

To install the latest pre-release development build continuously built from the `edge` branch:

```bash
pip install --upgrade "https://github.com/infinity6-ai/i6sdk/releases/download/edge/i6sdk-edge-py3-none-any.whl"
```

#### In `requirements.txt`

You can pin release wheels directly using standard PEP 508 direct URL specifications:

```text
# Specific release version
i6sdk @ https://github.com/infinity6-ai/i6sdk/releases/download/v0.1.0/i6sdk-0.1.0-py3-none-any.whl

# Or edge development build
i6sdk @ https://github.com/infinity6-ai/i6sdk/releases/download/edge/i6sdk-edge-py3-none-any.whl
```

### Direct from Git Repository

You can also install directly from the source repository using `pip`:

```bash
# Specific version tag (e.g. v0.1.0)
pip install "git+https://github.com/infinity6-ai/i6sdk.git@v0.1.0#subdirectory=sdkpy"

# Edge branch
pip install --upgrade "git+https://github.com/infinity6-ai/i6sdk.git@edge#subdirectory=sdkpy"
```

### Local Development

For local development inside the cloned mono-repo:

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
