package logzinit_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/logzinit"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

func TestManualTrace(t *testing.T) {
	ctx := t.Context()

	ctx, closer := logzinit.Init(ctx, logzinit.Options{
		App:       "a",
		Version:   "v5",
		ProjectId: "i6-rs-exp-main",
	})
	defer closer.Close()

	logger.Info(ctx, "aaaaa", nil)

}
