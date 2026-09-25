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

func TestUnitWrite(t *testing.T) {
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
 * @property {string} dataset
 * @property {ingestzGetUrlReqPayload} payload
 */

/**
 * ingestzGetUrl Response
 * @typedef {object} ingestzGetUrlResp
 * @property {ingestzGetUrlRespPayload} payload
 */

/**
 * @param {import("../services/services.js").Services} services
 * @param {ingestzGetUrlReq} params
 * @returns {Promise<ingestzGetUrlResp>}
 */
export async function ingestzGetUrl(services, params) {
  const resp = await services.invoker().invoke({
    method: "POST",
    path: ` + "`api/ingest/get-url/dataset/${encodeURIComponent(String(params.dataset))}`," + `
    json: params.payload,
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
 * @property {number} denominator
 * @property {number} numerator
 * @property {number} precision
 * @property {string} x_i6_trace_id
 * @property {samplefractionReqPayload} payload
 */

/**
 * samplefraction Response
 * @typedef {object} samplefractionResp
 * @property {string} x_i6_trace_message
 * @property {samplefractionRespPayload} payload
 */

/**
 * @param {import("../services/services.js").Services} services
 * @param {samplefractionReq} params
 * @returns {Promise<samplefractionResp>}
 */
export async function samplefraction(services, params) {
  const resp = await services.invoker().invoke({
    method: "POST",
    path: ` + "`api/gox/routez/sample/fraction/${encodeURIComponent(String(params.numerator))}/${encodeURIComponent(String(params.denominator))}`," + `
    json: params.payload,
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
			Req: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "object", Name: "healthCheckReq"},
				Properties: &jstypez.Tuples{},
			},
			Resp: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "object", Name: "healthCheckResp"},
				Properties: &jstypez.Tuples{},
			},
		}

		want := `/**
 * @typedef {object} healthCheckReq
 */

/**
 * @typedef {object} healthCheckResp
 */

/**
 * @param {import("../services/services.js").Services} services
 * @returns {Promise<healthCheckResp>}
 */
export async function healthCheck(services) {
  const resp = await services.invoker().invoke({
    method: "GET",
    path: ` + "`api/health`," + `
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
			Req: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "object", Name: "getUserReq"},
				Properties: &jstypez.Tuples{},
			},
			Resp: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "object", Name: "getUserResp"},
				Properties: &jstypez.Tuples{},
			},
		}
		api.Req.Properties.Add(jstypez.Tuple{Typ: "string", Name: "id"})
		api.Resp.Properties.Add(jstypez.Tuple{Typ: "getUserRespPayload", Name: "payload"})

		want := `/**
 * @typedef {object} getUserReq
 * @property {string} id
 */

/**
 * @typedef {object} getUserResp
 * @property {getUserRespPayload} payload
 */

/**
 * @param {import("../services/services.js").Services} services
 * @param {getUserReq} params
 * @returns {Promise<getUserResp>}
 */
export async function getUser(services, params) {
  const resp = await services.invoker().invoke({
    method: "GET",
    path: ` + "`api/user/${encodeURIComponent(String(params.id))}`," + `
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
			Req: jstypez.Typedef{
				Head: &jstypez.Tuple{Typ: "object", Name: "pingReq"},
			},
			Resp: jstypez.Typedef{
				Head: &jstypez.Tuple{Typ: "object", Name: "pingResp"},
			},
		}

		want := `/**
 * @typedef {object} pingReq
 */

/**
 * @typedef {object} pingResp
 */

/**
 * @param {import("../services/services.js").Services} services
 * @returns {Promise<pingResp>}
 */
export async function ping(services) {
  const resp = await services.invoker().invoke({
    method: "GET",
    path: ` + "`ping`," + `
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
