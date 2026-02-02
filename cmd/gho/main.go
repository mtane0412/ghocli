/**
 * main.go
 * gho - Ghost Admin API CLI
 *
 * gog-cliの使用感を備えたGhost Admin APIのCLIツール
 */

package main

import (
	"fmt"
	"os"

	"github.com/mtane0412/ghocli/internal/cmd"
)

var (
	// バージョン情報（ビルド時に-ldflagsで設定される）
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Execute関数を呼び出してコマンドを実行
	err := cmd.Execute(os.Args, cmd.ExecuteOptions{
		Version: buildVersion(),
	})

	// エラーがあればstderrに出力して終了コードを設定して終了
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(cmd.ExitCode(err))
	}
}

// buildVersion はバージョン文字列を構築します
func buildVersion() string {
	if version == "dev" {
		return "gho dev (commit: " + commit + ", built at: " + date + ")"
	}
	return "gho " + version
}
