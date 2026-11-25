package forgepyschema

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/infinity6-ai/gox/commonz/dataz"
	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/schemaz/schemamold"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
)

var pythonKeywords = map[string]bool{
	"False": true, "None": true, "True": true, "and": true, "as": true,
	"assert": true, "async": true, "await": true, "break": true, "class": true,
	"continue": true, "def": true, "del": true, "elif": true, "else": true,
	"except": true, "finally": true, "for": true, "from": true, "global": true,
	"if": true, "import": true, "in": true, "is": true, "lambda": true,
	"nonlocal": true, "not": true, "or": true, "pass": true, "raise": true,
	"return": true, "try": true, "while": true, "with": true, "yield": true,
}

func isValidIdentifier(s string) bool {
	if s == "" || pythonKeywords[s] {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}

// NewSdk returns a schemamold.Parser configured with SDK known types.
func NewSdk() *schemamold.Parser {
	return schemamold.New().
		AddKnownType(fsz.FileStat{}, typeastz.TypFileStat).
		AddKnownType(dataz.NullString{}, typeastz.TypString).
		AddKnownType(dataz.NullInt64{}, typeastz.TypNumber).
		AddKnownType(dataz.NullFloat64{}, typeastz.TypNumber)
}

var arrayRe = regexp.MustCompile(`^Array<(.+)>$`)
var recordRe = regexp.MustCompile(`^Record<([^,]+),\s*(.+)>$`)
var jsonEmbeddedRe = regexp.MustCompile(`^i6JsonEmbedded<(.+)>$`)

// ToPythonType converts a typeastz type string and optional Go type to a Python type annotation.
func ToPythonType(typ string, goTyp string) string {
	typ = strings.TrimSpace(typ)
	if typ == "" {
		return "Any"
	}

	switch typ {
	case typeastz.TypString:
		return "str"
	case typeastz.TypNumber:
		if strings.Contains(goTyp, "int") || strings.Contains(goTyp, "uint") {
			return "int"
		}
		if strings.Contains(goTyp, "float") {
			return "float"
		}
		return "Union[int, float]"
	case typeastz.TypBoolean:
		return "bool"
	case typeastz.TypObject:
		return "Dict[str, Any]"
	case typeastz.TypAny, "any":
		return "Any"
	case typeastz.TypTime:
		return "str"
	case typeastz.TypFileStat:
		return "FileStat"
	}

	if m := arrayRe.FindStringSubmatch(typ); len(m) > 1 {
		elemPy := ToPythonType(m[1], "")
		return fmt.Sprintf("List[%s]", elemPy)
	}

	if m := recordRe.FindStringSubmatch(typ); len(m) > 2 {
		keyPy := ToPythonType(m[1], "")
		valPy := ToPythonType(m[2], "")
		return fmt.Sprintf("Dict[%s, %s]", keyPy, valPy)
	}

	if m := jsonEmbeddedRe.FindStringSubmatch(typ); len(m) > 1 {
		elemPy := ToPythonType(m[1], "")
		return fmt.Sprintf("Union[str, %s]", elemPy)
	}

	// Handle union / enum types, e.g. "1 | 2 | 4 | 8 | number"
	if strings.Contains(typ, "|") {
		parts := strings.Split(typ, "|")
		allNumbers := true
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "number" && !isNumeric(trimmed) {
				allNumbers = false
				break
			}
		}
		if allNumbers {
			return "int"
		}

		var pyParts []string
		seen := make(map[string]bool)
		for _, part := range parts {
			pyPart := ToPythonType(strings.TrimSpace(part), "")
			if !seen[pyPart] {
				seen[pyPart] = true
				pyParts = append(pyParts, pyPart)
			}
		}
		if len(pyParts) == 1 {
			return pyParts[0]
		}
		return fmt.Sprintf("Union[%s]", strings.Join(pyParts, ", "))
	}

	return typ
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '-' && r != '+' && r != '.' {
			return false
		}
	}
	return true
}

func writeln(sb *strings.Builder, msg string, args ...any) {
	_, err := fmt.Fprintf(sb, msg, args...)
	errorz.Check(err)
	_, err = sb.WriteRune('\n')
	errorz.Check(err)
}

