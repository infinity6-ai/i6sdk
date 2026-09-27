package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
)

func TestUnitFractionSample(t *testing.T) {
	api := routezsamplefraction.NewService().Api()
	actual := jsparser.ParseApi(api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id: "aaaa",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "xxx", Name: "bbb"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "xxx", Name: "yyyy"},
				}),
			},
		}),
	}
	require.Equal(t, expected, actual)
}

func TestUnitIngestz(t *testing.T) {
	api := &ingestzgeturlapi.GetUrlApi{}
	actual := jsparser.ParseApi(api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id: "aaaa",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "xxx", Name: "bbb"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "xxx", Name: "yyyy"},
				}),
			},
		}),
	}
	require.Equal(t, expected, actual)
}
