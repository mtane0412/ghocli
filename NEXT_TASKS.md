# gho UX Improvements - Remaining Tasks (Phase 4 & 5)

## Overview

Continuation of the project to apply gogcli UX improvement features to gho. Phases 1-3 are complete, with Phases 4-5 remaining.

## Completed

- ✅ **Phase 1**: Command aliases (`gho p list`, `gho t list`, etc.)
- ✅ **Phase 2**: Shell completion (bash/zsh/fish/powershell support)
- ✅ **Phase 3**: Custom help (coloring, terminal width adjustment)

## Remaining Tasks

### Phase 4: Error Message Improvements

**Purpose**: Provide user-friendly error messages and suggest solutions

**Implementation**:

1. **Authentication Error Improvements**
   - Current: `No API key configured`
   - Improved:
     ```
     No API key configured for site "myblog".

     Add credentials:
       gho auth add myblog https://myblog.com
     ```

2. **Site Unspecified Error Improvements**
   - Current: `No site specified`
   - Improved:
     ```
     No site specified.

     Specify with --site flag or set default:
       gho config set default_site myblog
     ```

3. **Unknown Flag Error Improvements**
   - Current: `unknown flag --foo`
   - Improved:
     ```
     unknown flag --foo
     Run with --help to see available flags
     ```

**Modified Files**: `internal/errfmt/errfmt.go`

**Implementation Steps**:

1. **Follow TDD principles**: Write tests before implementation
2. Edit `internal/errfmt/errfmt.go`
3. Add/modify the following functions:
   - `FormatAuthError(site string) string` - Format authentication errors
   - `FormatSiteError() string` - Format site unspecified errors
   - `FormatFlagError(flag string) string` - Format unknown flag errors

**Reference Files**:
- `/Users/mtane0412/dev/gogcli/internal/errfmt/errfmt.go`

**Example Test Cases**:
```go
func TestFormatAuthError(t *testing.T) {
    msg := FormatAuthError("myblog")
    assert.Contains(t, msg, "No API key configured")
    assert.Contains(t, msg, "gho auth add")
}
```

---

### Phase 5: Flag Aliases

**Purpose**: Provide shorthand for commonly used flags

**Implementation**:

Add aliases to major flags:

| Flag | Aliases | Target Commands |
|------|---------|----------------|
| `--limit` | `--max`, `-n` | list commands |
| `--filter` | `--where`, `-w` | list commands |
| `--output` | `--format`, `-o` | All commands (future) |

**Modified Files**:
- `internal/cmd/posts.go`
- `internal/cmd/pages.go`
- `internal/cmd/tags.go`
- `internal/cmd/members.go`
- `internal/cmd/users.go`
- `internal/cmd/newsletters.go`
- `internal/cmd/tiers.go`
- `internal/cmd/offers.go`

**Implementation Steps**:

1. **Follow TDD principles**: Write tests before implementation
2. Edit each command's `ListCmd` structure
3. Example (`internal/cmd/posts.go`):
   ```go
   // Before
   Limit  int    `name:"limit" help:"Maximum number of posts to return"`
   Filter string `name:"filter" help:"Filter posts"`

   // After
   Limit  int    `name:"limit" aliases:"max,n" help:"Maximum number of posts to return"`
   Filter string `name:"filter" aliases:"where,w" help:"Filter posts"`
   ```

**Example Test Cases**:
```go
func TestPostsListCmd_LimitAliases(t *testing.T) {
    testCases := []struct {
        name string
        args []string
    }{
        {"--limit", []string{"posts", "list", "--limit=10"}},
        {"--max", []string{"posts", "list", "--max=10"}},
        {"-n", []string{"posts", "list", "-n=10"}},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            var cli CLI
            parser, err := kong.New(&cli)
            require.NoError(t, err)

            _, err = parser.Parse(tc.args)
            require.NoError(t, err)
        })
    }
}
```

---

## Implementation Guidelines

### TDD Principles (Strictly Applied)

1. **RED**: Write a failing test first
2. **GREEN**: Write minimal code to make the test pass
3. **REFACTOR**: Clean up the code

### Quality Checks

After implementation, always execute:

```bash
# Run tests
make test

# Type check
make type-check

# Lint
make lint

# Build
make build
```

### Commit Messages

Create commit for each Phase:

```bash
# Phase 4
git add <modified files>
git commit -m "feat: improve error messages (Phase 4)

<description of changes>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

# Phase 5
git add <modified files>
git commit -m "feat: add flag aliases (Phase 5)

<description of changes>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### Branch

Continue work on existing feature branch:

```bash
# Check current branch
git branch --show-current
# => feature/gogcli-ux-improvements

# Continue work on this branch
```

---

## Reference Information

### Existing Implementation (gogcli)

gogcli files to reference:
- `/Users/mtane0412/dev/gogcli/internal/errfmt/errfmt.go` - Error message formatting
- `/Users/mtane0412/dev/gogcli/internal/cmd/root.go` - Flag alias examples

### gho Current Status

- Command aliases: ✅ Implemented (`posts` → `post`, `p`)
- Shell completion: ✅ Implemented (bash/zsh/fish/powershell)
- Custom help: ✅ Implemented (coloring, terminal width adjustment)
- Error message improvements: ❌ Not implemented
- Flag aliases: ❌ Not implemented

### Priority

1. **Phase 5 (Flag Aliases)** - Simpler with smaller scope
2. **Phase 4 (Error Messages)** - Requires error handling design

Starting with Phase 5 is recommended.

---

## Success Criteria

- [ ] All Phase 4 tests pass
- [ ] All Phase 5 tests pass
- [ ] All existing tests pass (no regressions)
- [ ] `make lint` completes without errors
- [ ] `make type-check` completes without errors
- [ ] Can actually build and verify operation
- [ ] Commit messages are clear and changes understandable

---

## Questions & Unclear Points

For unclear points, refer to:

1. **TDD**: `@rules/tdd.md`
2. **Quality Checks**: `@rules/quality-checks.md`
3. **Git Workflow**: `@rules/git-workflow.md`
4. **Test Code**: `@rules/testing.md`

Or ask the user.
