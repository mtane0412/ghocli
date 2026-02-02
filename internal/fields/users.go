/**
 * users.go
 * User field definitions
 *
 * Defines all fields for Ghost Admin API User resources.
 */

package fields

// UserFields is the field set for User resources
var UserFields = FieldSet{
	// Default fields for list operations (used in table display)
	Default: []string{
		"id",
		"name",
		"slug",
		"email",
		"created_at",
	},

	// Detail fields for get operations (used in detail display)
	Detail: []string{
		"id",
		"name",
		"slug",
		"email",
		"bio",
		"location",
		"website",
		"profile_image",
		"cover_image",
		"roles",
		"created_at",
		"updated_at",
	},

	// All fields (all fields available in Ghost Admin API User resource)
	All: []string{
		"id",
		"name",
		"slug",
		"email",
		"bio",
		"location",
		"website",
		"profile_image",
		"cover_image",
		"roles",
		"created_at",
		"updated_at",
	},
}
