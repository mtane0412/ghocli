/**
 * members.go
 * Member field definitions
 *
 * Defines all fields for Ghost Admin API Member resources.
 */

package fields

// MemberFields is the field set for Member resources
var MemberFields = FieldSet{
	// Default fields for list operations (used in table display)
	Default: []string{
		"id",
		"email",
		"name",
		"status",
		"created_at",
	},

	// Detail fields for get operations (used in detail display)
	Detail: []string{
		"id",
		"uuid",
		"email",
		"name",
		"note",
		"status",
		"labels",
		"created_at",
		"updated_at",
	},

	// All fields (all fields available in Ghost Admin API Member resource)
	All: []string{
		"id",
		"uuid",
		"email",
		"name",
		"note",
		"status",
		"labels",
		"created_at",
		"updated_at",
	},
}
