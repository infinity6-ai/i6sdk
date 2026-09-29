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
 * apiSamplefraction Request
 * @typedef {object} apiSamplefractionReq
 * @property {apiSamplefractionparams} params
 * @property {apiSamplefractionquery} query
 * @property {apiSamplefractionreqHeaders} headers
 * @property {apiSamplefractionreqPayload} payload
 */

/**
 * apiSamplefraction Response
 * @typedef {object} apiSamplefractionResp
 * @property {apiSamplefractionrespHeaders} headers
 * @property {apiSamplefractionrespPayload} payload
 */

/**
 * @typedef {object} apiSamplefractionparams
 * @property {number} denominator
 * @property {number} numerator
 */

/**
 * @typedef {object} apiSamplefractionquery
 * @property {number} precision
 */

/**
 * @typedef {object} apiSamplefractionreqHeaders
 * @property {string} x_i6_trace_id
 */

/**
 * @typedef {object} apiSamplefractionreqPayload
 * @property {string} reason
 */

/**
 * @typedef {object} apiSamplefractionrespHeaders
 * @property {string} x_i6_trace_message
 */

/**
 * @typedef {object} apiSamplefractionrespPayload
 * @property {string} display
 * @property {string} result
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
   * @param {apiSamplefractionReq} req
   * @returns {Promise<apiSamplefractionResp>}
   */
  async samplefraction(req) {
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
