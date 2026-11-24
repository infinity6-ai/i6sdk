# Issue: Missing Silent Access Token Renewal via Refresh Token in `sdkpy`

**Component:** `sdkpy` (`i6sdk.tokens.tokens` & `i6sdk.invoker`)  
**Severity:** High (Breaks long-running sessions, causing unexpected 401 Unauthorized errors)  
**Status:** Open  

---

## Summary

In `sdkpy`, when a user access token expires, neither `Tokens.get_token()` nor `Invoker` attempts to refresh the access token using the stored `user_refresh_token`. Instead, `get_token()` simply returns `None`. This causes subsequent API requests and dataset token derivations to fail with authentication errors even though valid refresh credentials exist in storage.

---

## Root Cause Analysis

### 1. `get_token()` Drops Expired Tokens Without Triggering Refresh
In `sdkpy/src/i6sdk/tokens/tokens.py`:
```python
def get_token(self) -> Optional[str]:
    """Returns the active access token (user access token)."""
    ut = self.get_user_token()
    if ut and (access_token := ut.get("user_access_token")):
        expires_at = ut.get("expires_at", 0)
        now_ms = int(time.time() * 1000)
        if expires_at == 0 or expires_at > now_ms:
            return access_token
    return None  # <--- BUG: Drops token, never checks refresh_token or triggers auth.refresh()
```

### 2. Dataset Access Token Fetching Fails After Expiry
In `Tokens.get_valid_dataset_access_token()`:
```python
user_token = self.get_token()
if not user_token:
    return None  # <--- Silently aborts if access token expired, even if refresh token is valid!
```

### 3. Contrast with `sdkweb`
In `sdkweb`, `Tokens.getValidAccessToken()` and `Invoker.resolveToken()` perform silent renewal automatically:
```javascript
// sdkweb/src/main/tokens/tokens.js
async getValidAccessToken(skewSeconds = 30) {
  const token = this.getUserToken();
  if (!token) return null;

  if (!this.isAccessTokenExpired(token, skewSeconds)) {
    return token.user_access_token;
  }

  // Access token is expired or expiring soon; refresh session
  const refreshed = await this.refresh();
  return refreshed?.user_access_token || null;
}
```

---

## Steps to Reproduce

1. Authenticate with `client.auth().sign_in(...)` or device code flow.
2. Advance time past the access token expiration (`expires_at`), or manually set `expires_at` to a past timestamp in token storage.
3. Call an authenticated API endpoint (e.g. `client.auth().me()`, `client.apis().user_conf_settings()`).
4. **Expected:** SDK recognizes expired access token, calls `/api/userz/refresh-user` with the refresh token, updates storage, and proceeds.
5. **Actual:** SDK drops the token, sends no `Authorization` header, and the call fails with HTTP 401 Unauthorized.

---

## Proposed Remediation

1. Implement `get_valid_access_token(skew_seconds: float = 30.0) -> Optional[str]` in `Tokens`:
```python
def get_valid_access_token(self, skew_seconds: float = 30.0) -> Optional[str]:
    ut = self.get_user_token()
    if not ut:
        return None

    access_token = ut.get("user_access_token")
    expires_at = ut.get("expires_at", 0)
    now_ms = int(time.time() * 1000)

    # Return immediately if not expired
    if access_token and (expires_at == 0 or now_ms < (expires_at - int(skew_seconds * 1000))):
        return access_token

    # Attempt silent refresh
    if ut.get("user_refresh_token"):
        try:
            refreshed = self._services.auth().refresh()
            if refreshed and (new_token := refreshed.get("user_access_token")):
                return new_token
        except Exception:
            return None

    return None
```

2. Update `get_valid_dataset_access_token()` and `Invoker.invoke()` to use `get_valid_access_token()` instead of `get_token()`.
