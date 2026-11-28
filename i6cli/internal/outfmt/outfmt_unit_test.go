package outfmt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/stretchr/testify/require"
)

func TestUnitParseFormat(t *testing.T) {
	type testScenario struct {
		input       string
		defaultFmt  outfmt.Format
		wantFormat  outfmt.Format
		wantErrCode int
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var opts []outfmt.Format
		if s.defaultFmt != "" {
			opts = append(opts, s.defaultFmt)
		}
		got, err := outfmt.ParseFormat(s.input, opts...)
		if s.wantErrCode != 0 {
			require.Error(t, err)
			require.Equal(t, s.wantErrCode, exitcode.Code(err))
			return
		}
		require.NoError(t, err)
		require.Equal(t, s.wantFormat, got)
	}

	t.Run("Empty string defaults to table", func(t *testing.T) {
		check(t, testScenario{
			input:      "",
			wantFormat: outfmt.FormatTable,
		})
	})

	t.Run("Empty string uses provided defaultFormat", func(t *testing.T) {
		check(t, testScenario{
			input:      "",
			defaultFmt: outfmt.FormatRaw,
			wantFormat: outfmt.FormatRaw,
		})
	})

	t.Run("Valid formats parsed case-insensitively", func(t *testing.T) {
		for _, f := range outfmt.AllFormats {
			check(t, testScenario{
				input:      strings.ToUpper(string(f)),
				wantFormat: f,
			})
		}
	})

	t.Run("Invalid format returns Usage error", func(t *testing.T) {
		check(t, testScenario{
			input:       "xml",
			wantErrCode: exitcode.Usage,
		})
	})
}

func TestUnitRenderFormats(t *testing.T) {
	sampleData := outfmt.Data{
		Headers: []string{"Dataset", "Solution", "Status"},
		Rows: [][]string{
			{"fashion", "relevance", "active"},
			{"finance", "targeting", "pending"},
		},
		QuietId: []string{"fashion", "finance"},
		Raw:     "i6_token_sample_jwt_data",
		JSON: map[string]any{
			"total": 2,
			"items": []string{"fashion", "finance"},
		},
	}

	t.Run("Table format renders uppercase headers and rows", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := outfmt.Render(buf, outfmt.FormatTable, sampleData)
		require.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, "DATASET")
		require.Contains(t, output, "SOLUTION")
		require.Contains(t, output, "fashion")
		require.Contains(t, output, "finance")
	})

	t.Run("TSV format renders tab-delimited columns", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := outfmt.Render(buf, outfmt.FormatTSV, sampleData)
		require.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		require.Len(t, lines, 3)
		require.Equal(t, "Dataset\tSolution\tStatus", lines[0])
		require.Equal(t, "fashion\trelevance\tactive", lines[1])
	})

	t.Run("JSON format renders indented json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := outfmt.Render(buf, outfmt.FormatJSON, sampleData)
		require.NoError(t, err)
		output := buf.String()
		require.Contains(t, output, `"total": 2`)
		require.Contains(t, output, `"fashion"`)
	})

	t.Run("JSON fallback when JSON object is nil", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := outfmt.Render(buf, outfmt.FormatJSON, outfmt.Data{
			Headers: []string{"id", "val"},
			Rows:    [][]string{{"a", "1"}},
		})
		require.NoError(t, err)
		require.Contains(t, buf.String(), `"id": "a"`)
		require.Contains(t, buf.String(), `"val": "1"`)
	})

	t.Run("Quiet format renders primary identifiers only", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := outfmt.Render(buf, outfmt.FormatQuiet, sampleData)
		require.NoError(t, err)
		output := buf.String()
		require.Equal(t, "fashion\nfinance\n", output)
	})

	t.Run("Raw format renders unquoted raw string", func(t *testing.T) {
		buf := new(bytes.Buffer)
		err := outfmt.Render(buf, outfmt.FormatRaw, sampleData)
		require.NoError(t, err)
		require.Equal(t, "i6_token_sample_jwt_data\n", buf.String())
	})
}
