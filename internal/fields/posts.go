/**
 * posts.go
 * Post field definitions
 *
 * Defines all fields for Ghost Admin API Post/Page resources.
 */

package fields

// PostFields is the field set for Post/Page resources
var PostFields = FieldSet{
	// Default fields for list operations (used in table display)
	Default: []string{
		"id",
		"title",
		"status",
		"created_at",
		"published_at",
	},

	// Detail fields for get operations (used in detail display)
	Detail: []string{
		"id",
		"uuid",
		"title",
		"slug",
		"status",
		"url",
		"excerpt",
		"feature_image",
		"created_at",
		"updated_at",
		"published_at",
		"visibility",
		"featured",
	},

	// All fields (all fields available in Ghost Admin API Post resource)
	All: []string{
		// Basic information
		"id",
		"uuid",
		"title",
		"slug",
		"status",
		"url",

		// Content
		"html",
		"lexical",
		"excerpt",
		"custom_excerpt",

		// Images
		"feature_image",
		"feature_image_alt",
		"feature_image_caption",
		"og_image",
		"twitter_image",

		// SEO
		"meta_title",
		"meta_description",
		"og_title",
		"og_description",
		"twitter_title",
		"twitter_description",
		"canonical_url",

		// Timestamps
		"created_at",
		"updated_at",
		"published_at",

		// Control
		"visibility",
		"featured",
		"email_only",

		// Custom
		"codeinjection_head",
		"codeinjection_foot",
		"custom_template",

		// Relations
		"tags",
		"authors",
		"primary_author",
		"primary_tag",

		// Other
		"comment_id",
		"reading_time",

		// Email/Newsletter
		"email_segment",
		"newsletter_id",
		"send_email_when_published",
	},
}
