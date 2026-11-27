// inspect.js — Request recorder and inspect(x) runtime helper

(function () {
  const originalFetch = window.fetch;
  const recordedRequests = [];
  let lastCheckpointIndex = 0;
  let activeLogEntries = [];
  let methodsCatalog = null;

  // 1. Wrap window.fetch to record all HTTP traffic
  window.fetch = async function (input, init) {
    const startTime = performance.now();
    let method = "GET";
    let url = "";
    let reqHeaders = {};
    let reqBody = null;

    if (typeof input === "string") {
      url = input;
    } else if (input instanceof URL) {
      url = input.href;
    } else if (input && input.url) {
      url = input.url;
      method = input.method || "GET";
    }

    if (init) {
      if (init.method) method = init.method.toUpperCase();
      if (init.headers) {
        if (init.headers instanceof Headers) {
          init.headers.forEach((v, k) => (reqHeaders[k] = v));
        } else if (Array.isArray(init.headers)) {
          init.headers.forEach(([k, v]) => (reqHeaders[k] = v));
        } else {
          reqHeaders = { ...init.headers };
        }
      }
      if (init.body) {
        reqBody = typeof init.body === "string" ? init.body : "[Binary/FormData body]";
      }
    }

    const record = {
      method,
      url,
      reqHeaders,
      reqBody,
      startTime,
      endTime: 0,
      elapsedMs: 0,
      status: 0,
      statusText: "",
      respHeaders: {},
      respBody: null,
      error: null,
    };

    try {
      const response = await originalFetch.apply(window, arguments);
      record.endTime = performance.now();
      record.elapsedMs = Math.round(record.endTime - startTime);
      record.status = response.status;
      record.statusText = response.statusText;

      if (response.headers && response.headers.forEach) {
        response.headers.forEach((v, k) => (record.respHeaders[k] = v));
      }

      // Clone response to read body without consuming it for caller
      try {
        const clone = response.clone();
        record.respBody = await clone.text();
      } catch (e) {
        record.respBody = "[Unable to read response body]";
      }

      recordedRequests.push(record);
      return response;
    } catch (err) {
      record.endTime = performance.now();
      record.elapsedMs = Math.round(record.endTime - startTime);
      record.error = err;
      recordedRequests.push(record);
      throw err;
    }
  };

  // 2. Token masking helper
  function maskTokens(str, showTokens) {
    if (!str || showTokens) return str;
    const HIDDEN = "eyJh...(hidden)";
    return str
      .replace(/(Bearer\s+)[A-Za-z0-9._~+\/=-]+/gi, `$1${HIDDEN}`)
      .replace(/"(access_token|refresh_token|id_token|token)"(\s*):(\s*)"[^"]*"/gi, `"$1"$2:$3"${HIDDEN}"`)
      .replace(/eyJ[A-Za-z0-9_-]{8,}(?:\.[A-Za-z0-9_-]*)*/g, HIDDEN);
  }

  function formatJson(val) {
    try {
      return JSON.stringify(val, null, 2);
    } catch {
      return String(val);
    }
  }

  // Resp (SDK response wrapper) keeps its body in private fields, so
  // JSON.stringify would only show {status, headers, ok}. Duck-typed on purpose.
  function isResp(v) {
    return !!v && typeof v === "object" && typeof v.status === "number" && typeof v.ok === "boolean" && "text" in v;
  }

  function prettyBody(text) {
    if (!text) return "";
    try {
      return JSON.stringify(JSON.parse(text), null, 2);
    } catch {
      return text;
    }
  }

  function formatResp(r) {
    const lines = [`Resp ${r.status} (ok: ${r.ok})`];
    if (!r.ok) {
      let err = null;
      try {
        err = r.error;
      } catch {}
      if (err) lines.push(`error: ${err.code || ""} ${err.message || ""}`.trim());
    }
    const body = prettyBody(r.text);
    if (body) lines.push(body);
    return lines.join("\n");
  }

  function formatValue(v) {
    if (isResp(v)) return formatResp(v);
    if (v === undefined) return "undefined";
    return v && typeof v === "object" ? formatJson(v) : String(v);
  }

  // 3. Methods catalog lazy loader
  async function getMethodsCatalog() {
    if (methodsCatalog) return methodsCatalog;
    try {
      const res = await originalFetch("methods.json");
      if (res.ok) {
        methodsCatalog = await res.json();
      }
    } catch (e) {
      console.warn("Could not load methods.json:", e);
    }
    return methodsCatalog || {};
  }

  function findSdkMethodCallPath(fn) {
    if (!window.sdk || typeof fn !== "function") return null;

    const targets = [
      { prefix: "sdk.apis().", obj: window.sdk.apis?.() },
      { prefix: "sdk.funcs().", obj: window.sdk.funcs?.() },
      { prefix: "sdk.infer().", obj: window.sdk.infer?.() },
      { prefix: "sdk.dash().", obj: window.sdk.dash?.() },
      { prefix: "sdk.auth().", obj: window.sdk.auth?.() },
      { prefix: "sdk.tokens().", obj: window.sdk.tokens?.() },
      { prefix: "sdk.ingest().", obj: window.sdk.ingest?.() },
      { prefix: "sdk.", obj: window.sdk },
    ];

    for (const { prefix, obj } of targets) {
      if (!obj) continue;
      for (const key of Object.getOwnPropertyNames(Object.getPrototypeOf(obj) || obj)) {
        if (typeof obj[key] === "function" && obj[key] === fn) {
          return prefix + key;
        }
      }
    }
    return null;
  }

  function renderOutput() {
    const out = document.getElementById("result-output");
    const panel = document.getElementById("result-panel");
    const showTokens = document.getElementById("chk-show-tokens")?.checked;

    if (!out) return;
    if (panel) panel.style.display = "block";

    const text = activeLogEntries.map((e) => maskTokens(e, showTokens)).join("\n\n");
    out.textContent = text;
  }

  // 4. Main inspect(x) function
  async function inspect(target) {
    const startIndex = lastCheckpointIndex;
    lastCheckpointIndex = recordedRequests.length;

    // A. Function inspection (e.g. inspect(sdk.apis().dszCreateCustomerPortal))
    if (typeof target === "function") {
      const callPath = findSdkMethodCallPath(target);
      if (callPath) {
        const catalog = await getMethodsCatalog();
        const doc = catalog[callPath];
        if (doc) {
          const lines = [
            `Method: ${doc.callPath}`,
            doc.httpMethod ? `${doc.httpMethod} ${doc.httpPath} (token: ${doc.tokenType || "none"})` : "",
            doc.summary ? `\n${doc.summary}` : "",
            `\nDocumentation: ${doc.url}`,
          ].filter(Boolean);
          activeLogEntries.push(lines.join("\n"));
          renderOutput();
          return target;
        }
      }
      activeLogEntries.push(`[Function: ${target.name || "anonymous"}]`);
      renderOutput();
      return target;
    }

    // B. Promise / Call inspection
    if (target && typeof target.then === "function") {
      let resolved = null;
      let rejected = null;

      try {
        resolved = await target;
      } catch (err) {
        rejected = err;
      }

      const calls = recordedRequests.slice(startIndex);
      const outputParts = [];

      for (const call of calls) {
        const reqLines = [
          `> ${call.method} ${call.url}`,
          ...Object.entries(call.reqHeaders).map(([k, v]) => `> ${k}: ${v}`),
        ];
        if (call.reqBody) {
          reqLines.push(`>\n> ${call.reqBody}`);
        }

        const respLines = [
          `< ${call.status} ${call.statusText || ""} (${call.elapsedMs}ms)`,
          ...Object.entries(call.respHeaders).map(([k, v]) => `< ${k}: ${v}`),
        ];
        if (call.error) {
          respLines[0] = `< network error (${call.elapsedMs}ms): ${call.error.message || call.error}`;
        }
        if (call.respBody) {
          respLines.push(`<\n< ${prettyBody(call.respBody).split("\n").join("\n< ")}`);
        }

        outputParts.push(reqLines.join("\n") + "\n\n" + respLines.join("\n"));
      }

      if (rejected) {
        outputParts.push(`\nError: ${rejected.name || "Error"}: ${rejected.message}\n` + (rejected.stack || ""));
        activeLogEntries.push(outputParts.join("\n\n---\n\n"));
        renderOutput();
        throw rejected;
      } else {
        outputParts.push(`\nReturned:\n${formatValue(resolved)}`);
        activeLogEntries.push(outputParts.join("\n\n---\n\n"));
        renderOutput();
        return resolved;
      }
    }

    // C. Error inspection
    if (target instanceof Error) {
      const errText = `Error: ${target.name || "Error"}: ${target.message}` + (target.stack ? `\n${target.stack}` : "");
      activeLogEntries.push(errText);
      renderOutput();
      return target;
    }

    // D. Regular value or Resp
    activeLogEntries.push(formatValue(target));
    renderOutput();
    return target;
  }

  // Hook token toggle
  document.addEventListener("DOMContentLoaded", () => {
    const chk = document.getElementById("chk-show-tokens");
    if (chk) {
      chk.onchange = renderOutput;
    }
  });

  window.inspect = inspect;
  // Used by tryit.js to derive the real HTTP status of a run and to print
  // non-inspect() results through the same token masking.
  window.docsInspect = {
    requestCount: () => recordedRequests.length,
    requestsSince: (i) => recordedRequests.slice(i),
    appendEntry: (text) => {
      activeLogEntries.push(text);
      renderOutput();
    },
    entryCount: () => activeLogEntries.length,
    isResp,
    formatValue,
  };
  window.resetInspectLog = () => {
    activeLogEntries = [];
    lastCheckpointIndex = recordedRequests.length;
    const out = document.getElementById("result-output");
    if (out) out.textContent = "";
  };
})();
