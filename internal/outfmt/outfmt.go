/**
 * outfmt.go
 * Output format functionality
 *
 * Supports output in JSON, table, and plain (TSV) formats.
 */

package outfmt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/mattn/go-runewidth"
)

// Formatter is the output formatter
type Formatter struct {
	writer    io.Writer
	tabwriter *tabwriter.Writer
	mode      string // "json", "table", "plain"
}

// Mode represents the output format mode
type Mode struct {
	// JSON determines whether to output in JSON format
	JSON bool
	// Plain determines whether to output in plain format (TSV)
	Plain bool
}

// context key type
type contextKey int

const (
	modeKey contextKey = iota
)

// WithMode sets the output mode in the context
func WithMode(ctx context.Context, mode Mode) context.Context {
	return context.WithValue(ctx, modeKey, mode)
}

// getMode retrieves the output mode from the context
func getMode(ctx context.Context) Mode {
	if mode, ok := ctx.Value(modeKey).(Mode); ok {
		return mode
	}
	// Default is table mode
	return Mode{JSON: false, Plain: false}
}

// IsJSON returns whether to output in JSON format
func IsJSON(ctx context.Context) bool {
	return getMode(ctx).JSON
}

// IsPlain returns whether to output in plain format (TSV)
func IsPlain(ctx context.Context) bool {
	return getMode(ctx).Plain
}

// tableWriter wraps with tabwriter in table mode, returns as is otherwise
func tableWriter(ctx context.Context, w io.Writer) io.Writer {
	mode := getMode(ctx)
	// Return as is for JSON or Plain mode
	if mode.JSON || mode.Plain {
		return w
	}
	// Wrap with tabwriter for table mode
	return tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
}

// NewFormatter creates a new output formatter.
func NewFormatter(writer io.Writer, mode string) *Formatter {
	f := &Formatter{
		writer: writer,
		mode:   mode,
	}

	// Wrap with tabwriter for table format
	if mode == "table" {
		f.tabwriter = tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	}

	return f
}

// Flush flushes buffered output.
// Necessary to flush tabwriter for table format.
func (f *Formatter) Flush() error {
	if f.tabwriter != nil {
		return f.tabwriter.Flush()
	}
	return nil
}

// getWriter returns the destination writer.
func (f *Formatter) getWriter() io.Writer {
	if f.tabwriter != nil {
		return f.tabwriter
	}
	return f.writer
}

// Print outputs arbitrary data.
// Outputs as JSON for JSON format.
func (f *Formatter) Print(data interface{}) error {
	if f.mode == "json" {
		encoder := json.NewEncoder(f.writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	}

	// Default is standard output
	_, err := fmt.Fprintln(f.writer, data)
	return err
}

// PrintTable outputs data in table format.
func (f *Formatter) PrintTable(headers []string, rows [][]string) error {
	switch f.mode {
	case "plain":
		// Output in TSV format
		return f.printTSV(headers, rows)
	case "json":
		// Output as JSON array
		return f.printJSONTable(headers, rows)
	default:
		// Output in table format
		return f.printTableFormat(headers, rows)
	}
}

// printTSV outputs in TSV format (tab-separated).
func (f *Formatter) printTSV(headers []string, rows [][]string) error {
	// Output header row
	if _, err := fmt.Fprintln(f.writer, strings.Join(headers, "\t")); err != nil {
		return err
	}

	// Output data rows
	for _, row := range rows {
		if _, err := fmt.Fprintln(f.writer, strings.Join(row, "\t")); err != nil {
			return err
		}
	}

	return nil
}

// printJSONTable outputs in JSON array format.
func (f *Formatter) printJSONTable(headers []string, rows [][]string) error {
	// Convert each row to a map
	var data []map[string]string
	for _, row := range rows {
		item := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			}
		}
		data = append(data, item)
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printTableFormat outputs in table format (human-readable).
func (f *Formatter) printTableFormat(headers []string, rows [][]string) error {
	// Calculate maximum display width for each column (considering full-width characters)
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = runewidth.StringWidth(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				cellWidth := runewidth.StringWidth(cell)
				if cellWidth > colWidths[i] {
					colWidths[i] = cellWidth
				}
			}
		}
	}

	// Output header row
	for i, header := range headers {
		if i > 0 {
			fmt.Fprint(f.writer, "  ")
		}
		// Add padding based on display width
		fmt.Fprint(f.writer, header)
		padding := colWidths[i] - runewidth.StringWidth(header)
		if padding > 0 {
			fmt.Fprint(f.writer, strings.Repeat(" ", padding))
		}
	}
	fmt.Fprintln(f.writer)

	// Output data rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(f.writer, "  ")
			}
			if i < len(colWidths) {
				// Add padding based on display width
				fmt.Fprint(f.writer, cell)
				padding := colWidths[i] - runewidth.StringWidth(cell)
				if padding > 0 {
					fmt.Fprint(f.writer, strings.Repeat(" ", padding))
				}
			} else {
				fmt.Fprint(f.writer, cell)
			}
		}
		fmt.Fprintln(f.writer)
	}

	return nil
}

// PrintMessage outputs a message.
func (f *Formatter) PrintMessage(message string) {
	fmt.Fprintln(f.writer, message)
}

// PrintError outputs an error message.
func (f *Formatter) PrintError(message string) {
	fmt.Fprintln(f.writer, "Error:", message)
}

// PrintKeyValue outputs key/value pairs without headers.
// Used for displaying single item information.
func (f *Formatter) PrintKeyValue(rows [][]string) error {
	if f.mode == "json" {
		// Output as key/value map
		data := make(map[string]string)
		for _, row := range rows {
			if len(row) >= 2 {
				data[row[0]] = row[1]
			}
		}
		encoder := json.NewEncoder(f.writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	}

	// Common for Plain/table format: output tab-separated
	// Automatically aligned by tabwriter in table format
	w := f.getWriter()
	for _, row := range rows {
		if _, err := fmt.Fprintln(w, strings.Join(row, "\t")); err != nil {
			return err
		}
	}

	return nil
}
