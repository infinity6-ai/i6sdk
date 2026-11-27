package main_test

import (
	"context"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/infinity6-ai/gox/brwjsz/brwclient"
	"github.com/infinity6-ai/gox/brwjsz/brwjsz"
	"github.com/infinity6-ai/gox/brwjsz/brwqunit"
	"github.com/infinity6-ai/gox/commonz/pathz"
	"github.com/stretchr/testify/require"
)

var defaultImportMap = map[string]string{
	"sql.js": "https://cdn.jsdelivr.net/npm/sql.js@1.14.2/+esm",
}

func findSdkWebDir(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return startDir
		}
		dir = parent
	}
}

func TestWebSdk(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 120*time.Second)
	defer cancel()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	sdkDir := findSdkWebDir(cwd)

	opts := brwjsz.Options{
		Site: &brwjsz.Site{
			Disk: pathz.MustParse(sdkDir),
		},
	}
	opts.Chrome("", true, "")

	b := brwjsz.New(ctx, opts)
	defer b.Close()

	err = b.Start(ctx)
	require.NoError(t, err)

	client := b.Client()

	files, err := filepath.Glob(filepath.Join(sdkDir, "src", "test", "*-test.js"))
	require.NoError(t, err)
	require.NotEmpty(t, files, "no test files found matching src/test/*-test.js")

	sort.Strings(files)

	for _, file := range files {
		fileName := filepath.Base(file)
		testUrl := "/src/test/" + fileName
		t.Run(fileName, func(t *testing.T) {
			brwqunit.TestWithOptions(t, client, brwqunit.Options{
				TestUrls:  []string{testUrl},
				ImportMap: defaultImportMap,
			})
		})
	}
}

// docsPage is what TestWebDocs reads from a built docs page, once Monaco has
// loaded the SDK types. ExampleDiags has the type errors of each code given,
// checked the way the Try it editor checks it.
type docsPage struct {
	Title        string              `json:"title"`
	Fields       []string            `json:"fields"`
	SdkType      string              `json:"sdkType"`
	ExampleDiags map[string][]string `json:"exampleDiags"`
}

// Opens the page in an iframe, waits for the editor and asks its TypeScript
// worker for the type of `sdk`, then puts each example in the editor and
// reads its diagnostics. The ignored codes are the ones tryit.js ignores
// (top-level return/await).
const readDocsPageJs = `async (url, examples) => {
	const f = document.createElement("iframe");
	f.style.width = "1200px";
	f.style.height = "800px";
	f.src = url;
	document.body.appendChild(f);
	await new Promise((r) => (f.onload = r));
	const w = f.contentWindow;
	const start = Date.now();
	while (!w.monacoEditor) {
		if (Date.now() - start > 60000) throw new Error("Monaco editor did not load");
		await new Promise((r) => setTimeout(r, 200));
	}
	const model = w.monacoEditor.getModel();
	const worker = await (await w.monaco.languages.typescript.getJavaScriptWorker())(model.uri);
	const info = await worker.getQuickInfoAtPosition(model.uri.toString(), model.getValue().indexOf("sdk.") + 1);
	const exampleDiags = {};
	const uri = model.uri.toString();
	for (const [name, code] of Object.entries(examples || {})) {
		model.setValue(code);
		const diags = [...(await worker.getSyntacticDiagnostics(uri)), ...(await worker.getSemanticDiagnostics(uri))]
			.filter((d) => ![1108, 1375, 1378].includes(d.code))
			.map((d) => (typeof d.messageText === "string" ? d.messageText : d.messageText.messageText));
		if (diags.length > 0) exampleDiags[name] = diags;
	}
	const ret = {
		title: w.document.title,
		fields: [...w.document.querySelectorAll("tr[id^=field-]")].map((e) => e.id),
		sdkType: info ? info.displayParts.map((p) => p.text).join("") : "",
		exampleDiags,
	};
	f.remove();
	return ret;
}`

var exampleCodeRe = regexp.MustCompile(`data-default-code="([^"]*)"`)

// readExamples returns the Try it code of every built page, by page path.
func readExamples(t *testing.T, docsDir string) map[string]string {
	examples := map[string]string{}
	err := filepath.WalkDir(docsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if m := exampleCodeRe.FindSubmatch(b); m != nil {
			rel, _ := filepath.Rel(docsDir, path)
			examples[filepath.ToSlash(rel)] = html.UnescapeString(string(m[1]))
		}
		return nil
	})
	require.NoError(t, err)
	return examples
}

// TestWebDocs checks pages built by `./build.sh docs` (target/docs).
func TestWebDocs(t *testing.T) {
	type testScenario struct {
		page   string
		title  string
		fields []string
	}

	ctx, cancel := context.WithTimeout(t.Context(), 180*time.Second)
	defer cancel()

	cwd, err := os.Getwd()
	require.NoError(t, err)
	sdkDir := findSdkWebDir(cwd)
	docsDir := filepath.Join(sdkDir, "target", "docs")
	_, err = os.Stat(filepath.Join(docsDir, "types.json"))
	require.NoError(t, err, "docs are not built: run ./build.sh docs --allow-missing")

	opts := brwjsz.Options{
		Site: &brwjsz.Site{
			Disk: pathz.MustParse(sdkDir),
		},
	}
	opts.Chrome("", true, "")

	b := brwjsz.New(ctx, opts)
	defer b.Close()
	require.NoError(t, b.Start(ctx))
	client := b.Client()

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		page, err := brwclient.Call[docsPage](ctx, client, readDocsPageJs, "/target/docs/"+s.page, nil)
		require.NoError(t, err)

		require.Equal(t, s.title, page.Title)
		for _, f := range s.fields {
			require.Contains(t, page.Fields, f)
		}
		require.Contains(t, page.SdkType, "I6Sdk")
	}

	t.Run("API page has its fields and typed sdk", func(t *testing.T) {
		check(t, testScenario{
			page:   "apis/dsz-create-customer-portal.html",
			title:  "Create Customer Portal (sdk.apis().dszCreateCustomerPortal) — i6 SDK",
			fields: []string{"field-params-dataset", "field-payload-locale", "field-payload-return_url"},
		})
	})

	t.Run("Handwritten page has its fields and typed sdk", func(t *testing.T) {
		check(t, testScenario{
			page:   "infer-run.html",
			title:  "run (sdk.infer().run) — i6 SDK",
			fields: []string{"field-opts-modelName", "field-opts-partitions"},
		})
	})

	t.Run("Every example type checks in the editor", func(t *testing.T) {
		examples := readExamples(t, docsDir)
		require.NotEmpty(t, examples)

		page, err := brwclient.Call[docsPage](ctx, client, readDocsPageJs, "/target/docs/infer-run.html", examples)
		require.NoError(t, err)

		require.Empty(t, page.ExampleDiags)
	})
}
