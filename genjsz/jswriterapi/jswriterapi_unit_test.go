package jswriterapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jswriterapi"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   *jsparserapi.Apis
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var apisList []*jsparserapi.Api
		if s.in != nil {
			apisList = s.in.Apis
		}
		got := jswriterapi.WriteV2(ctx, apisList)
		println(got)

		require.Equal(t, s.want, got)
	}

	t.Run("nil apis", func(t *testing.T) {
		check(t, testScenario{
			in:   nil,
			want: "",
		})
	})

	t.Run("empty apis", func(t *testing.T) {
		check(t, testScenario{
			in: &jsparserapi.Apis{},
			want: ``,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := jsparserapi.ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
		want := `class Apis {

  /**
   * Constructs an Ingest client instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    /**
     * @private
     * @type {Services}
     */
    this._services = services;
  }

  /**
   * Sample Fraction
   * 
   * Sample API that performs a fraction operation
   * 
   * @param {apiSampleFractionReq} req
   * @returns {Promise<apiSampleFractionResp>}
   */
  async sampleFraction(req) {
    const resp = await this._services.invoker().invoke({
      method: "POST",
      path: ` + "`" + `api/gox/routez/sample/fraction/${encodeURIComponent(String(req.params.numerator))}/${encodeURIComponent(String(req.params.denominator))}` + "`" + `,
      query: req.query,
      headers: req.headers,
      json: req.payload,
    });
    const result = {status: resp.status, ok: resp.ok};
    if (resp.ok) {
      result.payload = await resp.json();
    }
    if (!resp.ok) {
      const text = await resp.text();
      try {
        result.error = JSON.parse(text);
      } catch {
        result.error = text;
      }
    }
    return result;
  }
}

export { Apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

}
