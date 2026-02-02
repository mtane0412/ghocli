/**
 * fields.go
 * Field definition infrastructure
 *
 * Provides field set definition, parsing, and validation functionality.
 */

package fields

import (
	"fmt"
	"strings"
)

// FieldSet represents a resource field set
type FieldSet struct {
	// Default is the default fields for list operations
	Default []string
	// Detail is the default fields for get operations
	Detail []string
	// All is all available fields
	All []string
}

// Parse parses a comma-separated field specification string
//
// Input examples:
//   - "id,title,status" -> []string{"id", "title", "status"}
//   - "all" -> fieldSet.All
//   - "" -> nil (indicates to use default fields)
func Parse(input string, fieldSet FieldSet) ([]string, error) {
	// Return nil for empty string (use default fields)
	if input == "" {
		return nil, nil
	}

	// Return all fields for "all"
	if input == "all" {
		result := make([]string, len(fieldSet.All))
		copy(result, fieldSet.All)
		return result, nil
	}

	// Parse comma-separated values
	fields := strings.Split(input, ",")

	// Trim each field
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	// Validate
	if err := Validate(fields, fieldSet.All); err != nil {
		return nil, err
	}

	return fields, nil
}

// Validate verifies whether the specified fields are available
func Validate(fields []string, available []string) error {
	// Convert available fields to map
	availableMap := make(map[string]bool)
	for _, field := range available {
		availableMap[field] = true
	}

	// Check if each field is available
	for _, field := range fields {
		if !availableMap[field] {
			return fmt.Errorf("unknown field '%s'. Available fields: %s", field, strings.Join(available, ", "))
		}
	}

	return nil
}

// ListAvailable returns a list of available fields as a string
func ListAvailable(fieldSet FieldSet) string {
	return "Specify fields with --fields. Available fields: " + strings.Join(fieldSet.All, ", ")
}
