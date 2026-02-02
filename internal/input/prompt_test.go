/**
 * prompt_test.go
 * Test code for user input prompt functionality
 */
package input_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/mtane0412/ghocli/internal/input"
	"github.com/mtane0412/ghocli/internal/ui"
)

// TestPromptLineFrom tests the PromptLineFrom function that reads input from io.Reader
func TestPromptLineFrom(t *testing.T) {
	// Precondition: prepare UI to output prompt to stderr
	var stderr bytes.Buffer
	output := ui.NewOutput(&stderr, &stderr)
	ctx := ui.WithUI(context.Background(), output)

	// Execute: call PromptLineFrom function (read from strings.Reader)
	line, err := input.PromptLineFrom(ctx, "Prompt: ", strings.NewReader("hello\n"))

	// Verify: no error occurs
	if err != nil {
		t.Fatalf("PromptLineFrom function returned an error: %v", err)
	}

	// Verify: correct line is read
	if line != "hello" {
		t.Errorf("Read line differs from expected: got %q, want %q", line, "hello")
	}

	// Verify: prompt is output to stderr
	if !strings.Contains(stderr.String(), "Prompt: ") {
		t.Errorf("Prompt not output to stderr: %q", stderr.String())
	}
}

// TestPromptLine tests the PromptLine function that reads input from os.Stdin
func TestPromptLine(t *testing.T) {
	// Precondition: save os.Stdin and restore it after test
	orig := os.Stdin
	defer func() {
		os.Stdin = orig
	}()

	// Precondition: create a pipe and replace os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	defer func() {
		_ = r.Close()
	}()
	os.Stdin = r

	// Precondition: write input data to the pipe
	_, writeErr := w.WriteString("world\n")
	if writeErr != nil {
		t.Fatalf("Failed to write to pipe: %v", writeErr)
	}
	_ = w.Close()

	// Execute: call PromptLine function
	line, err := input.PromptLine(context.Background(), "Prompt: ")

	// Verify: no error occurs
	if err != nil {
		t.Fatalf("PromptLine function returned an error: %v", err)
	}

	// Verify: correct line is read
	if line != "world" {
		t.Errorf("Read line differs from expected: got %q, want %q", line, "world")
	}
}
