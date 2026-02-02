/**
 * content.go
 * Content input utility
 *
 * Provides functionality to read content from files, stdin, and inline content
 */

package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ContentFormat represents the format of content
type ContentFormat string

const (
	// FormatUnknown represents an unknown format
	FormatUnknown ContentFormat = ""
	// FormatHTML represents HTML format
	FormatHTML ContentFormat = "html"
	// FormatMarkdown represents Markdown format
	FormatMarkdown ContentFormat = "markdown"
	// FormatLexical represents Lexical JSON format
	FormatLexical ContentFormat = "lexical"
)

// ReadContent reads content from a file or returns inline content
//
// Priority:
// 1. If filePath is specified, read content from the file
// 2. If filePath is empty, return inlineContent
func ReadContent(filePath string, inlineContent string) (string, error) {
	// Read from file if filePath is specified
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
		return string(data), nil
	}

	// Return inline content if filePath is empty
	return inlineContent, nil
}

// DetectFormat detects content format from file extension
//
// Parameters:
//   - filePath: file path
//
// Returns:
//   - ContentFormat: detected format
//
// Detection rules:
//   - .md, .markdown → FormatMarkdown
//   - .html, .htm → FormatHTML
//   - .json → FormatLexical
//   - others → FormatUnknown
func DetectFormat(filePath string) ContentFormat {
	if filePath == "" {
		return FormatUnknown
	}

	// Get extension (convert to lowercase)
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".md", ".markdown":
		return FormatMarkdown
	case ".html", ".htm":
		return FormatHTML
	case ".json":
		return FormatLexical
	default:
		return FormatUnknown
	}
}

// ReadContentWithFormat reads file or inline content and returns the format
//
// Parameters:
//   - filePath: file path (use inlineContent if empty)
//   - inlineContent: inline content
//
// Returns:
//   - content: content string
//   - format: detected format (FormatUnknown for inline content)
//   - error: error
//
// Priority:
//  1. If filePath is specified, read from file and detect format
//  2. If filePath is empty, return inlineContent with FormatUnknown
func ReadContentWithFormat(filePath string, inlineContent string) (content string, format ContentFormat, err error) {
	// If filePath is specified
	if filePath != "" {
		// Detect format
		format = DetectFormat(filePath)

		// Read from file
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", FormatUnknown, fmt.Errorf("failed to read file: %w", err)
		}

		return string(data), format, nil
	}

	// Return inline content (format is unknown)
	return inlineContent, FormatUnknown, nil
}
