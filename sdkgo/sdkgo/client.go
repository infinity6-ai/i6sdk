package sdkgo

import (
	"net/http"
	"os"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/infinity6-ai/i6sdk/versionz/version"
)

const (
	DefaultApp         = "default"
	DefaultDecsuiteUrl = "https://decsuite.sandbox.rs.infinity6.ai"
	DefaultTfinferUrl  = "https://tfinfer.sandbox.lb.infinity6.ai"
	DefaultUserzUrl    = "https://userz.sandbox.rs.infinity6.ai"
)

type Config struct {
	App      string
	ServUrls map[string]string
	Client   *http.Client
}

func LoadConfigFromEnv() Config {
	decsuiteUrl := DefaultDecsuiteUrl
	if env := os.Getenv("I6_SDK_DECSUITE_URL"); env != "" {
		decsuiteUrl = env
	}

	tfinferUrl := DefaultTfinferUrl
	if env := os.Getenv("I6_SDK_TFINFER_URL"); env != "" {
		tfinferUrl = env
	}

	userzUrl := DefaultUserzUrl
	if env := os.Getenv("I6_SDK_USERZ_URL"); env != "" {
		userzUrl = env
	}

	return Config{
		App:      DefaultApp,
		ServUrls: map[string]string{
			"decsuite": decsuiteUrl,
			"tfinfer":  tfinferUrl,
			"userz":    userzUrl,
		},
		Client: http.DefaultClient,
	}
}

func DefaultConfig() Config {
	return LoadConfigFromEnv()
}

type Client struct {
	config Config
}

func New(cfg ...Config) *Client {
	c := DefaultConfig()
	if len(cfg) > 0 {
		c = cfg[0]
	}
	return &Client{config: c}
}

func (c *Client) Version() string {
	return version.Version()
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) Apis() []apiz.Api {
	return apicollectionz.Apis()
}

func (c *Client) Api(id string) apiz.Api {
	return apicollectionz.Api(id)
}

func (c *Client) SdkFuncs() []sdkfuncz.Func {
	return apicollectionz.SdkFuncs()
}

func (c *Client) SdkFunc(id string) sdkfuncz.Func {
	return apicollectionz.SdkFunc(id)
}
