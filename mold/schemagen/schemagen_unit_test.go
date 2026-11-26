package schemagen_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/i6sdk/mold/schemagen"
	"github.com/stretchr/testify/require"
)

func TestUnitSchemagen(t *testing.T) {
	type testScenario struct {
		name          string
		goSource      string
		expectedInGen []string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		tmpDir, err := os.MkdirTemp("", "schemagen-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		srcFile := filepath.Join(tmpDir, "api.go")
		err = os.WriteFile(srcFile, []byte(s.goSource), 0644)
		require.NoError(t, err)

		code, pkgName, err := schemagen.ParseAndGenerate(tmpDir)
		require.NoError(t, err)
		require.Equal(t, "sampleapi", pkgName)

		codeStr := string(code)
		for _, exp := range s.expectedInGen {
			require.Contains(t, codeStr, exp)
		}
	}

	t.Run("Generates PathParams and Payload with slice and doc comments", func(t *testing.T) {
		source := `package sampleapi

// Params has path segment parameters.
type Params struct {
	// Unique identifier for user.
	UserId string ` + "`json:\"userId\"`" + `
}

// Item represents a single data record.
type Item struct {
	// Item unique id.
	Id string ` + "`json:\"id\"`" + `
	// Item name.
	Name string ` + "`json:\"name\"`" + `
}

// Payload is the response body.
type Payload struct {
	// List of items returned.
	Items []*Item ` + "`json:\"items\"`" + `
}

type Req struct {
	Params *Params ` + "`json:\"params,omitempty\"`" + `
}

type Resp struct {
	Payload *Payload ` + "`json:\"payload\"`" + `
}

type Api struct {
	Req  *Req  ` + "`json:\"req\"`" + `
	Resp *Resp ` + "`json:\"resp\"`" + `
}
`
		check(t, testScenario{
			name:     "Generates PathParams and Payload with slice and doc comments",
			goSource: source,
			expectedInGen: []string{
				"package sampleapi",
				"func (d *Api) ApiDataRefs() *apiz.DataRefs {",
				"PathParams: &schemaz.Schema{",
				"\"userId\": {",
				"Summary: \"Unique identifier for user.\"",
				"UrlValue: schemaz.ParseStr(&d.Req.Params.UserId)",
				"RespBody: &schemaz.Schema{",
				"\"items\": {",
				"Summary: \"List of items returned.\"",
				"Array: schemaz.NewArray(schemaz.ArrayOptions[[]*Item, *Item]{",
				"\"id\": {",
				"Summary: \"Item unique id.\"",
				"Field: func() any { return &(*element).Id }",
				"\"name\": {",
				"Summary: \"Item name.\"",
				"Field: func() any { return &(*element).Name }",
				"QueryParams: nil",
				"ReqHeaders:  nil",
				"ReqBody:     nil",
				"RespHeaders: nil",
			},
		})
	})

	t.Run("Generates enriched ApiSpec and ApiDataRefs when ApiRoute is defined", func(t *testing.T) {
		source := `package sampleapi

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Params struct {
	Id string ` + "`json:\"id\"`" + `
}

type Req struct {
	Params *Params ` + "`json:\"params\"`" + `
}

type Resp struct{}

// Sample Item Get
// Retrieves a single item record from the dataset.
type Api struct {
	Req  *Req  ` + "`json:\"req\"`" + `
	Resp *Resp ` + "`json:\"resp\"`" + `
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id:     "sample-item-get",
		Method: "GET",
		Path:   "/items/{id}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "sample"},
			}
		},
	}
}
`
		check(t, testScenario{
			name:     "Generates enriched ApiSpec and ApiDataRefs when ApiRoute is defined",
			goSource: source,
			expectedInGen: []string{
				"func (d *Api) ApiSpec() apiz.ApiSpec {",
				"spec := d.ApiRoute()",
				"origDesc := spec.Desc",
				"if desc.Name == \"\" {",
				"desc.Name = \"Sample Item Get\"",
				"if desc.Summary == \"\" {",
				"desc.Summary = \"Retrieves a single item record from the dataset.\"",
				"return desc",
				"return spec",
				"func (d *Api) ApiDataRefs() *apiz.DataRefs {",
			},
		})
	})

	t.Run("GenerateDir writes out .donotedit.go file and ignores it on subsequent runs", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "schemagen-donotedit-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		srcFile := filepath.Join(tmpDir, "api.go")
		source := `package testpkg

type Req struct{}
type Resp struct{}
type Api struct {
	Req  *Req  ` + "`json:\"req\"`" + `
	Resp *Resp ` + "`json:\"resp\"`" + `
}
`
		err = os.WriteFile(srcFile, []byte(source), 0644)
		require.NoError(t, err)

		outFile, err := schemagen.GenerateDir(tmpDir)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "testpkg.donotedit.go"), outFile)

		// Verify file exists
		_, err = os.Stat(outFile)
		require.NoError(t, err)

		// Re-run should cleanly parse and regenerate without failing due to testpkg.donotedit.go
		outFile2, err := schemagen.GenerateDir(tmpDir)
		require.NoError(t, err)
		require.Equal(t, outFile, outFile2)
	})
}
