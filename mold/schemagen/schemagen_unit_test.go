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
}
