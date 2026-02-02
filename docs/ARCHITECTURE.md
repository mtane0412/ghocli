# gho Architecture Design

## Overview

gho is a CLI tool for the Ghost Admin API. It adopts a simple and maintainable architecture based on the design patterns of gog-cli.

## Project Structure

```
gho/
├── cmd/gho/
│   └── main.go              # Entry point
├── internal/
│   ├── cmd/                  # CLI command definitions
│   │   ├── root.go          # CLI structure, RootFlags
│   │   ├── auth.go          # Authentication commands
│   │   ├── site.go          # Site information commands
│   │   ├── posts.go         # Posts commands (Phase 2)
│   │   ├── pages.go         # Pages commands (Phase 2)
│   │   ├── tags.go          # Tags commands (Phase 3)
│   │   └── images.go        # Images commands (Phase 3)
│   ├── config/              # Configuration file management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── secrets/             # Keyring integration
│   │   ├── store.go
│   │   └── store_test.go
│   ├── ghostapi/            # Ghost API client
│   │   ├── client.go        # HTTP client + JWT generation
│   │   ├── client_test.go
│   │   ├── jwt.go           # JWT generation
│   │   ├── jwt_test.go
│   │   ├── posts.go         # Posts API (Phase 2)
│   │   ├── posts_test.go
│   │   ├── pages.go         # Pages API (Phase 2)
│   │   ├── pages_test.go
│   │   ├── tags.go          # Tags API (Phase 3)
│   │   ├── tags_test.go
│   │   ├── images.go        # Images API (Phase 3)
│   │   └── images_test.go
│   ├── outfmt/              # Output formatting
│   │   ├── outfmt.go
│   │   └── outfmt_test.go
│   └── errfmt/              # Error formatting (to be implemented)
├── docs/                    # Documentation
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── .gitignore
└── README.md
```

## Layer Architecture

```
┌─────────────────────────────────────┐
│          CLI Layer (cmd/)           │  ← User Interface
│  - Command definitions              │
│  - Flag parsing                     │
│  - Input validation                 │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Business Logic Layer           │
│  - config/  : Configuration mgmt    │  ← Business Logic
│  - secrets/ : Credential mgmt       │
│  - ghostapi/: API operations        │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Infrastructure Layer           │
│  - outfmt/  : Output formatting     │  ← Infrastructure
│  - errfmt/  : Error formatting      │
│  - HTTP Client                      │
│  - OS Keyring                       │
└─────────────────────────────────────┘
```

## Component Design

### 1. CLI Layer (`internal/cmd/`)

**Responsibility**: Receive user input and invoke appropriate business logic

**Main Components**:

- **RootFlags**: Common flags for all commands
  ```go
  type RootFlags struct {
      Site    string // Site alias or URL
      JSON    bool   // Output in JSON format
      Plain   bool   // Output in TSV format
      Force   bool   // Skip confirmation
      Verbose bool   // Enable verbose logging
  }
  ```

- **CLI**: CLI structure defined by Kong
  ```go
  type CLI struct {
      RootFlags `embed:""`
      Version   kong.VersionFlag
      Auth      AuthCmd
      Site      SiteCmd
      Posts     PostsCmd
      Pages     PagesCmd
      Tags      TagsCmd
      Images    ImagesCmd
      // ...
  }
  ```

**Design Pattern**: Command Pattern (used internally by Kong)

### 2. Config Layer (`internal/config/`)

**Responsibility**: Read/write configuration files, site management

**Main Features**:

- Configuration file path: `~/.config/gho/config.json`
- Multi-site support (alias functionality)
- Default site management

**Configuration File Format**:
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

**Main Methods**:
- `Load(path string) (*Config, error)` - Load configuration
- `Save(path string) error` - Save configuration
- `AddSite(alias, url string)` - Add site
- `GetSiteURL(aliasOrURL string) (string, bool)` - Get URL

### 3. Secrets Layer (`internal/secrets/`)

**Responsibility**: Secure storage and retrieval of Admin API keys

**Keyring Backends**:
- macOS: Keychain
- Linux: Secret Service (GNOME Keyring, KWallet)
- Windows: Credential Manager
- Fallback: Encrypted file

**Main Methods**:
- `NewStore(backend, fileDir string) (*Store, error)` - Create store
- `Set(alias, apiKey string) error` - Save API key
- `Get(alias string) (string, error)` - Retrieve API key
- `Delete(alias string) error` - Delete API key
- `List() ([]string, error)` - List saved aliases
- `ParseAdminAPIKey(apiKey string) (id, secret string, err error)` - Parse API key

**Security**:
- API keys are stored in OS keyring (not in plain text files)
- Fallback (file backend) is password-protected
- Configuration file stores only URLs (no API keys)

### 4. Ghost API Layer (`internal/ghostapi/`)

**Responsibility**: Communication with Ghost Admin API

**Main Components**:

#### Client
Integrates HTTP client and JWT generation

```go
type Client struct {
    baseURL    string
    keyID      string
    secret     string
    httpClient *http.Client
}
```

**Main Methods**:
- `NewClient(baseURL, keyID, secret string) (*Client, error)`
- `doRequest(method, path string, body io.Reader) ([]byte, error)`
- `GetSite() (*Site, error)`

#### JWT Generation
Ghost Admin API requires JWT signed with HS256

```go
func GenerateJWT(keyID, secret string) (string, error)
```

**JWT Claims**:
```json
{
  "iat": 1234567890,      // Issued at (Unix time)
  "exp": 1234568190,      // Expiration (iat + 5 minutes)
  "aud": "/admin/"        // Ghost Admin API path
}
```

