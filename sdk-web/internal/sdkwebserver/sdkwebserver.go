package sdkwebserver

import (
	"context"
	"io"

	"go.code.infinity6.ai/platform/config"
	"go.code.infinity6.ai/platform/httpz/httpzserver"
	"go.code.infinity6.ai/platform/service/boxservice"
	"go.code.infinity6.ai/sdkweb/version"
	"go.code.infinity6.ai/web/api/webhandler"
	"go.code.infinity6.ai/web/internal/stfiles"
)

var I6_SDK_WE_PORT = config.Create("I6_SDK_WE_PORT")

type Options struct {
	Context            context.Context
	LocalAddress       string
	DisableAudit       bool
	AllowExternalReqId bool
}

type Server struct {
	server   *httpzserver.Server
	reverter io.Closer
}

func (me *Server) Start() {
	me.server.Start()
}

func (me *Server) Stop() {
	if me.reverter != nil {
		defer me.reverter.Close()
	}
	if me.server != nil {
		defer me.server.Stop()
	}
}

func (me *Server) Address() string {
	return me.server.Address()
}

func (me *Server) Base() string {
	return me.server.Base()
}

func (me *Server) Serve() {
	me.server.Serve()
}

func Listen(opts Options) *Server {
	if opts.LocalAddress == "" {
		if boxservice.Config().Id == "local" {
			// opts.LocalAddress = discovery.Service(datasetmanagercfg.ME).Bind()
		}
	}

	ret := &Server{}
	ret.server = &httpzserver.Server{
		Context:      opts.Context,
		LocalAddress: opts.LocalAddress,
		Audit:        !opts.DisableAudit,
		ServName:     "sdkweb",
		ServVersion:  version.Version(),
	}
	ret.server.Listen()
	ret.server.AllowExternalReqId(opts.AllowExternalReqId)
	return ret
}

func (me *Server) AddHandlers() {
	me.server.Health()

	webhandler.Add(me.server)
	me.server.AddRestPrefix("GET", "/static", stfiles.Handle())

	// datasetmanagerdiscoveryimpl.AddHandlers(me.server)

	// me.server.AddRestPrefix("GET", "/", stfiles.IndexRedirect)
}

func Start(opts Options) *Server {
	server := Listen(opts)
	server.AddHandlers()
	server.Start()
	return server
}

func (me *Server) Server() *httpzserver.Server {
	return me.server
}
