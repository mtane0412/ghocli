/**
 * helpers.go
 * Output format helper functions
 *
 * Provides common helper functions used in default display for posts/pages commands.
 */

package outfmt

import (
	"strings"

	"github.com/mtane0412/ghocli/internal/ghostapi"
)

// FormatAuthors formats a list of authors as comma-separated names
func FormatAuthors(authors []ghostapi.Author) string {
	// Return empty string if no authors
	if len(authors) == 0 {
		return ""
	}

	// Collect author names
	names := make([]string, len(authors))
	for i, author := range authors {
		names[i] = author.Name
	}

	// Join with comma separator
	return strings.Join(names, ", ")
}

// FormatTags formats a list of tags as comma-separated names
func FormatTags(tags []ghostapi.Tag) string {
	// Return empty string if no tags
	if len(tags) == 0 {
		return ""
	}

	// Collect tag names
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}

	// Join with comma separator
	return strings.Join(names, ", ")
}

// TruncateExcerpt truncates an excerpt to a specified number of characters
// If it exceeds maxLen, truncate to maxLen characters and append "..."
func TruncateExcerpt(excerpt string, maxLen int) string {
	// Return as is for empty string
	if excerpt == "" {
		return ""
	}

	// Convert string to rune (Unicode code point) slice
	// To correctly count multibyte characters like Japanese
	runes := []rune(excerpt)

	// Return as is if within maxLen
	if len(runes) <= maxLen {
		return excerpt
	}

	// Truncate to maxLen characters and append "..."
	return string(runes[:maxLen]) + "..."
}