// FormatTypedefs formats a typeastz.Typedefs collection into Python TypedDicts and type aliases.
func FormatTypedefs(typedefs *typeastz.Typedefs) string {
	if typedefs == nil {
		return ""
	}

	all := typedefs.All()
	if len(all) == 0 {
		return ""
	}

	var b strings.Builder
	seen := map[string]bool{
		"ServiceOptions": true,
		"FileStat":       true,
	}

	for _, td := range all {
		name := td.Head.Name
		if name == typeastz.TypFileStat {
			name = "FileStat"
		}
		if seen[name] {
			continue
		}
		seen[name] = true

		hasProps := td.Properties != nil && len(td.Properties.All()) > 0
		hasKeywordOrSpecial := false
		if hasProps {
			for _, p := range td.Properties.All() {
				if !isValidIdentifier(p.Name) {
					hasKeywordOrSpecial = true
					break
				}
			}
		}

		summary := td.Head.Summary
		goTyp := td.Head.GoTyp

		if hasProps {
			if hasKeywordOrSpecial {
				// Functional syntax for keywords or special characters
				writeln(&b, "%s = TypedDict(", name)
				writeln(&b, "    %q,", name)
				writeln(&b, "    {")
				for _, p := range td.Properties.All() {
					propPy := ToPythonType(p.Typ, p.GoTyp)
					writeln(&b, "        %q: %s,", p.Name, propPy)
				}
				writeln(&b, "    },")
				writeln(&b, "    total=False,")
				writeln(&b, ")")
				writeln(&b, "")
			} else {
				// Standard class syntax
				writeln(&b, "class %s(TypedDict, total=False):", name)
				var docLines []string
				if summary != "" {
					docLines = append(docLines, summary)
				}
				if goTyp != "" {
					docLines = append(docLines, fmt.Sprintf("[Go: %s]", goTyp))
				}
				if len(docLines) > 0 {
					writeln(&b, "    \"\"\"%s\"\"\"", strings.Join(docLines, "\n\n    "))
				}

				for _, p := range td.Properties.All() {
					propPy := ToPythonType(p.Typ, p.GoTyp)
					var comments []string
					if p.GoTyp != "" {
						comments = append(comments, fmt.Sprintf("[Go: %s]", p.GoTyp))
					}
					if p.Summary != "" {
						comments = append(comments, p.Summary)
					}
					if len(comments) > 0 {
						writeln(&b, "    # %s", strings.Join(comments, " - "))
					}
					writeln(&b, "    %s: %s", p.Name, propPy)
				}
				writeln(&b, "")
			}
		} else if td.Head.Typ == typeastz.TypObject {
			writeln(&b, "class %s(TypedDict, total=False):", name)
			if summary != "" || goTyp != "" {
				var docLines []string
				if summary != "" {
					docLines = append(docLines, summary)
				}
				if goTyp != "" {
					docLines = append(docLines, fmt.Sprintf("[Go: %s]", goTyp))
				}
				writeln(&b, "    \"\"\"%s\"\"\"", strings.Join(docLines, "\n\n    "))
			}
			writeln(&b, "    pass")
			writeln(&b, "")
		} else {
			// Type alias
			var comments []string
			if goTyp != "" {
				comments = append(comments, fmt.Sprintf("[Go: %s]", goTyp))
			}
			if summary != "" {
				comments = append(comments, summary)
			}
			if len(comments) > 0 {
				writeln(&b, "# %s", strings.Join(comments, " - "))
			}
			aliasPy := ToPythonType(td.Head.Typ, goTyp)
			writeln(&b, "%s = %s", name, aliasPy)
			writeln(&b, "")
		}
	}

	return b.String()
}

// WriteModuleTypedefs produces the complete types.py file contents including imports and base types.
func WriteModuleTypedefs(ctx context.Context, typedefs typeastz.Typedefs) string {
	var b strings.Builder
	writeln(&b, "\"\"\"Generated type definitions for i6 SDK.\"\"\"")
	writeln(&b, "")
	writeln(&b, "from __future__ import annotations")
	writeln(&b, "")
	writeln(&b, "from typing import Any, Dict, List, Optional, TypedDict, Union")
	writeln(&b, "")
	writeln(&b, "")
	writeln(&b, "class ServiceOptions(TypedDict, total=False):")
	writeln(&b, "    \"\"\"Execution options for service requests.\"\"\"")
	writeln(&b, "")
	writeln(&b, "    timeout: Optional[float]")
	writeln(&b, "    headers: Optional[Dict[str, str]]")
	writeln(&b, "    labels: Optional[Dict[str, str]]")
	writeln(&b, "    service_name: Optional[str]")
	writeln(&b, "")
	writeln(&b, "")
	writeln(&b, "class FileStat(TypedDict, total=False):")
	writeln(&b, "    \"\"\"File storage metadata.\"\"\"")
	writeln(&b, "")
	writeln(&b, "    url: Optional[str]")
	writeln(&b, "    content_type: Optional[str]")
	writeln(&b, "    md5: Optional[str]")
	writeln(&b, "    size: Optional[str]")
	writeln(&b, "    etag: Optional[str]")
	writeln(&b, "    created_at: Optional[str]")
	writeln(&b, "    updated_at: Optional[str]")
	writeln(&b, "")

	if formatted := FormatTypedefs(&typedefs); formatted != "" {
		writeln(&b, "")
		b.WriteString(strings.TrimSuffix(formatted, "\n"))
		writeln(&b, "")
	}

	return b.String()
}

// WriteTypedefs aliases WriteModuleTypedefs for consistency.
func WriteTypedefs(ctx context.Context, typedefs typeastz.Typedefs) string {
	return WriteModuleTypedefs(ctx, typedefs)
}
