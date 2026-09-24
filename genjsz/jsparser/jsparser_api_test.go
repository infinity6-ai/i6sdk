package jsparser_test

import (
	"fmt"
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
)

func TestManualFractionSample(t *testing.T) {
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
			"js": "object",
			"fields": {
			"denominator": {
				"js": "float"
			},
			"numerator": {
				"js": "float"
			},
			"payload": {
				"js": "object",
				"fields": {
				"reason": {
					"js": "string"
				}
				}
			},
			"precision": {
				"js": "int"
			},
			"x_i6_trace_id": {
				"js": "string"
			}
			}
		},
		"resp": {
			"js": "object",
			"fields": {
			"payload": {
				"js": "object",
				"fields": {
				"display": {
					"js": "string"
				},
				"result": {
					"js": "string"
				}
				}
			},
			"x_i6_trace_message": {
				"js": "string"
			}
			}
		}
		}
	`, jsonz.MustFormat(jsparser.ParseApi(api)).String())

}
