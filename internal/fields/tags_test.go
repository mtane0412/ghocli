/**
 * tags_test.go
 * Test code for Tag field definitions
 */

package fields

import "testing"

// TestTagFields_TotalFieldCount verifies that TagFields has the expected number of fields
func TestTagFields_TotalFieldCount(t *testing.T) {
	// Tag struct has 7 fields
	expectedCount := 7
	if len(TagFields.All) != expectedCount {
		t.Errorf("Incorrect number of TagFields.All fields. expected=%d, got=%d", expectedCount, len(TagFields.All))
	}
}

// TestTagFields_BasicFields verifies that TagFields contains basic fields
func TestTagFields_BasicFields(t *testing.T) {
	// Verify basic fields exist
	expectedFields := []string{"id", "name", "slug", "visibility"}

	fieldMap := make(map[string]bool)
	for _, field := range TagFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Basic field '%s' not found", expected)
		}
	}
}

// TestTagFields_DetailFields verifies that TagFields contains detail fields
func TestTagFields_DetailFields(t *testing.T) {
	// Verify detail fields exist
	expectedFields := []string{"description", "created_at", "updated_at"}

	fieldMap := make(map[string]bool)
	for _, field := range TagFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detail field '%s' not found", expected)
		}
	}
}

// TestTagFields_DefaultFields verifies that TagFields.Default contains the expected fields
func TestTagFields_DefaultFields(t *testing.T) {
	// Default fields should be 5
	expectedCount := 5
	if len(TagFields.Default) != expectedCount {
		t.Errorf("Incorrect number of TagFields.Default fields. expected=%d, got=%d", expectedCount, len(TagFields.Default))
	}

	// Verify required fields are included
	expectedFields := []string{"id", "name", "slug"}
	fieldMap := make(map[string]bool)
	for _, field := range TagFields.Default {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Default field '%s' not found", expected)
		}
	}
}

// TestTagFields_DetailDisplayFields verifies that TagFields.Detail contains the expected fields
func TestTagFields_DetailDisplayFields(t *testing.T) {
	// Detail fields should be 7 (all fields)
	expectedCount := 7
	if len(TagFields.Detail) != expectedCount {
		t.Errorf("Incorrect number of TagFields.Detail fields. expected=%d, got=%d", expectedCount, len(TagFields.Detail))
	}

	// Verify description and updated_at are included
	expectedFields := []string{"description", "updated_at"}
	fieldMap := make(map[string]bool)
	for _, field := range TagFields.Detail {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detail field '%s' not found", expected)
		}
	}
}
