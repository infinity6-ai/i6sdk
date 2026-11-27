package cmdsdkweb_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkweb/cmdsdkweb"
	"github.com/stretchr/testify/require"
)

func TestUnitRenderMarkdown(t *testing.T) {
	type testScenario struct {
		input        string
		want         string
		linkResolver func(string) string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := cmdsdkweb.RenderMarkdown(s.input, s.linkResolver)
		require.Equal(t, s.want, got)
	}

	t.Run("Empty input returns empty string", func(t *testing.T) {
		check(t, testScenario{
			input: "",
			want:  "",
		})
	})

	t.Run("Paragraphs separated by blank lines", func(t *testing.T) {
		check(t, testScenario{
			input: "First paragraph.\n\nSecond paragraph.",
			want:  "<p>First paragraph.</p>\n<p>Second paragraph.</p>",
		})
	})

	t.Run("Inline code", func(t *testing.T) {
		check(t, testScenario{
			input: "Use `sdk.apis()` to make calls.",
			want:  "<p>Use <code>sdk.apis()</code> to make calls.</p>",
		})
	})

	t.Run("Unordered list", func(t *testing.T) {
		check(t, testScenario{
			input: "- First item\n- Second item",
			want:  "<ul>\n<li>First item</li>\n<li>Second item</li>\n</ul>",
		})
	})

	t.Run("Ordered list", func(t *testing.T) {
		check(t, testScenario{
			input: "1. Step one\n2. Step two",
			want:  "<ol>\n<li>Step one</li>\n<li>Step two</li>\n</ol>",
		})
	})

	t.Run("Code blocks", func(t *testing.T) {
		check(t, testScenario{
			input: "```js\nconst sdk = new I6Sdk()\n```",
			want:  "<pre><code>const sdk = new I6Sdk()</code></pre>",
		})
	})

	t.Run("Bold text", func(t *testing.T) {
		check(t, testScenario{
			input: "This is **very important**.",
			want:  "<p>This is <strong>very important</strong>.</p>",
		})
	})

	t.Run("JSDoc link with custom resolver", func(t *testing.T) {
		check(t, testScenario{
			input: "See {@link Auth#handleCallback} for details.",
			linkResolver: func(target string) string {
				if target == "Auth#handleCallback" {
					return "auth-handle-callback.html"
				}
				return target
			},
			want: `<p>See <a href="auth-handle-callback.html">Auth#handleCallback</a> for details.</p>`,
		})
	})

	t.Run("HTML escaping prevents injection", func(t *testing.T) {
		check(t, testScenario{
			input: `<script>alert("xss")</script> & "quotes"`,
			want:  "<p>&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt; &amp; &#34;quotes&#34;</p>",
		})
	})
}
