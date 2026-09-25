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
				"method": "POST",
				"path": "/api/gox/routez/sample/fraction/{numerator}/{denominator}"
			},
			"req": {
				"head": {
				"typ": "object",
				"name": "samplefractionReq",
				"summary": "samplefraction Request"
				},
				"properties": [
				{
					"typ": "number",
					"name": "numerator",
					"summary": ""
				},
				{
					"typ": "number",
					"name": "denominator",
					"summary": ""
				},
				{
					"typ": "number",
					"name": "precision",
					"summary": ""
				},
				{
					"typ": "string",
					"name": "x_i6_trace_id",
					"summary": ""
				},
				{
					"typ": "samplefractionReqPayload",
					"name": "payload",
					"summary": ""
				}
				]
			},
			"resp": {
				"head": {
				"typ": "object",
				"name": "samplefractionResp",
				"summary": "samplefraction Response"
				},
				"properties": [
				{
					"typ": "string",
					"name": "x_i6_trace_message",
					"summary": ""
				},
				{
					"typ": "samplefractionRespPayload",
					"name": "payload",
					"summary": ""
				}
				]
			},
			"typedefs": {
				"samplefractionReqPayload": {
				"head": {
					"typ": "object",
					"name": "samplefractionReqPayload",
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
				"samplefractionRespPayload": {
				"head": {
					"typ": "object",
					"name": "samplefractionRespPayload",
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
				"method": "POST",
				"path": "/api/ingest/get-url/dataset/{dataset}"
			},
			"req": {
				"head": {
				"typ": "object",
				"name": "ingestzGetUrlReq",
				"summary": "ingestzGetUrl Request"
				},
				"properties": [
				{
					"typ": "string",
					"name": "dataset",
					"summary": ""
				},
				{
					"typ": "ingestzGetUrlReqPayload",
					"name": "payload",
					"summary": ""
				}
				]
			},
			"resp": {
				"head": {
				"typ": "object",
				"name": "ingestzGetUrlResp",
				"summary": "ingestzGetUrl Response"
				},
				"properties": [
				{
					"typ": "ingestzGetUrlRespPayload",
					"name": "payload",
					"summary": ""
				}
				]
			},
			"typedefs": {
				"ingestzGetUrlReqPayload": {
				"head": {
					"typ": "object",
					"name": "ingestzGetUrlReqPayload",
					"summary": ""
				},
				"properties": [
					{
					"typ": "string",
					"name": "table",
					"summary": ""
					},
					{
					"typ": "Record\u003cstring, string\u003e",
					"name": "partitions",
					"summary": ""
					},
					{
					"typ": "number",
					"name": "amount",
					"summary": ""
					}
				]
				},
				"ingestzGetUrlRespPayload": {
				"head": {
					"typ": "object",
					"name": "ingestzGetUrlRespPayload",
					"summary": ""
				},
				"properties": null
				}
			}
		}
	`, jsonz.MustFormat(jsparser.ParseApi(api)).String())

}
