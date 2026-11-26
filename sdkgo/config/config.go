package config

import (
	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/httpz/httpzclient"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
)

// DefaultUrlOptions contains the default production service endpoints ([service].app.infinity6.ai).
var DefaultUrlOptions = UrlOptions{
	Decsuite: "https://decsuite.app.infinity6.ai",
	Tfinfer:  "https://tfinfer.app.infinity6.ai",
	Userz:    "https://userz.app.infinity6.ai",
}

// SandboxUrlOptions contains the sandbox service endpoints.
var SandboxUrlOptions = UrlOptions{
	Decsuite: "https://decsuite.sandbox.rs.infinity6.ai",
	Tfinfer:  "https://tfinfer.sandbox.lb.infinity6.ai",
	Userz:    "https://userz.sandbox.rs.infinity6.ai",
}

// StagingUrlOptions contains the staging service endpoints.
var StagingUrlOptions = UrlOptions{
	Decsuite: "https://decsuite.staging.rs.infinity6.ai",
	Tfinfer:  "https://tfinfer.staging.lb.infinity6.ai",
	Userz:    "https://userz.staging.rs.infinity6.ai",
}

// Options holds all configuration parameters for creating an I6Client.
type Options struct {
	// Token is the initial authentication bearer token.
	Token string

	// Urls specifies custom service endpoints. If empty, DefaultUrlOptions is used.
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
