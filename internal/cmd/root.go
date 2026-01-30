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
	Plain   bool   `help:"Output stable TSV" env:"GHO_PLAIN"`
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
	Site        SiteCmd        `cmd:"" help:"Site information"`
	Posts       PostsCmd       `cmd:"" help:"Posts management"`
	Pages       PagesCmd       `cmd:"" help:"Pages management"`
	Tags        TagsCmd        `cmd:"" help:"Tags management"`
	Images      ImagesCmd      `cmd:"" help:"Images management"`
	Members     MembersCmd     `cmd:"" help:"Members management"`
	Users       UsersCmd       `cmd:"" help:"Users management"`
	Newsletters NewslettersCmd `cmd:"" help:"Newsletters management"`
	Tiers       TiersCmd       `cmd:"" help:"Tiers management"`
	Offers      OffersCmd      `cmd:"" help:"Offers management"`
	Themes      ThemesCmd      `cmd:"" help:"Themes management"`
	Webhooks    WebhooksCmd    `cmd:"" help:"Webhooks management"`
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
