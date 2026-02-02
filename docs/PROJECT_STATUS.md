# Project Status

## Overview

**gho** is a CLI tool for the Ghost Admin API. It provides a gog-cli-like user experience and allows operations on the Ghost Admin API from the command line.

## Implementation Phases

### ✅ Phase 1: Foundation (Completed)

**Completion Date**: 2026-01-29

**Implementation Details**:

1. **Project Initialization**
   - Go modules initialization
   - Dependencies added (Kong, Keyring, JWT)

2. **Configuration System** (`internal/config/`)
   - Configuration file management (`~/.config/gho/config.json`)
   - Multi-site support (alias feature)
   - Default site management

3. **Keyring Integration** (`internal/secrets/`)
   - Secure API key storage via OS keyring
   - Support for macOS Keychain, Linux Secret Service, Windows Credential Manager
   - API key parsing functionality

4. **Ghost API Client** (`internal/ghostapi/`)
   - JWT generation (HS256, 5-minute expiration)
   - HTTP client
   - Site information retrieval API

5. **Output Formats** (`internal/outfmt/`)
   - JSON format
   - Table format (human-friendly, gogcli-style)
   - Plain format (TSV, for program integration)

6. **Auth Commands** (`internal/cmd/auth.go`)
   ```
   gho auth add <site-url>      # Register API key
   gho auth list                # List registered sites
   gho auth remove <alias>      # Remove API key
   gho auth status              # Check authentication status
   ```

7. **Basic Commands**
   ```
   gho site                     # Get site information
   gho version                  # Display version
   ```

**Quality Checks**:
- ✅ All tests pass
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Commit**: `68b9340 Phase 1: Complete foundation implementation`

### ✅ Phase 2: Content Management (Posts/Pages) (Completed)

**Completion Date**: 2026-01-29

**Implementation Details**:

1. **Posts API** (`internal/ghostapi/posts.go`)
   - Post type definition (ID, Title, Slug, HTML, Status, PublishedAt, etc.)
   - ListOptions type definition (Limit, Status, Filter, etc.)
   - `ListPosts(options ListOptions) ([]Post, error)` implementation
   - `GetPost(idOrSlug string) (*Post, error)` implementation
   - `CreatePost(post *Post) (*Post, error)` implementation
   - `UpdatePost(id string, post *Post) (*Post, error)` implementation
   - `DeletePost(id string) error` implementation

2. **Pages API** (`internal/ghostapi/pages.go`)
   - Page type definition (ID, Title, Slug, HTML, Status, etc.)
   - `ListPages(options ListOptions) ([]Page, error)` implementation
   - `GetPage(idOrSlug string) (*Page, error)` implementation
   - `CreatePage(page *Page) (*Page, error)` implementation
   - `UpdatePage(id string, page *Page) (*Page, error)` implementation
   - `DeletePage(id string) error` implementation

3. **Posts Commands** (`internal/cmd/posts.go`)
   ```
   gho posts list [--status draft|published|scheduled] [--limit N]
   gho posts info <id-or-slug>   # Previously: gho posts get (backward compatible)
   gho posts create --title "..." [--html "..."] [--status draft|published]
   gho posts update <id> [--title "..."] [--html "..."]
   gho posts delete <id>
   gho posts publish <id>
   ```

4. **Pages Commands** (`internal/cmd/pages.go`)
   ```
   gho pages list [--status draft|published|scheduled] [--limit N]
   gho pages info <id-or-slug>   # Previously: gho pages get (backward compatible)
   gho pages create --title "..." [--html "..."]
   gho pages update <id> [--title "..."] [--html "..."]
   gho pages delete <id>
   ```

**Quality Checks**:
- ✅ All tests pass (Posts: 7 tests, Pages: 5 tests)
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Commits**:
- `40c33f2 feat(ghostapi): Implement Posts API`
- `016fe5c feat(ghostapi): Implement Pages API`
- `a84e3da feat(cmd): Implement Posts/Pages commands`

### ✅ Phase 3: Taxonomy + Media (Completed)

**Completion Date**: 2026-01-30

**Implementation Details**:

1. **Tags API** (`internal/ghostapi/tags.go`)
   - Tag type definition (ID, Name, Slug, Description, Visibility, etc.)
   - TagListOptions type definition (pagination, filter support)
   - `ListTags(options TagListOptions) (*TagListResponse, error)` implementation
   - `GetTag(idOrSlug string) (*Tag, error)` implementation ("slug:" prefix support)
   - `CreateTag(tag *Tag) (*Tag, error)` implementation
   - `UpdateTag(id string, tag *Tag) (*Tag, error)` implementation
   - `DeleteTag(id string) error` implementation

