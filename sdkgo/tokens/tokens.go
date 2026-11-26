package tokens

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccesscreate"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefresh"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapiauthrefresh"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/reqsign"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
)

const (
	// UserTokenKey is the storage key for user authentication tokens.
	UserTokenKey = "i6_user_user_token"

	// DatasetTokenPrefix is the key prefix for cached dataset access tokens.
	DatasetTokenPrefix = "i6_ds:"

	// DatasetRefreshPrefix is the key prefix for machine robot dataset refresh keys.
	DatasetRefreshPrefix = "i6_rf:"

	// Deprecated: Monolithic storage keys retained for backward compatibility constants.
	DatasetTokensKey        = "i6_dataset_access_tokens"
	DatasetRefreshTokensKey = "i6_dataset_refresh_tokens"

	// DefaultSkewSeconds is the default clock skew window applied before token expiration.
	DefaultSkewSeconds = 30.0
)

func datasetTokenKey(dataset string) string {
	return DatasetTokenPrefix + dataset
}

func datasetRefreshKey(dataset string) string {
	return DatasetRefreshPrefix + dataset
}

// UserToken represents authenticated human user credentials.
type UserToken struct {
	UserAccessToken       string `json:"user_access_token"`
	UserRefreshToken      string `json:"user_refresh_token,omitempty"`
	UserId                string `json:"user_id,omitempty"`
	ExpiresAt             int64  `json:"expires_at,omitempty"`                // epoch milliseconds
	RefreshTokenExpiresAt *int64 `json:"refresh_token_expires_at,omitempty"` // epoch milliseconds
}

// DatasetToken represents dataset authorization credentials.
type DatasetToken struct {
	Dataset     string       `json:"dataset"`
	AccessToken string       `json:"access_token"`
	ExpiresAt   int64        `json:"expires_at"` // epoch milliseconds
	Grants      grants.Grant `json:"grants,omitempty"`
	Solution    string       `json:"solution,omitempty"`
}

// DatasetTokenOptions configures dataset token derivation options.
type DatasetTokenOptions struct {
	Subsolution string
	Subgrants   grants.Grant
	ExpSeconds  int64
	Unbound     bool
}

// ManagerOptions holds configuration parameters for the Token Manager.
type ManagerOptions struct {
	Store        tokenstore.Store
	Apis         *apicollectionz.SdkApis
	UserzBaseUrl *urlz.Url
	InitialToken string
}

// Manager orchestrates storage, double-checked locking renewal, and derivation of access tokens.
type Manager struct {
	store        tokenstore.Store
	apis         *apicollectionz.SdkApis
	userzBaseUrl *urlz.Url
	initialToken string

	refreshMu sync.Mutex
	datasetMu sync.Mutex
}

// New creates a new token Manager instance.
func New(opts ManagerOptions) *Manager {
	store := opts.Store
	if store == nil {
		store = tokenstore.NewMemory()
	}
	return &Manager{
		store:        store,
		apis:         opts.Apis,
		userzBaseUrl: opts.UserzBaseUrl,
		initialToken: opts.InitialToken,
	}
}

// Store returns the underlying token persistence store.
func (m *Manager) Store() tokenstore.Store {
	return m.store
}

// SetStore updates the token persistence store.
func (m *Manager) SetStore(s tokenstore.Store) {
	m.store = s
}

// GetUserToken retrieves the stored user token if present.
func (m *Manager) GetUserToken(ctx context.Context) (*UserToken, error) {
	val, err := m.store.Get(ctx, UserTokenKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token from store: %w", err)
	}
	if strings.TrimSpace(val) == "" {
		return nil, nil
	}
	var token UserToken
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user token: %w", err)
	}
	return &token, nil
}

