# gho - Ghost Admin API CLI

A CLI tool for Ghost Admin API with the user experience of gog-cli.

## Features

- Execute Ghost Admin API operations from the command line
- Secure API key storage using OS keyrings (macOS Keychain, Linux Secret Service, Windows Credential Manager)
- Multi-site support (alias functionality)
- Output in JSON/Table/TSV formats
- Robust implementation through TDD (Test-Driven Development)

## Implementation Status

✅ **Phase 1-7 Complete** - All major Ghost Admin API features implemented

**Implemented Features**:
- Authentication Management (Auth)
- Site Information (Site)
- Post Management (Posts)
- Page Management (Pages)
- Tag Management (Tags)
- Image Management (Images)
- Member Management (Members)
- User Management (Users)
- Newsletters (Newsletters)
- Tiers (Tiers)
- Offers (Offers)
- Theme Management (Themes)
- Webhook Management (Webhooks)

See [`docs/PROJECT_STATUS.md`](./docs/PROJECT_STATUS.md) for details.

## Installation

```bash
# Build from source
git clone https://github.com/mtane0412/gho.git
cd gho
make build

# Or
go install github.com/mtane0412/gho/cmd/gho@latest
```

## Usage

### Authentication Setup

```bash
# Register Ghost Admin API key
gho auth add https://your-blog.ghost.io

# List registered sites
gho auth list

# Check authentication status
gho auth status

# Remove site authentication
gho auth remove <alias>
```

### Site Information

```bash
# Get site information
gho site

# Get information for a specific site
gho -s myblog site

# Output in JSON format
gho site --json
```

### Posts

```bash
# List posts
gho posts list

# Filter by status
gho posts list --status draft
gho posts list --status published
gho posts list --status scheduled

# Limit number of results
gho posts list --limit 10

# Get post details (by ID or Slug)
gho posts info <id-or-slug>

# Display post content
gho posts cat <id-or-slug>
gho posts cat <id-or-slug> --format text    # Display as text
gho posts cat <id-or-slug> --format lexical # Display as Lexical JSON

# Create new post
gho posts create --title "Title" --html "Content" --status draft

# Update post
gho posts update <id> --title "New Title"
gho posts update <id> --html "New Content"

# Delete post
gho posts delete <id>

# Publish post
gho posts publish <id>

# Copy post (create as new draft)
gho posts copy <id-or-slug>
gho posts copy <id-or-slug> --title "New Title"
```

### Pages

```bash
# List pages
gho pages list

# Filter by status
gho pages list --status draft
gho pages list --status published
gho pages list --status scheduled

# Limit number of results
gho pages list --limit 10

# Get page details (by ID or Slug)
gho pages info <id-or-slug>

# Display page content
gho pages cat <id-or-slug>
gho pages cat <id-or-slug> --format text    # Display as text
gho pages cat <id-or-slug> --format lexical # Display as Lexical JSON

# Create new page
gho pages create --title "Title" --html "Content"

# Update page
gho pages update <id> --title "New Title"
gho pages update <id> --html "New Content"

# Delete page
gho pages delete <id>

# Copy page (create as new draft)
gho pages copy <id-or-slug>
gho pages copy <id-or-slug> --title "New Title"
```

### Tags

```bash
# List tags
gho tags list

# Limit number of results
gho tags list --limit 10

# Get tag details (by ID or Slug)
gho tags info <id-or-slug>
gho tags info slug:technology

# Create new tag
gho tags create --name "Technology" --description "Technical articles"

# Specify tag visibility
gho tags create --name "Internal" --visibility internal

# Update tag
gho tags update <id> --name "Tech" --description "New description"

# Delete tag
gho tags delete <id>
```

### Images

```bash
# Upload image
gho images upload path/to/image.jpg

# Upload with purpose
gho images upload avatar.png --purpose profile_image
gho images upload icon.png --purpose icon

# Specify reference ID
gho images upload banner.jpg --ref post-123
```

### Members

```bash
# List members
gho members list

# Limit number of results
gho members list --limit 10

# Apply filter
gho members list --filter "status:paid"

# Get member details
gho members info <id>

# Create new member
gho members create --email "user@example.com" --name "Taro Yamada"

# Create member with labels
gho members create --email "user@example.com" --labels "VIP,Premium"

# Update member
gho members update <id> --name "Hanako Tanaka" --note "Important customer"

# Delete member
gho members delete <id>
```

### Users

```bash
# List users
gho users list

# Include role information
gho users list --include roles

# Include post count
gho users list --include count.posts

# Get user details (by ID or Slug)
gho users info <id-or-slug>
gho users info slug:john-doe

# Update user information
gho users update <id> --name "New Name" --bio "New bio"
gho users update <id> --location "Tokyo" --website "https://example.com"
```

### Newsletters

```bash
# List newsletters
gho newsletters list

# Apply filter
gho newsletters list --filter "status:active"

# Get newsletter details (by ID or Slug)
gho newsletters info <id-or-slug>
gho newsletters info slug:weekly-newsletter

# Create new newsletter
gho newsletters create --name "Weekly Newsletter" --description "Delivered every Friday"

# Create with sender information
gho newsletters create --name "Monthly Letter" --sender-name "Editorial Team" --sender-email "editor@example.com"

# Update newsletter
gho newsletters update <id> --name "New Name"
gho newsletters update <id> --visibility paid --subscribe-on-signup=false
```

