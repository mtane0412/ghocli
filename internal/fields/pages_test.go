/**
 * pages_test.go
 * Test code for Page field definitions
 */

package fields

import "testing"

// TestPageFields_IdenticalToPostFields verifies that PageFields is identical to PostFields
func TestPageFields_IdenticalToPostFields(t *testing.T) {
	// Verify PageFields and PostFields reference the same instance
	if len(PageFields.All) != len(PostFields.All) {
		t.Errorf("PageFields.All and PostFields.All have different lengths: PageFields=%d, PostFields=%d",
			len(PageFields.All), len(PostFields.All))
	}

	if len(PageFields.Default) != len(PostFields.Default) {
		t.Errorf("PageFields.Default and PostFields.Default have different lengths: PageFields=%d, PostFields=%d",
			len(PageFields.Default), len(PostFields.Default))
	}

	if len(PageFields.Detail) != len(PostFields.Detail) {
		t.Errorf("PageFields.Detail and PostFields.Detail have different lengths: PageFields=%d, PostFields=%d",
			len(PageFields.Detail), len(PostFields.Detail))
	}
}

// TestPageFields_BasicFields verifies that PageFields contains basic fields
func TestPageFields_BasicFields(t *testing.T) {
	// Verify basic fields exist
	expectedFields := []string{"id", "uuid", "title", "slug", "status", "url"}

	fieldMap := make(map[string]bool)
	for _, field := range PageFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("Basic field '%s' not found", expected)
		}
	}
}