// SetUserToken stores the given user token and evicts any robot keys and derived dataset access tokens.
func (m *Manager) SetUserToken(ctx context.Context, token *UserToken) error {
	if token == nil {
		return m.RemoveUserToken(ctx)
	}
	raw, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal user token: %w", err)
	}
	// Mutual exclusivity: clear any machine robot refresh keys
	_ = m.store.DeletePrefix(ctx, DatasetRefreshPrefix)
	// Clear any previous derived dataset tokens
	_ = m.store.DeletePrefix(ctx, DatasetTokenPrefix)

	if err := m.store.Set(ctx, UserTokenKey, string(raw)); err != nil {
		return fmt.Errorf("failed to store user token: %w", err)
	}
	return nil
}

// RemoveUserToken deletes the stored user token.
func (m *Manager) RemoveUserToken(ctx context.Context) error {
	if err := m.store.Delete(ctx, UserTokenKey); err != nil {
		return fmt.Errorf("failed to delete user token from store: %w", err)
	}
	return nil
}

// ClearTokens deletes all stored user and dataset tokens.
func (m *Manager) ClearTokens(ctx context.Context) error {
	if err := m.RemoveUserToken(ctx); err != nil {
		return err
	}
	if err := m.store.DeletePrefix(ctx, DatasetTokenPrefix); err != nil {
		return fmt.Errorf("failed to delete dataset tokens: %w", err)
	}
	if err := m.store.DeletePrefix(ctx, DatasetRefreshPrefix); err != nil {
		return fmt.Errorf("failed to delete dataset refresh tokens: %w", err)
	}
	return nil
}

// GetDatasetRefreshTokensMap returns the map of dataset IDs to machine robot refresh keys.
func (m *Manager) GetDatasetRefreshTokensMap(ctx context.Context) (map[string]string, error) {
	rawMap, err := m.store.ListPrefix(ctx, DatasetRefreshPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list dataset refresh tokens: %w", err)
	}
	tokens := make(map[string]string, len(rawMap))
	for k, v := range rawMap {
		ds := strings.TrimPrefix(k, DatasetRefreshPrefix)
		tokens[ds] = v
	}
	return tokens, nil
}

// GetDatasetRefreshToken returns the machine robot Ed25519 refresh key for the given dataset.
func (m *Manager) GetDatasetRefreshToken(ctx context.Context, dataset string) (string, error) {
	if dataset == "" {
		return "", fmt.Errorf("dataset identifier cannot be empty")
	}
	val, err := m.store.Get(ctx, datasetRefreshKey(dataset))
	if err != nil {
		return "", fmt.Errorf("failed to get dataset refresh token for %q: %w", dataset, err)
	}
	return val, nil
}

// SetDatasetRefreshToken stores a machine robot Ed25519 refresh key for the given dataset,
// enforcing mutual exclusivity by evicting any user session, previous robot keys, and derived tokens.
func (m *Manager) SetDatasetRefreshToken(ctx context.Context, dataset string, key string) error {
	if dataset == "" {
		return fmt.Errorf("dataset identifier cannot be empty")
	}
	if key == "" {
		return m.RemoveDatasetRefreshToken(ctx, dataset)
	}
	// Mutual exclusivity: clear any user session
	_ = m.store.Delete(ctx, UserTokenKey)
	// Exactly one robot key: clear any previous robot keys
	_ = m.store.DeletePrefix(ctx, DatasetRefreshPrefix)
	// Clear any previous derived dataset tokens
	_ = m.store.DeletePrefix(ctx, DatasetTokenPrefix)

	if err := m.store.Set(ctx, datasetRefreshKey(dataset), key); err != nil {
		return fmt.Errorf("failed to store dataset refresh token for %q: %w", dataset, err)
	}
	return nil
}

// RemoveDatasetRefreshToken deletes the stored refresh key for the given dataset.
func (m *Manager) RemoveDatasetRefreshToken(ctx context.Context, dataset string) error {
	if dataset == "" {
		return fmt.Errorf("dataset identifier cannot be empty")
	}
	if err := m.store.Delete(ctx, datasetRefreshKey(dataset)); err != nil {
		return fmt.Errorf("failed to delete dataset refresh token for %q: %w", dataset, err)
	}
	return nil
}

