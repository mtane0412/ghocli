/**
 * outfmt_test.go
 * Test code for output formatting functionality
 */

package outfmt

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

// TestPrintJSON_OutputsInJSONFormat tests JSON format output
func TestPrintJSON_OutputsInJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	data := map[string]interface{}{
		"title": "Test Blog",
		"url":   "https://test.ghost.io",
	}

	if err := formatter.Print(data); err != nil {
		t.Fatalf("Failed to print: %v", err)
	}

	output := buf.String()

	// Verify it can be parsed as JSON
	if !strings.Contains(output, `"title"`) {
		t.Error("JSON does not contain 'title' field")
	}
	if !strings.Contains(output, `"Test Blog"`) {
		t.Error("JSON does not contain 'Test Blog' value")
	}
}

// TestPrintTable_OutputsInTableFormat tests table format output
func TestPrintTable_OutputsInTableFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	headers := []string{"Name", "URL"}
	rows := [][]string{
		{"Site1", "https://site1.ghost.io"},
		{"Site2", "https://site2.ghost.io"},
	}

	if err := formatter.PrintTable(headers, rows); err != nil {
		t.Fatalf("Failed to print table: %v", err)
	}

	output := buf.String()

	// Verify headers and rows are included
	if !strings.Contains(output, "Name") {
		t.Error("Output does not contain header 'Name'")
	}
	if !strings.Contains(output, "Site1") {
		t.Error("Output does not contain 'Site1'")
	}
	if !strings.Contains(output, "Site2") {
		t.Error("Output does not contain 'Site2'")
	}
}

// TestPrintPlain_OutputsInPlainFormat tests plain format (TSV) output
func TestPrintPlain_OutputsInPlainFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "plain")

	headers := []string{"Name", "URL"}
	rows := [][]string{
		{"Site1", "https://site1.ghost.io"},
		{"Site2", "https://site2.ghost.io"},
	}

	if err := formatter.PrintTable(headers, rows); err != nil {
		t.Fatalf("Failed to print plain: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Verify header and data rows exist
	if len(lines) != 3 {
		t.Errorf("Line count = %d; want 3", len(lines))
	}

	// Verify TSV format (tab-separated)
	if !strings.Contains(lines[0], "\t") {
		t.Error("Header row is not tab-separated")
	}
	if !strings.Contains(lines[1], "\t") {
		t.Error("Data row 1 is not tab-separated")
	}
}

// TestPrintMessage_PrintsMessage tests message output
func TestPrintMessage_PrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	message := "Test message"
	formatter.PrintMessage(message)

	output := buf.String()
	if !strings.Contains(output, message) {
		t.Errorf("Output does not contain message: %s", output)
	}
}

// TestPrintError_PrintsError tests error output
func TestPrintError_PrintsError(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	errMsg := "Test error"
	formatter.PrintError(errMsg)

	output := buf.String()
	if !strings.Contains(output, errMsg) {
		t.Errorf("Output does not contain error message: %s", output)
	}
}

// TestPrintTable_TableWithJapaneseStrings tests table display with Japanese strings
func TestPrintTable_TableWithJapaneseStrings(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	headers := []string{"Title", "Status"}
	rows := [][]string{
		{"Non-engineer uncle's dev environment 2026", "published"},
		{"Making a cat tower with 1x4 boards and hemp rope", "published"},
		{"Test", "draft"},
	}

	if err := formatter.PrintTable(headers, rows); err != nil {
		t.Fatalf("Failed to print table: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Header row, 3 data rows, last empty line = 5 lines
	if len(lines) != 5 {
		t.Errorf("Line count = %d; want 5", len(lines))
	}

	// Verify all rows have the same display width
	headerLine := lines[0]
	headerWidth := runewidth.StringWidth(headerLine)
	for i := range rows {
		dataLine := lines[i+1] // After header
		// Verify data row display width matches header row
		dataWidth := runewidth.StringWidth(dataLine)
		if headerWidth != dataWidth {
			t.Errorf("Row %d display width differs from header (header=%d, data=%d)\n  Header: %q\n  Data:   %q",
				i, headerWidth, dataWidth, headerLine, dataLine)
		}
	}
}

// TestPrintKeyValue_OutputsInPlainFormat tests key-value output in plain format
func TestPrintKeyValue_OutputsInPlainFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "plain")

	// Key-value pairs (no headers)
	rows := [][]string{
		{"Title", "Story Seeds"},
		{"URL", "https://hanashinotane.com"},
		{"Version", "5.102"},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		t.Fatalf("Failed to print key-value: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// 3 data rows (no headers)
	if len(lines) != 3 {
		t.Errorf("Line count = %d; want 3", len(lines))
	}

	// Verify TSV format (tab-separated)
	if !strings.Contains(lines[0], "\t") {
		t.Error("Data row 1 is not tab-separated")
	}

	// Verify first row is "Title\tStory Seeds"
	expected := "Title\tStory Seeds"
	if lines[0] != expected {
		t.Errorf("First row = %q; want %q", lines[0], expected)
	}

	// Verify no "Field" or "Value" headers exist
	if strings.Contains(output, "Field") || strings.Contains(output, "Value") {
		t.Error("Output contains headers (Field/Value)")
	}
}

// TestPrintKeyValue_OutputsInJSONFormat tests key-value output in JSON format
func TestPrintKeyValue_OutputsInJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// Key-value pairs
	rows := [][]string{
		{"Title", "Story Seeds"},
		{"URL", "https://hanashinotane.com"},
		{"Version", "5.102"},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		t.Fatalf("Failed to print key-value: %v", err)
	}

	output := buf.String()

	// Verify output as JSON object (not array)
	if strings.HasPrefix(strings.TrimSpace(output), "[") {
		t.Error("JSON output is in array format (should be object format)")
	}

	// Verify each key is included
	if !strings.Contains(output, `"Title"`) {
		t.Error("JSON does not contain 'Title' field")
	}
	if !strings.Contains(output, `"Story Seeds"`) {
		t.Error("JSON does not contain 'Story Seeds' value")
	}
	if !strings.Contains(output, `"URL"`) {
		t.Error("JSON does not contain 'URL' field")
	}
}

// TestPrintKeyValue_OutputsInTableFormat tests key-value output in table format
func TestPrintKeyValue_OutputsInTableFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	// Key-value pairs
	rows := [][]string{
		{"Title", "Story Seeds"},
		{"URL", "https://hanashinotane.com"},
		{"Version", "5.102"},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		t.Fatalf("Failed to print key-value: %v", err)
	}

	// Flush tabwriter
	if err := formatter.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	// 3 data rows, last empty line = 4 lines (no header or separator rows)
	if len(lines) != 4 {
		t.Errorf("Line count = %d; want 4 (got: %v)", len(lines), lines)
	}

	// Verify first row contains "Title"
	if !strings.Contains(lines[0], "Title") {
		t.Error("First row does not contain 'Title'")
	}

	// Verify first row contains "Story Seeds"
	if !strings.Contains(lines[0], "Story Seeds") {
		t.Error("First row does not contain 'Story Seeds'")
	}

	// Verify no tab characters (converted to spaces by tabwriter)
	if strings.Contains(lines[0], "\t") {
		t.Error("Contains tab character (should be aligned by tabwriter)")
	}
}

