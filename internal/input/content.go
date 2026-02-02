/**
 * content.go
 * コンテンツ入力ユーティリティ
 *
 * ファイル、標準入力、インラインコンテンツからコンテンツを読み込む機能を提供する
 */

package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ContentFormat はコンテンツのフォーマットを表す型
type ContentFormat string

const (
	// FormatUnknown は不明なフォーマット
	FormatUnknown ContentFormat = ""
	// FormatHTML はHTML形式
	FormatHTML ContentFormat = "html"
	// FormatMarkdown はMarkdown形式
	FormatMarkdown ContentFormat = "markdown"
	// FormatLexical はLexical JSON形式
	FormatLexical ContentFormat = "lexical"
)

// ReadContent reads content from a file or returns inline content
//
// 優先順位:
// 1. filePathが指定されている場合、ファイルからコンテンツを読み込む
// 2. filePathが空の場合、inlineContentを返す
func ReadContent(filePath string, inlineContent string) (string, error) {
	// ファイルパスが指定されている場合はファイルから読み込み
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("ファイルの読み込みに失敗: %w", err)
		}
		return string(data), nil
	}

	// ファイルパスが空の場合はインラインコンテンツを返す
	return inlineContent, nil
}

// DetectFormat はファイルパスの拡張子からコンテンツフォーマットを検出する
//
// 引数:
//   - filePath: ファイルパス
//
// 戻り値:
//   - ContentFormat: 検出されたフォーマット
//
// 検出ルール:
//   - .md, .markdown → FormatMarkdown
//   - .html, .htm → FormatHTML
//   - .json → FormatLexical
//   - その他 → FormatUnknown
func DetectFormat(filePath string) ContentFormat {
	if filePath == "" {
		return FormatUnknown
	}

	// 拡張子を取得（小文字に変換）
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".md", ".markdown":
		return FormatMarkdown
	case ".html", ".htm":
		return FormatHTML
	case ".json":
		return FormatLexical
	default:
		return FormatUnknown
	}
}

// ReadContentWithFormat はファイルまたはインラインコンテンツを読み込み、フォーマットも返す
//
// 引数:
//   - filePath: ファイルパス（空の場合はinlineContentを使用）
//   - inlineContent: インラインコンテンツ
//
// 戻り値:
//   - content: コンテンツ文字列
//   - format: 検出されたフォーマット（インラインの場合はFormatUnknown）
//   - error: エラー
//
// 優先順位:
//  1. filePathが指定されている場合、ファイルから読み込み、フォーマットを検出
//  2. filePathが空の場合、inlineContentを返し、フォーマットはFormatUnknown
func ReadContentWithFormat(filePath string, inlineContent string) (content string, format ContentFormat, err error) {
	// ファイルパスが指定されている場合
	if filePath != "" {
		// フォーマットを検出
		format = DetectFormat(filePath)

		// ファイルから読み込み
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", FormatUnknown, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
		}

		return string(data), format, nil
	}

	// インラインコンテンツを返す（フォーマットは不明）
	return inlineContent, FormatUnknown, nil
}
