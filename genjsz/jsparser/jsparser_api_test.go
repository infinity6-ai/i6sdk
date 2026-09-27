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
	require.JSONEq(t, `
		{
			"spec": {
				"id": "samplefraction",
				"camel_id": "samplefraction",
				"method": "POST",
				"path": "/api/gox/routez/sample/fraction/{numerator}/{denominator}"
			},
			"typedefs": {
				"samplefractionparams": {
					"head": {
						"typ": "object",
						"name": "samplefractionparams",
						"summary": ""
					},
					"properties": [
						{
							"typ": "number",
							"name": "denominator",
							"summary": ""
						},
						{
							"typ": "number",
							"name": "numerator",
							"summary": ""
						}
					]
				},
				"samplefractionquery": {
					"head": {
						"typ": "object",
						"name": "samplefractionquery",
						"summary": ""
					},
					"properties": [
						{
							"typ": "number",
							"name": "precision",
							"summary": ""
						}
					]
				},
				"samplefractionreqHeaders": {
					"head": {
						"typ": "object",
						"name": "samplefractionreqHeaders",
						"summary": ""
					},
					"properties": [
						{
							"typ": "string",
							"name": "x_i6_trace_id",
							"summary": ""
						}
					]
				},
				"samplefractionreqPayload": {
					"head": {
						"typ": "object",
						"name": "samplefractionreqPayload",
						"summary": ""
					},
					"properties": [
						{
							"typ": "string",
							"name": "reason",
							"summary": ""
						}
					]
				},
				"samplefractionrespHeaders": {
					"head": {
						"typ": "object",
						"name": "samplefractionrespHeaders",
						"summary": ""
					},
					"properties": [
						{
							"typ": "string",
							"name": "x_i6_trace_message",
							"summary": ""
						}
					]
				},
				"samplefractionrespPayload": {
					"head": {
						"typ": "object",
						"name": "samplefractionrespPayload",
						"summary": ""
					},
					"properties": [
						{
							"typ": "string",
							"name": "display",
							"summary": ""
						},
						{
							"typ": "string",
							"name": "result",
							"summary": ""
						}
					]
				}
			},
			"req": {
				"params": {
					"head": {
						"typ": "samplefractionparams",
						"name": "params",
						"summary": ""
					},
					"properties": []
				},
				"query": {
					"head": {
						"typ": "samplefractionquery",
						"name": "query",
						"summary": ""
					},
					"properties": []
				},
				"headers": {
					"head": {
						"typ": "samplefractionreqHeaders",
						"name": "reqHeaders",
						"summary": ""
					},
					"properties": []
				},
				"payload": {
					"head": {
						"typ": "samplefractionreqPayload",
						"name": "reqPayload",
						"summary": ""
					},
					"properties": []
				}
			},
			"resp": {
				"headers": {
					"head": {
						"typ": "samplefractionrespHeaders",
						"name": "respHeaders",
						"summary": ""
					},
					"properties": []
				},
				"payload": {
					"head": {
						"typ": "samplefractionrespPayload",
						"name": "respPayload",
						"summary": ""
					},
					"properties": []
				}
			}
		}
	`, jsonz.MustFormat(jsparser.ParseApi(api)).String())
}

func TestUnitIngestz(t *testing.T) {
	api := &ingestzgeturlapi.GetUrlApi{}
	fmt.Printf("%s\n", jsonz.MustFormatPretty(jsparser.ParseApi(api)).String())
	require.JSONEq(t, `
		{
			"spec": {
				"id": "ingestz-get-url",
				"camel_id": "ingestzGetUrl",
				"method": "POST",
				"path": "/api/ingest/get-url/dataset/{dataset}"
			},
			"typedefs": {
				"ingestzGetUrlparams": {
					"head": {
						"typ": "object",
						"name": "ingestzGetUrlparams",
						"summary": ""
					},
					"properties": [
						{
							"typ": "string",
							"name": "dataset",
							"summary": ""
						}
					]
				},
				"ingestzGetUrlreqPayload": {
					"head": {
						"typ": "object",
						"name": "ingestzGetUrlreqPayload",
						"summary": ""
					},
					"properties": [
						{
							"typ": "number",
							"name": "amount",
							"summary": ""
						},
						{
							"typ": "Record<string, string>",
							"name": "partitions",
							"summary": ""
						},
						{
							"typ": "string",
							"name": "table",
							"summary": ""
						}
					]
				},
				"ingestzGetUrlrespPayload": {
					"head": {
						"typ": "object",
						"name": "ingestzGetUrlrespPayload",
						"summary": ""
					},
					"properties": null
				}
			},
			"req": {
				"params": {
					"head": {
						"typ": "ingestzGetUrlparams",
						"name": "params",
						"summary": ""
					},
					"properties": []
				},
				"query": {
					"head": null,
					"properties": null
				},
				"headers": {
					"head": null,
					"properties": null
				},
				"payload": {
					"head": {
						"typ": "ingestzGetUrlreqPayload",
						"name": "reqPayload",
						"summary": ""
					},
					"properties": []
				}
			},
			"resp": {
				"headers": {
					"head": null,
					"properties": null
				},
				"payload": {
					"head": {
						"typ": "ingestzGetUrlrespPayload",
						"name": "respPayload",
						"summary": ""
					},
					"properties": []
				}
			}
		}
	`, jsonz.MustFormat(jsparser.ParseApi(api)).String())
}
