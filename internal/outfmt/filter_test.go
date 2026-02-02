/**
 * filter_test.go
 * Test code for field filtering functionality
 *
 * Provides tests for functionality that outputs only specified fields.
 */

package outfmt

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestFilterFields_JSONOutput verifies outputting only specified fields in JSON format
func TestFilterFields_JSONOutput(t *testing.T) {
	// Test data
	data := map[string]interface{}{
		"id":     "abc123",
		"title":  "Test Article",
		"status": "published",
		"html":   "<p>HTML Content</p>",
		"slug":   "test-post",
	}

	// Output buffer
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// Output with field specification
	fields := []string{"id", "title", "status"}
	err := FilterFields(formatter, data, fields)
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Verify result
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify only specified fields are included
	if result["id"] != "abc123" {
		t.Errorf("id is not included")
	}
	if result["title"] != "Test Article" {
		t.Errorf("title is not included")
	}
	if result["status"] != "published" {
		t.Errorf("status is not included")
	}

	// Verify unspecified fields are not included
	if _, ok := result["html"]; ok {
		t.Errorf("html is included (should be excluded)")
	}
	if _, ok := result["slug"]; ok {
		t.Errorf("slug is included (should be excluded)")
	}
}

// TestFilterFields_SliceJSONOutput verifies outputting only specified fields in slice data
func TestFilterFields_SliceJSONOutput(t *testing.T) {
	// Test data (multiple items)
	data := []map[string]interface{}{
		{
			"id":     "abc123",
			"title":  "Article 1",
			"status": "published",
			"html":   "<p>HTML1</p>",
		},
		{
			"id":     "def456",
			"title":  "Article 2",
			"status": "draft",
			"html":   "<p>HTML2</p>",
		},
	}

	// Output buffer
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// Output with field specification
	fields := []string{"id", "title"}
	err := FilterFields(formatter, data, fields)
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Verify result
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify element count
	if len(result) != 2 {
		t.Fatalf("Invalid element count: got=%d, want=2", len(result))
	}

	// Verify first element
	if result[0]["id"] != "abc123" {
		t.Errorf("result[0].id is invalid")
	}
	if result[0]["title"] != "Article 1" {
		t.Errorf("result[0].title is invalid")
	}
	if _, ok := result[0]["status"]; ok {
		t.Errorf("result[0].status is included (should be excluded)")
	}

	// Verify second element
	if result[1]["id"] != "def456" {
		t.Errorf("result[1].id is invalid")
	}
	if result[1]["title"] != "Article 2" {
		t.Errorf("result[1].title is invalid")
	}
}

// TestFilterFields_PlainOutput verifies outputting only specified fields in Plain format (TSV)
func TestFilterFields_PlainOutput(t *testing.T) {
	// Test data
	data := []map[string]interface{}{
		{
			"id":     "abc123",
			"title":  "Article 1",
			"status": "published",
		},
		{
			"id":     "def456",
			"title":  "Article 2",
			"status": "draft",
		},
	}

	// Output buffer
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "plain")

	// Output with field specification
	fields := []string{"id", "title"}
	err := FilterFields(formatter, data, fields)
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Verify result
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("Invalid line count: got=%d, want=3 (header + 2 data rows)", len(lines))
	}

	// Verify header row
	header := lines[0]
	if !strings.Contains(header, "id") || !strings.Contains(header, "title") {
		t.Errorf("Invalid header: %s", header)
	}
	if strings.Contains(header, "status") {
		t.Errorf("Header contains status (should be excluded): %s", header)
	}

	// Verify data rows
	if !strings.Contains(lines[1], "abc123") || !strings.Contains(lines[1], "Article 1") {
		t.Errorf("Invalid data in row 1: %s", lines[1])
	}
	if !strings.Contains(lines[2], "def456") || !strings.Contains(lines[2], "Article 2") {
		t.Errorf("Invalid data in row 2: %s", lines[2])
	}
}

// TestFilterFields_NoFieldsSpecified verifies outputting all fields when no fields are specified
func TestFilterFields_NoFieldsSpecified(t *testing.T) {
	// Test data
	data := map[string]interface{}{
		"id":    "abc123",
		"title": "Test Article",
	}

	// Output buffer
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// Output without field specification (nil or empty slice)
	err := FilterFields(formatter, data, nil)
	if err != nil {
		t.Fatalf("FilterFields failed: %v", err)
	}

	// Verify result
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify all fields are included
	if result["id"] != "abc123" {
		t.Errorf("id is not included")
	}
	if result["title"] != "Test Article" {
		t.Errorf("title is not included")
	}
}
