package sdkgo

import (
	"net/http"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
	"github.com/infinity6-ai/i6sdk/versionz/version"
)

type Config struct {
	ServUrls map[string]string
	Client   *http.Client
}

func DefaultConfig() Config {
	return Config{
		ServUrls: map[string]string{
			"decsuite": "https://decsuite.sandbox.rs.infinity6.ai",
			"tfinfer":  "https://tfinfer.sandbox.lb.infinity6.ai",
			"userz":    "https://userz.sandbox.rs.infinity6.ai",
		},
		Client: http.DefaultClient,
	}
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
