package dszapisignedurl_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapisignedurl"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakepathz"
)

func TestUnitApiSpec(t *testing.T) {
	api := &dszapisignedurl.Api{}
	spec := api.ApiSpec()

	require.Equal(t, "dsz-signed-url", spec.Id)
	require.Equal(t, "POST", spec.Method)
	require.Equal(t, "/api/dsz/signed-url/dataset/{dataset}", spec.Path)

	desc := spec.Desc()
	require.NotNil(t, desc)
	require.Equal(t, "Get Dataset Signed URLs", desc.Name)
	require.Equal(t, "Generates presigned URLs for accessing or mutating files in the specified dataset.", desc.Summary)
}

func TestUnitGetDataRefs(t *testing.T) {
	api := &dszapisignedurl.Api{
		Req: dszapisignedurl.Req{
			Dataset: lakedsz.MustParse("sample-dataset"),
			Paths: []dszapisignedurl.Path{
				{
					Method: "GET",
					Path:   lakepathz.MustParse("path/to/file.parquet"),
				},
			},
		},
	}

	dataRefs := api.GetDataRefs()
	require.NotNil(t, dataRefs)
	require.NotNil(t, dataRefs.PathParams)
	require.NotNil(t, dataRefs.ReqBody)
	require.NotNil(t, dataRefs.RespBody)

	pathParams := dataRefs.PathParams.Object(true)
	require.Contains(t, pathParams, "dataset")
	require.Equal(t, "Name of the dataset to generate signed URLs for.", pathParams["dataset"].Desc.Summary)

	reqBody := dataRefs.ReqBody.Object(true)
	require.Contains(t, reqBody, "paths")
	require.Equal(t, "List of paths and HTTP methods to generate signed URLs for.", reqBody["paths"].Desc.Summary)

	pathsArr := reqBody["paths"].Array()
	itemSchema := pathsArr.Get(0, true)
	require.Equal(t, "Path and HTTP method specification for a signed URL.", itemSchema.Desc.Summary)

	itemObj := itemSchema.Object(true)
	require.Equal(t, "HTTP method (e.g. GET, PUT) for the signed URL.", itemObj["method"].Desc.Summary)
	require.Equal(t, "File path within the dataset.", itemObj["path"].Desc.Summary)

	require.Equal(t, "Presigned URLs generated for the requested paths.", dataRefs.RespBody.Desc.Summary)
}
