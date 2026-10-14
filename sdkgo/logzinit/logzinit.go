package logzinit

import (
	"context"
	"io"

	"github.com/infinity6-ai/gox/commonz/ioz"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/commonz/logz/tracerz"
	"go.code.infinity6.ai/sdkgo/sdkgo/logzinit/internal/logztracer"
	"google.golang.org/grpc"
)

var done = false

type Options struct {
	App       string
	Version   string
	ProjectId string
}

func Init(ctx context.Context, opts Options) (context.Context, io.Closer) {
	if done {
		panic("already configured")
	}
	if opts.App == "" {
		panic("app is required")
	}
	if opts.Version == "" {
		panic("version is required")
	}
	ctx = logz.SetAttrs(ctx, map[string]string{
		"app.name":    opts.App,
		"app.version": opts.Version,
	})
	grpc.EnableTracing = false
	tracerClose := logztracer.InitTracer(ctx, opts.App, opts.Version, &logztracer.Options{ProjectId: opts.ProjectId})
	// metricsClose := logztracer.InitMetric(ctx, opts.App, opts.Version, tracerOptions)
	if opts.ProjectId != "" {
		tracerz.Init(opts.App)
		// meterz.Init(app.App)
	}
	done = true
	return ctx, ioz.Closer(func() error {
		// metricsClose.Close()
		return tracerClose.Close()
	})
}
