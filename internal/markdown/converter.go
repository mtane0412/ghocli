/**
 * Markdown to HTML conversion functionality
 *
 * This package provides functionality to convert Markdown text to HTML format.
 * Uses the goldmark library for safe and fast conversion.
 */
package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// ConvertToHTML converts a Markdown string to HTML
//
// Parameters:
//   - markdown: source Markdown string
//
// Returns:
//   - string: converted HTML string
//   - error: conversion error (usually returns nil)
//
// Example usage:
//
//	html, err := ConvertToHTML("# Heading\n\nThis is a paragraph.")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(html)
func ConvertToHTML(markdown string) (string, error) {
	// Return as is for empty string
	if markdown == "" {
		return "", nil
	}

	// Prepare buffer
	var buf bytes.Buffer

	// Convert Markdown to HTML using goldmark
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
