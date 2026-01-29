/**
 * main.go
 * gho - Ghost Admin API CLI
 *
 * gog-cliの使用感を備えたGhost Admin APIのCLIツール
 */

package main

import (
	"github.com/alecthomas/kong"
	"github.com/mtane0412/gho/internal/cmd"
)

var (
	// バージョン情報（ビルド時に-ldflagsで設定される）
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// CLIを初期化
	cli := &cmd.CLI{}

	// Kongでパース
	ctx := kong.Parse(cli,
		kong.Name("gho"),
		kong.Description("Ghost Admin API CLI"),
		kong.UsageOnError(),
		kong.Vars{
			"version": buildVersion(),
		},
	)

	// コマンドを実行
	err := ctx.Run(&cli.RootFlags)
	ctx.FatalIfErrorf(err)
}

// buildVersion はバージョン文字列を構築します
func buildVersion() string {
	if version == "dev" {
		return "gho dev (commit: " + commit + ", built at: " + date + ")"
	}
	return "gho " + version
}
