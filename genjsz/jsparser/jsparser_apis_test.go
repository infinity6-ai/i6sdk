package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
)

func TestUnitParseApis(t *testing.T) {
	apis := []apiz.Api{
		routezsamplefraction.NewService().Api(),
		&ingestzgeturlapi.GetUrlApi{},
	}
	actual := jsparser.ParseApis(apis)
	require.NotNil(t, actual)
}