// GetDatasetTokensMap returns all cached dataset access tokens.
func (m *Manager) GetDatasetTokensMap(ctx context.Context) (map[string]*DatasetToken, error) {
	rawMap, err := m.store.ListPrefix(ctx, DatasetTokenPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list dataset access tokens: %w", err)
	}
	tokens := make(map[string]*DatasetToken, len(rawMap))
	for k, v := range rawMap {
		ds := strings.TrimPrefix(k, DatasetTokenPrefix)
		var dt DatasetToken
		if err := json.Unmarshal([]byte(v), &dt); err == nil {
			tokens[ds] = &dt
		}
	}
	return tokens, nil
}

// GetDatasetToken returns the cached dataset access token for the target dataset.
func (m *Manager) GetDatasetToken(ctx context.Context, dataset string) (*DatasetToken, error) {
	if dataset == "" {
		return nil, fmt.Errorf("dataset identifier cannot be empty")
	}
	val, err := m.store.Get(ctx, datasetTokenKey(dataset))
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset access token for %q: %w", dataset, err)
	}
	if strings.TrimSpace(val) == "" {
		return nil, nil
	}
	var dt DatasetToken
	if err := json.Unmarshal([]byte(val), &dt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dataset token for %q: %w", dataset, err)
	}
	return &dt, nil
}

// SetDatasetToken stores a cached dataset access token.
func (m *Manager) SetDatasetToken(ctx context.Context, dataset string, token *DatasetToken) error {
	if dataset == "" {
		return fmt.Errorf("dataset identifier cannot be empty")
	}
	if token == nil {
		return m.RemoveDatasetToken(ctx, dataset)
	}
	raw, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal dataset token for %q: %w", dataset, err)
	}
	if err := m.store.Set(ctx, datasetTokenKey(dataset), string(raw)); err != nil {
		return fmt.Errorf("failed to store dataset token for %q: %w", dataset, err)
	}
	return nil
}

// RemoveDatasetToken removes the cached access token for the target dataset.
func (m *Manager) RemoveDatasetToken(ctx context.Context, dataset string) error {
	if dataset == "" {
		return fmt.Errorf("dataset identifier cannot be empty")
	}
	if err := m.store.Delete(ctx, datasetTokenKey(dataset)); err != nil {
		return fmt.Errorf("failed to delete dataset token for %q: %w", dataset, err)
	}
	return nil
}

// IsAccessTokenExpired checks whether the user access token is expired or expiring within skewSeconds.
func (m *Manager) IsAccessTokenExpired(token *UserToken, skewSeconds ...float64) bool {
	if token == nil || token.UserAccessToken == "" {
		return true
	}
	if token.ExpiresAt == 0 {
		return false
	}
	skew := DefaultSkewSeconds
	if len(skewSeconds) > 0 {
		skew = skewSeconds[0]
	}
	nowMs := time.Now().UnixMilli()
	return nowMs >= (token.ExpiresAt - int64(skew*1000))
}

// IsRefreshTokenExpired checks whether the user refresh token is missing or expired.
func (m *Manager) IsRefreshTokenExpired(token *UserToken) bool {
	if token == nil || token.UserRefreshToken == "" {
		return true
	}
	if token.RefreshTokenExpiresAt != nil && *token.RefreshTokenExpiresAt > 0 {
		return time.Now().UnixMilli() >= *token.RefreshTokenExpiresAt
	}
	return false
}

// IsDatasetTokenExpired checks whether a dataset access token is expired or expiring within skewSeconds.
func (m *Manager) IsDatasetTokenExpired(token *DatasetToken, skewSeconds ...float64) bool {
	if token == nil || token.AccessToken == "" {
		return true
	}
	if token.ExpiresAt == 0 {
		return false
	}
	skew := DefaultSkewSeconds
	if len(skewSeconds) > 0 {
		skew = skewSeconds[0]
	}
	nowMs := time.Now().UnixMilli()
	return nowMs >= (token.ExpiresAt - int64(skew*1000))
}