2. **Images API** (`internal/ghostapi/images.go`)
   - Image type definition (URL, Ref)
   - `UploadImage(file io.Reader, filename string, opts ImageUploadOptions) (*Image, error)` implementation
   - multipart/form-data upload support
   - Purpose specification support (image/profile_image/icon)

3. **Tags Commands** (`internal/cmd/tags.go`)
   ```
   gho tags list [--limit N] [--page N]
   gho tags info <id-or-slug>       # Previously: gho tags get (backward compatible)
   gho tags create --name "..." [--description "..."] [--visibility public|internal]
   gho tags update <id> [--name "..."] [--description "..."]
   gho tags delete <id>
   ```

4. **Images Commands** (`internal/cmd/images.go`)
   ```
   gho images upload <file-path> [--purpose image|profile_image|icon] [--ref <ref-id>]
   ```

**Quality Checks**:
- ✅ All tests pass (Tags: 6 tests, Images: 2 tests)
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Commits**:
- `b5299e8 feat(api): Implement Tags API and Images API`

### ✅ Phase 4: Members Management (Completed)

**Completion Date**: 2026-01-30

**Implementation Details**:

1. **Members API** (`internal/ghostapi/members.go`)
   - Member type definition (ID, UUID, Email, Name, Note, Status, Labels, etc.)
   - Label type definition (ID, Name, Slug)
   - MemberListOptions type definition (pagination, filter, order support)
   - `ListMembers(options MemberListOptions) (*MemberListResponse, error)` implementation
   - `GetMember(id string) (*Member, error)` implementation
   - `CreateMember(member *Member) (*Member, error)` implementation
   - `UpdateMember(id string, member *Member) (*Member, error)` implementation
   - `DeleteMember(id string) error` implementation

2. **Members Commands** (`internal/cmd/members.go`)
   ```
   gho members list [--limit N] [--page N] [--filter "..."] [--order "..."]
   gho members info <id>            # Previously: gho members get (backward compatible)
   gho members create --email "..." [--name "..."] [--note "..."] [--labels "..."]
   gho members update <id> [--name "..."] [--note "..."] [--labels "..."]
   gho members delete <id>
   ```

**Quality Checks**:
- ✅ All tests pass (Members: 6 tests)
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Commits**:
- `3a935e6 feat(api): Implement Members API`

### ✅ Phase 5: Users Management (Completed)

**Completion Date**: 2026-01-30

**Implementation Details**:

1. **Users API** (`internal/ghostapi/users.go`)
   - User type definition (ID, Name, Slug, Email, Bio, Location, Website, ProfileImage, CoverImage, Roles, etc.)
   - Role type definition (ID, Name)
   - UserListOptions type definition (pagination, include, filter support)
   - `ListUsers(options UserListOptions) (*UserListResponse, error)` implementation
   - `GetUser(idOrSlug string) (*User, error)` implementation ("slug:" prefix support)
   - `UpdateUser(id string, user *User) (*User, error)` implementation
   - **Note**: Create/Delete operations not supported (use Ghost dashboard invite feature)

2. **Users Commands** (`internal/cmd/users.go`)
   ```
   gho users list [--limit N] [--page N] [--include roles,count.posts]
   gho users info <id-or-slug>      # Previously: gho users get (backward compatible)
   gho users update <id> [--name "..."] [--slug "..."] [--bio "..."] [--location "..."] [--website "..."]
   ```

**Quality Checks**:
- ✅ All tests pass (Users: 7 tests)
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Commits**:
- `1884ff0 feat(api): Implement Users API`

### ✅ Phase 6: Newsletters/Tiers/Offers (Completed)

**Completion Date**: 2026-01-30

**Implementation Details**:

1. **Newsletters API** (`internal/ghostapi/newsletters.go`)
   - Newsletter type definition (ID, Name, Slug, Description, Status, SubscribeOnSignup, etc.)
   - NewsletterListOptions type definition (pagination, filter support)
   - `ListNewsletters(options NewsletterListOptions) (*NewsletterListResponse, error)` implementation
   - `GetNewsletter(idOrSlug string) (*Newsletter, error)` implementation ("slug:" prefix support)
   - `CreateNewsletter(newsletter *Newsletter) (*Newsletter, error)` implementation
   - `UpdateNewsletter(id string, newsletter *Newsletter) (*Newsletter, error)` implementation

