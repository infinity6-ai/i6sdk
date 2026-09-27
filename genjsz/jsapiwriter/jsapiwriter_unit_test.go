package jsapiwriter_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapiwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
)

func TestManualWrite(t *testing.T) {
	type testScenario struct {
		in   *jsparser.Api
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := jsapiwriter.Write(s.in)
		require.Equal(t, s.want, got)
	}

	t.Run("with path param, request payload, and response payload", func(t *testing.T) {
		api := jsparser.ParseApi(&ingestzgeturlapi.GetUrlApi{})
		want := `/**
 * ingestzGetUrl Request
 * @typedef {object} ingestzGetUrlReq
 * @property {ingestzGetUrlparams} params
 * @property {ingestzGetUrlreqPayload} payload
 */

/**
 * ingestzGetUrl Response
 * @typedef {object} ingestzGetUrlResp
 * @property {ingestzGetUrlrespPayload} payload
 */

/**
 * @typedef {object} ingestzGetUrlparams
 * @property {string} dataset
 */

/**
 * @typedef {object} ingestzGetUrlreqPayload
 * @property {number} amount
 * @property {Record<string, string>} partitions
 * @property {string} table
 */

/**
 * @typedef {object} ingestzGetUrlrespPayload
 */

/**
 * @param {import("../services/services.js").Services} services
 * @param {ingestzGetUrlReq} req
 * @returns {Promise<ingestzGetUrlResp>}
 */
export async function ingestzGetUrl(services, req) {
  const resp = await services.invoker().invoke({
    method: "POST",
    path: ` + "`" + `api/ingest/get-url/dataset/${encodeURIComponent(String(req.params.dataset))}` + "`" + `,
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
`
		check(t, testScenario{
			in:   api,
			want: want,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		api := jsparser.ParseApi(routezsamplefraction.NewService().Api())
		want := `/**
 * samplefraction Request
 * @typedef {object} samplefractionReq
 * @property {samplefractionparams} params
 * @property {samplefractionquery} query
 * @property {samplefractionreqHeaders} headers
 * @property {samplefractionreqPayload} payload
 */

/**
 * samplefraction Response
 * @typedef {object} samplefractionResp
 * @property {samplefractionrespHeaders} headers
 * @property {samplefractionrespPayload} payload
 */

/**
 * @typedef {object} samplefractionparams
 * @property {number} denominator
 * @property {number} numerator
 */

/**
 * @typedef {object} samplefractionquery
 * @property {number} precision
 */

/**
 * @typedef {object} samplefractionreqHeaders
 * @property {string} x_i6_trace_id
 */

/**
 * @typedef {object} samplefractionreqPayload
 * @property {string} reason
 */

/**
 * @typedef {object} samplefractionrespHeaders
 * @property {string} x_i6_trace_message
 */

/**
 * @typedef {object} samplefractionrespPayload
 * @property {string} display
 * @property {string} result
 */

/**
 * @param {import("../services/services.js").Services} services
 * @param {samplefractionReq} req
 * @returns {Promise<samplefractionResp>}
 */
export async function samplefraction(services, req) {
  const resp = await services.invoker().invoke({
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
`
		check(t, testScenario{
			in:   api,
			want: want,
		})
	})

	t.Run("without request params and response body fallback", func(t *testing.T) {
		api := &jsparser.Api{
			Spec: jsparser.ApiSpec{
				Id:     "health-check",
				Method: "GET",
				Path:   "/api/health",
			},
		}

		want := `/**
 * healthCheck Response
 * @typedef {object} healthCheckResp
 */

/**
 * @param {import("../services/services.js").Services} services
 * @returns {Promise<healthCheckResp>}
 */
export async function healthCheck(services) {
  const resp = await services.invoker().invoke({
    method: "GET",
    path: ` + "`" + `api/health` + "`" + `,
  });
  const result = {status: resp.status, ok: resp.ok};
  if (resp.ok) {
    result.body = await resp.json();
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
`
		check(t, testScenario{
			in:   api,
			want: want,
		})
	})

	t.Run("with request params but without payload", func(t *testing.T) {
		api := &jsparser.Api{
			Spec: jsparser.ApiSpec{
				Id:     "get-user",
				Method: "GET",
				Path:   "/api/user/{id}",
			},
		}
		api.Req.Params.Head = &jstypez.Tuple{Typ: "getUserParams", Name: "params"}
		api.Resp.Payload.Head = &jstypez.Tuple{Typ: "getUserRespPayload", Name: "payload"}

		want := `/**
 * getUser Request
 * @typedef {object} getUserReq
 * @property {getUserParams} params
 */

/**
 * getUser Response
 * @typedef {object} getUserResp
 * @property {getUserRespPayload} payload
 */

/**
 * @param {import("../services/services.js").Services} services
 * @param {getUserReq} req
 * @returns {Promise<getUserResp>}
 */
export async function getUser(services, req) {
  const resp = await services.invoker().invoke({
    method: "GET",
    path: ` + "`" + `api/user/${encodeURIComponent(String(req.params.id))}` + "`" + `,
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
`
		check(t, testScenario{
			in:   api,
			want: want,
		})
	})

	t.Run("with nil properties", func(t *testing.T) {
		api := &jsparser.Api{
			Spec: jsparser.ApiSpec{
				Id:     "ping",
				Method: "GET",
				Path:   "ping",
			},
		}

		want := `/**
 * ping Response
 * @typedef {object} pingResp
 */

/**
 * @param {import("../services/services.js").Services} services
 * @returns {Promise<pingResp>}
 */
export async function ping(services) {
  const resp = await services.invoker().invoke({
    method: "GET",
    path: ` + "`" + `ping` + "`" + `,
  });
  const result = {status: resp.status, ok: resp.ok};
  if (resp.ok) {
    result.body = await resp.json();
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
`
		check(t, testScenario{
			in:   api,
			want: want,
		})
	})
}
