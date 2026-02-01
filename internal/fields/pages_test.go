/**
 * pages_test.go
 * Page用フィールド定義のテストコード
 */

package fields

import "testing"

// TestPageFields_PostFieldsと同一 はPageFieldsがPostFieldsと同一であることを確認します
func TestPageFields_PostFieldsと同一(t *testing.T) {
	// PageFieldsとPostFieldsが同じインスタンスを参照していることを確認
	if len(PageFields.All) != len(PostFields.All) {
		t.Errorf("PageFields.AllとPostFields.Allの長さが異なります: PageFields=%d, PostFields=%d",
			len(PageFields.All), len(PostFields.All))
	}

	if len(PageFields.Default) != len(PostFields.Default) {
		t.Errorf("PageFields.DefaultとPostFields.Defaultの長さが異なります: PageFields=%d, PostFields=%d",
			len(PageFields.Default), len(PostFields.Default))
	}

	if len(PageFields.Detail) != len(PostFields.Detail) {
		t.Errorf("PageFields.DetailとPostFields.Detailの長さが異なります: PageFields=%d, PostFields=%d",
			len(PageFields.Detail), len(PostFields.Detail))
	}
}

// TestPageFields_基本フィールド はPageFieldsが基本フィールドを含むことを確認します
func TestPageFields_基本フィールド(t *testing.T) {
	// 基本フィールドが存在することを確認
	expectedFields := []string{"id", "uuid", "title", "slug", "status", "url"}

	fieldMap := make(map[string]bool)
	for _, field := range PageFields.All {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("基本フィールド '%s' が見つかりません", expected)
		}
	}
}
