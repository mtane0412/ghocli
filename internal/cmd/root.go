/**
 * root.go
 * gho CLIのルート定義
 *
 * Kongを使用したCLI構造の定義
 */

package cmd

import (
	"github.com/alecthomas/kong"
)

// RootFlags はすべてのコマンドで共通のフラグです
type RootFlags struct {
	Site    string `help:"Site alias or URL" short:"s" env:"GHO_SITE"`
	JSON    bool   `help:"Output JSON" env:"GHO_JSON"`
	Plain   bool   `help:"Output stable TSV"`
	Force   bool   `help:"Skip confirmations" short:"f"`
	Verbose bool   `help:"Enable verbose logging" short:"v"`
}

// CLI はgho CLIのルート構造体です
type CLI struct {
	RootFlags `embed:""`
	Version   kong.VersionFlag `help:"Print version"`

	Auth    AuthCmd    `cmd:"" help:"Authentication management"`
	Site    SiteCmd    `cmd:"" help:"Site information"`
	Posts   PostsCmd   `cmd:"" help:"Posts management"`
	Pages   PagesCmd   `cmd:"" help:"Pages management"`
	Tags    TagsCmd    `cmd:"" help:"Tags management"`
	Images  ImagesCmd  `cmd:"" help:"Images management"`
	Members MembersCmd `cmd:"" help:"Members management"`
	Users   UsersCmd   `cmd:"" help:"Users management"`
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
