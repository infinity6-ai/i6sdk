package jsparser_test

import (
	"fmt"
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
)

func TestUnitFractionSample(t *testing.T) {
	api := routezsamplefraction.NewService().Api()
	fmt.Printf("%s\n", jsonz.MustFormatPretty(jsparser.ParseApi(api)).String())
	require.Equal(t, jsparser.Api{}, jsparser.ParseApi(api))
}

func TestUnitIngestz(t *testing.T) {
	api := &ingestzgeturlapi.GetUrlApi{}
	fmt.Printf("%s\n", jsonz.MustFormatPretty(jsparser.ParseApi(api)).String())
}
