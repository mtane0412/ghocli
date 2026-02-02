# gho Implementation Plan

## Overview

Create a Ghost Admin API CLI tool with the user experience of gog-cli.

## Technology Stack

- **Language**: Go 1.22+
- **CLI Framework**: Kong (`github.com/alecthomas/kong`)
- **Credential Management**: 99designs/keyring (OS keyring integration)
- **JWT**: golang-jwt/jwt/v5

## Implementation Phases

### Phase 1: Foundation ✅

**Goal**: Build project foundation and implement authentication and site information retrieval

**Implementation**:

1. **Project Initialization**
   - `go mod init github.com/mtane0412/gho`
   - Add dependencies

2. **Configuration System** (`internal/config/`)
   - Configuration file: `~/.config/gho/config.json`
   - Multi-site support (alias functionality)
   - Structure:
     ```json
     {
       "keyring_backend": "auto",
       "default_site": "myblog",
       "sites": {
         "myblog": "https://myblog.ghost.io",
         "company": "https://blog.company.com"
       }
     }
     ```

3. **Keyring Integration** (`internal/secrets/`)
   - Secure storage of Admin API keys
   - Backends: auto/file/keychain/secretservice/wincred

4. **Ghost API Client** (`internal/ghostapi/`)
   - JWT generation (HS256, 5-minute expiration)
   - HTTP client
   - Site information retrieval API

5. **Output Format** (`internal/outfmt/`)
   - JSON format
   - Table format
   - TSV format (plain)

6. **Authentication Commands** (`internal/cmd/auth.go`)
   - `gho auth add <site-url>` - Register API key
   - `gho auth list` - List registered sites
   - `gho auth remove <alias>` - Delete API key
   - `gho auth status` - Check authentication status

7. **Basic Commands**
   - `gho site` - Get site information
   - `gho version` - Display version

**Verification**:
```bash
make build
./gho auth add https://your-ghost-site.ghost.io
./gho auth list
./gho site
```

**Completed**: ✅ 2026-01-29

---

### Phase 2: Content Management (Posts/Pages) ✅

**Goal**: Implement Posts/Pages create, update, delete, and publish functionality

**Completed**: 2026-01-29

**Implementation**:

1. **Posts API** (`internal/ghostapi/posts.go`)
   - `ListPosts(options ListOptions) ([]Post, error)`
   - `GetPost(idOrSlug string) (*Post, error)`
   - `CreatePost(post *Post) (*Post, error)`
   - `UpdatePost(id string, post *Post) (*Post, error)`
   - `DeletePost(id string) error`

2. **Posts Type Definition**
   ```go
   type Post struct {
       ID          string     `json:"id"`
       Title       string     `json:"title"`
       Slug        string     `json:"slug"`
       HTML        string     `json:"html,omitempty"`
       MobileDoc   string     `json:"mobiledoc,omitempty"`
       Status      string     `json:"status"` // draft/published/scheduled
       CreatedAt   time.Time  `json:"created_at"`
       PublishedAt *time.Time `json:"published_at,omitempty"`
       Tags        []Tag      `json:"tags,omitempty"`
       Authors     []Author   `json:"authors,omitempty"`
   }
   ```

3. **Posts Commands** (`internal/cmd/posts.go`)
   - `gho posts list [--status draft|published|scheduled] [--limit N]`
   - `gho posts get <id-or-slug>`
   - `gho posts create --title "..." [--html "..."]`
   - `gho posts update <id> --title "..."`
   - `gho posts delete <id>`
   - `gho posts publish <id>`

4. **Pages API** (`internal/ghostapi/pages.go`)
   - `ListPages(options ListOptions) ([]Page, error)`
   - `GetPage(idOrSlug string) (*Page, error)`
   - `CreatePage(page *Page) (*Page, error)`
   - `UpdatePage(id string, page *Page) (*Page, error)`
   - `DeletePage(id string) error`

5. **Pages Commands** (`internal/cmd/pages.go`)
   - `gho pages list`
   - `gho pages get <id-or-slug>`
   - `gho pages create --title "..."`
   - `gho pages update <id> ...`
   - `gho pages delete <id>`

**Tests**:
- Posts API tests (`internal/ghostapi/posts_test.go`)
- Posts command tests (`internal/cmd/posts_test.go`)
- Pages API tests (`internal/ghostapi/pages_test.go`)
- Pages command tests (`internal/cmd/pages_test.go`)

**Verification**:
```bash
./gho posts list
./gho posts get <slug>
./gho posts create --title "Test Post" --status draft
./gho posts publish <id>
./gho posts delete <id>

./gho pages list
./gho pages create --title "Test Page"
```

---

### Phase 3: Taxonomy + Media ✅

**Goal**: Implement Tags management and Images upload functionality

