/**
 * root.go
 * gho CLIのルート定義
 *
 * Kongを使用したCLI構造の定義
 */

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/mtane0412/ghocli/internal/outfmt"
	"github.com/mtane0412/ghocli/internal/ui"
)

// RootFlags はすべてのコマンドで共通のフラグです
type RootFlags struct {
	Site    string `help:"Site alias or URL" short:"s" env:"GHO_SITE"`
	JSON    bool   `help:"Output JSON" env:"GHO_JSON"`
	Plain   bool   `help:"Output stable TSV" env:"GHO_PLAIN"`
	Fields  string `help:"Fields to output (comma-separated)" short:"F" env:"GHO_FIELDS"`
	Force   bool   `help:"Skip confirmations" short:"f"`
	NoInput bool   `help:"Never prompt; fail instead (useful for CI)" env:"GHO_NO_INPUT"`
	Verbose bool   `help:"Enable verbose logging" short:"v" env:"GHO_VERBOSE"`
	Color   string `help:"Color output (auto, always, never)" enum:"auto,always,never" default:"auto" env:"GHO_COLOR"`
}

// CLI はgho CLIのルート構造体です
type CLI struct {
	RootFlags `embed:""`
	Version   kong.VersionFlag `help:"Print version"`

	Auth        AuthCmd        `cmd:"" help:"Authentication management"`
	Config      ConfigCmd      `cmd:"" help:"Configuration management"`
	Site        SiteCmd        `cmd:"" help:"Site information"`
	Posts       PostsCmd       `cmd:"" aliases:"post,p" help:"Posts management"`
	Pages       PagesCmd       `cmd:"" aliases:"page" help:"Pages management"`
	Tags        TagsCmd        `cmd:"" aliases:"tag,t" help:"Tags management"`
	Images      ImagesCmd      `cmd:"" aliases:"image,img" help:"Images management"`
	Members     MembersCmd     `cmd:"" aliases:"member,m" help:"Members management"`
	Users       UsersCmd       `cmd:"" aliases:"user,u" help:"Users management"`
	Newsletters NewslettersCmd `cmd:"" aliases:"newsletter,nl" help:"Newsletters management"`
	Tiers       TiersCmd       `cmd:"" aliases:"tier" help:"Tiers management"`
	Offers      OffersCmd      `cmd:"" aliases:"offer" help:"Offers management"`
	Themes      ThemesCmd      `cmd:"" aliases:"theme" help:"Themes management"`
	Webhooks    WebhooksCmd    `cmd:"" aliases:"webhook,wh" help:"Webhooks management"`
	Settings    SettingsCmd    `cmd:"" aliases:"setting" help:"Settings management"`

	Completion         CompletionCmd         `cmd:"" help:"Generate shell completion script"`
	CompletionInternal CompletionInternalCmd `cmd:"" name:"__complete" hidden:"" help:""`
}

// GetOutputMode はRootFlagsから出力モードを決定します
func (r *RootFlags) GetOutputMode() string {
	if r.JSON {
		return "json"
	}
	if r.Plain {
		return "plain"
	}
	return "table"
}

// ExecuteOptions はExecute関数のオプションです
type ExecuteOptions struct {
	// Version はバージョン文字列（省略時は"dev"）
	Version string
}

// Execute はgho CLIのエントリーポイントです
// argsにはコマンドライン引数を渡します（os.Argsなど）
func Execute(args []string, opts ...ExecuteOptions) (err error) {
	// オプションの取得
	version := "dev"
	if len(opts) > 0 && opts[0].Version != "" {
		version = opts[0].Version
	}

	// CLIを初期化
	cli := &CLI{}

	// Kongパーサーを作成
	parser, err := kong.New(cli,
		kong.Name("gho"),
		kong.Description("Ghost Admin API CLI"),
		kong.UsageOnError(),
		kong.Writers(os.Stdout, os.Stderr),
		kong.Help(helpPrinter),
		kong.Vars{
			"version": version,
		},
	)
	if err != nil {
		return err
	}

	// コマンドライン引数をパース（最初の要素はプログラム名なのでスキップ）
	var parseArgs []string
	if len(args) > 0 {
		parseArgs = args[1:]
	}

	// Kongでパース
	kctx, err := parser.Parse(parseArgs)
	if err != nil {
		// パースエラーをstderrに出力
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// contextを初期化
	ctx := context.Background()

	// 出力モードを設定
	mode := outfmt.Mode{
		JSON:  cli.JSON,
		Plain: cli.Plain,
	}
	ctx = outfmt.WithMode(ctx, mode)

	// UI出力を設定
	uiOutput := ui.NewOutput(os.Stdout, os.Stderr)
	ctx = ui.WithUI(ctx, uiOutput)

	// contextをKongにバインド
	kctx.BindTo(ctx, (*context.Context)(nil))

	// RootFlagsをバインド
	kctx.Bind(&cli.RootFlags)

	// コマンドを実行
	return kctx.Run()
}
