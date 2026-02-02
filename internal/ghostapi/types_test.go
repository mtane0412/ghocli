/**
 * types_test.go
 * Tests for common type definitions
 *
 * Provides tests for common types such as Author and Tag.
 */

package ghostapi

import (
	"encoding/json"
	"testing"
)

// TestAuthor_JSONConversion verifies Author struct can be converted to/from JSON
func TestAuthor_JSONConversion(t *testing.T) {
	// Test data
	author := Author{
		ID:    "abc123",
		Name:  "John Smith",
		Slug:  "john-smith",
		Email: "john@example.com",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(author)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Restore from JSON
	var restored Author
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify values match
	if restored.ID != author.ID {
		t.Errorf("ID does not match: got=%s, want=%s", restored.ID, author.ID)
	}
	if restored.Name != author.Name {
		t.Errorf("Name does not match: got=%s, want=%s", restored.Name, author.Name)
	}
	if restored.Slug != author.Slug {
		t.Errorf("Slug does not match: got=%s, want=%s", restored.Slug, author.Slug)
	}
	if restored.Email != author.Email {
		t.Errorf("Email does not match: got=%s, want=%s", restored.Email, author.Email)
	}
}

// TestTag_JSONConversion verifies Tag struct can be converted to/from JSON
func TestTag_JSONConversion(t *testing.T) {
	// Test data
	tag := Tag{
		ID:          "tag123",
		Name:        "Technology",
		Slug:        "tech",
		Description: "Technology-related articles",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Restore from JSON
	var restored Tag
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify values match
	if restored.ID != tag.ID {
		t.Errorf("ID does not match: got=%s, want=%s", restored.ID, tag.ID)
	}
	if restored.Name != tag.Name {
		t.Errorf("Name does not match: got=%s, want=%s", restored.Name, tag.Name)
	}
	if restored.Slug != tag.Slug {
		t.Errorf("Slug does not match: got=%s, want=%s", restored.Slug, tag.Slug)
	}
	if restored.Description != tag.Description {
		t.Errorf("Description does not match: got=%s, want=%s", restored.Description, tag.Description)
	}
}
