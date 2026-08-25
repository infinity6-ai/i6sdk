package appstack

import (
	"context"

	"go.code.infinity6.ai/platform/httpz"
	"go.code.infinity6.ai/platform/util/deferz"
	"go.code.infinity6.ai/sdkweb/internal/sdkwebserver"
	"go.code.infinity6.ai/web/api/webhandler"
)

type Options struct {
	Context            context.Context
	SdkWebLocalAddress string
	WebCheckerEnabled  bool
}

type Server struct {
	dfz    *deferz.Deferz
	client *httpz.Client
}

func (s *Server) Client() *httpz.Client {
	return s.client
}

func (me *Server) Stop() {
	me.dfz.Close()
}

func StartWithOptions(opts Options) *Server {
	ctx := opts.Context
	ret := &Server{}
	dfz := deferz.New(opts.Context)
	defer dfz.Close()

	if opts.WebCheckerEnabled {
		dfz.AddCloser("I6_WEB_CHECKER_ENABLED", webhandler.I6_WEB_CHECKER_ENABLED.Set("true"))
	}
	if opts.SdkWebLocalAddress != "" {
		dfz.AddCloser("I6_SDK_WEB_SERVER_LOCAL_ADDRESS", sdkwebserver.I6_SDK_WEB_SERVER_LOCAL_ADDRESS.Set("%s", opts.SdkWebLocalAddress))
	}

	s := sdkwebserver.Listen(sdkwebserver.Options{
		Context: ctx,
	})

	dfz.Add("sdkwebserver", s.Stop)

	s.AddHandlers()
	s.Start()

	ret.client = httpz.CreateClient(s.Base())

	ret.dfz = dfz.Detach()

	return ret
}

func StartForTest(ctx context.Context) *Server {
	return StartWithOptions(Options{
		Context:            ctx,
		SdkWebLocalAddress: "127.0.0.1:0",
		WebCheckerEnabled:  false,
	})
}

func StartForTestWeb(ctx context.Context) *Server {
	return StartWithOptions(Options{
		Context:            ctx,
		SdkWebLocalAddress: "127.0.0.1:0",
		WebCheckerEnabled:  true,
	})
}
