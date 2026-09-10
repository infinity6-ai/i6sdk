# i6sdk

The i6 SDK for the web, and the Go server that builds, serves, and tests it.

## Install

Requires Go and Node.js.

```bash
cd i6sdk-web
npm ci
```

## Run

```bash
cd i6sdk-web
./build.sh run_server
```

Starts the server at `127.0.0.1:8043` (override with the `I6_SDK_WEB_SERVER_LOCAL_ADDRESS`
env var). The SDK demo page is served at `/static/index.html`.

## Test

```bash
cd i6sdk-web
npm test
```

Runs the Go unit tests and the browser (QUnit) tests for the SDK.
