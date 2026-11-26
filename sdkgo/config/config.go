package config

import (
	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/httpz/httpzclient"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
)

const (
	// Production URLs ([service].app.infinity6.ai) - DEFAULT
	ProdDecsuiteUrl = "https://decsuite.app.infinity6.ai"
	ProdTfinferUrl  = "https://tfinfer.app.infinity6.ai"
	ProdUserzUrl    = "https://userz.app.infinity6.ai"

	// Staging URLs
	StagingDecsuiteUrl = "https://decsuite.staging.rs.infinity6.ai"
	StagingTfinferUrl  = "https://tfinfer.staging.lb.infinity6.ai"
	StagingUserzUrl    = "https://userz.staging.rs.infinity6.ai"

	// Sandbox URLs
	SandboxDecsuiteUrl = "https://decsuite.sandbox.rs.infinity6.ai"
	SandboxTfinferUrl  = "https://tfinfer.sandbox.lb.infinity6.ai"
	SandboxUserzUrl    = "https://userz.sandbox.rs.infinity6.ai"

	// Defaults point to Production
	DefaultDecsuiteUrl = ProdDecsuiteUrl
	DefaultTfinferUrl  = ProdTfinferUrl
	DefaultUserzUrl    = ProdUserzUrl
)

// Options holds all configuration parameters for creating an I6Client.
type Options struct {
	// Token is the initial authentication bearer token.
	Token string

	// Urls specifies custom service endpoints. If empty, production defaults are used.
	Urls UrlOptions

	// Cache configures SQLite analytical bundle storage.
	Cache CacheOptions

	// TokenStore customizes token persistence (defaults to in-memory).
	TokenStore tokenstore.Store

	// Http configures options for the underlying httpzclient instances.
	Http httpzclient.Options
}

// UrlOptions configures individual service endpoints.
type UrlOptions struct {
	Decsuite string // Defaults to https://decsuite.app.infinity6.ai
	Tfinfer  string // Defaults to https://tfinfer.app.infinity6.ai
	Userz    string // Defaults to https://userz.app.infinity6.ai
}

// CacheOptions configures analytical bundle caching.
type CacheOptions struct {
	Dir    string    // Local filesystem cache directory (converted to file://)
	Url    *urlz.Url // Custom fsz URL scheme
	Memory bool      // Ephemeral in-memory SQLite execution
}

// ProdUrls contains the production service endpoints (default).
var ProdUrls = UrlOptions{
	Decsuite: ProdDecsuiteUrl,
	Tfinfer:  ProdTfinferUrl,
	Userz:    ProdUserzUrl,
}

// StagingUrls contains the staging service endpoints.
var StagingUrls = UrlOptions{
	Decsuite: StagingDecsuiteUrl,
	Tfinfer:  StagingTfinferUrl,
	Userz:    StagingUserzUrl,
}

// SandboxUrls contains the sandbox service endpoints.
var SandboxUrls = UrlOptions{
	Decsuite: SandboxDecsuiteUrl,
	Tfinfer:  SandboxTfinferUrl,
	Userz:    SandboxUserzUrl,
}