// GetValidAccessToken returns an active user access token, silently refreshing it if expired.
// Employs double-checked locking to avoid concurrent refresh storms.
func (m *Manager) GetValidAccessToken(ctx context.Context, skewSeconds ...float64) (string, error) {
	// 1. Fast read without mutex
	ut, err := m.GetUserToken(ctx)
	if err != nil {
		return "", err
	}
	if ut != nil && !m.IsAccessTokenExpired(ut, skewSeconds...) {
		return ut.UserAccessToken, nil
	}

	// If no user token in store, fall back to initial static token
	if ut == nil {
		if m.initialToken != "" {
			return m.initialToken, nil
		}
		return "", nil
	}

	// 2. Slow path: synchronize under mutex
	m.refreshMu.Lock()
	defer m.refreshMu.Unlock()

	// Re-check under lock in case another goroutine refreshed it
	ut, err = m.GetUserToken(ctx)
	if err != nil {
		return "", err
	}
	if ut != nil && !m.IsAccessTokenExpired(ut, skewSeconds...) {
		return ut.UserAccessToken, nil
	}

	// Check if refresh token is present and valid
	if ut == nil || ut.UserRefreshToken == "" || m.IsRefreshTokenExpired(ut) {
		if m.initialToken != "" {
			return m.initialToken, nil
		}
		return "", nil
	}

	// 3. Perform silent renewal via SdkApis
	if m.apis == nil {
		if m.initialToken != "" {
			return m.initialToken, nil
		}
		return "", fmt.Errorf("cannot refresh user token: apis client not wired")
	}

	refreshApi := &userapiauthrefresh.Api{
		Req: &userapiauthrefresh.Req{
			Payload: &userapiauthrefresh.ReqPayload{
				RefreshToken: ut.UserRefreshToken,
			},
		},
	}

	status, err := m.apis.UserAuthRefresh(ctx, refreshApi)
	if err != nil || status < 200 || status >= 300 || refreshApi.Resp == nil || refreshApi.Resp.Payload == nil {
		if status == 401 || status == 403 {
			_ = m.RemoveUserToken(ctx)
		}
		if m.initialToken != "" {
			return m.initialToken, nil
		}
		return "", fmt.Errorf("failed to refresh user token (status %d): %w", status, err)
	}

	p := refreshApi.Resp.Payload
	expAt := p.ExpiresAt
	if 0 < expAt && expAt < 10_000_000_000 {
		expAt *= 1000
	}

	var rfExp *int64
	if p.RefreshTokenExpiresAt > 0 {
		val := p.RefreshTokenExpiresAt
		if val < 10_000_000_000 {
			val *= 1000
		}
		rfExp = &val
	} else {
		rfExp = ut.RefreshTokenExpiresAt
	}

	refreshed := &UserToken{
		UserAccessToken:       p.UserAccessToken,
		UserRefreshToken:      p.UserRefreshToken,
		UserId:                p.UserId,
		ExpiresAt:             expAt,
		RefreshTokenExpiresAt: rfExp,
	}
	if refreshed.UserRefreshToken == "" {
		refreshed.UserRefreshToken = ut.UserRefreshToken
	}

	if err := m.SetUserToken(ctx, refreshed); err != nil {
		return "", fmt.Errorf("failed to save refreshed user token: %w", err)
	}

	return refreshed.UserAccessToken, nil
}

