"""HTTP client invoker and response container using Python Standard Library."""

from __future__ import annotations

import json
import urllib.error
import urllib.parse
import urllib.request
from typing import (
    TYPE_CHECKING,
    Any,
    Callable,
    Dict,
    Generic,
    Optional,
    Tuple,
    TypeVar,
    TypedDict,
    Union,
)

if TYPE_CHECKING:
    from .services import Services

from .config import I6SdkConfig
from .exceptions import I6ApiError, I6ConnectionError

T = TypeVar("T")


def code_from_status(status: int) -> str:
    """Maps an HTTP status code to an error code matching sdkweb."""
    if status == 401:
        return "unauthenticated"
    if status == 403:
        return "forbidden"
    if status == 404:
        return "not_found"
    if status >= 500:
        return "unavailable"
    return "bad_request"


class ApiErrorDetail(TypedDict, total=False):
    """Structured error details when an API request fails."""

    status: int
    code: str
    message: str
    cause: Optional[str]


class ServiceOptions(TypedDict, total=False):
    """Optional per-request execution overrides."""

    timeout: Optional[float]
    headers: Optional[Dict[str, str]]
    labels: Optional[Dict[str, str]]
    service_name: Optional[str]


class InvokeOptions(TypedDict, total=False):
    """Options passed to the low-level HTTP invoker."""

    method: str
    url: Optional[str]
    path: Optional[str]
    query: Optional[Dict[str, Any]]
    headers: Optional[Dict[str, str]]
    json: Optional[Any]
    body: Optional[bytes]
    opts: Optional[Union[ServiceOptions, Dict[str, Any]]]


class Resp(Generic[T]):
    """Response container matching sdkweb's Resp class."""

    def __init__(
        self,
        status: int,
        headers: Dict[str, str],
        parser: Callable[[], T],
        text: str,
    ) -> None:
        self.status = status
        self.headers = headers
        self._parser = parser
        self.text = text
        self._value: Optional[T] = None
        self._has_value = False

    @property
    def ok(self) -> bool:
        """True if the HTTP status code is in the 2xx success range."""
        return 200 <= self.status < 300

    @property
    def error(self) -> Optional[ApiErrorDetail]:
        """Structured error details if request failed, None otherwise."""
        if self.ok:
            return None
        return {
            "status": self.status,
            "code": code_from_status(self.status),
            "message": self.text or f"HTTP {self.status}",
            "cause": self.text,
        }

    def parse(self) -> T:
        """Parses and returns the response payload. Raises I6ApiError on non-2xx status."""
        if not self._has_value:
            self._value = self._parser()
            self._has_value = True
        return self._value  # type: ignore[return-value]


class Invoker:
    """HTTP dispatch engine executing requests using urllib.request."""

    def __init__(self, services: Services) -> None:
        self._services = services

    def invoke(self, opts: Union[InvokeOptions, Dict[str, Any]]) -> Tuple[int, Dict[str, str], str]:
        """Executes the HTTP call and returns (status, headers, text)."""
        config: I6SdkConfig = self._services.config()

        # 1. Target URL
        if raw_url := opts.get("url"):
            target_url = raw_url
        else:
            path = (opts.get("path") or "").lstrip("/")
            decsuite_url = config.serv_urls.get("decsuite", "")
            target_url = f"{decsuite_url.rstrip('/')}/{path}"

        # 2. Query string
        if query := opts.get("query"):
            filtered_query = {k: v for k, v in query.items() if v is not None}
            if filtered_query:
                qs = urllib.parse.urlencode(filtered_query, doseq=True)
                separator = "&" if "?" in target_url else "?"
                target_url = f"{target_url}{separator}{qs}"

        # 3. Headers
        headers: Dict[str, str] = {
            "User-Agent": "i6sdk-python/0.0.1",
            "Accept": "application/json",
        }

        # Auth token injection
        if token := self._services.tokens().get_token():
            headers["Authorization"] = f"Bearer {token}"

        # Custom headers
        if custom_headers := opts.get("headers"):
            headers.update({str(k): str(v) for k, v in custom_headers.items() if v is not None})

        req_opts = opts.get("opts")
        if req_opts and isinstance(req_opts, dict):
            if extra_headers := req_opts.get("headers"):
                headers.update({str(k): str(v) for k, v in extra_headers.items() if v is not None})

        print("invoker.debug", opts)

        # 4. Payload data
        body_bytes: Optional[bytes] = None
        if (json_data := opts.get("json")) is not None:
            body_bytes = json.dumps(json_data).encode("utf-8")
            headers["Content-Type"] = "application/json"
        elif (raw_body := opts.get("body")) is not None:
            if isinstance(raw_body, str):
                body_bytes = raw_body.encode("utf-8")
            elif isinstance(raw_body, (bytes, bytearray)):
                body_bytes = bytes(raw_body)

        # 5. Timeout
        timeout = config.timeout
        if req_opts and isinstance(req_opts, dict):
            if (t := req_opts.get("timeout")) is not None:
                timeout = float(t)

        method = (opts.get("method") or "GET").upper()

        # 6. Pluggable fetcher support
        if config.fetcher is not None:
            return config.fetcher(target_url, method, headers, body_bytes, timeout)

        # 7. Standard library urllib execution
        req = urllib.request.Request(
            url=target_url,
            data=body_bytes,
            headers=headers,
            method=method,
        )

        try:
            with urllib.request.urlopen(req, timeout=timeout) as resp:
                resp_status = resp.status
                resp_headers = {k.lower(): v for k, v in resp.headers.items()}
                resp_text = resp.read().decode("utf-8", errors="replace")
                return resp_status, resp_headers, resp_text
        except urllib.error.HTTPError as e:
            err_headers = {k.lower(): v for k, v in e.headers.items()} if e.headers else {}
            err_text = ""
            try:
                err_text = e.read().decode("utf-8", errors="replace")
            except Exception:
                pass
            return e.code, err_headers, err_text
        except urllib.error.URLError as e:
            raise I6ConnectionError(f"Connection failed to {target_url}: {e.reason}") from e
        except TimeoutError as e:
            raise I6ConnectionError(f"Request timed out after {timeout}s: {target_url}") from e
        except Exception as e:
            raise I6ConnectionError(f"HTTP request error: {e}") from e

    def invoke_result(self, **kwargs: Any) -> Resp[Any]:
        """Executes the request and wraps it in a Resp container with lazy parsing."""
        status, headers, text = self.invoke(kwargs)

        def parser() -> Any:
            if not (200 <= status < 300):
                code = code_from_status(status)
                msg = text or f"HTTP {status}"
                raise I6ApiError(message=msg, status=status, code=code, cause=text)
            if not text:
                return None
            return json.loads(text)

        return Resp(status, headers, parser, text)
