package logzinit

import (
	"context"
	"io"

	"github.com/infinity6-ai/gox/commonz/ioz"
	"github.com/infinity6-ai/gox/commonz/logz"
	"go.code.infinity6.ai/sdkgo/sdkgo/logzinit/internal/logztracer"
	"google.golang.org/grpc"
)

var done = false

type App struct{ App string }
type Version struct{ Version string }

func Init(ctx context.Context, app App, version Version, tracerOptions *logztracer.Options) (context.Context, io.Closer) {
	if done {
		panic("already configured")
	}
	if app.App == "" {
		panic("app is required")
	}
	if version.Version == "" {
		panic("version is required")
	}
	ctx = logz.SetAttrs(ctx, map[string]string{
		"app.name":    app.App,
		"app.version": version.Version,
	})
	grpc.EnableTracing = false
	tracerClose := logztracer.InitTracer(ctx, app.App, version.Version, tracerOptions)
	metricsClose := logztracer.InitMetric(ctx, app.App, version.Version, tracerOptions)
	if tracerOptions.ProjectId != "" {
		// tracerz.Init(app.App)
		// meterz.Init(app.App)
	}
	done = true
	return ctx, ioz.Closer(func() error {
		metricsClose.Close()
		return tracerClose.Close()
	})
}