// RefreshDatasetAccessToken mints a fresh dataset access token using an Ed25519 refresh key.
func (m *Manager) RefreshDatasetAccessToken(ctx context.Context, dataset string, opts ...DatasetTokenOptions) (*DatasetToken, error) {
	if dataset == "" {
		return nil, fmt.Errorf("dataset identifier cannot be empty")
	}
	rfToken, err := m.GetDatasetRefreshToken(ctx, dataset)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(rfToken) == "" {
		return nil, fmt.Errorf("no dataset refresh token found for %q", dataset)
	}

	privKey, err := parseEd25519PrivateKey(rfToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ed25519 private key for dataset %q: %w", dataset, err)
	}

	var opt DatasetTokenOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	expSeconds := opt.ExpSeconds
	if expSeconds <= 0 {
		expSeconds = 300 // default 5-minute token
	}

	// Sign standalone RefreshAssertion (Type 3)
	var nonce [16]byte
	exp := uint64(time.Now().Unix() + expSeconds)
	assertion, err := reqsign.SignAssertion(privKey, "POST", "/api/userz/refresh-dataset", "", nil, exp, nonce, false)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh assertion: %w", err)
	}
	sigHeader := assertion.Pack()

	if m.apis == nil {
		return nil, fmt.Errorf("cannot refresh dataset token: apis client not wired")
	}

	var subSol solution.Solution
	if opt.Subsolution != "" {
		if err := checked.Set(&subSol, opt.Subsolution); err != nil {
			return nil, fmt.Errorf("invalid subsolution %q: %w", opt.Subsolution, err)
		}
	}

	var ds lakedsz.Dataset
	if err := checked.Set(&ds, dataset); err != nil {
		return nil, fmt.Errorf("invalid dataset %q: %w", dataset, err)
	}

	refreshApi := &userapidatasettokenrefresh.Api{
		Req: &userapidatasettokenrefresh.Req{
			Headers: &userapidatasettokenrefresh.ReqHeaders{
				SignatureHeader: sigHeader,
			},
			Payload: &userapidatasettokenrefresh.ReqPayload{
				Dataset:     ds,
				Subsolution: subSol,
				Subgrants:   opt.Subgrants,
				ExpSeconds:  expSeconds,
				Unbound:     true,
			},
		},
	}

	status, err := m.apis.UserDatasetTokenRefresh(ctx, refreshApi)
	if err != nil || status < 200 || status >= 300 || refreshApi.Resp == nil || refreshApi.Resp.Payload == nil {
		return nil, fmt.Errorf("failed to refresh dataset token for %q (status %d): %w", dataset, status, err)
	}

	p := refreshApi.Resp.Payload
	expAt := p.ExpiresAt
	if 0 < expAt && expAt < 10_000_000_000 {
		expAt *= 1000
	}

	dt := &DatasetToken{
		Dataset:     p.Dataset.Optional().Or(""),
		AccessToken: p.AccessToken,
		ExpiresAt:   expAt,
		Grants:      p.Grants,
		Solution:    p.Solution.Optional().Or(""),
	}

	if err := m.SetDatasetToken(ctx, dataset, dt); err != nil {
		return nil, fmt.Errorf("failed to store refreshed dataset token: %w", err)
	}

	return dt, nil
}

