package jsapiwriter_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapiwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   *jsparser.Apis
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := jsapiwriter.Write(ctx, s.in)
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
			in: &jsparser.Apis{},
			want: `/**
 * @typedef {import("../services/services.js").Services} Services
 */

`,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := jsparser.ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
		want := `/**
 * @typedef {import("../services/services.js").Services} Services
 */

/**
 * apiSampleFraction Request
 * @typedef {object} apiSampleFractionReq
 * @property {apiSampleFractionparams} params
 * @property {apiSampleFractionquery} query
 * @property {apiSampleFractionreqHeaders} headers
 * @property {apiSampleFractionreqPayload} payload
 */

/**
 * apiSampleFraction Response
 * @typedef {object} apiSampleFractionResp
 * @property {apiSampleFractionrespHeaders} headers
 * @property {apiSampleFractionrespPayload} payload
 */

/**
 * @typedef {object} apiSampleFractionparams
 * @property {number} denominator - Fraction Denominator
 * @property {number} numerator - Fraction Numerator
 */

/**
 * @typedef {object} apiSampleFractionquery
 * @property {number} precision - Fraction max precision
 */

/**
 * @typedef {object} apiSampleFractionreqHeaders
 * @property {string} x_i6_trace_id - Sample of Trace Id
 */

/**
 * @typedef {object} apiSampleFractionreqPayload
 * @property {string} reason - Sample reason of this fraction
 */

/**
 * @typedef {object} apiSampleFractionrespHeaders
 * @property {string} x_i6_trace_message - Sample of any trace message
 */

/**
 * @typedef {object} apiSampleFractionrespPayload
 * @property {string} display - Human Fraction Representation
 * @property {string} result - Fraction Result
 */

class Apis {

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
