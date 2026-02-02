/**
 * posts_extended_test.go
 * Tests for extended fields of Post structure
 *
 * Provides tests for Post structure that supports all fields of Ghost Admin API.
 */

package ghostapi

import (
	"encoding/json"
	"testing"
)

// TestPost_AllFieldsJSONConversion verifies all fields of Post struct can be converted to/from JSON
func TestPost_AllFieldsJSONConversion(t *testing.T) {
	// Test data: Post with all fields
	post := Post{
		// Basic information
		ID:     "post123",
		UUID:   "uuid123",
		Title:  "Test Article",
		Slug:   "test-post",
		Status: "published",
		URL:    "https://example.com/test-post",

		// Content
		HTML:          "<p>HTML Content</p>",
		Lexical:       "{}",
		Excerpt:       "Excerpt",
		CustomExcerpt: "Custom Excerpt",

		// Images
		FeatureImage:        "https://example.com/image.jpg",
		FeatureImageAlt:     "Image Description",
		FeatureImageCaption: "Image Caption",
		OGImage:             "https://example.com/og.jpg",
		TwitterImage:        "https://example.com/twitter.jpg",

		// SEO
		MetaTitle:          "Meta Title",
		MetaDescription:    "Meta Description",
		OGTitle:            "OG Title",
		OGDescription:      "OG Description",
		TwitterTitle:       "Twitter Title",
		TwitterDescription: "Twitter Description",
		CanonicalURL:       "https://example.com/canonical",

		// Control
		Visibility: "public",
		Featured:   true,
		EmailOnly:  false,

		// Custom
		CodeinjectionHead: "<script>head</script>",
		CodeinjectionFoot: "<script>foot</script>",
		CustomTemplate:    "custom-template",

		// Other
		CommentID:   "comment123",
		ReadingTime: 5,

		// Email/Newsletter
		EmailSegment:           "all",
		NewsletterID:           "newsletter123",
		SendEmailWhenPublished: true,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Restore from JSON
	var restored Post
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify main fields match
	if restored.ID != post.ID {
		t.Errorf("ID does not match: got=%s, want=%s", restored.ID, post.ID)
	}
	if restored.Title != post.Title {
		t.Errorf("Title does not match: got=%s, want=%s", restored.Title, post.Title)
	}
	if restored.FeatureImage != post.FeatureImage {
		t.Errorf("FeatureImage does not match: got=%s, want=%s", restored.FeatureImage, post.FeatureImage)
	}
	if restored.MetaTitle != post.MetaTitle {
		t.Errorf("MetaTitle does not match: got=%s, want=%s", restored.MetaTitle, post.MetaTitle)
	}
	if restored.Visibility != post.Visibility {
		t.Errorf("Visibility does not match: got=%s, want=%s", restored.Visibility, post.Visibility)
	}
	if restored.Featured != post.Featured {
		t.Errorf("Featured does not match: got=%v, want=%v", restored.Featured, post.Featured)
	}
}

// TestPost_RelatedFieldsJSONConversion verifies related fields (tags, authors) of Post struct can be converted to/from JSON
func TestPost_RelatedFieldsJSONConversion(t *testing.T) {
	// Test data: Post with related fields
	post := Post{
		ID:    "post123",
		Title: "Test Article",
		Tags: []Tag{
			{ID: "tag1", Name: "Technology", Slug: "tech"},
			{ID: "tag2", Name: "Go Language", Slug: "golang"},
		},
		Authors: []Author{
			{ID: "author1", Name: "John Smith", Slug: "john-smith"},
		},
		PrimaryAuthor: &Author{
			ID:   "author1",
			Name: "John Smith",
			Slug: "john-smith",
		},
		PrimaryTag: &Tag{
			ID:   "tag1",
			Name: "Technology",
			Slug: "tech",
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Restore from JSON
	var restored Post
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify Tags match
	if len(restored.Tags) != len(post.Tags) {
		t.Errorf("Tags length does not match: got=%d, want=%d", len(restored.Tags), len(post.Tags))
	}
	if len(restored.Tags) > 0 && restored.Tags[0].Name != post.Tags[0].Name {
		t.Errorf("Tags[0].Name does not match: got=%s, want=%s", restored.Tags[0].Name, post.Tags[0].Name)
	}

	// Verify Authors match
	if len(restored.Authors) != len(post.Authors) {
		t.Errorf("Authors length does not match: got=%d, want=%d", len(restored.Authors), len(post.Authors))
	}
	if len(restored.Authors) > 0 && restored.Authors[0].Name != post.Authors[0].Name {
		t.Errorf("Authors[0].Name does not match: got=%s, want=%s", restored.Authors[0].Name, post.Authors[0].Name)
	}

	// Verify PrimaryAuthor matches
	if restored.PrimaryAuthor == nil || restored.PrimaryAuthor.Name != post.PrimaryAuthor.Name {
		t.Errorf("PrimaryAuthor.Name does not match")
	}

	// Verify PrimaryTag matches
	if restored.PrimaryTag == nil || restored.PrimaryTag.Name != post.PrimaryTag.Name {
		t.Errorf("PrimaryTag.Name does not match")
	}
}