**Completed**: 2026-01-30

**Implementation**:

1. **Tags API** (`internal/ghostapi/tags.go`)
   - `ListTags() ([]Tag, error)`
   - `GetTag(idOrSlug string) (*Tag, error)`
   - `CreateTag(tag *Tag) (*Tag, error)`
   - `UpdateTag(id string, tag *Tag) (*Tag, error)`
   - `DeleteTag(id string) error`

2. **Tags Type Definition**
   ```go
   type Tag struct {
       ID          string `json:"id"`
       Name        string `json:"name"`
       Slug        string `json:"slug"`
       Description string `json:"description,omitempty"`
   }
   ```

3. **Tags Commands** (`internal/cmd/tags.go`)
   - `gho tags list`
   - `gho tags get <id-or-slug>`
   - `gho tags create --name "..."`
   - `gho tags update <id> --name "..."`
   - `gho tags delete <id>`

4. **Images API** (`internal/ghostapi/images.go`)
   - `UploadImage(filePath string) (*ImageUploadResponse, error)`

5. **Images Commands** (`internal/cmd/images.go`)
   - `gho images upload <file-path>`

**Tests**:
- Tags API tests
- Tags command tests
- Images API tests
- Images command tests

**Verification**:
```bash
./gho tags list
./gho tags create --name "Technology"
./gho images upload ./image.png
```

---

### Phase 4: Members Management

**Goal**: Implement Members (subscribers) management functionality

**Implementation**:

1. **Members API** (`internal/ghostapi/members.go`)
   - `ListMembers(options ListOptions) ([]Member, error)`
   - `GetMember(id string) (*Member, error)`
   - `CreateMember(member *Member) (*Member, error)`
   - `UpdateMember(id string, member *Member) (*Member, error)`
   - `DeleteMember(id string) error`

2. **Members Commands** (`internal/cmd/members.go`)
   - `gho members list`
   - `gho members get <id>`
   - `gho members create --email "..."`
   - `gho members update <id> ...`
   - `gho members delete <id>`

---

### Phase 5: Users Management

**Goal**: Implement Users (administrators/editors) management functionality

**Implementation**:

1. **Users API** (`internal/ghostapi/users.go`)
   - `ListUsers() ([]User, error)`
   - `GetUser(id string) (*User, error)`
   - `UpdateUser(id string, user *User) (*User, error)`

2. **Users Commands** (`internal/cmd/users.go`)
   - `gho users list`
   - `gho users get <id>`
   - `gho users update <id> ...`

---

### Phase 6: Newsletters/Tiers/Offers

**Goal**: Implement Newsletter, Tiers (subscription plans), and Offers (benefits) management functionality

**Implementation**:

1. **Newsletters API**
   - `ListNewsletters() ([]Newsletter, error)`
   - `GetNewsletter(id string) (*Newsletter, error)`

2. **Tiers API**
   - `ListTiers() ([]Tier, error)`
   - `GetTier(id string) (*Tier, error)`

3. **Offers API**
   - `ListOffers() ([]Offer, error)`
   - `GetOffer(id string) (*Offer, error)`

---

### Phase 7: Themes/Webhooks

**Goal**: Implement Themes management and Webhooks management functionality

**Implementation**:

1. **Themes API**
   - `ListThemes() ([]Theme, error)`
   - `UploadTheme(filePath string) error`
   - `ActivateTheme(name string) error`
   - `DeleteTheme(name string) error`

2. **Webhooks API**
   - `ListWebhooks() ([]Webhook, error)`
   - `CreateWebhook(webhook *Webhook) (*Webhook, error)`
   - `DeleteWebhook(id string) error`

---

## Development Workflow

### TDD Principles

Follow the TDD cycle for all implementations:

1. **RED** - Write a failing test first
2. **GREEN** - Write minimal code to make the test pass
3. **REFACTOR** - Clean up the code

### Quality Checks

Execute the following after completing each phase:

```bash
# Run tests
make test

# Type check
make type-check

# Lint
make lint

# Build verification
make build
```

### Git Workflow

```bash
# Create feature branch for each phase
git checkout -b feature/phase2-content-management

# Pre-commit check
make test
make type-check
make lint

# Create commit
git add .
git commit -m "Phase 2: Implement content management features

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Reference Resources

### Ghost Admin API Documentation
- https://ghost.org/docs/admin-api/

### Reference Project (gog-cli)
- `../gogcli/internal/cmd/root.go` - CLI structure pattern
- `../gogcli/internal/cmd/auth.go` - Authentication command implementation
- `../gogcli/internal/secrets/store.go` - Keyring integration
- `../gogcli/internal/config/config.go` - Configuration file management
- `../gogcli/internal/outfmt/outfmt.go` - Output formatting
