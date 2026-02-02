# gho — Ghost in your terminal.

Fast, script-friendly CLI for Ghost Admin API. JSON-first output, multiple sites, and secure credential storage built in.

> **Inspired by [gogcli](https://github.com/steipete/gogcli)** — This project is heavily influenced by gogcli's design philosophy and user experience. The motivation was simple: "I wanted a Ghost version of gogcli."

## Features

**Content Management**
- **Posts** — list, search, create, update, delete, publish, unpublish, schedule, drafts, copy
- **Pages** — full CRUD operations, content export
- **Tags** — manage tags with visibility control (public/internal)

**User & Member Management**
- **Members** — manage subscribers with filters, labels, and notes
- **Users** — view and update staff users with role information
- **Newsletters** — create and manage newsletters with sender configuration

**Monetization**
- **Tiers** — manage membership tiers (free/paid) with pricing
- **Offers** — create discount codes and promotions (percentage/fixed amount)

**Site Management**
- **Images** — upload images with purpose specification (profile_image, icon, etc.)
- **Themes** — list, upload, and activate themes
- **Webhooks** — create, update, and delete webhooks for events
- **Settings** — view site settings and configuration

**Developer Experience**
- **Multiple sites** — manage multiple Ghost sites with aliases
- **Secure storage** — OS keyring integration (Keychain, Secret Service, Credential Manager)
- **Output formats** — JSON, Table, TSV (Plain) for scripting
- **Batch operations** — combine with `jq` for powerful automation

## Installation

### Homebrew (Recommended)

```bash
brew tap mtane0412/ghocli
brew install gho
```

### Build from Source

```bash
git clone https://github.com/mtane0412/ghocli.git
cd ghocli
make build
```

The binary will be created as `./gho` in the project directory.

### Go Install

```bash
go install github.com/mtane0412/ghocli/cmd/gho@latest
```

Make sure `$GOPATH/bin` is in your `PATH`.

## Quick Start

### 1. Get Admin API Key

Navigate to your Ghost Admin panel:
1. Go to **Settings** → **Integrations**
2. Click **Add custom integration**
3. Name it (e.g., "CLI Access")
4. Copy the **Admin API Key**

### 2. Add Authentication

```bash
gho auth add https://your-blog.ghost.io
# Paste your Admin API key when prompted
```

### 3. Run Your First Command

```bash
gho posts list
```

That's it! You're ready to manage your Ghost site from the terminal.

## Authentication & Secrets

### Keyring Backends

gho stores credentials securely using your operating system's keyring:

- **macOS**: Keychain
- **Linux**: Secret Service (GNOME Keyring, KWallet)
- **Windows**: Credential Manager

### Environment Variables

Configure keyring behavior with environment variables:

```bash
# Force a specific backend
export GHO_KEYRING_BACKEND=keychain  # or: file, secretservice, wincred, auto

# Set password for file-based backend
export GHO_KEYRING_PASSWORD="your-secure-password"
```

### Multiple Sites

Manage multiple Ghost sites with aliases:

```bash
# Add sites with custom aliases
gho auth add https://blog.example.com --alias myblog
gho auth add https://news.example.com --alias news

# List all registered sites
gho auth list

# Use specific site
gho -s myblog posts list
gho -s news posts list

# Check authentication status
gho auth status
gho auth status -s myblog
```

## Configuration

### Config File

Configuration is stored at `~/.config/gho/config.json`.

### Available Settings

```bash
# View all configuration keys
gho config keys

# Get a specific value
gho config get default_site

# Set a value
gho config set default_site myblog

# Unset a value
gho config unset default_site

# List all current settings
gho config list

# Show config file path
gho config path
```

### Site Selection Priority

When running commands, gho selects the site in this order:
1. `--site` flag
2. `GHO_SITE` environment variable
3. `default_site` config value

## Commands

### Authentication

```bash
gho auth add <url>              # Add new site
gho auth list                   # List registered sites
gho auth remove <alias>         # Remove site
gho auth status                 # Check authentication status
gho auth tokens                 # List stored tokens
gho auth credentials            # Show credentials (admin-only)
```

### Configuration

```bash
gho config get <key>            # Get config value
gho config set <key> <value>    # Set config value
gho config unset <key>          # Unset config value
gho config list                 # List all settings
gho config path                 # Show config file path
gho config keys                 # Show available keys
```

### Site Information

```bash
gho site                        # Get site information
gho site --json                 # Output as JSON
```

### Posts

```bash
# List & Search
gho posts list                  # List all posts
gho posts list --status draft   # Filter by status (draft/published/scheduled)
gho posts list --limit 10       # Limit results
gho posts search <query>        # Search posts by keyword
gho posts drafts                # List draft posts only
gho posts url <url>             # Get post by URL

# View
gho posts get <id-or-slug>      # Get post details
gho posts cat <id-or-slug>      # Display post content
gho posts cat <id> --format text    # Display as plain text
gho posts cat <id> --format lexical # Display as Lexical JSON

# Create & Update
gho posts create --title "Title" --html "Content"
gho posts create --title "Title" --markdown "# Heading\n\nContent"
gho posts create --title "Title" --lexical '{"root":{"children":[...]}}'
gho posts create --title "Title" --file article.md   # Auto-detect format (.md, .html, .json)
gho posts create --title "New Post" --status draft
gho posts update <id> --title "New Title"
gho posts update <id> --html "New Content"
gho posts update <id> --markdown "# Updated Content"
gho posts update <id> --file updated.md
gho posts copy <id-or-slug>     # Copy post as new draft
gho posts copy <id> --title "Copy of Original"

# Publishing
gho posts publish <id>          # Publish immediately
gho posts unpublish <id>        # Unpublish to draft
gho posts schedule <id> "2026-12-31T23:59:59Z"  # Schedule publication

# Delete
gho posts delete <id>           # Delete post
gho posts delete <id> --force   # Skip confirmation

# Batch Operations
gho posts batch publish <id1> <id2> ...
gho posts batch delete <id1> <id2> ...
```

### Pages

```bash
gho pages list                  # List all pages
gho pages list --status draft   # Filter by status
gho pages info <id-or-slug>     # Get page details
gho pages cat <id-or-slug>      # Display page content
gho pages create --title "Title" --html "Content"
gho pages create --title "Title" --markdown "# Page Content"
gho pages create --title "Title" --lexical '{"root":{"children":[...]}}'
gho pages create --title "Title" --file page.md      # Auto-detect format (.md, .html, .json)
gho pages update <id> --title "New Title"
gho pages update <id> --markdown "# Updated Content"
gho pages update <id> --file updated.md
gho pages delete <id>           # Delete page
gho pages copy <id-or-slug>     # Copy page as new draft
```

### Tags

```bash
gho tags list                   # List all tags
gho tags info <id-or-slug>      # Get tag details
gho tags info slug:technology   # Get tag by slug
gho tags create --name "Tech" --description "Technical posts"
gho tags create --name "Internal" --visibility internal
gho tags update <id> --name "Technology"
gho tags delete <id>            # Delete tag
```

### Members

```bash
gho members list                # List all members
gho members list --limit 10     # Limit results
gho members list --filter "status:paid"  # Apply filter
gho members info <id>           # Get member details
gho members create --email "user@example.com" --name "John Doe"
gho members create --email "user@example.com" --labels "VIP,Premium"
gho members update <id> --name "Jane Doe" --note "Important customer"
gho members delete <id>         # Delete member
```

### Users

```bash
gho users list                  # List all users
gho users list --include roles  # Include role information
gho users list --include count.posts  # Include post count
gho users info <id-or-slug>     # Get user details
gho users info slug:john-doe    # Get user by slug
gho users update <id> --name "New Name" --bio "New bio"
gho users update <id> --location "Tokyo" --website "https://example.com"
```

### Newsletters

```bash
gho newsletters list            # List all newsletters
gho newsletters list --filter "status:active"
gho newsletters info <id-or-slug>
gho newsletters info slug:weekly-newsletter
gho newsletters create --name "Weekly" --description "Delivered every Friday"
gho newsletters create --name "Monthly" --sender-name "Editorial" --sender-email "editor@example.com"
gho newsletters update <id> --name "New Name"
gho newsletters update <id> --visibility paid --subscribe-on-signup=false
```

### Tiers

```bash
gho tiers list                  # List all tiers
gho tiers list --include monthly_price,yearly_price
gho tiers info <id-or-slug>
gho tiers info slug:premium
gho tiers create --name "Free Plan" --type free
gho tiers create --name "Premium" --type paid --monthly-price 1000 --yearly-price 10000 --currency JPY
gho tiers create --name "VIP" --type paid --monthly-price 3000 --benefits "Priority Support" --benefits "Exclusive Content"
gho tiers update <id> --name "New Premium"
gho tiers update <id> --monthly-price 1200 --yearly-price 12000
```

### Offers

```bash
gho offers list                 # List all offers
gho offers list --filter "status:active"
gho offers info <id>
gho offers create --name "Welcome" --code "WELCOME2024" --type percent --amount 20 --tier-id <tier-id>
gho offers create --name "500 Off" --code "SAVE500" --type fixed --amount 500 --currency JPY --tier-id <tier-id>
gho offers create --name "3 Month" --code "TRIAL3M" --type percent --amount 50 --duration repeating --duration-in-months 3 --tier-id <tier-id>
gho offers update <id> --name "New Campaign"
gho offers update <id> --amount 30
```

### Images

```bash
gho images upload path/to/image.jpg
gho images upload avatar.png --purpose profile_image
gho images upload icon.png --purpose icon
gho images upload banner.jpg --ref post-123
```

### Themes

```bash
gho themes list                 # List installed themes
gho themes upload theme.zip     # Upload and install theme
gho themes activate casper      # Activate theme by name
```

### Webhooks

```bash
gho webhooks create --event post.published --target-url https://example.com/webhook
gho webhooks create --event member.added --target-url https://example.com/webhook --name "Member notification"
gho webhooks update <id> --target-url https://new-example.com/webhook
gho webhooks delete <id>
```

### Settings

```bash
gho settings list               # List all settings
gho settings info <key>         # Get specific setting
```

## Output Formats

gho supports three output formats optimized for different use cases.

### Table (Default)

Human-readable format with aligned columns.

**Single item:**
```bash
$ gho site
title        My Blog
description  A technical blog about programming
url          https://blog.example.com/
version      6.8
```

**Multiple items:**
```bash
$ gho posts list --limit 3
ID                        TITLE                     STATUS     CREATED     PUBLISHED
697b61d44921c40001f01aa3  Cannot/Don't Use CLI      draft      2026-01-29
696ce7244921c40001f017ed  Dev Environment 2026      published  2026-01-18  2026-01-28
```

### Plain (TSV)

Tab-separated values for scripting and pipelines.

```bash
$ gho posts list --plain --limit 2
ID	TITLE	STATUS	CREATED	PUBLISHED
697b61d44921c40001f01aa3	Cannot/Don't Use CLI	draft	2026-01-29
696ce7244921c40001f017ed	Dev Environment 2026	published	2026-01-18	2026-01-28

$ gho posts list --plain | cut -f2  # Extract titles only
```

### JSON

Structured data for programmatic processing.

```bash
$ gho site --json
{
  "site": {
    "title": "My Blog",
    "description": "A technical blog about programming",
    "url": "https://blog.example.com/",
    "version": "6.8"
  }
}

$ gho posts list --json --limit 1 | jq '.posts[0].title'
"Cannot/Don't Use CLI"
```

## Examples

### Search and List Posts

```bash
# Find posts about "docker"
gho posts search docker --json | jq '.posts[] | {title, slug, status}'

# List all draft posts
gho posts drafts --json | jq '.posts[] | .title'

# Get published posts from last 30 days
gho posts list --status published --json | jq '.posts[] | select(.published_at > (now - 2592000))'
```

### Batch Operations with jq

```bash
# Get all post IDs
gho posts list --json | jq -r '.posts[].id'

# Export all published post titles and URLs
gho posts list --status published --json | \
  jq -r '.posts[] | [.title, .url] | @tsv' > posts.tsv

# Bulk update: add tag to all draft posts
gho posts drafts --json | \
  jq -r '.posts[].id' | \
  xargs -I {} gho posts update {} --tags "draft-review"

# Create posts from JSON file
cat posts.json | \
  jq -c '.[]' | \
  while read post; do
    gho posts create \
      --title "$(echo $post | jq -r .title)" \
      --html "$(echo $post | jq -r .content)" \
      --status draft
  done
```

### Multi-Site Management

```bash
# Add multiple sites
gho auth add https://blog.example.com --alias blog
gho auth add https://news.example.com --alias news

# Set default site
gho config set default_site blog

# Run commands on different sites
gho posts list                    # Uses default (blog)
gho -s news posts list            # Uses news site
GHO_SITE=news gho members list    # Environment variable

# Compare post counts across sites
echo "Blog: $(gho -s blog posts list --json | jq '.posts | length')"
echo "News: $(gho -s news posts list --json | jq '.posts | length')"
```

### Automated Publishing Workflow

```bash
#!/bin/bash
# Publish all scheduled posts that are ready

# Get current timestamp
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Find and publish ready posts
gho posts list --status scheduled --json | \
  jq -r --arg now "$NOW" '.posts[] | select(.published_at <= $now) | .id' | \
  xargs -I {} gho posts publish {} --force

echo "✅ Published all ready posts"
```

## Global Flags

These flags work with all commands:

| Flag | Short | Environment Variable | Description |
|------|-------|---------------------|-------------|
| `--site <alias>` | `-s` | `GHO_SITE` | Specify site alias or URL |
| `--json` | | `GHO_JSON=1` | Output as JSON |
| `--plain` | | `GHO_PLAIN=1` | Output as TSV (tab-separated) |
| `--fields <list>` | `-F` | `GHO_FIELDS` | Select specific fields to display |
| `--force` | `-f` | | Skip confirmation prompts |
| `--no-input` | | `GHO_NO_INPUT=1` | Non-interactive mode (fail if input required) |
| `--verbose` | `-v` | `GHO_VERBOSE=1` | Enable verbose logging |
| `--color <mode>` | | `GHO_COLOR` | Color output (auto/always/never) |

### Examples

```bash
# Combine flags
gho -s myblog posts list --json --fields id,title,status

# Use environment variables
export GHO_SITE=myblog
export GHO_JSON=1
gho posts list

# Force mode for automation
gho posts delete abc123 --force
gho newsletters create --name "Test" --force
```

## Shell Completions

Generate shell completion scripts for faster typing.

### Bash

```bash
# Generate completion script
gho completion bash > /etc/bash_completion.d/gho

# Or for user-level installation
gho completion bash > ~/.local/share/bash-completion/completions/gho
```

### Zsh

```bash
# Generate completion script
gho completion zsh > "${fpath[1]}/_gho"

# Then reload your shell
source ~/.zshrc
```

### Fish

```bash
gho completion fish > ~/.config/fish/completions/gho.fish
```

### PowerShell

```powershell
gho completion powershell | Out-String | Invoke-Expression
```

## Development

### Prerequisites

- Go 1.23 or later
- Make

### Build

```bash
make build
```

The binary will be created as `./gho`.

### Run Tests

```bash
make test
```

### Run Linter

```bash
make lint
```

Requires [golangci-lint](https://golangci-lint.run/) to be installed.

### Type Check

```bash
make type-check
```

### Test Coverage

```bash
make test-coverage
```

Opens coverage report in your browser at `coverage.html`.

### All Quality Checks

```bash
make lint
make type-check
make test
```

## Architecture

### Design Philosophy

gho is **heavily inspired by [gogcli](https://github.com/steipete/gogcli)**, a CLI tool for Google services. The entire architecture, coding patterns, and user experience design follow gogcli's proven approach.

**Why gogcli?**
- Excellent UX with JSON/Table/Plain output modes
- Clean architecture with context propagation
- Robust error handling and exit code management
- Script-friendly design that works well with Unix pipes

gho applies these same principles to the Ghost Admin API, aiming to provide the same level of polish and usability that gogcli brings to Google services.

### Core Design Principles

- **Context propagation**: All commands use Go context for clean dependency injection
- **Exit code management**: Proper exit codes through `ExitError` type
- **Output format abstraction**: Unified interface for JSON/Table/Plain formats
- **TDD**: Comprehensive test coverage with Test-Driven Development
- **Type safety**: Leveraging Go's type system for correctness
- **Script-friendly**: Designed to work seamlessly in automation and pipelines

### Project Structure

- `internal/cmd/` — Command implementations
- `internal/outfmt/` — Output format handlers (JSON/Table/Plain)
- `internal/ui/` — UI output (stdout/stderr separation)
- `internal/ghostapi/` — Ghost Admin API client
- `internal/secrets/` — OS keyring integration
- `internal/config/` — Configuration management

For more details:
- [Architecture Documentation](./docs/ARCHITECTURE.md)
- [Development Guide](./docs/DEVELOPMENT_GUIDE.md)

## License

MIT

## Links

- [GitHub Repository](https://github.com/mtane0412/gho)
- [Ghost Admin API Documentation](https://ghost.org/docs/admin-api/)
- [Issue Tracker](https://github.com/mtane0412/gho/issues)
- [gogcli](https://github.com/steipete/gogcli) — The project that inspired gho's design
