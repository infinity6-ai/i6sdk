package ingestzapistream_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzapistream"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbnamez"
)

type MockSdkObject struct {
	IdValue string `json:"id"`
	Name    string `json:"name"`
}

func (m *MockSdkObject) Table() laketbnamez.TableName {
	return laketbnamez.MustParse("mock-object")
}

func (m *MockSdkObject) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Raw: func() any { return &m.IdValue },
				},
				"name": {
					Raw: func() any { return &m.Name },
				},
			}
		},
	}
}

func TestUnitIngestzStreamApi(t *testing.T) {
	api := &ingestzapistream.Api[*MockSdkObject]{
		Req: ingestzapistream.StreamReq[*MockSdkObject]{
			Dataset: lakedsz.MustParse("my-dataset"),
			Data: []*MockSdkObject{
				{
					IdValue: "obj1",
					Name:    "Object 1",
				},
			},
		},
	}

	// Verify standard JSON copy / two-way binding
	apiCopy := jsonz.MustCopy(api, &ingestzapistream.Api[*MockSdkObject]{})
	require.Equal(t, api.Req.Dataset.Get(), apiCopy.Req.Dataset.Get())
	require.Len(t, apiCopy.Req.Data, 1)
	require.Equal(t, api.Req.Data[0].IdValue, apiCopy.Req.Data[0].IdValue)
	require.Equal(t, api.Req.Data[0].Name, apiCopy.Req.Data[0].Name)

	// Verify API specifications
	spec := api.ApiSpec()
	require.Equal(t, "ingestz-stream-mock-object", spec.Id)
	require.Equal(t, "/api/ingest/dataset/{dataset}/stream/mock-object", spec.Path)

	// Verify Schema parsing and mapping
	parsed := jsparserapi.ParseApi(api)
	require.NotNil(t, parsed)
	require.Equal(t, "ingestz-stream-mock-object", parsed.Spec.Id)
}
