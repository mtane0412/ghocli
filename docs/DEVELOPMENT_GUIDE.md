# gho Development Guide

## Development Environment Setup

### Required Tools

- **Go**: 1.22 or later
- **Make**: Build automation
- **golangci-lint**: Lint execution (optional)
- **Git**: Version control

### Setup Steps

```bash
# Clone repository
git clone https://github.com/mtane0412/gho.git
cd gho

# Install dependencies
go mod download

# Build
make build

# Run tests
make test
```

## Development Workflow

### TDD Principles

Follow the TDD cycle for all implementations:

1. **RED** - Write a failing test first
2. **GREEN** - Write minimal code to make the test pass
3. **REFACTOR** - Clean up the code

### Implementation Example

```go
// 1. RED: Write a failing test first
func TestGenerateJWT_GeneratesCorrectFormatToken(t *testing.T) {
    token, err := GenerateJWT("keyid", "secret")
    if err != nil {
        t.Fatalf("Failed to generate JWT: %v", err)
    }
    if token == "" {
        t.Error("Generated token is empty")
    }
}

// 2. GREEN: Write minimal code to pass the test
func GenerateJWT(keyID, secret string) (string, error) {
    // Minimal implementation
}

// 3. REFACTOR: Clean up the code
func GenerateJWT(keyID, secret string) (string, error) {
    // Refactored implementation
}
```

## Coding Conventions

### File Header Comments

Add specifications at the beginning of each file:

```go
/**
 * jwt.go
 * JWT generation for Ghost Admin API
 *
 * Ghost Admin API requires JWT signed with HS256 algorithm.
 * Token expiration is 5 minutes.
 */

package ghostapi
```

### Function Comments

Describe purpose, content, and notes in detail:

```go
// GenerateJWT generates a JWT token for the Ghost Admin API.
// keyID: ID part of the Admin API key
// secret: Secret part of the Admin API key
func GenerateJWT(keyID, secret string) (string, error) {
    // ...
}
```

### Test Function Naming

Test function names should describe the specific content:

```go
// ✅ Good examples
func TestGenerateJWT_GeneratesCorrectFormatToken(t *testing.T) { }
func TestGenerateJWT_ErrorOnEmptyKeyID(t *testing.T) { }

// ❌ Bad examples
func TestGenerateJWT(t *testing.T) { }
func TestJWT1(t *testing.T) { }
```

### Error Handling

Error messages should be specific and descriptive:

```go
// ✅ Good example
if keyID == "" {
    return "", errors.New("key ID is empty")
}

// ❌ Bad example
if keyID == "" {
    return "", errors.New("invalid key")
}
```

Use `fmt.Errorf` with `%w` for error wrapping:

```go
if err := store.Set(alias, apiKey); err != nil {
    return fmt.Errorf("failed to save API key: %w", err)
}
```

### Struct Tags

Add appropriate tags to JSON structs:

```go
type Site struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Version     string `json:"version"`
}
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Verbose output
go test ./... -v

# Specific package only
go test ./internal/config/... -v

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Writing Tests

#### Unit Tests

```go
func TestStore_SetAndGetBasicOperation(t *testing.T) {
    // Create test store (using file backend)
    store, err := NewStore("file", t.TempDir())
    if err != nil {
        t.Fatalf("Failed to create store: %v", err)
    }

    // Save API key
    testKey := "64fac5417c4c6b0001234567:89abcdef..."
    if err := store.Set("testsite", testKey); err != nil {
        t.Fatalf("Failed to save API key: %v", err)
    }

    // Retrieve API key
    retrieved, err := store.Get("testsite")
    if err != nil {
        t.Fatalf("Failed to retrieve API key: %v", err)
    }

    // Verify saved and retrieved keys match
    if retrieved != testKey {
        t.Errorf("Retrieved key = %q; want %q", retrieved, testKey)
    }
}
```

#### HTTP Client Testing

Use `httptest` package to create mock server:

```go
func TestGetSite_RetrievesSiteInfo(t *testing.T) {
    // Create test HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Return response
        response := map[string]interface{}{
            "site": map[string]interface{}{
                "title": "Test Blog",
                "url":   "https://test.ghost.io",
            },
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }))
    defer server.Close()

    // Create client
    client, err := NewClient(server.URL, "keyid", "secret")
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }

    // Get site info
    site, err := client.GetSite()
    if err != nil {
        t.Fatalf("Failed to get site info: %v", err)
    }

    // Verify response
    if site.Title != "Test Blog" {
        t.Errorf("Title = %q; want %q", site.Title, "Test Blog")
    }
}
```

#### Table-Driven Tests

Efficiently test multiple test cases:

```go
func TestParseAdminAPIKey(t *testing.T) {
    testCases := []struct {
        name       string
        input      string
        wantID     string
        wantSecret string
        wantErr    bool
    }{
        {
            name:       "correct format",
            input:      "64fac5417c4c6b0001234567:89abcdef...",
            wantID:     "64fac5417c4c6b0001234567",
            wantSecret: "89abcdef...",
            wantErr:    false,
        },
        {
            name:    "no colon",
            input:   "64fac5417c4c6b000123456789abcdef",
            wantErr: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            id, secret, err := ParseAdminAPIKey(tc.input)

            if tc.wantErr {
                if err == nil {
                    t.Error("Expected error but got nil")
                }
                return
            }

            if err != nil {
                t.Fatalf("Unexpected error: %v", err)
            }

            if id != tc.wantID {
                t.Errorf("id = %q; want %q", id, tc.wantID)
            }
            if secret != tc.wantSecret {
                t.Errorf("secret = %q; want %q", secret, tc.wantSecret)
            }
        })
    }
}
```

## Quality Checks

### Pre-Commit Checks

Always execute the following before committing:

```bash
# Run tests
make test

