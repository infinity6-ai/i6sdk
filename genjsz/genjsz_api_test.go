package genjsz_test

import (
	"fmt"
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"go.code.infinity6.ai/sdkgo/genjsz"
)

func TestManualFractionSample(t *testing.T) {
	api := routezsamplefraction.NewService().Api()
	fmt.Printf("%s\n", jsonz.MustFormatPretty(genjsz.ParseApi(api)).String())

}
