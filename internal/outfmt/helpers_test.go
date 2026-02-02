/**
 * helpers_test.go
 * Test code for output format helper functions
 */

package outfmt

import (
	"testing"

	"github.com/mtane0412/ghocli/internal/ghostapi"
)

// TestFormatAuthors_FormatsAuthorList tests formatting author list
func TestFormatAuthors_FormatsAuthorList(t *testing.T) {
	// Test case: multiple authors
	authors := []ghostapi.Author{
		{ID: "1", Name: "Taro Yamada"},
		{ID: "2", Name: "Hanako Suzuki"},
	}

	// Execute formatting
	result := FormatAuthors(authors)

	// Verify expected value
	expected := "Taro Yamada, Hanako Suzuki"
	if result != expected {
		t.Errorf("FormatAuthors() = %q; want %q", result, expected)
	}
}

// TestFormatAuthors_SingleAuthor tests formatting single author
func TestFormatAuthors_SingleAuthor(t *testing.T) {
	// Test case: single author
	authors := []ghostapi.Author{
		{ID: "1", Name: "Taro Yamada"},
	}

	// Execute formatting
	result := FormatAuthors(authors)

	// Verify expected value
	expected := "Taro Yamada"
	if result != expected {
		t.Errorf("FormatAuthors() = %q; want %q", result, expected)
	}
}

// TestFormatAuthors_NoAuthors tests formatting with no authors
func TestFormatAuthors_NoAuthors(t *testing.T) {
	// Test case: empty slice
	authors := []ghostapi.Author{}

	// Execute formatting
	result := FormatAuthors(authors)

	// Verify expected value (empty string)
	expected := ""
	if result != expected {
		t.Errorf("FormatAuthors() = %q; want %q", result, expected)
	}
}

// TestFormatTags_FormatsTagList tests formatting tag list
func TestFormatTags_FormatsTagList(t *testing.T) {
	// Test case: multiple tags
	tags := []ghostapi.Tag{
		{ID: "1", Name: "Travel"},
		{ID: "2", Name: "Hokkaido"},
		{ID: "3", Name: "Gourmet"},
	}

	// Execute formatting
	result := FormatTags(tags)

	// Verify expected value
	expected := "Travel, Hokkaido, Gourmet"
	if result != expected {
		t.Errorf("FormatTags() = %q; want %q", result, expected)
	}
}

// TestFormatTags_SingleTag tests formatting single tag
func TestFormatTags_SingleTag(t *testing.T) {
	// Test case: single tag
	tags := []ghostapi.Tag{
		{ID: "1", Name: "Travel"},
	}

	// Execute formatting
	result := FormatTags(tags)

	// Verify expected value
	expected := "Travel"
	if result != expected {
		t.Errorf("FormatTags() = %q; want %q", result, expected)
	}
}

// TestFormatTags_NoTags tests formatting with no tags
func TestFormatTags_NoTags(t *testing.T) {
	// Test case: empty slice
	tags := []ghostapi.Tag{}

	// Execute formatting
	result := FormatTags(tags)

	// Verify expected value (empty string)
	expected := ""
	if result != expected {
		t.Errorf("FormatTags() = %q; want %q", result, expected)
	}
}

// TestTruncateExcerpt_TruncatesExcerpt tests truncating excerpt
func TestTruncateExcerpt_TruncatesExcerpt(t *testing.T) {
	// Test case: long string (actual string exceeding 140 characters)
	excerpt := "This is a very long excerpt text. This excerpt exceeds 140 characters and needs to be properly truncated. An ellipsis (...) will be added to the truncated part. This is processing to make it easy for humans and LLMs to read. Adding more characters to exceed 140 characters. Almost reaching 140 characters. Adding a bit more. This should definitely exceed 140 characters now."

	// Execute formatting (max 140 characters)
	result := TruncateExcerpt(excerpt, 140)

	// Verify expected value
	// Verify by rune count (character count)
	resultRunes := []rune(result)
	excerptRunes := []rune(excerpt)

	// Verify original string exceeds 140 characters
	if len(excerptRunes) <= 140 {
		t.Errorf("Test case excerpt is 140 characters or less. len = %d", len(excerptRunes))
	}

	// Verify result is 143 characters (140 + "...")
	if len(resultRunes) != 143 {
		t.Errorf("TruncateExcerpt() length = %d; want %d", len(resultRunes), 143)
	}

	// Verify ending with "..."
	if len(resultRunes) >= 3 {
		suffix := string(resultRunes[len(resultRunes)-3:])
		if suffix != "..." {
			t.Errorf("TruncateExcerpt() ending = %q; want %q", suffix, "...")
		}
	}
}

// TestTruncateExcerpt_ShortString tests with short string
func TestTruncateExcerpt_ShortString(t *testing.T) {
	// Test case: short string
	excerpt := "This is a short excerpt."

	// Execute formatting (max 140 characters)
	result := TruncateExcerpt(excerpt, 140)

	// Verify expected value (returned as is)
	if result != excerpt {
		t.Errorf("TruncateExcerpt() = %q; want %q", result, excerpt)
	}
}

// TestTruncateExcerpt_ExactlyMaxLength tests with exactly maximum length
func TestTruncateExcerpt_ExactlyMaxLength(t *testing.T) {
	// Test case: exactly 140 characters (Japanese 1 character = 3 bytes, but counted by character count)
	excerpt := "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"

	// Execute formatting (max 140 characters)
	result := TruncateExcerpt(excerpt, 140)

	// Verify expected value (returned as is)
	if result != excerpt {
		t.Errorf("TruncateExcerpt() = %q; want %q", result, excerpt)
	}
}

// TestTruncateExcerpt_EmptyString tests with empty string
func TestTruncateExcerpt_EmptyString(t *testing.T) {
	// Test case: empty string
	excerpt := ""

	// Execute formatting
	result := TruncateExcerpt(excerpt, 140)

	// Verify expected value (empty string)
	if result != "" {
		t.Errorf("TruncateExcerpt() = %q; want %q", result, "")
	}
}
