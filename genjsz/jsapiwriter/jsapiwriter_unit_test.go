package jsapiwriter_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapiwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
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
			in:   &jsparser.Apis{},
			want: "",
		})
	})

	t.Run("with path param, request payload, and response payload", func(t *testing.T) {
		apis := jsparser.ParseApis(ctx, []apiz.Api{&ingestzgeturlapi.GetUrlApi{}})
		want := `/**
 * apiIngestzGetUrl Request
 * @typedef {object} apiIngestzGetUrlReq
 * @property {apiIngestzGetUrlparams} params
 * @property {apiIngestzGetUrlreqPayload} payload
 */

/**
 * apiIngestzGetUrl Response
 * @typedef {object} apiIngestzGetUrlResp
 * @property {apiIngestzGetUrlrespPayload} payload
 */

/**
 * @typedef {object} apiIngestzGetUrlparams
 * @property {string} dataset
 */

/**
 * @typedef {object} apiIngestzGetUrlreqPayload
 * @property {number} amount
 * @property {Record<string, string>} partitions
 * @property {string} table
 */

/**
 * @typedef {object} apiIngestzGetUrlrespPayload
 */

const apis = {
  /**
   * @param {import("../services/services.js").Services} services
   * @param {apiIngestzGetUrlReq} req
   * @returns {Promise<apiIngestzGetUrlResp>}
   */
  async apiIngestzGetUrl(services, req) {
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
  },
};

export { apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := jsparser.ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
		want := `/**
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

const apis = {
  /**
   * @param {import("../services/services.js").Services} services
   * @param {apiSamplefractionReq} req
   * @returns {Promise<apiSamplefractionResp>}
   */
  async apiSamplefraction(services, req) {
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
  },
};

export { apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("without request params and response body fallback", func(t *testing.T) {
		api := &jsparser.Api{
			Spec: jsparser.ApiSpec{
				Id:      "health-check",
				CamelId: "healthCheck",
				Method:  "GET",
				Path:    "/api/health",
			},
		}
		apis := &jsparser.Apis{
			Apis: []*jsparser.Api{api},
		}

		want := `/**
 * healthCheck Response
 * @typedef {object} healthCheckResp
 */

const apis = {
  /**
   * @param {import("../services/services.js").Services} services
   * @returns {Promise<healthCheckResp>}
   */
  async healthCheck(services) {
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
  },
};

export { apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("with request params but without payload", func(t *testing.T) {
		api := &jsparser.Api{
			Spec: jsparser.ApiSpec{
				Id:      "get-user",
				CamelId: "getUser",
				Method:  "GET",
				Path:    "/api/user/{id}",
			},
		}
		api.Req.Params.Head = &jstypez.Tuple{Typ: "getUserParams", Name: "params"}
		api.Resp.Payload.Head = &jstypez.Tuple{Typ: "getUserRespPayload", Name: "payload"}
		apis := &jsparser.Apis{
			Apis: []*jsparser.Api{api},
		}

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

const apis = {
  /**
   * @param {import("../services/services.js").Services} services
   * @param {getUserReq} req
   * @returns {Promise<getUserResp>}
   */
  async getUser(services, req) {
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
  },
};

export { apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("with nil properties", func(t *testing.T) {
		api := &jsparser.Api{
			Spec: jsparser.ApiSpec{
				Id:      "ping",
				CamelId: "ping",
				Method:  "GET",
				Path:    "ping",
			},
		}
		apis := &jsparser.Apis{
			Apis: []*jsparser.Api{api},
		}

		want := `/**
 * ping Response
 * @typedef {object} pingResp
 */

const apis = {
  /**
   * @param {import("../services/services.js").Services} services
   * @returns {Promise<pingResp>}
   */
  async ping(services) {
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
  },
};

export { apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("with multiple apis: all typedefs first then all api declarations", func(t *testing.T) {
		apis := jsparser.ParseApis(ctx, []apiz.Api{
			routezsamplefraction.NewService().Api(),
			&ingestzgeturlapi.GetUrlApi{},
		})

		want := `/**
 * apiIngestzGetUrl Request
 * @typedef {object} apiIngestzGetUrlReq
 * @property {apiIngestzGetUrlparams} params
 * @property {apiIngestzGetUrlreqPayload} payload
 */

/**
 * apiIngestzGetUrl Response
 * @typedef {object} apiIngestzGetUrlResp
 * @property {apiIngestzGetUrlrespPayload} payload
 */

/**
 * @typedef {object} apiIngestzGetUrlparams
 * @property {string} dataset
 */

/**
 * @typedef {object} apiIngestzGetUrlreqPayload
 * @property {number} amount
 * @property {Record<string, string>} partitions
 * @property {string} table
 */

/**
 * @typedef {object} apiIngestzGetUrlrespPayload
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

const apis = {
  /**
   * @param {import("../services/services.js").Services} services
   * @param {apiSamplefractionReq} req
   * @returns {Promise<apiSamplefractionResp>}
   */
  async apiSamplefraction(services, req) {
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
  },

  /**
   * @param {import("../services/services.js").Services} services
   * @param {apiIngestzGetUrlReq} req
   * @returns {Promise<apiIngestzGetUrlResp>}
   */
  async apiIngestzGetUrl(services, req) {
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
  },
};

export { apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})
}
