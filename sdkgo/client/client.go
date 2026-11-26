package client

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/httpz/httpzclient"
	"github.com/infinity6-ai/gox/httpz/httpzrequest"
	"github.com/infinity6-ai/gox/routez/apiclientz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/infinity6-ai/i6sdk/sdkgo/dash"
	"github.com/infinity6-ai/i6sdk/sdkgo/infer"
	"github.com/infinity6-ai/i6sdk/sdkgo/ingest"
	"github.com/infinity6-ai/i6sdk/sdkgo/internal/fszcache"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokens"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
)

type authSkipKey struct{}

func extractDatasetFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, p := range parts {
		if p == "dataset" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// I6Client is the primary client for the Infinity6 SDK.
type I6Client struct {
	tokenStore   tokenstore.Store
	tokens       *tokens.Manager
	cacheUrl     *urlz.Url
	memoryCache  bool
	initialToken string

	decsuiteUrl *urlz.Url
	tfinferUrl  *urlz.Url
	userzUrl    *urlz.Url

	decsuiteClient *httpzclient.Client
	tfinferClient  *httpzclient.Client
	userzClient    *httpzclient.Client

	apis   *apicollectionz.SdkApis
	funcs  *apicollectionz.SdkFuncs
	dash   *dash.Client
	infer  *infer.Client
	ingest *ingest.Client
}

// New constructs an I6Client configured with the provided Options.
func New(opts config.Options) (*I6Client, error) {
	ctx := context.Background()
	c := &I6Client{
		initialToken: opts.Token,
	}

	// 1. Resolve URLs
	decUrlStr := opts.Urls.Decsuite
	if decUrlStr == "" {
		decUrlStr = config.DefaultUrlOptions.Decsuite
	}
	decUrl, err := urlz.Parse(decUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decsuite url %q: %w", decUrlStr, err)
	}
	c.decsuiteUrl = decUrl

	tfUrlStr := opts.Urls.Tfinfer
	if tfUrlStr == "" {
		tfUrlStr = config.DefaultUrlOptions.Tfinfer
	}
	tfUrl, err := urlz.Parse(tfUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tfinfer url %q: %w", tfUrlStr, err)
	}
	c.tfinferUrl = tfUrl

	userUrlStr := opts.Urls.Userz
	if userUrlStr == "" {
		userUrlStr = config.DefaultUrlOptions.Userz
	}
	userUrl, err := urlz.Parse(userUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userz url %q: %w", userUrlStr, err)
	}
	c.userzUrl = userUrl

	// 2. Auth filter for bearer token injection (supports dynamic user and dataset tokens)
	authFilter := func(ctx context.Context, req *httpzrequest.Req, next httpzclient.Handler) (*httpzclient.Resp, error) {
		if ctx.Value(authSkipKey{}) != true && req.Headers.Get("Authorization") == "" {
			var tokenToInject string
			var pathStr string
			if req.Path != nil {
				pathStr = req.Path.String()
			} else if req.Url != nil && req.Url.Path != nil {
				pathStr = req.Url.Path.String()
			}

			if ds := extractDatasetFromPath(pathStr); ds != "" && c.tokens != nil {
				if dsTok, err := c.tokens.GetValidDatasetAccessToken(ctx, ds); err == nil && dsTok != "" {
					tokenToInject = dsTok
				}
			}

			if tokenToInject == "" && c.tokens != nil {
				if usrTok, err := c.tokens.GetValidAccessToken(ctx); err == nil && usrTok != "" {
					tokenToInject = usrTok
				}
			}

			if tokenToInject == "" && c.initialToken != "" {
				tokenToInject = c.initialToken
			}

			if tokenToInject != "" {
				req.Headers.Set("Authorization", "Bearer "+tokenToInject)
			}
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
		defUrl, err := fszcache.DefaultBaseUrl()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve default cache url: %w", err)
		}
		c.cacheUrl = defUrl
	}

	// 5. Token store
	if opts.TokenStore != nil {
		c.tokenStore = opts.TokenStore
	} else {
		c.tokenStore = tokenstore.NewMemory()
	}

	// 6. Initialize SdkApis with apiclientz invoker dispatching by service path
	c.apis = &apicollectionz.SdkApis{
		Invoker: func(ctx context.Context, api apiz.Api) (int, error) {
			spec := api.ApiSpec()
			targetClient := c.ClientForPath(spec.Path)
			if spec.Desc != nil {
				if desc := spec.Desc(); desc != nil && desc.Labels != nil && desc.Labels[speclabeltokentype.Key] == speclabeltokentype.None {
					ctx = context.WithValue(ctx, authSkipKey{}, true)
				}
			}
			return apiclientz.Do(ctx, targetClient, api)
		},
	}

	// 7. Token manager
	c.tokens = tokens.New(tokens.ManagerOptions{
		Store:        c.tokenStore,
		Apis:         c.apis,
		UserzBaseUrl: c.userzUrl,
		InitialToken: c.initialToken,
	})

	// 8. Dash analytical query runner & Infer model runner
	fsStore, err := fszcache.NewFs(c.cacheUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize fszcache store: %w", err)
	}
	downloader := fszcache.NewDownloader(fsStore)
	c.dash = dash.New(dash.Options{
		Apis:       c.apis,
		Downloader: downloader,
	})
	c.infer = infer.New(infer.Options{
		Apis:       c.apis,
		Downloader: downloader,
	})
	c.ingest = ingest.New(ingest.Options{
		Apis: c.apis,
	})

	// 9. Initialize SdkFuncs with Dash & Infer delegators
	c.funcs = &apicollectionz.SdkFuncs{
		Invoker: func(ctx context.Context, fn sdkfuncz.Func) (int, error) {
			spec := fn.SdkFuncSpec()
			switch spec.Typ {
			case sdkfuncz.TypClientQuery:
				return c.dash.Execute(ctx, fn)
			case sdkfuncz.TypClientInfer:
				return c.infer.Execute(ctx, fn)
			default:
				return 501, fmt.Errorf("sdkfunc typ %q not yet implemented", spec.Typ)
			}
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
func (c *I6Client) TokenStore() tokenstore.Store {
	return c.tokenStore
}

// InitialToken returns the initial token provided to the client.
func (c *I6Client) InitialToken() string {
	return c.initialToken
}

// Tokens returns the Token Manager instance.
func (c *I6Client) Tokens() *tokens.Manager {
	return c.tokens
}

// Dash returns the Dash analytical query runner.
func (c *I6Client) Dash() *dash.Client {
	return c.dash
}

// Infer returns the Infer model execution runner.
func (c *I6Client) Infer() *infer.Client {
	return c.infer
}

// Ingest returns the Ingest client runner.
func (c *I6Client) Ingest() *ingest.Client {
	return c.ingest
}
