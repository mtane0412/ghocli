/**
 * tags.go
 * Tag field definitions
 *
 * Defines all fields for Ghost Admin API Tag resources.
 */

package fields

// TagFields is the field set for Tag resources
var TagFields = FieldSet{
	// Default fields for list operations (used in table display)
	Default: []string{
		"id",
		"name",
		"slug",
		"visibility",
		"created_at",
	},

	// Detail fields for get operations (used in detail display)
	Detail: []string{
		"id",
		"name",
		"slug",
		"description",
		"visibility",
		"created_at",
		"updated_at",
	},

	// All fields (all fields available in Ghost Admin API Tag resource)
	All: []string{
		"id",
		"name",
		"slug",
		"description",
		"visibility",
		"created_at",
		"updated_at",
	},
}
