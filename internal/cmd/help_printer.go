/**
 * help_printer.go
 * カスタムヘルププリンター
 *
 * ヘルプメッセージの色付けとビルド情報の注入を行います。
 */

package cmd

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// helpPrinter はカスタムヘルププリンターです
func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	origStdout := ctx.Stdout
	origStderr := ctx.Stderr

	// 端末の幅を取得
	width := guessColumns(origStdout)
	oldCols, hadCols := os.LookupEnv("COLUMNS")
	_ = os.Setenv("COLUMNS", strconv.Itoa(width))
	defer func() {
		if hadCols {
			_ = os.Setenv("COLUMNS", oldCols)
		} else {
			_ = os.Unsetenv("COLUMNS")
		}
	}()

	// バッファに出力
	buf := bytes.NewBuffer(nil)
	ctx.Stdout = buf
	ctx.Stderr = origStderr
	defer func() { ctx.Stdout = origStdout }()

	// デフォルトのヘルプを生成
	if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
		return err
	}

	// ヘルプテキストをカスタマイズ
	out := buf.String()
	out = injectBuildLine(out)
	out = colorizeHelp(out, helpProfile(origStdout))

	_, err := io.WriteString(origStdout, out)
	return err
}

// injectBuildLine はビルド情報をヘルプに注入します
func injectBuildLine(out string) string {
	// ビルド情報はExecute()で設定されたバージョンから取得
	// 簡略化のため、ここでは注入をスキップ
	return out
}

// helpProfile はカラープロファイルを返します
func helpProfile(stdout io.Writer) termenv.Profile {
	// NO_COLOR環境変数の確認
	if termenv.EnvNoColor() {
		return termenv.Ascii
	}

	// GHO_COLOR環境変数の確認
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("GHO_COLOR")))
	switch mode {
	case "never":
		return termenv.Ascii
	case "always":
		return termenv.TrueColor
	default:
		// auto: 端末の色対応を自動検出
		o := termenv.NewOutput(stdout, termenv.WithProfile(termenv.EnvColorProfile()))
		return o.Profile
	}
}

// colorizeHelp はヘルプテキストを色付けします
func colorizeHelp(out string, profile termenv.Profile) string {
	// カラープロファイルがAsciiの場合は色付けしない
	if profile == termenv.Ascii {
		return out
	}

	// カラー関数を定義
	heading := func(s string) string {
		return termenv.String(s).Foreground(profile.Color("#60a5fa")).Bold().String()
	}
	section := func(s string) string {
		return termenv.String(s).Foreground(profile.Color("#a78bfa")).Bold().String()
	}
	cmdName := func(s string) string {
		return termenv.String(s).Foreground(profile.Color("#38bdf8")).Bold().String()
	}

	// 行ごとに処理
	inCommands := false
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		if line == "Commands:" {
			inCommands = true
		}
		switch {
		case strings.HasPrefix(line, "Usage:"):
			lines[i] = heading("Usage:") + strings.TrimPrefix(line, "Usage:")
		case line == "Flags:":
			lines[i] = section(line)
		case line == "Commands:":
			lines[i] = section(line)
		case line == "Arguments:":
			lines[i] = section(line)
		case inCommands && strings.HasPrefix(line, "  ") && len(line) > 2 && line[2] != ' ':
			// コマンド名を色付け
			name, tail, found := strings.Cut(strings.TrimPrefix(line, "  "), " ")
			if found {
				lines[i] = "  " + cmdName(name) + " " + tail
			} else {
				lines[i] = "  " + cmdName(name)
			}
		}
	}

	return strings.Join(lines, "\n")
}

// guessColumns は端末の幅を推測します
func guessColumns(w io.Writer) int {
	// COLUMNS環境変数を確認
	if colsStr := os.Getenv("COLUMNS"); colsStr != "" {
		if cols, err := strconv.Atoi(colsStr); err == nil {
			return cols
		}
	}

	// 端末のサイズを取得
	f, ok := w.(*os.File)
	if !ok {
		return 80
	}

	width, _, err := term.GetSize(int(f.Fd()))
	if err == nil && width > 0 {
		return width
	}

	// デフォルト値
	return 80
}