// TestWithMode_EmbedsModeInContext tests embedding Mode in context
func TestWithMode_EmbedsModeInContext(t *testing.T) {
	tests := []struct {
		name      string
		mode      Mode
		wantJSON  bool
		wantPlain bool
	}{
		{
			name:      "JSON mode",
			mode:      Mode{JSON: true, Plain: false},
			wantJSON:  true,
			wantPlain: false,
		},
		{
			name:      "Plain mode",
			mode:      Mode{JSON: false, Plain: true},
			wantJSON:  false,
			wantPlain: true,
		},
		{
			name:      "Table mode (default)",
			mode:      Mode{JSON: false, Plain: false},
			wantJSON:  false,
			wantPlain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = WithMode(ctx, tt.mode)

			if got := IsJSON(ctx); got != tt.wantJSON {
				t.Errorf("IsJSON() = %v, want %v", got, tt.wantJSON)
			}
			if got := IsPlain(ctx); got != tt.wantPlain {
				t.Errorf("IsPlain() = %v, want %v", got, tt.wantPlain)
			}
		})
	}
}

// TestIsJSON_WhenModeNotSet tests IsJSON when Mode is not set
func TestIsJSON_WhenModeNotSet(t *testing.T) {
	ctx := context.Background()
	if IsJSON(ctx) {
		t.Error("IsJSON() = true, want false (Mode not set)")
	}
}

// TestIsPlain_WhenModeNotSet tests IsPlain when Mode is not set
func TestIsPlain_WhenModeNotSet(t *testing.T) {
	ctx := context.Background()
	if IsPlain(ctx) {
		t.Error("IsPlain() = true, want false (Mode not set)")
	}
}

// TestTableWriter_ReturnsTabwriterInTableMode tests tabwriter return in table mode
func TestTableWriter_ReturnsTabwriterInTableMode(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	// Table mode (JSON=false, Plain=false)
	ctx = WithMode(ctx, Mode{JSON: false, Plain: false})

	w := tableWriter(ctx, &buf)

	// Verify tabwriter is returned
	// In table mode, writer should differ from buf (wrapped by tabwriter)
	if w == &buf {
		t.Error("tableWriter() returned original writer, want tabwriter wrapper")
	}

	// Test write and flush
	_, err := w.Write([]byte("test\tdata\n"))
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
}

// TestTableWriter_ReturnsOriginalWriterInJSONMode tests original writer return in JSON mode
func TestTableWriter_ReturnsOriginalWriterInJSONMode(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	// JSON mode
	ctx = WithMode(ctx, Mode{JSON: true, Plain: false})

	w := tableWriter(ctx, &buf)

	// Verify original writer is returned
	if w != &buf {
		t.Error("tableWriter() in JSON mode should return original writer")
	}
}

// TestTableWriter_ReturnsOriginalWriterInPlainMode tests original writer return in plain mode
func TestTableWriter_ReturnsOriginalWriterInPlainMode(t *testing.T) {
	var buf bytes.Buffer
	ctx := context.Background()
	// Plain mode
	ctx = WithMode(ctx, Mode{JSON: false, Plain: true})

	w := tableWriter(ctx, &buf)

	// Verify original writer is returned
	if w != &buf {
		t.Error("tableWriter() in Plain mode should return original writer")
	}
}
