/**
 * content_test.go
 * Test code for content input utilities
 */

package input

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadContent_ReadFromFile tests reading content from a file
func TestReadContent_ReadFromFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<p>Test Content</p>"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Read content from file
	content, err := ReadContent(tmpFile, "")
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	// Verify content
	if content != expectedContent {
		t.Errorf("Content = %q; want %q", content, expectedContent)
	}
}

// TestReadContent_ReturnInlineContent tests returning inline content
func TestReadContent_ReturnInlineContent(t *testing.T) {
	expectedContent := "<p>Inline Content</p>"

	// Read inline content (empty file path)
	content, err := ReadContent("", expectedContent)
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	// Verify content
	if content != expectedContent {
		t.Errorf("Content = %q; want %q", content, expectedContent)
	}
}

// TestReadContent_FileTakesPrecedence tests that file takes precedence
func TestReadContent_FileTakesPrecedence(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<p>File Content</p>"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Specify both file and inline content (file takes precedence)
	content, err := ReadContent(tmpFile, "<p>Inline Content</p>")
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	// Verify file content is returned
	if content != expectedContent {
		t.Errorf("Content = %q; want %q", content, expectedContent)
	}
}

// TestReadContent_ErrorWhenFileNotFound tests error when file does not exist
func TestReadContent_ErrorWhenFileNotFound(t *testing.T) {
	// Specify non-existent file
	_, err := ReadContent("/path/to/nonexistent/file.html", "")
	if err == nil {
		t.Error("Error should be returned")
	}
}

// TestReadContent_ReturnsEmptyStringWhenBothEmpty tests returning empty string when both are empty
func TestReadContent_ReturnsEmptyStringWhenBothEmpty(t *testing.T) {
	// Both file path and inline content are empty
	content, err := ReadContent("", "")
	if err != nil {
		t.Fatalf("Error should not be returned: %v", err)
	}

	// Verify empty string is returned
	if content != "" {
		t.Errorf("Content = %q; want %q", content, "")
	}
}

// TestDetectFormat_Markdown tests format detection for Markdown files
func TestDetectFormat_Markdown(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     ".md extension",
			filePath: "test.md",
			expected: FormatMarkdown,
		},
		{
			name:     ".markdown extension",
			filePath: "article.markdown",
			expected: FormatMarkdown,
		},
		{
			name:     ".md with path",
			filePath: "/path/to/file.md",
			expected: FormatMarkdown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestDetectFormat_HTML tests format detection for HTML files
func TestDetectFormat_HTML(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     ".html extension",
			filePath: "page.html",
			expected: FormatHTML,
		},
		{
			name:     ".htm extension",
			filePath: "index.htm",
			expected: FormatHTML,
		},
		{
			name:     ".html with path",
			filePath: "/var/www/page.html",
			expected: FormatHTML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestDetectFormat_Lexical tests format detection for Lexical JSON files
func TestDetectFormat_Lexical(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     ".json extension",
			filePath: "content.json",
			expected: FormatLexical,
		},
		{
			name:     ".json with path",
			filePath: "/data/lexical.json",
			expected: FormatLexical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestDetectFormat_Unknown tests detection of unknown file formats
func TestDetectFormat_Unknown(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected ContentFormat
	}{
		{
			name:     "no extension",
			filePath: "noextension",
			expected: FormatUnknown,
		},
		{
			name:     ".txt extension",
			filePath: "text.txt",
			expected: FormatUnknown,
		},
		{
			name:     "empty string",
			filePath: "",
			expected: FormatUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DetectFormat(tt.filePath)
			if format != tt.expected {
				t.Errorf("DetectFormat(%q) = %q; want %q", tt.filePath, format, tt.expected)
			}
		})
	}
}

// TestReadContentWithFormat_Markdown tests reading and format detection for Markdown files
func TestReadContentWithFormat_Markdown(t *testing.T) {
	// Create temporary Markdown file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	expectedContent := "# Heading\n\nThis is Markdown."
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Read from file and detect format
	content, format, err := ReadContentWithFormat(tmpFile, "")
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	// Verify content and format
	if content != expectedContent {
		t.Errorf("content = %q; want %q", content, expectedContent)
	}
	if format != FormatMarkdown {
		t.Errorf("format = %q; want %q", format, FormatMarkdown)
	}
}

// TestReadContentWithFormat_HTML tests reading and format detection for HTML files
func TestReadContentWithFormat_HTML(t *testing.T) {
	// Create temporary HTML file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.html")
	expectedContent := "<h1>Heading</h1><p>This is HTML.</p>"
	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	// Read from file and detect format
	content, format, err := ReadContentWithFormat(tmpFile, "")
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	// Verify content and format
	if content != expectedContent {
		t.Errorf("content = %q; want %q", content, expectedContent)
	}
	if format != FormatHTML {
		t.Errorf("format = %q; want %q", format, FormatHTML)
	}
}

// TestReadContentWithFormat_InlineContent tests reading inline content
func TestReadContentWithFormat_InlineContent(t *testing.T) {
	expectedContent := "<p>Inline Content</p>"

	// Read inline content
	content, format, err := ReadContentWithFormat("", expectedContent)
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	// Verify content and format
	if content != expectedContent {
		t.Errorf("content = %q; want %q", content, expectedContent)
	}
	// Format is unknown for inline content
	if format != FormatUnknown {
		t.Errorf("format = %q; want %q", format, FormatUnknown)
	}
}
