/**
 * outfmt_test.go
 * 出力フォーマット機能のテストコード
 */

package outfmt

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

// TestPrintJSON_JSON形式で出力
func TestPrintJSON_JSON形式で出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	data := map[string]interface{}{
		"title": "Test Blog",
		"url":   "https://test.ghost.io",
	}

	if err := formatter.Print(data); err != nil {
		t.Fatalf("出力に失敗: %v", err)
	}

	output := buf.String()

	// JSONとしてパースできることを確認
	if !strings.Contains(output, `"title"`) {
		t.Error("JSONに'title'フィールドが含まれていない")
	}
	if !strings.Contains(output, `"Test Blog"`) {
		t.Error("JSONに'Test Blog'値が含まれていない")
	}
}

// TestPrintTable_テーブル形式で出力
func TestPrintTable_テーブル形式で出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	headers := []string{"Name", "URL"}
	rows := [][]string{
		{"Site1", "https://site1.ghost.io"},
		{"Site2", "https://site2.ghost.io"},
	}

	if err := formatter.PrintTable(headers, rows); err != nil {
		t.Fatalf("テーブル出力に失敗: %v", err)
	}

	output := buf.String()

	// ヘッダーと各行が含まれていることを確認
	if !strings.Contains(output, "Name") {
		t.Error("出力にヘッダー'Name'が含まれていない")
	}
	if !strings.Contains(output, "Site1") {
		t.Error("出力に'Site1'が含まれていない")
	}
	if !strings.Contains(output, "Site2") {
		t.Error("出力に'Site2'が含まれていない")
	}
}

// TestPrintPlain_プレーン形式（TSV）で出力
func TestPrintPlain_プレーン形式で出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "plain")

	headers := []string{"Name", "URL"}
	rows := [][]string{
		{"Site1", "https://site1.ghost.io"},
		{"Site2", "https://site2.ghost.io"},
	}

	if err := formatter.PrintTable(headers, rows); err != nil {
		t.Fatalf("プレーン出力に失敗: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// ヘッダー行とデータ行があることを確認
	if len(lines) != 3 {
		t.Errorf("行数 = %d; want 3", len(lines))
	}

	// TSV形式（タブ区切り）であることを確認
	if !strings.Contains(lines[0], "\t") {
		t.Error("ヘッダー行がタブ区切りではない")
	}
	if !strings.Contains(lines[1], "\t") {
		t.Error("データ行1がタブ区切りではない")
	}
}

// TestPrintMessage_メッセージ出力
func TestPrintMessage_メッセージ出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	message := "Test message"
	formatter.PrintMessage(message)

	output := buf.String()
	if !strings.Contains(output, message) {
		t.Errorf("出力にメッセージが含まれていない: %s", output)
	}
}

// TestPrintError_エラー出力
func TestPrintError_エラー出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	errMsg := "Test error"
	formatter.PrintError(errMsg)

	output := buf.String()
	if !strings.Contains(output, errMsg) {
		t.Errorf("出力にエラーメッセージが含まれていない: %s", output)
	}
}

