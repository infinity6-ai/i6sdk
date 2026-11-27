package cmdsdkweb

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var (
	reCodeBlock = regexp.MustCompile("(?s)```(?:[a-zA-Z0-9_-]+)?\\n?(.*?)```")
	reInlineCode = regexp.MustCompile("`([^`]+)`")
	reBold       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	reJsDocLink  = regexp.MustCompile(`\{@link\s+([^}]+)\}`)
)

// RenderMarkdown renders a subset of Markdown used in SDK JSDocs into safe HTML.
// Input text is escaped first, so Markdown text cannot inject raw HTML.
func RenderMarkdown(input string, linkResolver func(target string) string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}

	// 1. Extract and placeholder code blocks so their contents are not touched by inline formatting
	type blockPlaceholder struct {
		placeholder string
		html        string
	}
	var placeholders []blockPlaceholder
	blockIdx := 0

	escapedInput := reCodeBlock.ReplaceAllStringFunc(trimmed, func(match string) string {
		sub := reCodeBlock.FindStringSubmatch(match)
		code := ""
		if len(sub) > 1 {
			code = sub[1]
		}
		ph := fmt.Sprintf("@@DOCSMD_BLOCK_%d@@", blockIdx)
		blockIdx++
		placeholders = append(placeholders, blockPlaceholder{
			placeholder: ph,
			html:        "<pre><code>" + html.EscapeString(strings.TrimRight(code, "\n")) + "</code></pre>",
		})
		return "\n\n" + ph + "\n\n"
	})

	// 2. Split into blocks by blank lines
	rawBlocks := strings.Split(escapedInput, "\n\n")
	var renderedBlocks []string

	for _, rawBlock := range rawBlocks {
		block := strings.TrimSpace(rawBlock)
		if block == "" {
			continue
		}

		// Check if it's a code block placeholder
		isPlaceholder := false
		for _, ph := range placeholders {
			if block == ph.placeholder {
				renderedBlocks = append(renderedBlocks, ph.html)
				isPlaceholder = true
				break
			}
		}
		if isPlaceholder {
			continue
		}

		// Check if block is a list
		lines := strings.Split(block, "\n")
		isUnorderedList := true
		isOrderedList := true

		for _, l := range lines {
			tl := strings.TrimSpace(l)
			if !strings.HasPrefix(tl, "- ") && !strings.HasPrefix(tl, "* ") {
				isUnorderedList = false
			}
			matchedOrdered, _ := regexp.MatchString(`^\d+\.\s+`, tl)
			if !matchedOrdered {
				isOrderedList = false
			}
		}

		if isUnorderedList && len(lines) > 0 {
			var listItems []string
			for _, l := range lines {
				tl := strings.TrimSpace(l)
				itemText := strings.TrimPrefix(strings.TrimPrefix(tl, "- "), "* ")
				listItems = append(listItems, "<li>"+renderInline(itemText, linkResolver)+"</li>")
			}
			renderedBlocks = append(renderedBlocks, "<ul>\n"+strings.Join(listItems, "\n")+"\n</ul>")
		} else if isOrderedList && len(lines) > 0 {
			var listItems []string
			reNum := regexp.MustCompile(`^\d+\.\s+`)
			for _, l := range lines {
				tl := strings.TrimSpace(l)
				itemText := reNum.ReplaceAllString(tl, "")
				listItems = append(listItems, "<li>"+renderInline(itemText, linkResolver)+"</li>")
			}
			renderedBlocks = append(renderedBlocks, "<ol>\n"+strings.Join(listItems, "\n")+"\n</ol>")
		} else {
			// Regular paragraph
			renderedBlocks = append(renderedBlocks, "<p>"+renderInline(block, linkResolver)+"</p>")
		}
	}

	result := strings.Join(renderedBlocks, "\n")

	// Restore any placeholders if they were nested
	for _, ph := range placeholders {
		result = strings.ReplaceAll(result, ph.placeholder, ph.html)
	}

	return result
}

func renderInline(text string, linkResolver func(target string) string) string {
	// Escape first
	escaped := html.EscapeString(text)

	// Replace code `...`
	escaped = reInlineCode.ReplaceAllString(escaped, "<code>$1</code>")

	// Replace bold **...**
	escaped = reBold.ReplaceAllString(escaped, "<strong>$1</strong>")

	// Replace JSDoc {@link Target}
	escaped = reJsDocLink.ReplaceAllStringFunc(escaped, func(match string) string {
		sub := reJsDocLink.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		rawTarget := strings.TrimSpace(sub[1])
		target := html.UnescapeString(rawTarget)
		href := target
		if linkResolver != nil {
			href = linkResolver(target)
		}
		return fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(target))
	})

	return escaped
}
