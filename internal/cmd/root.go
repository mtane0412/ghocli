/**
 * root.go
 * Root definition for gho CLI
 *
 * Defines CLI structure using Kong
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

// RootFlags are common flags for all commands
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

// CLI is the root structure for gho CLI
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

// GetOutputMode determines the output mode from RootFlags
func (r *RootFlags) GetOutputMode() string {
	if r.JSON {
		return "json"
	}
	if r.Plain {
		return "plain"
	}
	return "table"
}

// ExecuteOptions are options for the Execute function
type ExecuteOptions struct {
	// Version is the version string (defaults to "dev" if omitted)
	Version string
}

// Execute is the entry point for gho CLI
// args should contain command line arguments (such as os.Args)
func Execute(args []string, opts ...ExecuteOptions) (err error) {
	// Get options
	version := "dev"
	if len(opts) > 0 && opts[0].Version != "" {
		version = opts[0].Version
	}

	// Initialize CLI
	cli := &CLI{}

	// Create Kong parser
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

	// Parse command line arguments (skip first element as it's the program name)
	var parseArgs []string
	if len(args) > 0 {
		parseArgs = args[1:]
	}

	// Parse with Kong
	kctx, err := parser.Parse(parseArgs)
	if err != nil {
		// Output parse error to stderr
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// Initialize context
	ctx := context.Background()

	// Set output mode
	mode := outfmt.Mode{
		JSON:  cli.JSON,
		Plain: cli.Plain,
	}
	ctx = outfmt.WithMode(ctx, mode)

	// Set UI output
	uiOutput := ui.NewOutput(os.Stdout, os.Stderr)
	ctx = ui.WithUI(ctx, uiOutput)

	// Bind context to Kong
	kctx.BindTo(ctx, (*context.Context)(nil))

	// Bind RootFlags
	kctx.Bind(&cli.RootFlags)

	// Execute command
	return kctx.Run()
}
