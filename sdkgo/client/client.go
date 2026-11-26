package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/httpz/httpzclient"
	"github.com/infinity6-ai/gox/httpz/httpzrequest"
	"github.com/infinity6-ai/gox/routez/apiclientz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

const (
	// Default Sandbox URLs
	DefaultDecsuiteUrl = "https://decsuite.sandbox.rs.infinity6.ai"
	DefaultTfinferUrl  = "https://tfinfer.sandbox.lb.infinity6.ai"
	DefaultUserzUrl    = "https://userz.sandbox.rs.infinity6.ai"

	// Staging URLs
	StagingDecsuiteUrl = "https://decsuite.staging.rs.infinity6.ai"
	StagingTfinferUrl  = "https://tfinfer.staging.lb.infinity6.ai"
	StagingUserzUrl    = "https://userz.staging.rs.infinity6.ai"
)

// I6Client is the primary client for the Infinity6 SDK.
type I6Client struct {
	tokenStore   TokenStore
	cacheUrl     *urlz.Url
	memoryCache  bool
	initialToken string

	decsuiteUrl *urlz.Url
	tfinferUrl  *urlz.Url
	userzUrl    *urlz.Url

	decsuiteClient *httpzclient.Client
	tfinferClient  *httpzclient.Client
	userzClient    *httpzclient.Client

	apis  *apicollectionz.SdkApis
	funcs *apicollectionz.SdkFuncs
}

// New constructs an I6Client configured with the provided Options.
func New(opts Options) (*I6Client, error) {
	ctx := context.Background()
	c := &I6Client{
		initialToken: opts.Token,
	}

	// 1. Resolve URLs
	decUrlStr := opts.Urls.Decsuite
	if decUrlStr == "" {
		decUrlStr = DefaultDecsuiteUrl
	}
	decUrl, err := urlz.Parse(decUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decsuite url %q: %w", decUrlStr, err)
	}
	c.decsuiteUrl = decUrl

	tfUrlStr := opts.Urls.Tfinfer
	if tfUrlStr == "" {
		tfUrlStr = DefaultTfinferUrl
	}
	tfUrl, err := urlz.Parse(tfUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tfinfer url %q: %w", tfUrlStr, err)
	}
	c.tfinferUrl = tfUrl

	userUrlStr := opts.Urls.Userz
	if userUrlStr == "" {
		userUrlStr = DefaultUserzUrl
	}
	userUrl, err := urlz.Parse(userUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userz url %q: %w", userUrlStr, err)
	}
	c.userzUrl = userUrl

	// 2. Auth filter for bearer token injection
	authFilter := func(ctx context.Context, req *httpzrequest.Req, next httpzclient.Handler) (*httpzclient.Resp, error) {
		if c.initialToken != "" && req.Headers.Get("Authorization") == "" {
			req.Headers.Set("Authorization", "Bearer "+c.initialToken)
		}
		return next(ctx, req)
	}

	// 3. Initialize httpzclient instances with BaseUrl and auth filter
	decOpts := opts.Http
	decOpts.BaseUrl = decUrl
	c.decsuiteClient = httpzclient.New(ctx, decOpts)
	c.decsuiteClient.AddFilter(authFilter)

	tfOpts := opts.Http
	tfOpts.BaseUrl = tfUrl
	c.tfinferClient = httpzclient.New(ctx, tfOpts)
	c.tfinferClient.AddFilter(authFilter)

	userOpts := opts.Http
	userOpts.BaseUrl = userUrl
	c.userzClient = httpzclient.New(ctx, userOpts)
	c.userzClient.AddFilter(authFilter)

	// 4. Resolve cache settings
	if opts.Cache.Memory {
		c.memoryCache = true
	} else if opts.Cache.Url != nil {
		c.cacheUrl = opts.Cache.Url
	} else if opts.Cache.Dir != "" {
		absDir, err := filepath.Abs(opts.Cache.Dir)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve absolute cache dir %q: %w", opts.Cache.Dir, err)
		}
		c.cacheUrl = urlz.MustParse("file://" + filepath.Clean(absDir))
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = os.TempDir()
		}
		defaultCache := filepath.Join(homeDir, ".cache", "i6", "sdkgo", "fs")
		c.cacheUrl = urlz.MustParse("file://" + filepath.Clean(defaultCache))
	}

	// 5. Token store
	if opts.TokenStore != nil {
		c.tokenStore = opts.TokenStore
	} else {
		c.tokenStore = newMemoryTokenStore()
	}

	// 6. Initialize SdkApis with apiclientz invoker dispatching by service path
	c.apis = &apicollectionz.SdkApis{
		Invoker: func(ctx context.Context, api apiz.Api) (int, error) {
			targetClient := c.ClientForPath(api.ApiSpec().Path)
			return apiclientz.Do(ctx, targetClient, api)
		},
	}

	// 7. Initialize SdkFuncs with fallback invoker (wired in Step 4.7)
	c.funcs = &apicollectionz.SdkFuncs{
		Invoker: func(ctx context.Context, fn sdkfuncz.Func) (int, error) {
			return 501, fmt.Errorf("sdkfunc invoker not yet wired: step 4.7")
		},
	}

	return c, nil
}

// ClientForPath routes an API path to the appropriate service httpzclient.Client.
func (c *I6Client) ClientForPath(path string) *httpzclient.Client {
	if strings.HasPrefix(path, "/api/userz/") {
		return c.userzClient
	}
	if strings.HasPrefix(path, "/api/metadata/") || strings.HasPrefix(path, "/api/model/") {
		return c.tfinferClient
	}
	return c.decsuiteClient
}

// Apis returns the typed SdkApis client containing all API endpoints.
func (c *I6Client) Apis() *apicollectionz.SdkApis {
	return c.apis
}

// Funcs returns the typed SdkFuncs client containing all workflow functions.
func (c *I6Client) Funcs() *apicollectionz.SdkFuncs {
	return c.funcs
}

// DecsuiteUrl returns the configured DecSuite base URL.
func (c *I6Client) DecsuiteUrl() *urlz.Url {
	return c.decsuiteUrl
}

// TfinferUrl returns the configured TensorFlow Inference base URL.
func (c *I6Client) TfinferUrl() *urlz.Url {
	return c.tfinferUrl
}

// UserzUrl returns the configured Userz authentication base URL.
func (c *I6Client) UserzUrl() *urlz.Url {
	return c.userzUrl
}

// BaseUrl returns the default Decsuite base URL.
func (c *I6Client) BaseUrl() *urlz.Url {
	return c.decsuiteUrl
}

// DecsuiteClient returns the httpzclient for the Decsuite service.
func (c *I6Client) DecsuiteClient() *httpzclient.Client {
	return c.decsuiteClient
}

// TfinferClient returns the httpzclient for the TensorFlow Inference service.
func (c *I6Client) TfinferClient() *httpzclient.Client {
	return c.tfinferClient
}

// UserzClient returns the httpzclient for the Userz service.
func (c *I6Client) UserzClient() *httpzclient.Client {
	return c.userzClient
}

// CacheUrl returns the configured fsz cache URL.
func (c *I6Client) CacheUrl() *urlz.Url {
	return c.cacheUrl
}

// IsMemoryCache returns true if in-memory cache is active.
func (c *I6Client) IsMemoryCache() bool {
	return c.memoryCache
}

// TokenStore returns the active token persistence store.
func (c *I6Client) TokenStore() TokenStore {
	return c.tokenStore
}

// InitialToken returns the initial token provided to the client.
func (c *I6Client) InitialToken() string {
	return c.initialToken
}
