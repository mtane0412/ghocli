/**
 * filter_test.go
 * フィールドフィルタリング機能のテスト
 *
 * 指定されたフィールドのみを出力する機能のテストを提供します。
 */

package outfmt

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestFilterFields_JSON出力 はJSON形式で指定フィールドのみ出力できることを確認します
func TestFilterFields_JSON出力(t *testing.T) {
	// テストデータ
	data := map[string]interface{}{
		"id":     "abc123",
		"title":  "テスト記事",
		"status": "published",
		"html":   "<p>HTMLコンテンツ</p>",
		"slug":   "test-post",
	}

	// 出力先バッファ
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// フィールド指定で出力
	fields := []string{"id", "title", "status"}
	err := FilterFields(formatter, data, fields)
	if err != nil {
		t.Fatalf("FilterFieldsに失敗: %v", err)
	}

	// 結果を検証
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSONパースに失敗: %v", err)
	}

	// 指定したフィールドのみ含まれることを確認
	if result["id"] != "abc123" {
		t.Errorf("idが含まれていません")
	}
	if result["title"] != "テスト記事" {
		t.Errorf("titleが含まれていません")
	}
	if result["status"] != "published" {
		t.Errorf("statusが含まれていません")
	}

	// 指定していないフィールドが含まれないことを確認
	if _, ok := result["html"]; ok {
		t.Errorf("htmlが含まれています（除外されるべき）")
	}
	if _, ok := result["slug"]; ok {
		t.Errorf("slugが含まれています（除外されるべき）")
	}
}

// TestFilterFields_スライスJSON出力 はスライスデータで指定フィールドのみ出力できることを確認します
func TestFilterFields_スライスJSON出力(t *testing.T) {
	// テストデータ（複数のアイテム）
	data := []map[string]interface{}{
		{
			"id":     "abc123",
			"title":  "記事1",
			"status": "published",
			"html":   "<p>HTML1</p>",
		},
		{
			"id":     "def456",
			"title":  "記事2",
			"status": "draft",
			"html":   "<p>HTML2</p>",
		},
	}

	// 出力先バッファ
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// フィールド指定で出力
	fields := []string{"id", "title"}
	err := FilterFields(formatter, data, fields)
	if err != nil {
		t.Fatalf("FilterFieldsに失敗: %v", err)
	}

	// 結果を検証
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSONパースに失敗: %v", err)
	}

	// 要素数を確認
	if len(result) != 2 {
		t.Fatalf("要素数が不正: got=%d, want=2", len(result))
	}

	// 1つ目の要素を確認
	if result[0]["id"] != "abc123" {
		t.Errorf("result[0].idが不正")
	}
	if result[0]["title"] != "記事1" {
		t.Errorf("result[0].titleが不正")
	}
	if _, ok := result[0]["status"]; ok {
		t.Errorf("result[0].statusが含まれています（除外されるべき）")
	}

	// 2つ目の要素を確認
	if result[1]["id"] != "def456" {
		t.Errorf("result[1].idが不正")
	}
	if result[1]["title"] != "記事2" {
		t.Errorf("result[1].titleが不正")
	}
}

// TestFilterFields_Plain出力 はPlain形式（TSV）で指定フィールドのみ出力できることを確認します
func TestFilterFields_Plain出力(t *testing.T) {
	// テストデータ
	data := []map[string]interface{}{
		{
			"id":     "abc123",
			"title":  "記事1",
			"status": "published",
		},
		{
			"id":     "def456",
			"title":  "記事2",
			"status": "draft",
		},
	}

	// 出力先バッファ
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "plain")

	// フィールド指定で出力
	fields := []string{"id", "title"}
	err := FilterFields(formatter, data, fields)
	if err != nil {
		t.Fatalf("FilterFieldsに失敗: %v", err)
	}

	// 結果を検証
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("行数が不正: got=%d, want=3（ヘッダー+データ2行）", len(lines))
	}

	// ヘッダー行を確認
	header := lines[0]
	if !strings.Contains(header, "id") || !strings.Contains(header, "title") {
		t.Errorf("ヘッダーが不正: %s", header)
	}
	if strings.Contains(header, "status") {
		t.Errorf("ヘッダーにstatusが含まれています（除外されるべき）: %s", header)
	}

	// データ行を確認
	if !strings.Contains(lines[1], "abc123") || !strings.Contains(lines[1], "記事1") {
		t.Errorf("1行目のデータが不正: %s", lines[1])
	}
	if !strings.Contains(lines[2], "def456") || !strings.Contains(lines[2], "記事2") {
		t.Errorf("2行目のデータが不正: %s", lines[2])
	}
}

// TestFilterFields_フィールド未指定 はフィールド未指定の場合に全フィールドを出力することを確認します
func TestFilterFields_フィールド未指定(t *testing.T) {
	// テストデータ
	data := map[string]interface{}{
		"id":    "abc123",
		"title": "テスト記事",
	}

	// 出力先バッファ
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// フィールド未指定（nilまたは空スライス）で出力
	err := FilterFields(formatter, data, nil)
	if err != nil {
		t.Fatalf("FilterFieldsに失敗: %v", err)
	}

	// 結果を検証
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSONパースに失敗: %v", err)
	}

	// 全フィールドが含まれることを確認
	if result["id"] != "abc123" {
		t.Errorf("idが含まれていません")
	}
	if result["title"] != "テスト記事" {
		t.Errorf("titleが含まれていません")
	}
}
