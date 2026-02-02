/**
 * fields_test.go
 * Tests for field definitions
 *
 * Provides tests for field parser and validation functionality.
 */

package fields

import (
	"strings"
	"testing"
)

// TestParse_Success_CommaSeparated verifies that comma-separated field specifications can be parsed
func TestParse_Success_CommaSeparated(t *testing.T) {
	// Test data: available fields
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status", "html"},
		All:     []string{"id", "title", "status", "html", "slug", "url"},
	}

	// Parse field specification
	result, err := Parse("id,title,status", fieldSet)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Expected value
	expected := []string{"id", "title", "status"}

	// Verify result
	if len(result) != len(expected) {
		t.Errorf("Invalid result length: got=%d, want=%d", len(result), len(expected))
	}
	for i, field := range expected {
		if result[i] != field {
			t.Errorf("Invalid result[%d]: got=%s, want=%s", i, result[i], field)
		}
	}
}

// TestParse_Success_AllSpecification verifies that "all" specification returns all fields
func TestParse_Success_AllSpecification(t *testing.T) {
	// Test data
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status"},
		All:     []string{"id", "title", "status", "html", "slug"},
	}

	// Specify "all"
	result, err := Parse("all", fieldSet)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify all fields are returned
	if len(result) != len(fieldSet.All) {
		t.Errorf("Invalid result length: got=%d, want=%d", len(result), len(fieldSet.All))
	}
	for i, field := range fieldSet.All {
		if result[i] != field {
			t.Errorf("Invalid result[%d]: got=%s, want=%s", i, result[i], field)
		}
	}
}

// TestParse_Error_InvalidField verifies that an error is returned for invalid field specification
func TestParse_Error_InvalidField(t *testing.T) {
	// Test data
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status"},
		All:     []string{"id", "title", "status"},
	}

	// Specify invalid field
	_, err := Parse("id,invalid_field", fieldSet)
	if err == nil {
		t.Fatal("No error was returned")
	}

	// Verify error message contains "invalid_field"
	if !strings.Contains(err.Error(), "invalid_field") {
		t.Errorf("Error message does not contain invalid field name: %v", err)
	}
}

// TestValidate_Success verifies that no error is returned for a valid field list
func TestValidate_Success(t *testing.T) {
	// Test data
	available := []string{"id", "title", "status", "html"}
	fields := []string{"id", "title"}

	// Validation
	err := Validate(fields, available)
	if err != nil {
		t.Errorf("Validation error: %v", err)
	}
}

// TestValidate_Error_InvalidField verifies that an error is returned for invalid fields
func TestValidate_Error_InvalidField(t *testing.T) {
	// Test data
	available := []string{"id", "title", "status"}
	fields := []string{"id", "invalid"}

	// Validation
	err := Validate(fields, available)
	if err == nil {
		t.Fatal("No error was returned")
	}

	// Verify error message contains "invalid"
	if !strings.Contains(err.Error(), "invalid") {
		t.Errorf("Error message does not contain invalid field name: %v", err)
	}
}

// TestListAvailable verifies that the field list can be retrieved as a string
func TestListAvailable(t *testing.T) {
	// Test data
	fieldSet := FieldSet{
		Default: []string{"id", "title"},
		Detail:  []string{"id", "title", "status"},
		All:     []string{"id", "title", "status", "html", "slug"},
	}

	// Get field list
	result := ListAvailable(fieldSet)

	// Verify result is not empty
	if result == "" {
		t.Error("Field list is empty")
	}

	// Verify all fields are included
	for _, field := range fieldSet.All {
		if !strings.Contains(result, field) {
			t.Errorf("Field list does not contain '%s': %s", field, result)
		}
	}
}
