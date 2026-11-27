// tryit.js — Monaco editor loader, code runner, and autocomplete

(function () {
  const MONACO_CDN_BASE = "https://cdn.jsdelivr.net/npm/monaco-editor@0.52.2/min/vs";
  let monacoLoadingPromise = null;

  // Docs root (where types.json lives), from this script's own URL so it
  // is right on every page depth.
  const DOCS_ROOT_URL = new URL(".", (document.currentScript && document.currentScript.src) || window.location.href).href;

  function loadMonaco() {
    if (monacoLoadingPromise) return monacoLoadingPromise;

    monacoLoadingPromise = new Promise((resolve, reject) => {
      if (window.monaco) {
        return resolve(window.monaco);
      }

      // Browsers refuse to start a Worker from a cross-origin (CDN) URL; Monaco
      // then runs its TypeScript worker on the main thread, where completions
      // like `sdk.` do not work. Start the worker from a same-origin blob that
      // loads the CDN worker code.
      window.MonacoEnvironment = {
        getWorkerUrl() {
          const src =
            `self.MonacoEnvironment = { baseUrl: "${MONACO_CDN_BASE}/../" };\n` +
            `importScripts("${MONACO_CDN_BASE}/base/worker/workerMain.js");`;
          return URL.createObjectURL(new Blob([src], { type: "text/javascript" }));
        },
      };

      const script = document.createElement("script");
      script.src = `${MONACO_CDN_BASE}/loader.js`;
      script.onload = () => {
        window.require.config({ paths: { vs: MONACO_CDN_BASE } });
        window.require(["vs/editor/editor.main"], async () => {
          try {
            // types.json holds the SDK .d.ts files ({path: content}). Each is
            // added as if installed in node_modules/i6sdk-web, so their
            // relative imports and import("i6sdk-web") resolve.
            const res = await fetch(new URL("types.json", DOCS_ROOT_URL));
            if (res.ok) {
              for (const [p, content] of Object.entries(await res.json())) {
                monaco.languages.typescript.javascriptDefaults.addExtraLib(content, `file:///node_modules/i6sdk-web/${p}`);
              }
            } else {
              console.warn("types.json not found, no autocomplete:", res.status);
            }

            // Globals in Try It scope (the runner passes sdk and inspect).
            monaco.languages.typescript.javascriptDefaults.addExtraLib(
              [
                'declare const sdk: import("i6sdk-web").I6Sdk;',
                'declare const I6Sdk: typeof import("i6sdk-web").I6Sdk;',
                "declare function inspect<T>(x: T): Promise<Awaited<T>>;",
              ].join("\n") + "\n",
              "file:///globals.d.ts"
            );

            monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
              noSemanticValidation: false,
              noSyntaxValidation: false,
              // The snippet runs as an async function body: top-level
              // return/await are fine there.
              diagnosticCodesToIgnore: [1108, 1375, 1378],
            });

            monaco.languages.typescript.javascriptDefaults.setCompilerOptions({
              target: monaco.languages.typescript.ScriptTarget.ESNext,
              allowNonTsExtensions: true,
              checkJs: true,
              moduleResolution: monaco.languages.typescript.ModuleResolutionKind.NodeJs,
            });

            resolve(window.monaco);
          } catch (e) {
            console.warn("Error configuring Monaco TypeScript libs:", e);
            resolve(window.monaco);
          }
        });
      };
      script.onerror = reject;
      document.head.appendChild(script);
    });

    return monacoLoadingPromise;
  }

  async function initEditor() {
    const container = document.getElementById("editor-container");
    if (!container) return;

    const defaultCode = container.dataset.defaultCode || container.textContent.trim();
    const fallbackPre = container.querySelector(".fallback-code");

    try {
      const monaco = await loadMonaco();
      if (fallbackPre) fallbackPre.style.display = "none";

      const isDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
      window.monacoEditor = monaco.editor.create(container, {
        value: defaultCode,
        language: "javascript",
        theme: isDark ? "vs-dark" : "vs",
        automaticLayout: true,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        fontSize: 13,
        lineNumbers: "on",
        // The code box clips its children (rounded corners); render the
        // suggest/hover widgets as fixed so they can escape it.
        fixedOverflowWidgets: true,
        padding: { top: 8, bottom: 8 },
        scrollbar: { alwaysConsumeMouseWheel: false },
      });

      // Grow the editor with its content so snippets never scroll inside a
      // tiny box and the suggest widget has room.
      const MIN_H = 200;
      const MAX_H = 600;
      const fit = () => {
        const h = Math.min(MAX_H, Math.max(MIN_H, window.monacoEditor.getContentHeight()));
        if (container.style.height !== `${h}px`) container.style.height = `${h}px`;
      };
      window.monacoEditor.onDidContentSizeChange(fit);
      fit();
    } catch (err) {
      console.warn("Failed to load Monaco editor, using static fallback:", err);
    }

    // Bind Run button
    const runBtn = document.getElementById("btn-run-code");
    if (runBtn) {
      runBtn.onclick = runCode;
    }

    // Bind Reset button
    const resetBtn = document.getElementById("btn-reset-code");
    if (resetBtn) {
      resetBtn.onclick = () => {
        if (window.monacoEditor) {
          window.monacoEditor.setValue(defaultCode);
        } else if (fallbackPre) {
          fallbackPre.querySelector("code").textContent = defaultCode;
        }
        if (window.resetInspectLog) window.resetInspectLog();
        const panel = document.getElementById("result-panel");
        if (panel) panel.style.display = "none";
      };
    }
  }

  function setStatus(text, ok) {
    const statusEl = document.getElementById("result-status");
    if (!statusEl) return;
    statusEl.textContent = text;
    statusEl.className = "result-status " + (ok === null ? "" : ok ? "status-200" : "status-err");
  }

  // Status of a run: the returned Resp wins; otherwise the last HTTP call the
  // run made; otherwise there is no HTTP status to show.
  function statusOf(result, calls) {
    const di = window.docsInspect;
    if (di && di.isResp(result)) return { code: result.status, ok: result.ok };
    const last = calls[calls.length - 1];
    if (last) {
      if (last.error) return { code: 0, ok: false };
      return { code: last.status, ok: last.status >= 200 && last.status < 300 };
    }
    return null;
  }

  function statusText(st, elapsed, fallback) {
    if (!st) return `${fallback} (${elapsed}ms)`;
    if (st.code === 0) return `Network error (${elapsed}ms)`;
    return `${st.code} ${st.ok ? "OK" : "Error"} (${elapsed}ms)`;
  }

  async function runCode() {
    const container = document.getElementById("editor-container");
    if (!container) return;

    const code = window.monacoEditor
      ? window.monacoEditor.getValue()
      : container.dataset.defaultCode;

    const panel = document.getElementById("result-panel");
    const out = document.getElementById("result-output");
    const di = window.docsInspect;

    if (panel) panel.style.display = "block";
    if (window.resetInspectLog) window.resetInspectLog();
    if (out) out.textContent = "Running...";
    setStatus("...", null);

    if (!window.sdk && typeof I6Sdk !== "undefined") {
      const cfg = window.getSavedSdkConfig ? window.getSavedSdkConfig() : {};
      window.sdk = new I6Sdk(cfg);
    }
    if (!window.sdk) {
      if (out) out.textContent = "SDK instance is not available. Please ensure i6sdk-web.min.js is loaded.";
      setStatus("Error", false);
      return;
    }

    const reqStart = di ? di.requestCount() : 0;
    const startTime = performance.now();
    try {
      const AsyncFunction = Object.getPrototypeOf(async function () {}).constructor;
      const runner = new AsyncFunction("sdk", "inspect", code);
      const result = await runner(window.sdk, window.inspect);
      const elapsed = Math.round(performance.now() - startTime);
      const calls = di ? di.requestsSince(reqStart) : [];

      // Nothing printed via inspect(): show the returned value.
      if (di && di.entryCount() === 0) {
        di.appendEntry(result === undefined ? "[Completed with no return value]" : di.formatValue(result));
      } else if (!di && out) {
        out.textContent = String(result);
      }

      const st = statusOf(result, calls);
      setStatus(statusText(st, elapsed, "Done"), st ? st.ok : true);
    } catch (err) {
      console.error("Run error:", err);
      const elapsed = Math.round(performance.now() - startTime);
      const calls = di ? di.requestsSince(reqStart) : [];
      const msg = `Error: ${err.name || "Error"}: ${err.message}\n` + (err.stack || "");
      if (di && di.entryCount() > 0) di.appendEntry(msg);
      else if (out) out.textContent = msg;

      const st = statusOf(null, calls) || (typeof err.status === "number" ? { code: err.status, ok: false } : null);
      setStatus(statusText(st && { code: st.code, ok: false }, elapsed, "Error"), false);
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initEditor);
  } else {
    initEditor();
  }
})();
