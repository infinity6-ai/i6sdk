"""Ed25519 request signing and refresh assertion utilities matching sdkcorez/reqsign."""

from __future__ import annotations

import base64
import hashlib
import json
import struct
import time
from typing import Any, Dict, Optional, Tuple, Union

try:
    from cryptography.hazmat.primitives import serialization
    from cryptography.hazmat.primitives.asymmetric import ed25519

    HAS_CRYPTOGRAPHY = True
except ImportError:
    HAS_CRYPTOGRAPHY = False


TYPE_DATASET_ACCESS = 1
TYPE_SIG_ENVELOPE = 2
TYPE_REFRESH_ASSERTION = 3
FLAG_UNSIGNED_BODY = 0x80
ASSERTION_LEN = 121


def b64url_decode(s: str) -> bytes:
    """Decodes a base64url string with or without padding."""
    s = s.strip()
    rem = len(s) % 4
    if rem > 0:
        s += "=" * (4 - rem)
    return base64.urlsafe_b64decode(s)


def b64url_encode(b: bytes) -> str:
    """Encodes bytes to unpadded base64url string."""
    return base64.urlsafe_b64encode(b).decode("utf-8").rstrip("=")


def parse_ed25519_private_key(key_input: Union[str, Dict[str, Any], bytes]) -> Any:
    """Extracts and instantiates an Ed25519PrivateKey from a base64url string, dict, or bytes."""
    if not HAS_CRYPTOGRAPHY:
        raise RuntimeError("cryptography package is required for Ed25519 signing")

    if isinstance(key_input, dict):
        raw_key = key_input.get("priv_key") or key_input.get("private_key") or ""
    elif isinstance(key_input, bytes):
        raw_key = key_input
    elif isinstance(key_input, str):
        key_input = key_input.strip()
        if key_input.startswith("{") and key_input.endswith("}"):
            try:
                d = json.loads(key_input)
                raw_key = d.get("priv_key") or d.get("private_key") or ""
            except Exception:
                raw_key = key_input
        else:
            raw_key = key_input
    else:
        raise ValueError(f"Invalid private key input type: {type(key_input)}")

    if isinstance(raw_key, str):
        raw_bytes = b64url_decode(raw_key)
    else:
        raw_bytes = bytes(raw_key)

    if len(raw_bytes) == 64:
        seed = raw_bytes[:32]
    elif len(raw_bytes) == 32:
        seed = raw_bytes
    else:
        raise ValueError(
            f"Invalid Ed25519 private key length: expected 32 or 64 bytes, got {len(raw_bytes)}"
        )

    return ed25519.Ed25519PrivateKey.from_private_bytes(seed)


def build_refresh_assertion(
    priv_key: Any,
    dataset: str,
    method: str = "POST",
    path: str = "/api/userz/refresh-dataset",
    body_bytes: bytes = b"",
    ttl_seconds: int = 300,
    unsigned_body: bool = False,
) -> str:
    """Builds and signs a 121-byte packed RefreshAssertion string (162 Base64URL chars)."""
    if isinstance(priv_key, (str, dict, bytes)):
        priv_key = parse_ed25519_private_key(priv_key)

    norm_method = method.strip().upper()
    norm_path = path.strip()
    if norm_path != "/" and norm_path.endswith("/"):
        norm_path = norm_path.rstrip("/")

    norm_query = ""
    norm_body_hash = ""
    if not unsigned_body:
        norm_body_hash = hashlib.sha256(body_bytes).hexdigest().lower()

    token_hash = ""
    exp = int(time.time()) + ttl_seconds
    nonce = bytes(16)  # 16 zero bytes for time-windowed replay mode
    nonce_hex = nonce.hex()

    canonical = f"{norm_method}\n{norm_path}\n{norm_query}\n{norm_body_hash}\n{token_hash}\n{exp}\n{nonce_hex}"
    sig = priv_key.sign(canonical.encode("utf-8"))

    raw_pub = priv_key.public_key().public_bytes(
        serialization.Encoding.Raw, serialization.PublicFormat.Raw
    )

    buf = bytearray(ASSERTION_LEN)
    v = TYPE_REFRESH_ASSERTION
    if unsigned_body:
        v |= FLAG_UNSIGNED_BODY
    buf[0] = v
    struct.pack_into(">Q", buf, 1, exp)
    buf[9:25] = nonce
    buf[25:57] = raw_pub
    buf[57:121] = sig

    return b64url_encode(bytes(buf))
