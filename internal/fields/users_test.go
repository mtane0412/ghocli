/**
 * users_test.go
 * Test code for User field definitions
 */

package fields

import "testing"

// TestUserFields_TotalFieldCount verifies that UserFields has the expected number of fields
func TestUserFields_TotalFieldCount(t *testing.T) {
	// User struct has 12 fields
	expectedCount := 12
	if len(UserFields.All) != expectedCount {
		t.Errorf("Incorrect number of UserFields.All fields. expected=%d, got=%d", expectedCount, len(UserFields.All))
	}
}

// TestUserFields_BasicFields verifies that UserFields contains basic fields
func TestUserFields_BasicFields(t *testing.T) {
	// Verify basic fields exist
	expectedFields := []string{"id", "name", "slug", "email"}

	fieldMap := make(map[string]bool)
	for _, field := range UserFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Basic field '%s' not found", expected)
		}
	}
}

// TestUserFields_ProfileFields verifies that UserFields contains profile fields
func TestUserFields_ProfileFields(t *testing.T) {
	// Verify profile fields exist
	expectedFields := []string{"bio", "location", "website", "profile_image", "cover_image"}

	fieldMap := make(map[string]bool)
	for _, field := range UserFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Profile field '%s' not found", expected)
		}
	}
}

// TestUserFields_RolesField verifies that UserFields contains the roles field
func TestUserFields_RolesField(t *testing.T) {
	// Verify roles field exists
	fieldMap := make(map[string]bool)
	for _, field := range UserFields.All {
		fieldMap[field] = true
	}

	if !fieldMap["roles"] {
		t.Errorf("Roles field not found")
	}
}

// TestUserFields_DefaultFields verifies that UserFields.Default contains the expected fields
func TestUserFields_DefaultFields(t *testing.T) {
	// Default fields should be 5
	expectedCount := 5
	if len(UserFields.Default) != expectedCount {
		t.Errorf("Incorrect number of UserFields.Default fields. expected=%d, got=%d", expectedCount, len(UserFields.Default))
	}

	// Verify required fields are included
	expectedFields := []string{"id", "name", "email"}
	fieldMap := make(map[string]bool)
	for _, field := range UserFields.Default {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Default field '%s' not found", expected)
		}
	}
}

// TestUserFields_DetailFields verifies that UserFields.Detail contains the expected fields
func TestUserFields_DetailFields(t *testing.T) {
	// Detail fields should be 12 (all fields)
	expectedCount := 12
	if len(UserFields.Detail) != expectedCount {
		t.Errorf("Incorrect number of UserFields.Detail fields. expected=%d, got=%d", expectedCount, len(UserFields.Detail))
	}

	// Verify bio, location, roles are included
	expectedFields := []string{"bio", "location", "roles"}
	fieldMap := make(map[string]bool)
	for _, field := range UserFields.Detail {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Detail field '%s' not found", expected)
		}
	}
}
