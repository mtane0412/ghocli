/**
 * outfmt.go
 * 出力フォーマット機能
 *
 * JSON、テーブル、プレーン（TSV）形式での出力をサポートします。
 */

package outfmt

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Formatter は出力フォーマッターです
type Formatter struct {
	writer io.Writer
	mode   string // "json", "table", "plain"
}

// NewFormatter は新しい出力フォーマッターを作成します。
func NewFormatter(writer io.Writer, mode string) *Formatter {
	return &Formatter{
		writer: writer,
		mode:   mode,
	}
}

// Print は任意のデータを出力します。
// JSON形式の場合はJSONとして出力します。
func (f *Formatter) Print(data interface{}) error {
	if f.mode == "json" {
		encoder := json.NewEncoder(f.writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	}

	// デフォルトは標準出力
	_, err := fmt.Fprintln(f.writer, data)
	return err
}

// PrintTable はテーブル形式でデータを出力します。
func (f *Formatter) PrintTable(headers []string, rows [][]string) error {
	switch f.mode {
	case "plain":
		// TSV形式で出力
		return f.printTSV(headers, rows)
	case "json":
		// JSON配列として出力
		return f.printJSONTable(headers, rows)
	default:
		// テーブル形式で出力
		return f.printTableFormat(headers, rows)
	}
}

// printTSV はTSV形式（タブ区切り）で出力します。
func (f *Formatter) printTSV(headers []string, rows [][]string) error {
	// ヘッダー行を出力
	if _, err := fmt.Fprintln(f.writer, strings.Join(headers, "\t")); err != nil {
		return err
	}

	// データ行を出力
	for _, row := range rows {
		if _, err := fmt.Fprintln(f.writer, strings.Join(row, "\t")); err != nil {
			return err
		}
	}

	return nil
}

// printJSONTable はJSON配列形式で出力します。
func (f *Formatter) printJSONTable(headers []string, rows [][]string) error {
	// 各行をマップに変換
	var data []map[string]string
	for _, row := range rows {
		item := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				item[header] = row[i]
			}
		}
		data = append(data, item)
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printTableFormat はテーブル形式で出力します（人間向け）。
func (f *Formatter) printTableFormat(headers []string, rows [][]string) error {
	// 各列の最大幅を計算
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// ヘッダー行を出力
	for i, header := range headers {
		if i > 0 {
			fmt.Fprint(f.writer, "  ")
		}
		fmt.Fprintf(f.writer, "%-*s", colWidths[i], header)
	}
	fmt.Fprintln(f.writer)

	// セパレーター行を出力
	for i, width := range colWidths {
		if i > 0 {
			fmt.Fprint(f.writer, "  ")
		}
		fmt.Fprint(f.writer, strings.Repeat("-", width))
	}
	fmt.Fprintln(f.writer)

	// データ行を出力
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(f.writer, "  ")
			}
			if i < len(colWidths) {
				fmt.Fprintf(f.writer, "%-*s", colWidths[i], cell)
			} else {
				fmt.Fprint(f.writer, cell)
			}
		}
		fmt.Fprintln(f.writer)
	}

	return nil
}

// PrintMessage はメッセージを出力します。
func (f *Formatter) PrintMessage(message string) {
	fmt.Fprintln(f.writer, message)
}

// PrintError はエラーメッセージを出力します。
func (f *Formatter) PrintError(message string) {
	fmt.Fprintln(f.writer, "Error:", message)
}
