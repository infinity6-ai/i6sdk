package pipezapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/pipezapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/domain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
)

func TestUnitPipezApi(t *testing.T) {
	api := &pipezapi.PipezStartApi{
		Req: &pipezapi.PipezStartReq{
			Dataset: lakedsz.MustParse("myds"),
			SolutionDomain: solutiondomain.MustParse(solutiondomain.Value{
				Solution: solution.MustParse("raas"),
				Domain:   domain.MustParse("fashion"),
			}),
		},
	}

	api2 := jsonz.MustCopy(api, &pipezapi.PipezStartApi{})

	require.Equal(t, *api.Req, *api2.Req)
}
