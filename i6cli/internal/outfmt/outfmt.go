package outfmt

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
)

// Format represents an output format option.
type Format string

const (
	// FormatTable renders aligned ASCII tables for terminal readability.
	FormatTable Format = "table"

	// FormatTSV renders tab-separated values with headers for Unix stream pipelines.
	FormatTSV Format = "tsv"

	// FormatJSON renders indented JSON for programmatic processing and jq.
	FormatJSON Format = "json"

	// FormatQuiet outputs only primary identifiers, suppressing tables and headers.
	FormatQuiet Format = "quiet"

	// FormatRaw outputs unquoted string or byte data (e.g. raw JWT or private key).
	FormatRaw Format = "raw"
)

// AllFormats contains all recognized output formats.
var AllFormats = []Format{
	FormatTable,
	FormatTSV,
	FormatJSON,
	FormatQuiet,
	FormatRaw,
}

// ParseFormat validates and parses the format string. If empty, defaultFormat is used.
func ParseFormat(s string, defaultFormat ...Format) (Format, error) {
	trimmed := strings.ToLower(strings.TrimSpace(s))
	if trimmed == "" {
		if len(defaultFormat) > 0 && defaultFormat[0] != "" {
			return defaultFormat[0], nil
		}
		return FormatTable, nil
	}

	f := Format(trimmed)
	switch f {
	case FormatTable, FormatTSV, FormatJSON, FormatQuiet, FormatRaw:
		return f, nil
	default:
		return "", exitcode.New(
			exitcode.Usage,
			fmt.Sprintf("invalid format %q: expected one of table, tsv, json, quiet, raw", s),
		)
	}
}

// Data encapsulates the structured content to be rendered across different formats.
type Data struct {
	// Headers defines column titles for tabular displays.
	Headers []string

	// Rows contains the matrix of row data corresponding to Headers.
	Rows [][]string

	// JSON holds an arbitrary structured object for json format. If nil, Rows and Headers are converted.
	JSON any

	// QuietId holds the minimal identifier values printed in quiet mode.
	QuietId []string

	// Raw holds the unformatted string printed in raw mode.
	Raw string
}

// Render formats Data according to Format and writes the result to w.
func Render(w io.Writer, f Format, d Data) error {
	switch f {
	case FormatTable:
		return renderTable(w, d)
	case FormatTSV:
		return renderTSV(w, d)
	case FormatJSON:
		return renderJSON(w, d)
	case FormatQuiet:
		return renderQuiet(w, d)
	case FormatRaw:
		return renderRaw(w, d)
	default:
		return exitcode.New(exitcode.Usage, fmt.Sprintf("unsupported format %q", f))
	}
}

func renderTable(w io.Writer, d Data) error {
	if len(d.Headers) == 0 && len(d.Rows) == 0 {
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 4, 3, ' ', 0)
	if len(d.Headers) > 0 {
		upperHeaders := make([]string, len(d.Headers))
		for i, h := range d.Headers {
			upperHeaders[i] = strings.ToUpper(h)
		}
		if _, err := fmt.Fprintln(tw, strings.Join(upperHeaders, "\t")); err != nil {
			return fmt.Errorf("failed to write table headers: %w", err)
		}
	}
	for _, row := range d.Rows {
		if _, err := fmt.Fprintln(tw, strings.Join(row, "\t")); err != nil {
			return fmt.Errorf("failed to write table row: %w", err)
		}
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("failed to flush table writer: %w", err)
	}
	return nil
}

func renderTSV(w io.Writer, d Data) error {
	if len(d.Headers) > 0 {
		if _, err := fmt.Fprintln(w, strings.Join(d.Headers, "\t")); err != nil {
			return fmt.Errorf("failed to write tsv headers: %w", err)
		}
	}
	for _, row := range d.Rows {
		if _, err := fmt.Fprintln(w, strings.Join(row, "\t")); err != nil {
			return fmt.Errorf("failed to write tsv row: %w", err)
		}
	}
	return nil
}

func renderJSON(w io.Writer, d Data) error {
	target := d.JSON
	if target == nil {
		// Fallback: build an array of maps from Headers and Rows
		list := make([]map[string]any, 0, len(d.Rows))
		for _, row := range d.Rows {
			item := make(map[string]any, len(d.Headers))
			for i, h := range d.Headers {
				if i < len(row) {
					item[h] = row[i]
				} else {
					item[h] = ""
				}
			}
			list = append(list, item)
		}
		target = list
	}

	bytes, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json output: %w", err)
	}
	if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
		return fmt.Errorf("failed to write json output: %w", err)
	}
	return nil
}

func renderQuiet(w io.Writer, d Data) error {
	if len(d.QuietId) > 0 {
		for _, id := range d.QuietId {
			if _, err := fmt.Fprintln(w, id); err != nil {
				return fmt.Errorf("failed to write quiet output: %w", err)
			}
		}
		return nil
	}
	for _, row := range d.Rows {
		if len(row) > 0 {
			if _, err := fmt.Fprintln(w, row[0]); err != nil {
				return fmt.Errorf("failed to write quiet row id: %w", err)
			}
		}
	}
	return nil
}

func renderRaw(w io.Writer, d Data) error {
	if d.Raw != "" {
		if _, err := fmt.Fprintln(w, d.Raw); err != nil {
			return fmt.Errorf("failed to write raw output: %w", err)
		}
		return nil
	}
	if len(d.QuietId) > 0 {
		if _, err := fmt.Fprintln(w, d.QuietId[0]); err != nil {
			return fmt.Errorf("failed to write raw quiet id: %w", err)
		}
		return nil
	}
	if len(d.Rows) > 0 && len(d.Rows[0]) > 0 {
		if _, err := fmt.Fprintln(w, d.Rows[0][0]); err != nil {
			return fmt.Errorf("failed to write raw row value: %w", err)
		}
		return nil
	}
	return nil
}