// GetValidDatasetAccessToken returns an active dataset access token for the given dataset.
// If missing or expired:
// 1. Uses a robot Ed25519 refresh token if stored.
// 2. Otherwise derives a short-lived token using the valid user access token.
func (m *Manager) GetValidDatasetAccessToken(ctx context.Context, dataset string, skewSeconds ...float64) (string, error) {
	if dataset == "" {
		return "", fmt.Errorf("dataset identifier cannot be empty")
	}

	// 1. Fast read without lock
	dt, err := m.GetDatasetToken(ctx, dataset)
	if err != nil {
		return "", err
	}
	if dt != nil && !m.IsDatasetTokenExpired(dt, skewSeconds...) {
		return dt.AccessToken, nil
	}

	// 2. Slow path: synchronize under mutex
	m.datasetMu.Lock()
	defer m.datasetMu.Unlock()

	// Re-check under lock
	dt, err = m.GetDatasetToken(ctx, dataset)
	if err != nil {
		return "", err
	}
	if dt != nil && !m.IsDatasetTokenExpired(dt, skewSeconds...) {
		return dt.AccessToken, nil
	}

	// Mode A: Machine robot Ed25519 refresh token
	rfToken, err := m.GetDatasetRefreshToken(ctx, dataset)
	if err == nil && strings.TrimSpace(rfToken) != "" {
		refreshed, err := m.RefreshDatasetAccessToken(ctx, dataset)
		if err != nil {
			return "", fmt.Errorf("failed to refresh dataset token using robot key for %q: %w", dataset, err)
		}
		if refreshed != nil && refreshed.AccessToken != "" {
			return refreshed.AccessToken, nil
		}
		return "", fmt.Errorf("empty access token returned after refreshing dataset %q", dataset)
	}

	// Mode B: Derive via user access token
	userToken, err := m.GetValidAccessToken(ctx, skewSeconds...)
	if err != nil {
		return "", fmt.Errorf("failed to obtain user token for dataset derivation: %w", err)
	}
	if userToken == "" {
		return "", fmt.Errorf("no user token or robot key available to derive access token for dataset %q", dataset)
	}

	if m.apis == nil {
		return "", fmt.Errorf("cannot derive dataset token: apis client not wired")
	}

	var ds lakedsz.Dataset
	if err := checked.Set(&ds, dataset); err != nil {
		return "", fmt.Errorf("invalid dataset %q: %w", dataset, err)
	}

	deriveApi := &userapidatasettokenaccesscreate.Api{
		Req: &userapidatasettokenaccesscreate.Req{
			Params: &userapidatasettokenaccesscreate.ReqParams{
				Dataset: ds,
			},
			Headers: &userapidatasettokenaccesscreate.ReqHeaders{
				Authorization: "Bearer " + userToken,
			},
			Payload: &userapidatasettokenaccesscreate.ReqPayload{
				Grants: grants.All(),
			},
		},
	}

	status, err := m.apis.UserDatasetTokenAccessCreate(ctx, deriveApi)
	if err != nil || status < 200 || status >= 300 || deriveApi.Resp == nil || deriveApi.Resp.Payload == nil {
		return "", fmt.Errorf("failed to derive dataset token for %q (status %d): %w", dataset, status, err)
	}

	p := deriveApi.Resp.Payload
	expAt := p.ExpiresAt
	if 0 < expAt && expAt < 10_000_000_000 {
		expAt *= 1000
	}

	dt = &DatasetToken{
		Dataset:     p.Dataset.Optional().Or(""),
		AccessToken: p.AccessToken,
		ExpiresAt:   expAt,
		Grants:      p.Grants,
		Solution:    p.Solution.Optional().Or(""),
	}

	if err := m.SetDatasetToken(ctx, dataset, dt); err != nil {
		return "", fmt.Errorf("failed to cache derived dataset token: %w", err)
	}

	return dt.AccessToken, nil
}

// parseEd25519PrivateKey parses an Ed25519 private key from various input representations:
// - JSON string containing {"priv_key": ...} or {"private_key": ...}
// - Base64URL string (raw 32-byte seed or 64-byte private key)
// - Hex string
func parseEd25519PrivateKey(keyInput string) (ed25519.PrivateKey, error) {
	keyInput = strings.TrimSpace(keyInput)
	if keyInput == "" {
		return nil, fmt.Errorf("private key cannot be empty")
	}

	// Try JSON unmarshal
	if strings.HasPrefix(keyInput, "{") && strings.HasSuffix(keyInput, "}") {
		var d map[string]any
		if err := json.Unmarshal([]byte(keyInput), &d); err == nil {
			if pk, ok := d["priv_key"].(string); ok && pk != "" {
				keyInput = strings.TrimSpace(pk)
			} else if pk, ok := d["private_key"].(string); ok && pk != "" {
				keyInput = strings.TrimSpace(pk)
			}
		}
	}

	// Try Base64URL decode
	rawBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(keyInput, "="))
	if err != nil {
		// Try standard Base64
		rawBytes, err = base64.StdEncoding.DecodeString(keyInput)
		if err != nil {
			// Try hex
			rawBytes, err = hex.DecodeString(keyInput)
			if err != nil {
				// Treat as raw bytes
				rawBytes = []byte(keyInput)
			}
		}
	}

	if len(rawBytes) == ed25519.PrivateKeySize {
		return ed25519.PrivateKey(rawBytes), nil
	}
	if len(rawBytes) == ed25519.SeedSize {
		return ed25519.NewKeyFromSeed(rawBytes), nil
	}

	return nil, fmt.Errorf("invalid ed25519 private key length: expected %d (seed) or %d bytes, got %d",
		ed25519.SeedSize, ed25519.PrivateKeySize, len(rawBytes))
}
