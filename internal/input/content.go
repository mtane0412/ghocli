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