2. **Tiers API** (`internal/ghostapi/tiers.go`)
   - Tier type definition (ID, Name, Slug, Type, MonthlyPrice, YearlyPrice, etc.)
   - TierListOptions type definition (pagination, include support)
   - `ListTiers(options TierListOptions) (*TierListResponse, error)` implementation
   - `GetTier(idOrSlug string) (*Tier, error)` implementation ("slug:" prefix support)
   - `CreateTier(tier *Tier) (*Tier, error)` implementation
   - `UpdateTier(id string, tier *Tier) (*Tier, error)` implementation

3. **Offers API** (`internal/ghostapi/offers.go`)
   - Offer type definition (ID, Name, Code, Tier, DiscountType, DiscountAmount, etc.)
   - OfferListOptions type definition (pagination, filter support)
   - `ListOffers(options OfferListOptions) (*OfferListResponse, error)` implementation
   - `GetOffer(id string) (*Offer, error)` implementation
   - `CreateOffer(offer *Offer) (*Offer, error)` implementation
   - `UpdateOffer(id string, offer *Offer) (*Offer, error)` implementation

4. **Newsletters Commands** (`internal/cmd/newsletters.go`)
   ```
   gho newsletters list [--limit N] [--page N] [--filter "..."]
   gho newsletters info <id-or-slug>   # Previously: gho newsletters get (backward compatible)
   gho newsletters create --name "..." [--description "..."] [--visibility members|paid]
   gho newsletters update <id> [--name "..."] [--visibility "..."] [--sender-name "..."]
   ```

5. **Tiers Commands** (`internal/cmd/tiers.go`)
   ```
   gho tiers list [--limit N] [--page N] [--include monthly_price,yearly_price]
   gho tiers info <id-or-slug>         # Previously: gho tiers get (backward compatible)
   gho tiers create --name "..." [--type free|paid] [--monthly-price N] [--yearly-price N]
   gho tiers update <id> [--name "..."] [--monthly-price N] [--yearly-price N]
   ```

6. **Offers Commands** (`internal/cmd/offers.go`)
   ```
   gho offers list [--limit N] [--page N] [--filter "..."]
   gho offers info <id>                # Previously: gho offers get (backward compatible)
   gho offers create --name "..." --code "..." --type percent|fixed --amount N --tier-id <tier-id>
   gho offers update <id> [--name "..."] [--amount N]
   ```

7. **Destructive Operation Confirmation** (`internal/cmd/helpers.go`)
   - Confirmation prompts displayed for Create/Update operations
   - `--force` flag to skip confirmation

**Quality Checks**:
- ✅ All tests pass (Newsletters: 6 tests, Tiers: 6 tests, Offers: 6 tests)
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Commits**:
- `4545035 feat(api): Implement Newsletters, Tiers, Offers APIs`
- `eed5ff2 feat(newsletters): Implement Newsletters write operations`
- `8b158df feat(tiers): Implement Tiers write operations`
- `2874e8d feat(offers): Implement Offers write operations`
- `013086c feat(cmd): Implement destructive operation confirmation mechanism`

### ✅ Phase 7: Themes/Webhooks API (Completed)

**Completion Date**: 2026-01-30

**Implementation Details**:

1. **Themes API** (`internal/ghostapi/themes.go`)
   - Theme type definition (Name, Package, Active, Templates, etc.)
   - ThemePackage type definition (Name, Description, Version)
   - ThemeTemplate type definition (Filename)
   - `ListThemes() (*ThemeListResponse, error)` implementation
   - `UploadTheme(file io.Reader, filename string) (*Theme, error)` implementation (multipart upload)
   - `ActivateTheme(name string) (*Theme, error)` implementation

2. **Webhooks API** (`internal/ghostapi/webhooks.go`)
   - Webhook type definition (ID, Event, TargetURL, Name, Secret, APIVersion, IntegrationID, Status, LastTriggeredAt, CreatedAt, UpdatedAt, etc.)
   - `CreateWebhook(webhook *Webhook) (*Webhook, error)` implementation
   - `UpdateWebhook(id string, webhook *Webhook) (*Webhook, error)` implementation
   - `DeleteWebhook(id string) error` implementation
   - **Note**: Ghost API does not support Webhook List/Get

