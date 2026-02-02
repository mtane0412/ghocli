/**
 * output_test.go
 * Tests for UI output functionality
 */

package ui

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewOutput_StructExists verifies that NewOutput creates an Output struct
func TestNewOutput_StructExists(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)
	assert.NotNil(t, output)
}

// TestOutput_PrintData tests outputting data to stdout
func TestOutput_PrintData(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	err := output.PrintData("test data")
	assert.NoError(t, err)

	// Data is output to stdout
	assert.Contains(t, stdout.String(), "test data")
	// Nothing is output to stderr
	assert.Empty(t, stderr.String())
}

// TestOutput_PrintMessage tests outputting progress messages to stderr
func TestOutput_PrintMessage(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	output.PrintMessage("test message")

	// Message is output to stderr
	assert.Contains(t, stderr.String(), "test message")
	// Nothing is output to stdout
	assert.Empty(t, stdout.String())
}

// TestOutput_PrintError tests outputting error messages to stderr
func TestOutput_PrintError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	output.PrintError("test error")

	// Error message is output to stderr
	assert.Contains(t, stderr.String(), "test error")
	// Nothing is output to stdout
	assert.Empty(t, stdout.String())
}

// TestWithUI_EmbedUIInContext tests embedding UI in context
func TestWithUI_EmbedUIInContext(t *testing.T) {
	ctx := context.Background()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	ctx = WithUI(ctx, output)

	// Verify UI can be retrieved from context
	retrieved := FromContext(ctx)
	assert.NotNil(t, retrieved)
	assert.Equal(t, output, retrieved)
}

// TestFromContext_ReturnsNilWhenUINotSet tests that FromContext returns nil when UI is not set
func TestFromContext_ReturnsNilWhenUINotSet(t *testing.T) {
	ctx := context.Background()

	retrieved := FromContext(ctx)
	assert.Nil(t, retrieved)
}

// TestFromContext_GetUIFromContextAndOutput tests retrieving UI from context and outputting
func TestFromContext_GetUIFromContextAndOutput(t *testing.T) {
	ctx := context.Background()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	output := NewOutput(stdout, stderr)

	ctx = WithUI(ctx, output)

	// Retrieve UI from context
	ui := FromContext(ctx)
	assert.NotNil(t, ui)

	// Output using UI
	err := ui.PrintData("test data from context")
	assert.NoError(t, err)
	assert.Contains(t, stdout.String(), "test data from context")

	ui.PrintMessage("test message from context")
	assert.Contains(t, stderr.String(), "test message from context")
}
