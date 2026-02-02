/**
 * filter.go
 * Field filtering functionality
 *
 * Provides functionality to output only specified fields.
 */

package outfmt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// FilterFields filters and outputs only the specified fields
//
// If data is a map[string]interface{} or a slice thereof, it extracts only the specified fields.
// If fields is nil or empty, all fields are output as is.
func FilterFields(formatter *Formatter, data interface{}, fields []string) error {
	// If no field specification, output as is
	if len(fields) == 0 {
		return formatter.Print(data)
	}

	// Filter based on data type
	switch v := data.(type) {
	case map[string]interface{}:
		// Filter a single map
		filtered := filterMap(v, fields)
		return formatter.Print(filtered)

	case []map[string]interface{}:
		// Filter each element of the slice
		filtered := make([]map[string]interface{}, len(v))
		for i, item := range v {
			filtered[i] = filterMap(item, fields)
		}

		// Output based on mode
		if formatter.mode == "plain" {
			return formatter.PrintTable(fields, mapSliceToRows(filtered, fields))
		}
		return formatter.Print(filtered)

	case []interface{}:
		// Convert interface{} slice to map[string]interface{} slice
		var mapSlice []map[string]interface{}
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				mapSlice = append(mapSlice, m)
			}
		}
		if len(mapSlice) > 0 {
			filtered := make([]map[string]interface{}, len(mapSlice))
			for i, item := range mapSlice {
				filtered[i] = filterMap(item, fields)
			}

			// Output based on mode
			if formatter.mode == "plain" {
				return formatter.PrintTable(fields, mapSliceToRows(filtered, fields))
			}
			return formatter.Print(filtered)
		}

		// Output as is if conversion fails
		return formatter.Print(v)

	default:
		// Filter struct using reflection
		return filterStruct(formatter, data, fields)
	}
}

// filterMap extracts only the specified fields from a map
func filterMap(m map[string]interface{}, fields []string) map[string]interface{} {
	filtered := make(map[string]interface{})
	for _, field := range fields {
		if value, ok := m[field]; ok {
			filtered[field] = value
		}
	}
	return filtered
}

// mapSliceToRows converts a map slice to table rows
func mapSliceToRows(data []map[string]interface{}, fields []string) [][]string {
	rows := make([][]string, len(data))
	for i, item := range data {
		row := make([]string, len(fields))
		for j, field := range fields {
			if value, ok := item[field]; ok {
				row[j] = fmt.Sprintf("%v", value)
			} else {
				row[j] = ""
			}
		}
		rows[i] = row
	}
	return rows
}

// filterStruct extracts only the specified fields from a struct
func filterStruct(formatter *Formatter, data interface{}, fields []string) error {
	// Convert struct to map
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	var m interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Filter if conversion to map succeeded
	switch v := m.(type) {
	case map[string]interface{}:
		filtered := filterMap(v, fields)
		return formatter.Print(filtered)
	case []interface{}:
		// For slices
		var mapSlice []map[string]interface{}
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				mapSlice = append(mapSlice, itemMap)
			}
		}
		if len(mapSlice) > 0 {
			filtered := make([]map[string]interface{}, len(mapSlice))
			for i, item := range mapSlice {
				filtered[i] = filterMap(item, fields)
			}

			// Output based on mode
			if formatter.mode == "plain" {
				return formatter.PrintTable(fields, mapSliceToRows(filtered, fields))
			}
			return formatter.Print(filtered)
		}
	}

	// Output as is if filtering fails
	return formatter.Print(data)
}

// StructToMap converts a struct to map[string]interface{}
func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Get field name from JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Remove options like "omitempty"
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		// Skip if value is empty (for omitempty)
		if strings.Contains(jsonTag, "omitempty") && value.IsZero() {
			continue
		}

		result[fieldName] = value.Interface()
	}

	return result, nil
}
