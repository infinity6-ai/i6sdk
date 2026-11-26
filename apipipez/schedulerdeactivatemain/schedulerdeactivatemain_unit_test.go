package schedulerdeactivatemain_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/schedulerdeactivatemain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/stretchr/testify/require"
)

func TestUnitSchedulerDeactivateMainApi(t *testing.T) {
	api := &schedulerdeactivatemain.Api{
		Req: &schedulerdeactivatemain.Req{
			Params: &schedulerdeactivatemain.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &relevancefashionpipez.Main{},
		},
	}

	apiCopy := jsonz.MustCopy(api, &schedulerdeactivatemain.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())

	spec := api.ApiSpec()
	require.Equal(t, "pipez-scheduler-deactivate-relevance-fashion-relevance-fashion-main", spec.Id)
	require.Equal(t, "/api/pipe/dataset/{dataset}/scheduler/deactivate/relevance-fashion/relevance-fashion-main", spec.Path)
}
