package scheduleractivatemain_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apipipez/scheduleractivatemain"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/pipez/relevancefashionpipez"
	"github.com/stretchr/testify/require"
)

func TestUnitSchedulerActivateMainApi(t *testing.T) {
	api := &scheduleractivatemain.Api{
		Req: &scheduleractivatemain.Req{
			Params: &scheduleractivatemain.ReqParams{
				Dataset: lakedsz.MustParse("my-dataset"),
			},
			Payload: &relevancefashionpipez.Main{},
		},
	}

	apiCopy := jsonz.MustCopy(api, &scheduleractivatemain.Api{})
	require.Equal(t, api.Req.Params.Dataset.Get(), apiCopy.Req.Params.Dataset.Get())

	spec := api.ApiSpec()
	require.Equal(t, "pipez-scheduler-activate-relevance-fashion-relevance-fashion-main", spec.Id)
	require.Equal(t, "/api/pipe/dataset/{dataset}/scheduler/activate/relevance-fashion/relevance-fashion-main", spec.Path)
}