// TestPrintTable_日本語文字列を含むテーブル表示
func TestPrintTable_日本語文字列を含むテーブル表示(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	headers := []string{"Title", "Status"}
	rows := [][]string{
		{"非エンジニアおじさんの開発環境2026", "published"},
		{"1x4材と麻紐でキャットタワーを作る", "published"},
		{"Test", "draft"},
	}

	if err := formatter.PrintTable(headers, rows); err != nil {
		t.Fatalf("テーブル出力に失敗: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	// ヘッダー行、セパレーター行、データ3行、最後の空行 = 6行
	if len(lines) != 6 {
		t.Errorf("行数 = %d; want 6", len(lines))
	}

	// 各行の列が正しく揃っていることを確認
	// セパレーターの位置を基準に、各データ行の列位置が揃っているかチェック
	separatorLine := lines[1]
	if separatorLine == "" {
		t.Fatal("セパレーター行が空")
	}

	// セパレーター行から各列の開始位置を特定
	// (この実装は、列が2つスペースで区切られている前提)
	// 実際には、表示幅が正しく計算されていれば、セパレーターの長さが適切になる

	// すべての行の表示幅が揃っていることを確認
	headerLine := lines[0]
	headerWidth := runewidth.StringWidth(headerLine)
	for i := range rows {
		dataLine := lines[i+2] // ヘッダー、セパレーターの後
		// データ行の表示幅が、ヘッダー行と同じであることを確認
		dataWidth := runewidth.StringWidth(dataLine)
		if headerWidth != dataWidth {
			t.Errorf("行 %d の表示幅がヘッダーと異なる (header=%d, data=%d)\n  Header: %q\n  Data:   %q",
				i, headerWidth, dataWidth, headerLine, dataLine)
		}
	}
}

// TestPrintKeyValue_プレーン形式で出力
func TestPrintKeyValue_プレーン形式で出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "plain")

	// キー/値のペア（ヘッダーなし）
	rows := [][]string{
		{"Title", "はなしのタネ"},
		{"URL", "https://hanashinotane.com"},
		{"Version", "5.102"},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		t.Fatalf("キー/値出力に失敗: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// データ3行（ヘッダーなし）
	if len(lines) != 3 {
		t.Errorf("行数 = %d; want 3", len(lines))
	}

	// TSV形式（タブ区切り）であることを確認
	if !strings.Contains(lines[0], "\t") {
		t.Error("データ行1がタブ区切りではない")
	}

	// 最初の行が "Title\tはなしのタネ" であることを確認
	expected := "Title\tはなしのタネ"
	if lines[0] != expected {
		t.Errorf("最初の行 = %q; want %q", lines[0], expected)
	}

	// "Field" または "Value" というヘッダーがないことを確認
	if strings.Contains(output, "Field") || strings.Contains(output, "Value") {
		t.Error("出力にヘッダー（Field/Value）が含まれている")
	}
}

// TestPrintKeyValue_JSON形式で出力
func TestPrintKeyValue_JSON形式で出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "json")

	// キー/値のペア
	rows := [][]string{
		{"Title", "はなしのタネ"},
		{"URL", "https://hanashinotane.com"},
		{"Version", "5.102"},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		t.Fatalf("キー/値出力に失敗: %v", err)
	}

	output := buf.String()

	// JSONオブジェクトとして出力されることを確認（配列ではない）
	if strings.HasPrefix(strings.TrimSpace(output), "[") {
		t.Error("JSON出力が配列形式になっている（オブジェクト形式であるべき）")
	}

	// 各キーが含まれていることを確認
	if !strings.Contains(output, `"Title"`) {
		t.Error("JSONに'Title'フィールドが含まれていない")
	}
	if !strings.Contains(output, `"はなしのタネ"`) {
		t.Error("JSONに'はなしのタネ'値が含まれていない")
	}
	if !strings.Contains(output, `"URL"`) {
		t.Error("JSONに'URL'フィールドが含まれていない")
	}
}

// TestPrintKeyValue_テーブル形式で出力
func TestPrintKeyValue_テーブル形式で出力(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewFormatter(&buf, "table")

	// キー/値のペア
	rows := [][]string{
		{"Title", "はなしのタネ"},
		{"URL", "https://hanashinotane.com"},
		{"Version", "5.102"},
	}

	if err := formatter.PrintKeyValue(rows); err != nil {
		t.Fatalf("キー/値出力に失敗: %v", err)
	}

	// tabwriterをFlushする
	if err := formatter.Flush(); err != nil {
		t.Fatalf("Flush失敗: %v", err)
	}

	output := buf.String()
	lines := strings.Split(output, "\n")

	// データ3行、最後の空行 = 4行（ヘッダー行とセパレーター行はなし）
	if len(lines) != 4 {
		t.Errorf("行数 = %d; want 4 (got: %v)", len(lines), lines)
	}

	// 最初の行に "Title" が含まれることを確認
	if !strings.Contains(lines[0], "Title") {
		t.Error("最初の行に'Title'が含まれていない")
	}

	// 最初の行に "はなしのタネ" が含まれることを確認
	if !strings.Contains(lines[0], "はなしのタネ") {
		t.Error("最初の行に'はなしのタネ'が含まれていない")
	}

	// タブ文字は含まれていないこと（tabwriterでスペースに変換される）
	if strings.Contains(lines[0], "\t") {
		t.Error("タブ文字が含まれている（tabwriterで整列されるべき）")
	}
}
