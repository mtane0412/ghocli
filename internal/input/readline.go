/**
 * readline.go
 * Line reading functionality
 *
 * Provides functionality to read lines from io.Reader.
 * Supports both Unix (\n) and Windows (\r\n) line endings.
 */
package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ReadLine reads a single line from io.Reader
//
// Supports both Unix (\n) and Windows (\r\n) line endings.
// A standalone \r is also treated as a line ending.
//
// If EOF is reached before a line ending, the buffered content is returned
// if available (error is nil). If EOF is reached with an empty buffer,
// io.EOF is returned.
//
// r: The input io.Reader
//
// Returns:
//   - The read line (without line ending characters)
//   - An error if reading fails
func ReadLine(r io.Reader) (string, error) {
	br := bufio.NewReader(r)

	var sb strings.Builder

	for {
		b, err := br.ReadByte()
		if err != nil {
			// When EOF is reached
			if errors.Is(err, io.EOF) {
				// Return buffered content if available
				if sb.Len() > 0 {
					return sb.String(), nil
				}

				// Return EOF if buffer is empty
				return "", io.EOF
			}

			// Other errors
			return "", fmt.Errorf("read line: %w", err)
		}

		// Process line ending characters
		if b == '\n' || b == '\r' {
			// For \r, skip the next \n if present (Windows format)
			if b == '\r' {
				if next, _ := br.Peek(1); len(next) == 1 && next[0] == '\n' {
					_, _ = br.ReadByte()
				}
			}

			// Return the line without line ending characters
			return sb.String(), nil
		}

		// Add regular character to buffer
		sb.WriteByte(b)
	}
}