# Type check
make type-check

# Lint (requires golangci-lint)
make lint

# Build verification
make build
```

### Type Checking

```bash
# Type check with go vet
go vet ./...

# Or
make type-check
```

### Lint

Install golangci-lint:

```bash
# macOS
brew install golangci-lint

# Linux/Windows
# https://golangci-lint.run/usage/install/
```

Run lint:

```bash
golangci-lint run

# Or
make lint
```

## Git Workflow

### Branching Strategy

```bash
# Never commit directly to main branch
# Always work on feature branches

# Check current branch before starting work
git branch --show-current

# Create feature branch
git checkout -b feature/phase2-content-management
```

### Commit Messages

```bash
git commit -m "$(cat <<'EOF'
Phase 2: Implement content management features

Implemented Posts/Pages create, update, delete, and publish functionality.

Main implementation:
- Posts API (list/get/create/update/delete/publish)
- Pages API (list/get/create/update/delete)
- Posts/Pages commands
- Added tests

All tests passing.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
EOF
)"
```

## Adding New API Resources

### 1. Add API Type Definitions

`internal/ghostapi/posts.go`:

```go
package ghostapi

import "time"

// Post represents a Ghost post
type Post struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    HTML        string     `json:"html,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    PublishedAt *time.Time `json:"published_at,omitempty"`
}
```

### 2. Write Tests First (RED)

`internal/ghostapi/posts_test.go`:

```go
func TestListPosts_RetrievesPostList(t *testing.T) {
    // Create test HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Return mock response
    }))
    defer server.Close()

    // Write test
}
```

### 3. Add Implementation (GREEN)

`internal/ghostapi/posts.go`:

```go
// ListPosts retrieves a list of posts
func (c *Client) ListPosts(options ListOptions) ([]Post, error) {
    // Implementation
}
```

### 4. Add Commands

`internal/cmd/posts.go`:

```go
type PostsCmd struct {
    List   PostsListCmd   `cmd:"" help:"List posts"`
    Get    PostsGetCmd    `cmd:"" help:"Get a post"`
    Create PostsCreateCmd `cmd:"" help:"Create a post"`
}
```

### 5. Register in root.go

`internal/cmd/root.go`:

```go
type CLI struct {
    RootFlags `embed:""`
    Version   kong.VersionFlag
    Auth      AuthCmd
    Site      SiteCmd
    Posts     PostsCmd  // Add
}
```

## Debugging

### Logging

Enable logging with `--verbose` flag:

```go
if root.Verbose {
    log.Printf("API request: %s %s", method, url)
}
```

### JWT Debugging

Decode token at jwt.io:

```bash
# Get token
./gho site --verbose

# Visit jwt.io and paste the token
```

### HTTP Request Debugging

Enable detailed logging with environment variable:

```bash
# HTTP debug
export GODEBUG=http2debug=1
./gho site
```

## Troubleshooting

### Tests Failing

```bash
# Clear cache
go clean -testcache

# Re-run
go test ./...
```

### Build Errors

```bash
# Update dependencies
go mod tidy

# Rebuild
make build
```

### Keyring Errors

```bash
# Test with file backend
export GHO_KEYRING_BACKEND=file
./gho auth add https://test.ghost.io
```

## Release

### Version Tags

```bash
# Create tag
git tag -a v0.1.0 -m "Release v0.1.0"

# Push tag
git push origin v0.1.0
```

### Build

```bash
# Build with version
make build VERSION=0.1.0
```

## Reference Resources

### Go Language

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Ghost Admin API

- [Ghost Admin API Documentation](https://ghost.org/docs/admin-api/)
- [Ghost API Client Examples](https://github.com/TryGhost/Ghost/tree/main/ghost/admin-api)

### Testing

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