3. **Themes Commands** (`internal/cmd/themes.go`)
   ```
   gho themes list                    # List themes
   gho themes upload <file.zip>       # Upload theme
   gho themes activate <name>         # Activate theme
   ```

4. **Webhooks Commands** (`internal/cmd/webhooks.go`)
   ```
   gho webhooks create --event <event> --target-url <url> [--name <name>]
   gho webhooks update <id> [--event <event>] [--target-url <url>] [--name <name>]
   gho webhooks delete <id>
   ```

**Quality Checks**:
- ✅ All tests pass (Themes: 3 tests, Webhooks: 3 tests)
- ✅ Type check (`go vet`) successful
- ✅ Lint (golangci-lint) successful
- ✅ Build successful

**Commits**:
- `3a6a9ed feat(api): Implement Themes/Webhooks APIs`

## Current Structure

```
gho/
├── cmd/gho/
│   └── main.go              # Entry point
├── internal/
│   ├── cmd/                  # CLI command definitions
│   │   ├── root.go          # CLI struct, RootFlags
│   │   ├── auth.go          # Auth commands
│   │   ├── site.go          # Site info commands
│   │   ├── posts.go         # Posts commands
│   │   ├── pages.go         # Pages commands
│   │   ├── tags.go          # Tags commands
│   │   ├── images.go        # Images commands
│   │   ├── members.go       # Members commands
│   │   ├── users.go         # Users commands
│   │   ├── newsletters.go   # Newsletters commands
│   │   ├── tiers.go         # Tiers commands
│   │   ├── offers.go        # Offers commands
│   │   ├── themes.go        # Themes commands
│   │   └── webhooks.go      # Webhooks commands
│   ├── config/              # Configuration file management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── secrets/             # Keyring integration
│   │   ├── store.go
│   │   └── store_test.go
│   ├── ghostapi/            # Ghost API client
│   │   ├── client.go        # HTTP client
│   │   ├── client_test.go
│   │   ├── jwt.go           # JWT generation
│   │   ├── jwt_test.go
│   │   ├── posts.go         # Posts API
│   │   ├── posts_test.go
│   │   ├── pages.go         # Pages API
│   │   ├── pages_test.go
│   │   ├── tags.go          # Tags API
│   │   ├── tags_test.go
│   │   ├── images.go        # Images API
│   │   ├── images_test.go
│   │   ├── members.go       # Members API
│   │   ├── members_test.go
│   │   ├── users.go         # Users API
│   │   ├── users_test.go
│   │   ├── newsletters.go   # Newsletters API
│   │   ├── newsletters_test.go
│   │   ├── tiers.go         # Tiers API
│   │   ├── tiers_test.go
│   │   ├── offers.go        # Offers API
│   │   ├── offers_test.go
│   │   ├── themes.go        # Themes API
│   │   ├── themes_test.go
│   │   ├── webhooks.go      # Webhooks API
│   │   └── webhooks_test.go
│   └── outfmt/              # Output formatting
│       ├── outfmt.go
│       └── outfmt_test.go
├── docs/                    # Documentation
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── .gitignore
└── README.md
```

## Test Coverage

All core components are tested:

- `internal/config/` - Configuration file management (6 tests)
- `internal/secrets/` - Keyring integration (8 tests)
- `internal/ghostapi/` - API client (47 tests)
  - `client.go`, `jwt.go` - 11 tests
  - `posts.go` - 7 tests
  - `pages.go` - 5 tests
  - `tags.go` - 7 tests
  - `images.go` - 2 tests
  - `members.go` - 6 tests
  - `users.go` - 7 tests
  - `newsletters.go` - 6 tests (List, Get, Create, Update)
  - `tiers.go` - 6 tests (List, Get, Create, Update)
  - `offers.go` - 6 tests (List, Get, Create, Update)
  - `themes.go` - 3 tests
  - `webhooks.go` - 3 tests
- `internal/outfmt/` - Output formatting (5 tests)

Total: 66 tests, all passing

## Dependencies

```
github.com/alecthomas/kong v1.13.0        # CLI framework
github.com/99designs/keyring v1.2.2       # Keyring integration
github.com/golang-jwt/jwt/v5 v5.3.1       # JWT generation
github.com/k3a/html2text v1.3.0           # HTML→text conversion
```

## Quality Check Commands

```bash
# Run tests
make test

# Type check
make type-check

# Run lint (requires golangci-lint)
make lint

# Build
make build
```

