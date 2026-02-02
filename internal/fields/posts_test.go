/**
 * posts_test.go
 * Tests for Post field definitions
 *
 * Provides tests for Post field set definitions.
 */

package fields

import (
	"testing"
)

// TestPostFields_TotalFieldCount verifies the total number of Post fields
func TestPostFields_TotalFieldCount(t *testing.T) {
	// Verify PostFields is defined
	if PostFields.All == nil {
		t.Fatal("PostFields.All is not defined")
	}

	// Verify there are at least 40 fields (as planned)
	if len(PostFields.All) < 40 {
		t.Errorf("Insufficient total fields in PostFields: got=%d, want>=40", len(PostFields.All))
	}
}

// TestPostFields_BasicFields verifies that basic fields are included
func TestPostFields_BasicFields(t *testing.T) {
	// Basic fields
	basicFields := []string{"id", "uuid", "title", "slug", "status", "url"}

	// Verify each basic field is included in All
	for _, field := range basicFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.All does not contain '%s'", field)
		}
	}
}

// TestPostFields_ContentFields verifies that content-related fields are included
func TestPostFields_ContentFields(t *testing.T) {
	// Content fields
	contentFields := []string{"html", "lexical", "excerpt", "custom_excerpt"}

	// Verify each content field is included in All
	for _, field := range contentFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.All does not contain '%s'", field)
		}
	}
}

// TestPostFields_ImageFields verifies that image-related fields are included
func TestPostFields_ImageFields(t *testing.T) {
	// Image fields
	imageFields := []string{"feature_image", "feature_image_alt", "feature_image_caption", "og_image", "twitter_image"}

	// Verify each image field is included in All
	for _, field := range imageFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.All does not contain '%s'", field)
		}
	}
}

// TestPostFields_SEOFields verifies that SEO-related fields are included
func TestPostFields_SEOFields(t *testing.T) {
	// SEO fields
	seoFields := []string{
		"meta_title", "meta_description",
		"og_title", "og_description",
		"twitter_title", "twitter_description",
		"canonical_url",
	}

	// Verify each SEO field is included in All
	for _, field := range seoFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.All does not contain '%s'", field)
		}
	}
}

// TestPostFields_DateFields verifies that date-related fields are included
func TestPostFields_DateFields(t *testing.T) {
	// Date fields
	dateFields := []string{"created_at", "updated_at", "published_at"}

	// Verify each date field is included in All
	for _, field := range dateFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.All does not contain '%s'", field)
		}
	}
}

// TestPostFields_ControlFields verifies that control-related fields are included
func TestPostFields_ControlFields(t *testing.T) {
	// Control fields
	controlFields := []string{"visibility", "featured", "email_only"}

	// Verify each control field is included in All
	for _, field := range controlFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.All does not contain '%s'", field)
		}
	}
}

// TestPostFields_DefaultFields verifies that default fields are properly configured
func TestPostFields_DefaultFields(t *testing.T) {
	// Verify Default fields are set
	if PostFields.Default == nil {
		t.Fatal("PostFields.Default is not defined")
	}

	// Verify default field count is appropriate (5-10 fields)
	if len(PostFields.Default) < 3 || len(PostFields.Default) > 10 {
		t.Errorf("Inappropriate number of PostFields.Default fields: got=%d", len(PostFields.Default))
	}

	// Verify basic fields are included in Default
	requiredDefaults := []string{"id", "title", "status"}
	for _, field := range requiredDefaults {
		if !contains(PostFields.Default, field) {
			t.Errorf("PostFields.Default does not contain '%s'", field)
		}
	}
}

// TestPostFields_DetailFields verifies that detail fields are properly configured
func TestPostFields_DetailFields(t *testing.T) {
	// Verify Detail fields are set
	if PostFields.Detail == nil {
		t.Fatal("PostFields.Detail is not defined")
	}

	// Verify detail field count is greater than Default
	if len(PostFields.Detail) <= len(PostFields.Default) {
		t.Errorf("PostFields.Detail is not greater than Default: Detail=%d, Default=%d",
			len(PostFields.Detail), len(PostFields.Default))
	}
}

// Helper function: Check if an item is contained in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
