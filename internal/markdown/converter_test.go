/**
 * Test code for Markdown→HTML conversion functionality
 */
package markdown

import (
	"strings"
	"testing"
)

// TestConvertToHTML_Paragraph tests paragraph conversion
func TestConvertToHTML_Paragraph(t *testing.T) {
	markdown := "This is a paragraph."
	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("Conversion error occurred: %v", err)
	}

	// goldmark wraps with <p> tags and adds newline at the end
	expected := "<p>This is a paragraph.</p>\n"
	if html != expected {
		t.Errorf("Does not match expected value.\nExpected: %q\nActual: %q", expected, html)
	}
}

// TestConvertToHTML_Headings tests heading (h1-h6) conversion
func TestConvertToHTML_Headings(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "h1 heading",
			markdown: "# Heading 1",
			expected: "<h1>Heading 1</h1>\n",
		},
		{
			name:     "h2 heading",
			markdown: "## Heading 2",
			expected: "<h2>Heading 2</h2>\n",
		},
		{
			name:     "h3 heading",
			markdown: "### Heading 3",
			expected: "<h3>Heading 3</h3>\n",
		},
		{
			name:     "h6 heading",
			markdown: "###### Heading 6",
			expected: "<h6>Heading 6</h6>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ConvertToHTML(tt.markdown)
			if err != nil {
				t.Fatalf("Conversion error occurred: %v", err)
			}
			if html != tt.expected {
				t.Errorf("Does not match expected value.\nExpected: %q\nActual: %q", tt.expected, html)
			}
		})
	}
}

// TestConvertToHTML_Lists tests list (bullet points, numbered) conversion
func TestConvertToHTML_Lists(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string // Check elements that should be included rather than exact match
	}{
		{
			name:     "bullet list",
			markdown: "- Item 1\n- Item 2\n- Item 3",
			contains: []string{"<ul>", "<li>Item 1</li>", "<li>Item 2</li>", "<li>Item 3</li>", "</ul>"},
		},
		{
			name:     "numbered list",
			markdown: "1. First\n2. Second\n3. Last",
			contains: []string{"<ol>", "<li>First</li>", "<li>Second</li>", "<li>Last</li>", "</ol>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ConvertToHTML(tt.markdown)
			if err != nil {
				t.Fatalf("Conversion error occurred: %v", err)
			}

			// Verify each element is included
			for _, elem := range tt.contains {
				if !strings.Contains(html, elem) {
					t.Errorf("Expected element not included: %q\nHTML: %q", elem, html)
				}
			}
		})
	}
}

// TestConvertToHTML_CodeBlock tests code block conversion
func TestConvertToHTML_CodeBlock(t *testing.T) {
	markdown := "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```"
	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("Conversion error occurred: %v", err)
	}

	// Code blocks are wrapped with <pre><code>
	if !strings.Contains(html, "<pre>") {
		t.Errorf("<pre> tag not included: %q", html)
	}
	if !strings.Contains(html, "<code") {
		t.Errorf("<code> tag not included: %q", html)
	}
	if !strings.Contains(html, "func main()") {
		t.Errorf("Code content not included: %q", html)
	}
}

// TestConvertToHTML_LinksAndImages tests link and image conversion
func TestConvertToHTML_LinksAndImages(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string
	}{
		{
			name:     "link",
			markdown: "[Link to Google](https://google.com)",
			contains: []string{"<a href=\"https://google.com\"", "Link to Google</a>"},
		},
		{
			name:     "image",
			markdown: "![Alternative text](/path/to/image.png)",
			contains: []string{"<img src=\"/path/to/image.png\"", "alt=\"Alternative text\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := ConvertToHTML(tt.markdown)
			if err != nil {
				t.Fatalf("Conversion error occurred: %v", err)
			}

			for _, elem := range tt.contains {
				if !strings.Contains(html, elem) {
					t.Errorf("Expected element not included: %q\nHTML: %q", elem, html)
				}
			}
		})
	}
}

// TestConvertToHTML_EmptyString tests behavior with empty string input
func TestConvertToHTML_EmptyString(t *testing.T) {
	markdown := ""
	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("Error occurred with empty string: %v", err)
	}

	// Empty string should return empty string
	if html != "" {
		t.Errorf("Expected empty string, but got: %q", html)
	}
}

// TestConvertToHTML_ComplexMarkdown tests complex Markdown conversion
func TestConvertToHTML_ComplexMarkdown(t *testing.T) {
	markdown := `# Heading

This is a paragraph with **bold** and *italic* text.

- List item 1
- List item 2

It also includes a [link](https://example.com).`

	html, err := ConvertToHTML(markdown)

	if err != nil {
		t.Fatalf("Conversion error occurred: %v", err)
	}

	// Verify each element is included
	expectedElements := []string{
		"<h1>Heading</h1>",
		"<strong>bold</strong>",
		"<em>italic</em>",
		"<ul>",
		"<li>List item 1</li>",
		"<a href=\"https://example.com\"",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(html, elem) {
			t.Errorf("Expected element not included: %q\nHTML: %q", elem, html)
		}
	}
}