## Bug Fix History

### 2026-01-30: JWT Signature Error and Posts/Pages Update Lock Error Fix

**Problems**:
1. `Invalid token: invalid signature` error during Ghost Admin API communication
2. `Someone else is editing this post` error when updating posts/pages

**Causes**:
1. Ghost Admin API secret key is provided as a hexadecimal string, but was being used directly without binary decoding during JWT signing
2. When updating posts/pages, a new timestamp was generated and sent instead of the original `updated_at` timestamp retrieved from the server (violating Ghost API's optimistic locking mechanism)

**Fixes**:
- `internal/ghostapi/jwt.go:46-50` - Decode secret from hexadecimal to binary before signing
- `internal/ghostapi/jwt_test.go:9,58-63` - Update test code to support hexadecimal decoding
- All test files - Standardize test secrets to hexadecimal format
- `internal/cmd/posts.go:202` - Change `UpdatedAt: time.Now()` to `UpdatedAt: existingPost.UpdatedAt`
- `internal/cmd/posts.go:310` - Same fix for publish command
- `internal/cmd/pages.go:201` - Same fix for page updates

**Tests Added**:
- `TestUpdatePost_updated_atを保持して更新` - Test to verify original `updated_at` is sent during updates

**Verification**:
- All read operations (Posts, Tags, Users, Newsletters, Tiers, Pages, etc.) working correctly
- Write operations (create, update, delete) working correctly
- All 81 tests passing

**References**:
- [Ghost Admin API Overview](https://docs.ghost.org/admin-api)
- [Bash Example of Ghost JWT Auth](https://gist.github.com/ErisDS/6334f0e70ec7390ec08530d5ef9bd0d5)

### ✅ Phase 8: Command Design Improvements (In Progress)

**Start Date**: 2026-01-31

**Purpose**: Improve gho's command system based on gogcli's command design patterns

**Implementation Details**:

#### 8.1: get → info Rename (Completed)
- Renamed `get` command to `info` for all resources (posts, pages, members, tags, users, newsletters, tiers, offers)
- `get` maintained as alias (deprecated) for backward compatibility
- Added deprecation message to help

**Examples**:
```bash
# New command
gho posts info <id>
gho pages info <slug>
gho members info <id>

# Old command (works with deprecation warning)
gho posts get <id>
```

#### 8.2: cat Command Addition (Completed)
- Added `cat` command for posts/pages to display body content to stdout
- `--format` option to select html/text/lexical format
- Text format uses k3a/html2text library to convert HTML to plain text

**Usage Examples**:
```bash
gho posts cat <id>                      # Output in HTML format
gho posts cat <id> --format text        # Output in text format (HTML tags removed)
gho posts cat <id> --format lexical     # Output in Lexical JSON format
gho pages cat <slug> --format html      # Output page body in HTML format
```

#### 8.3: copy Command Addition (Completed)
- Added `copy` command to copy posts/pages
- Excludes ID/UUID/Slug/URL/timestamps and creates as new
- Status always set to `draft` on creation
- `--title` option to specify new title (defaults to "Original Title (Copy)")

**Usage Examples**:
```bash
gho posts copy <id-or-slug>                      # Copy post (title: "Original Title (Copy)")
gho posts copy <id-or-slug> --title "New Title"  # Copy with custom title
gho pages copy <slug>                            # Copy page
gho pages copy <slug> --title "New Title"        # Copy with custom title
```

#### 8.5: Complete Migration to gogcli-style Output Format (Completed)
- Changed info command key names to snake_case (lowercase)
- Removed separator lines in list commands
- Unified table and plain formats (auto-alignment via tabwriter)
- All info commands now use PrintKeyValue

**Changes**:
- info commands: No header, lowercase keys, tab-separated
- list commands: Uppercase headers, no separators, tab-separated
- Changed members/users/tags info commands from PrintTable to PrintKeyValue

**Output Examples**:
```bash
# info commands (table format)
$ gho site
title        はなしのタネ
description  技術・学問・ゲーム・田舎暮らしを中心に...
url          https://hanatane.net/
version      6.8

# info commands (plain format)
$ gho site --plain
title	はなしのタネ
description	技術・学問・ゲーム...
url	https://hanatane.net/
version	6.8

# list commands (table format)
$ gho posts list --limit 2
ID                        TITLE                               STATUS     CREATED     PUBLISHED
697b61d44921c40001f01aa3  CLIを使えない/使わない              draft      2026-01-29
696ce7244921c40001f017ed  非エンジニアおじさんの開発環境2026  published  2026-01-18  2026-01-28
```

#### 8.4: cat Command text Format Implementation (Completed)
- Implemented `--format text` for posts/pages `cat` command to work correctly
- Uses k3a/html2text library to convert HTML to plain text
- Judged that dedicated export command is unnecessary as shell redirection (`gho posts cat <id> --format html > output.html`) allows exporting

**Implementation Locations**:
- `internal/cmd/posts.go:939` - HTML to text conversion implementation
- `internal/cmd/pages.go:497` - Same as above
- `go.mod` - Added k3a/html2text v1.3.0

**Quality Checks (Phase 8.5 completion)**:
- ✅ All tests pass (164 tests)
- ✅ Type check (`go vet`) successful
- ✅ Lint (golangci-lint) successful (0 issues)
- ✅ Build successful

**Commits**:
- `dec99de feat(cmd): Phase 8.1 - get → info rename (all resources)`
- `a1d6f61 feat(cmd): Phase 8.2 - Add cat command (posts, pages)`
- `18ab842 feat(cmd): Phase 8.3 - Add copy command (posts, pages)`
- `bf260b7 feat(cmd): Phase 8.4 - Implement text format for cat command`
- `6f1a8b5 feat(outfmt): Remove headers from info commands and adopt gogcli style`
- `043ac98 feat(outfmt): Complete migration to gogcli style`

**Reference**: gogcli command design patterns (`gog docs info/cat/copy/export`)

#### 8.6: Field Selection Feature Implementation (In Progress)

**Completion Date**: 2026-02-01 (Phase 1-6 completed)

Implemented field selection feature supporting all Ghost Admin API fields with gh CLI-style `--fields` option.

**Completed Implementation**:

1. **Foundation Setup**
   - `internal/fields/` - Field definition package
     - `fields.go`: Parse, Validate, ListAvailable functions
     - `posts.go`: Post field definitions (40 fields)
   - Implemented unit tests for all functions

2. **Struct Extensions**
   - `internal/ghostapi/types.go` - Added Author struct
   - `internal/ghostapi/posts.go` - Extended Post struct to 40+ fields
     - Basic info, content, images, SEO, timestamps, control, custom, related, other, email/newsletter
   - Implemented JSON conversion tests

3. **Command Layer Preparation**
   - `internal/cmd/root.go` - Added `Fields` field to RootFlags
     - `-F, --fields` option
     - `GHO_FIELDS` environment variable support
   - Implemented environment variable and flag priority tests

4. **Output Layer Extensions**
   - `internal/outfmt/filter.go` - Field filtering functionality
     - `FilterFields()`: Extract and output only specified fields
     - `StructToMap()`: Convert struct to map[string]interface{}
   - Implemented tests for JSON/Plain/Table formats

5. **posts list Command Implementation**
   - `internal/cmd/posts.go` - Extended PostsListCmd.Run
     - JSON alone (without --fields): Display available fields list
     - With field specification: Output only specified fields

**Usage Examples**:
```bash
# Display field list
gho posts list --json

# Get only specified fields (JSON)
gho posts list --json --fields id,title,status,excerpt

# Specify fields in Plain format (TSV)
gho posts list --plain --fields id,title,url

# Field specification also works in table format
gho posts list --fields id,title,status,feature_image

# Specify via environment variable
export GHO_FIELDS="id,title,status"
gho posts list --json

# Short option
gho posts list --json -F id,title,url
```

**Quality Checks**:
- ✅ All tests pass
- ✅ Type check (`go vet`) successful
- ✅ Build successful

**Remaining Tasks**:
- posts get command support
- Horizontal expansion to other resources (pages, tags, members, users, etc.)
- Field definition creation for each resource

**Details**: See `docs/fields-feature-implementation.md`

## Next Steps

Phase 7 is complete, and all major Ghost Admin API feature implementations are finished.
Currently proceeding with Phase 8 (Command Design Improvements).

### Phase 8.6: Continue Field Selection Feature Implementation

Checklist for next work session:
1. Add --fields support for posts get command
2. Horizontal expansion to other resources (pages, tags, members, users, etc.)
3. Create field definitions for each resource

See `docs/fields-feature-implementation.md` for details.

For future extension features, see `docs/NEXT_STEPS.md`.
