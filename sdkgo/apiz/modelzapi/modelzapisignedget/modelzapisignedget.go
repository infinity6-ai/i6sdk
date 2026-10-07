package modelzapisignedget

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type Req struct {
	Dataset         lakedsz.Dataset   `json:"dataset"`          // PathParam
	ModelName       string            `json:"model_name"`       // PathParam
	ModelVersion    string            `json:"model_version"`    // Body
	ModelPartitions map[string]string `json:"model_partitions"` // Body
}

type ModelFile struct {
	Name string `json:"name"` // Resp Body
	Url  string `json:"url"`  // Resp Body
}

func (u *ModelFile) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"name": {
					Desc: &schemaz.Desc{Summary: "File Name"},
					Raw:  func() any { return &u.Name },
				},
				"url": {
					Desc: &schemaz.Desc{Summary: "signed URL."},
					Raw:  func() any { return &u.Url },
				},
			}
		},
	}
}

type Resp struct {
	Urls []*ModelFile `json:"urls"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "modelz-signed-get",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Get Dataset Signed URLs",
				Summary: "Generates presigned URLs for accessing model files in the specified dataset.",
			}
		},
		Method: "POST",
		Path:   "/api/model/dataset/{dataset}/signed-get/{model_name}",
	}
}

func (d *Api) GetDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc: &schemaz.Desc{Summary: "Name of the dataset to generate signed URLs for."},
						Str:  lakedsz.ParseStr(&d.Req.Dataset),
					},
					"model_name": {
						Desc: &schemaz.Desc{Summary: "Model name"},
						Str:  schemaz.ParseStr(&d.Req.ModelName),
					},
				}
			},
		},
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				panic("implement it")
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Presigned URLs generated for the requested paths."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"urls": {
						Desc: &schemaz.Desc{Summary: "Signed URLs"},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]*ModelFile, *ModelFile]{
							NewElement: func(idx int) *ModelFile { return &ModelFile{} },
							Slice:      func() *[]*ModelFile { return &d.Resp.Urls },
							Get:        func(idx int, element **ModelFile) *schemaz.Schema { return (*element).Schema() },
						}),
					},
				}
			},
		},
	}
}