### Tiers

```bash
# List tiers
gho tiers list

# Include price information
gho tiers list --include monthly_price,yearly_price

# Get tier details (by ID or Slug)
gho tiers info <id-or-slug>
gho tiers info slug:premium

# Create new tier (free plan)
gho tiers create --name "Free Plan" --type free

# Create paid tier
gho tiers create --name "Premium" --type paid --monthly-price 1000 --yearly-price 10000 --currency JPY

# Create tier with benefits
gho tiers create --name "VIP" --type paid --monthly-price 3000 --benefits "Priority Support" --benefits "Exclusive Content"

# Update tier
gho tiers update <id> --name "New Premium"
gho tiers update <id> --monthly-price 1200 --yearly-price 12000
```

### Offers

```bash
# List offers
gho offers list

# Apply filter
gho offers list --filter "status:active"

# Get offer details
gho offers info <id>

# Create percentage discount offer
gho offers create --name "New Member Discount" --code "WELCOME2024" --type percent --amount 20 --tier-id <tier-id>

# Create fixed amount discount offer
gho offers create --name "500 Yen Off" --code "SAVE500" --type fixed --amount 500 --currency JPY --tier-id <tier-id>

# Create limited time offer
gho offers create --name "3 Month Discount" --code "TRIAL3M" --type percent --amount 50 --duration repeating --duration-in-months 3 --tier-id <tier-id>

# Update offer
gho offers update <id> --name "New Registration Campaign"
gho offers update <id> --amount 30
```

### Themes

```bash
# List themes
gho themes list

# Upload theme
gho themes upload path/to/theme.zip

# Activate theme
gho themes activate casper
```

### Webhooks

```bash
# Create webhook
gho webhooks create --event post.published --target-url https://example.com/webhook

# Create webhook with name
gho webhooks create --event member.added --target-url https://example.com/webhook --name "Member notification"

# Update webhook
gho webhooks update <id> --target-url https://new-example.com/webhook

# Delete webhook
gho webhooks delete <id>
```

## Output Formats

gho supports three output formats:

### Table Format (Default)

Outputs in human-readable format.

**info commands (single item)**:
```bash
$ gho site
title        My Blog
description  A technical blog about...
url          https://hanatane.net/
version      6.8
```

**list commands (multiple items)**:
```bash
$ gho posts list --limit 3
ID                        TITLE                     STATUS     CREATED     PUBLISHED
697b61d44921c40001f01aa3  Cannot/Don't Use CLI      draft      2026-01-29
696ce7244921c40001f017ed  Dev Environment 2026      published  2026-01-18  2026-01-28
```

### Plain Format (TSV)

Tab-separated format suitable for scripting and pipelines.

```bash
$ gho site --plain
title	My Blog
description	A technical blog about...
url	https://hanatane.net/
version	6.8

$ gho posts list --plain --limit 2
ID	TITLE	STATUS	CREATED	PUBLISHED
697b61d44921c40001f01aa3	Cannot/Don't Use CLI	draft	2026-01-29
696ce7244921c40001f017ed	Dev Environment 2026	published	2026-01-18	2026-01-28
```

### JSON Format

Format suitable for programmatic processing and API integration.

```bash
$ gho site --json
{
  "site": {
    "title": "My Blog",
    "description": "A technical blog about...",
    "url": "https://hanatane.net/",
    "version": "6.8"
  }
}
```

## Global Options

The following options are available for all commands:

```bash
# Output in JSON format
gho posts list --json

# Output in TSV format (for script integration)
gho posts list --plain

# Specify site
gho -s myblog posts list

# Skip confirmation (create/update/delete commands)
gho posts delete <id> --force
gho newsletters create --name "Test" --force
gho tiers update <id> --name "New Name" --force

# Display verbose logs
gho -v posts list
```

## Architecture

gho is implemented based on gogcli's design patterns with the following characteristics:

### Design Principles

- **Context propagation**: All commands use context to safely propagate output mode and UI instances
- **Exit code management**: Proper exit code control through ExitError type
- **TDD**: Robust implementation through Test-Driven Development
- **Type safety**: Maximum utilization of Go's type system

### Main Components

- **internal/cmd**: Command implementations (all commands have `Run(ctx context.Context, root *RootFlags) error` signature)
- **internal/outfmt**: Output format management (JSON/Table/Plain)
- **internal/ui**: UI output management (stdout/stderr separation)
- **internal/ghostapi**: Ghost Admin API client
- **internal/secrets**: OS keyring integration
- **internal/config**: Configuration file management

For details, see:
- [Design Unification Progress](./docs/gogcli-alignment-status.md)
- [Remaining Tasks Implementation Guide](./docs/remaining-tasks-guide.md)

## Development

### Run Tests

```bash
make test
```

### Run Lint

```bash
make lint
```

### Type Check

```bash
make type-check
```

### Build

```bash
make build
```

## License

MIT
