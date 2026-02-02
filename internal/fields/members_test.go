/**
 * members_test.go
 * Test code for Member field definitions
 */

package fields

import "testing"

// TestMemberFields_TotalFieldCount verifies that MemberFields has the expected number of fields
func TestMemberFields_TotalFieldCount(t *testing.T) {
	// Member struct has 9 fields
	expectedCount := 9
	if len(MemberFields.All) != expectedCount {
		t.Errorf("Incorrect number of MemberFields.All fields. expected=%d, got=%d", expectedCount, len(MemberFields.All))
	}
}

// TestMemberFields_BasicFields verifies that MemberFields contains basic fields
func TestMemberFields_BasicFields(t *testing.T) {
	// Verify basic fields exist
	expectedFields := []string{"id", "email", "name", "status"}

	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Basic field '%s' not found", expected)
		}
	}
}

// TestMemberFields_DetailFields verifies that MemberFields contains detail fields
func TestMemberFields_DetailFields(t *testing.T) {
	// Verify detail fields exist
	expectedFields := []string{"uuid", "note", "labels", "created_at", "updated_at"}

	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detail field '%s' not found", expected)
		}
	}
}

// TestMemberFields_DefaultFields verifies that MemberFields.Default contains the expected fields
func TestMemberFields_DefaultFields(t *testing.T) {
	// Default fields should be 5
	expectedCount := 5
	if len(MemberFields.Default) != expectedCount {
		t.Errorf("Incorrect number of MemberFields.Default fields. expected=%d, got=%d", expectedCount, len(MemberFields.Default))
	}

	// Verify required fields are included
	expectedFields := []string{"id", "email", "name"}
	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.Default {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Default field '%s' not found", expected)
		}
	}
}

// TestMemberFields_DetailDisplayFields verifies that MemberFields.Detail contains the expected fields
func TestMemberFields_DetailDisplayFields(t *testing.T) {
	// Detail fields should be 9 (all fields)
	expectedCount := 9
	if len(MemberFields.Detail) != expectedCount {
		t.Errorf("Incorrect number of MemberFields.Detail fields. expected=%d, got=%d", expectedCount, len(MemberFields.Detail))
	}

	// Verify uuid, note, labels are included
	expectedFields := []string{"uuid", "note", "labels"}
	fieldMap := make(map[string]bool)
	for _, field := range MemberFields.Detail {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detail field '%s' not found", expected)
		}
	}
}
