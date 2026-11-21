package sdkwebservice

import (
	"context"
	"net/http"
	"os"

	"github.com/infinity6-ai/gox/commonz/configz"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"github.com/infinity6-ai/gox/httpz/httpzrequest"
	"github.com/infinity6-ai/gox/httpz/httpzserver"
	"github.com/infinity6-ai/gox/httpz/httpzserver/httpzserverdisk"
	"github.com/infinity6-ai/gox/httpz/httpzserver/httpzserverstatic"
	"github.com/infinity6-ai/i6sdk/sdkweb/internal/stzfiles"
)

var I6_SDK_WEB_SERVER_LOCAL_ADDRESS = configz.Create("I6_SDK_WEB_SERVER_LOCAL_ADDRESS", "127.0.0.1:8043")

type Options struct {
	Context      context.Context
	LocalAddress string
}

type Server struct {
	server *httpzserver.Server
}

func Listen(opts Options) *Server {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}
	addr := opts.LocalAddress
	if addr == "" {
		addr, _ = I6_SDK_WEB_SERVER_LOCAL_ADDRESS.Get(ctx)
		if addr == "" {
			addr = "127.0.0.1:8043"
		}
	}

	srv := httpzserver.New(ctx, httpzserver.Options{
		LocalAddress: addr,
	})

	root := "."
	if _, err := os.Stat("src/main"); err != nil {
		if _, err := os.Stat("sdkweb/src/main"); err == nil {
			root = "sdkweb"
		}
	}
	basePath := pathz.MustParse(root)

	srv.AddHandler("GET", "/site/*", httpzserverstatic.Handler(stzfiles.Name, httpzserverstatic.Options{
		Prefix: "/site",
	}))

	srv.AddHandler("GET", "/src/*", httpzserverdisk.Handler(httpzserverdisk.Options{
		Base:   basePath,
		Prefix: "/src",
	}))

	srv.AddHandler("POST", "/api/stripe/webhook", func(ctx context.Context, resp httpzserver.Resp, req *httpzrequest.Req, params map[string]string) {
		resp(http.StatusNoContent, nil)
	})

	srv.AddHandler("GET", "/", func(ctx context.Context, resp httpzserver.Resp, req *httpzrequest.Req, params map[string]string) {
		resp.Redirect(nil, "/src/site/index.html")
	})

	srv.Listen()

	return &Server{server: srv}
}

func (s *Server) Start() {
	s.server.Start()
}

func (s *Server) Stop() {
	s.server.Close()
}

func (s *Server) Address() string {
	return s.server.Addr().String()
}

func (s *Server) Base() string {
	return s.server.Base().String()
}

func (s *Server) Serve() {
	s.server.Serve()
}

func (s *Server) Server() *httpzserver.Server {
	return s.server
}

func Start(opts Options) *Server {
	s := Listen(opts)
	s.Start()
	return s
}