**JWT Header**:
```json
{
  "alg": "HS256",         // Signature algorithm
  "typ": "JWT",
  "kid": "64fac5417..."   // API key ID
}
```

#### API Type Definitions

Define types corresponding to each API resource

```go
type Site struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Version     string `json:"version"`
}

type Post struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    HTML        string     `json:"html,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    PublishedAt *time.Time `json:"published_at,omitempty"`
    Tags        []Tag      `json:"tags,omitempty"`
    Authors     []Author   `json:"authors,omitempty"`
}
```

### 5. Output Format Layer (`internal/outfmt/`)

**Responsibility**: Unified output format management

**Supported Formats**:

| Mode | Flag | Use Case | Format |
|------|------|---------|--------|
| Table | (default) | Human-readable | Column-aligned, with headers |
| JSON | `--json` | Programmatic | JSON format |
| Plain | `--plain` | Pipe processing | TSV format |

**Main Methods**:
- `NewFormatter(writer io.Writer, mode string) *Formatter`
- `Print(data interface{}) error` - Output arbitrary data
- `PrintTable(headers []string, rows [][]string) error` - Table output
- `PrintMessage(message string)` - Message output
- `PrintError(message string)` - Error output

**Output Examples**:

**Table Format**:
```
Alias   URL                           Default
------  ----------------------------  -------
myblog  https://myblog.ghost.io       *
work    https://blog.company.com
```

**JSON Format**:
```json
[
  {
    "Alias": "myblog",
    "URL": "https://myblog.ghost.io",
    "Default": "*"
  },
  {
    "Alias": "work",
    "URL": "https://blog.company.com",
    "Default": ""
  }
]
```

**Plain Format (TSV)**:
```
Alias	URL	Default
myblog	https://myblog.ghost.io	*
work	https://blog.company.com
```

## Authentication Flow

```
1. User creates Custom Integration in Ghost Admin
   ↓
2. Execute `gho auth add https://myblog.ghost.io`
   ↓
3. Enter API key (id:secret format)
   ↓
4. Parse API key (secrets.ParseAdminAPIKey)
   ↓
5. Verify with `/ghost/api/admin/site/`
   - Generate JWT (jwt.GenerateJWT)
   - Execute HTTP request (client.GetSite)
   ↓
6. Save to keyring (secrets.Store.Set)
   ↓
7. Add site to config (config.Config.AddSite)
   ↓
8. Save configuration file (config.Config.Save)
```

## API Request Flow

```
1. User executes command (e.g., gho site)
   ↓
2. Determine site from RootFlags
   - Site specified by -s flag
   - Or default_site from config
   ↓
3. Get URL from config (config.Config.GetSiteURL)
   ↓
4. Get API key from keyring (secrets.Store.Get)
   ↓
5. Parse API key (secrets.ParseAdminAPIKey)
   ↓
6. Create API client (ghostapi.NewClient)
   ↓
7. Generate JWT (ghostapi.GenerateJWT)
   ↓
8. Execute HTTP request (ghostapi.Client.doRequest)
   - Authorization: Ghost <JWT>
   - Accept: application/json
   ↓
9. Parse response
   ↓
10. Display with output format (outfmt.Formatter)
```

## Error Handling

### Error Types

1. **Configuration Errors**
   - Configuration file not found
   - Site not registered
   - Default site not set

2. **Authentication Errors**
   - Invalid API key
   - Invalid API key format
   - Keyring access error

3. **API Errors**
   - HTTP errors (401, 404, 500, etc.)
   - Response parsing errors
   - Network errors

4. **Input Errors**
   - Missing required parameters
   - Invalid parameter format

### Error Message Design

All error messages follow this format:

```
Error: <error description>
```

Examples:
```
Error: Site 'myblog' not found
Error: API key verification failed: Unauthorized
Error: Failed to load configuration: open /Users/user/.config/gho/config.json: no such file or directory
```

## Testing Strategy

### Unit Testing

Each component is independently testable:

- **config**: Configuration file read/write
- **secrets**: Keyring operations (tested with file backend)
- **ghostapi**: HTTP client (mocked with httptest)
- **outfmt**: Output formatting (verified with bytes.Buffer)

### Test Coverage Goals

- Core logic: 80%+
- API layer: 70%+
- CLI layer: 60%+ (covered by manual testing)

## Performance Considerations

### JWT Generation

- Generate JWT for each API request
- Validity period: 5 minutes (Ghost Admin API requirement)
- No caching needed (low generation cost)

### HTTP Connections

- Timeout: 30 seconds
- Keep-Alive: Enabled by default
- Connection reuse for multiple requests

### Keyring Access

- Open keyring on first access
- Reusable across multiple operations
- Password input depends on backend

## Security Considerations

### API Key Storage

- Stored in OS keyring (not in plain text files)
- Configuration file stores only URLs (no API keys)
- File backend is password-protected

### JWT

- Short-lived (5 minutes)
- HS256 signature (Ghost Admin API requirement)
- Includes kid (key ID) in header

### File Permissions

- Configuration file: 0600 (owner read/write only)
- Keyring file: 0600

## Extensibility

### Adding New API Resources

1. Add type definitions to `internal/ghostapi/`
2. Add API functions to `internal/ghostapi/`
3. Add commands to `internal/cmd/`
4. Add tests
5. Update documentation

### Adding New Output Formats

1. Add new formatter to `internal/outfmt/`
2. Add new flag to `RootFlags`
3. Modify `GetOutputMode()` to return new mode

### Adding New Keyring Backends

Backends supported by 99designs/keyring are automatically available
