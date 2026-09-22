package generatorjsz_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/infinity6-ai/gox/commonz/slicez"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/generatorjsz"
)

type capturedRequest struct {
	rawQuery string
	header   http.Header
	body     string
}

func serveCapture(t *testing.T) (baseUrl string, got func() capturedRequest) {
	t.Helper()
	var mu sync.Mutex
	var last capturedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		last = capturedRequest{rawQuery: r.URL.RawQuery, header: r.Header.Clone(), body: string(body)}
		mu.Unlock()
	}))
	t.Cleanup(srv.Close)
	return srv.URL, func() capturedRequest {
		mu.Lock()
		defer mu.Unlock()
		return last
	}
}

func rawField(ptr any) *schemaz.Schema {
	return &schemaz.Schema{Raw: func() any { return ptr }}
}

func strField(native any) *schemaz.Schema {
	return &schemaz.Schema{Str: func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse:  func(string) {},
			Format: func() (string, any) { return "", native },
		}
	}}
}

func strsField(native any) *schemaz.Schema {
	return &schemaz.Schema{Strs: func() *schemaz.Parser[[]string] {
		return &schemaz.Parser[[]string]{
			Parse:  func([]string) {},
			Format: func() ([]string, any) { return nil, native },
		}
	}}
}

type wireTestScenario struct {
	refs        *apiz.DataRefs
	params      map[string]any
	wantCode    []string
	wantQuery   string
	wantHeaders map[string]string
	wantBody    string // JSON; "" quando não deve haver corpo
}

func checkJsClientWireFormat(t *testing.T, s wireTestScenario) {
	t.Helper()
	baseUrl, got := serveCapture(t)
	api := &fakeApi{spec: apiz.ApiSpec{Id: "wire", Method: "POST", Path: "/api/x"}, refs: s.refs}

	fnName, code, err := generatorjsz.BuildJsFetch(api.ApiSpec(), api.GetDataRefs())
	require.NoError(t, err)
	for _, want := range s.wantCode {
		require.Contains(t, code, want)
	}

	result := runJs(t, fnName, code, invokerPath, baseUrl, s.params)
	require.Equal(t, true, result["ok"], "the capture server answers 200")

	req := got()
	require.Equal(t, s.wantQuery, req.rawQuery)
	for name, want := range s.wantHeaders {
		require.Equal(t, want, req.header.Get(name), "header %s", name)
	}
	if s.wantBody == "" {
		require.Empty(t, req.body)
	} else {
		require.JSONEq(t, s.wantBody, req.body)
	}
}

func TestUnitJsClientWireFormatQueryStrsSliceBecomesRepeatedKeys(t *testing.T) {
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			QueryParams: objectSchema(map[string]*schemaz.Schema{"ids": strsField([]int{})}),
		},
		params:    map[string]any{"ids": []int{1, 2, 3}},
		wantCode:  []string{"@property {number[]} ids"},
		wantQuery: "ids=1&ids=2&ids=3",
	})
}

func TestUnitJsClientWireFormatQueryRawSliceBecomesRepeatedKeys(t *testing.T) {
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			QueryParams: objectSchema(map[string]*schemaz.Schema{"tags": rawField(&[]string{})}),
		},
		params:    map[string]any{"tags": []string{"a", "b"}},
		wantCode:  []string{"@property {string[]} tags"},
		wantQuery: "tags=a&tags=b",
	})
}

func TestUnitJsClientWireFormatQueryScalarSentAndOmittedParamIsNot(t *testing.T) {
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			QueryParams: objectSchema(map[string]*schemaz.Schema{
				"page": strField(0),
				"q":    strField(""),
			}),
		},
		params:    map[string]any{"page": 2}, // q left out: must not become q=undefined
		wantCode:  []string{"@property {number} page", "@property {string} q"},
		wantQuery: "page=2",
	})
}

func TestUnitJsClientWireFormatHeaderStrsJoinedAndStrIsSingleValue(t *testing.T) {
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			ReqHeaders: objectSchema(map[string]*schemaz.Schema{
				"x_tags":     strsField([]string{}),
				"x_trace_id": strField(""),
			}),
		},
		params:   map[string]any{"x_tags": []string{"a", "b"}, "x_trace_id": "t1"},
		wantCode: []string{"@property {string[]} x_tags", `"x-tags": params.x_tags.map(String).join(", ")`},
		wantHeaders: map[string]string{
			"x-tags":     "a, b",
			"x-trace-id": "t1",
		},
	})
}

func TestUnitJsClientWireFormatBodyStrAndStrsTravelAsJsonArraysOfStrings(t *testing.T) {
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			ReqBody: objectSchema(map[string]*schemaz.Schema{
				"a": strField(""),          // Str -> ["x"]
				"b": strsField([]int{}),    // Strs with slice native -> ["1","2"]
				"c": strsField(float64(0)), // Strs with scalar native -> ["30.5"], not an array on the JS side
			}),
		},
		params:      map[string]any{"a": "x", "b": []int{1, 2}, "c": 30.5},
		wantCode:    []string{"@property {string} a", "@property {number[]} b", "@property {number} c"},
		wantHeaders: map[string]string{"Content-Type": "application/json"},
		wantBody:    `{"a":["x"],"b":["1","2"],"c":["30.5"]}`,
	})
}

func TestUnitJsClientWireFormatBodyRawValuesTravelUntouched(t *testing.T) {
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			ReqBody: objectSchema(map[string]*schemaz.Schema{
				"nums": rawField(&[]int{}),
				"flag": rawField(new(bool)),
				"name": rawField(new(string)),
			}),
		},
		params:   map[string]any{"nums": []int{1, 2}, "flag": true, "name": "n"},
		wantCode: []string{"@property {number[]} nums", "@property {boolean} flag", "@property {string} name"},
		wantBody: `{"nums":[1,2],"flag":true,"name":"n"}`,
	})
}

func TestUnitJsClientWireFormatBodyNestedObjectAndArrayOfObjectsTravelUntouched(t *testing.T) {
	var street string
	var items []struct{ Name string }
	checkJsClientWireFormat(t, wireTestScenario{
		refs: &apiz.DataRefs{
			ReqBody: objectSchema(map[string]*schemaz.Schema{
				"address": {
					Object: func(read bool) map[string]*schemaz.Schema {
						return map[string]*schemaz.Schema{"street": rawField(&street)}
					},
				},
				"items": {
					Array: func() *schemaz.Array {
						return &schemaz.Array{
							Len: len(items),
							Get: func(idx int, read bool) *schemaz.Schema {
								items = slicez.GrowLenTo(items, idx+1)
								return &schemaz.Schema{
									Object: func(read bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{"name": rawField(&items[idx].Name)}
									},
								}
							},
						}
					},
				},
			}),
		},
		params: map[string]any{
			"address": map[string]any{"street": "s"},
			"items":   []map[string]any{{"name": "a"}, {"name": "b"}},
		},
		wantCode: []string{
			"@property {{street: string}} address",
			"@property {{name: string}[]} items",
		},
		wantBody: `{"address":{"street":"s"},"items":[{"name":"a"},{"name":"b"}]}`,
	})
}
