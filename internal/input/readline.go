/**
 * readline.go
 * 行読み取り機能
 *
 * io.Readerから1行ずつ読み取る機能を提供する。
 * Unix (\n) と Windows (\r\n) の改行形式に対応。
 */
package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ReadLine は、io.Readerから1行を読み取る
//
// Unix (\n) と Windows (\r\n) の改行形式に対応する。
// 単独の\rも改行として扱う。
//
// 改行前にEOFに到達した場合、バッファに内容があればその内容を返す（エラーはnil）。
// バッファに内容がない状態でEOFに到達した場合は、io.EOFを返す。
//
// r: 入力元のio.Reader
//
// 戻り値:
//   - 読み取った行（改行文字を除く）
//   - エラー（読み取りに失敗した場合）
func ReadLine(r io.Reader) (string, error) {
	br := bufio.NewReader(r)

	var sb strings.Builder

	for {
		b, err := br.ReadByte()
		if err != nil {
			// EOFに到達した場合
			if errors.Is(err, io.EOF) {
				// バッファに内容があればそれを返す
				if sb.Len() > 0 {
					return sb.String(), nil
				}

				// バッファが空の場合はEOFを返す
				return "", io.EOF
			}

			// その他のエラー
			return "", fmt.Errorf("read line: %w", err)
		}

		// 改行文字の処理
		if b == '\n' || b == '\r' {
			// \rの場合、次の文字が\nならそれも読み飛ばす（Windows形式）
			if b == '\r' {
				if next, _ := br.Peek(1); len(next) == 1 && next[0] == '\n' {
					_, _ = br.ReadByte()
				}
			}

			// 改行文字を除いた行を返す
			return sb.String(), nil
		}

		// 通常の文字をバッファに追加
		sb.WriteByte(b)
	}
}
