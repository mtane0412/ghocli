# Fields Selection Feature Implementation

## Overview

Implement support for all Ghost Admin API fields and gh CLI-style `--fields` option for field selection.

## Implementation Status (2026-02-01)

### Phase 1: Foundation ✅

#### 1.1 Fields Definition Package Creation

**New Files**:
- `internal/fields/fields.go` - Fields definition foundation
  - `FieldSet`: Field set structure (Default, Detail, All)
  - `Parse(input, fieldSet)`: Parse comma-separated field specification
  - `Validate(fields, available)`: Field validation
  - `ListAvailable(fieldSet)`: Display list of available fields

- `internal/fields/posts.go` - Post field definitions (40 fields)
  - Default: Default fields for list (id, title, status, created_at, published_at)
  - Detail: Default fields for get (13 fields)
  - All: All fields (40 fields)
    - Basic info: id, uuid, title, slug, status, url
    - Content: html, lexical, excerpt, custom_excerpt
    - Images: feature_image, feature_image_alt, feature_image_caption, og_image, twitter_image
    - SEO: meta_title, meta_description, og_*, twitter_*, canonical_url
    - Dates: created_at, updated_at, published_at
    - Control: visibility, featured, email_only
    - Custom: codeinjection_head, codeinjection_foot, custom_template
    - Related: tags, authors, primary_author, primary_tag
    - Other: comment_id, reading_time
    - Email/Newsletter: email_segment, newsletter_id, send_email_when_published

**Tests**:
- `internal/fields/fields_test.go` - Foundation functionality tests
- `internal/fields/posts_test.go` - Post field definition tests

### Phase 2: Structure Extension ✅

#### 2.1 Common Type Definitions

**New Files**:
- `internal/ghostapi/types.go` - Common type definitions
  - `Author`: Author information (ID, Name, Slug, Email, Bio, Location, Website)
  - Tag structure already defined in existing tags.go

**Tests**:
- `internal/ghostapi/types_test.go` - Author/Tag JSON conversion tests

#### 2.2 Post Structure All-Field Support

**Modified Files**:
- `internal/ghostapi/posts.go` - Extend Post structure to 40+ fields
  - Basic info, content, images, SEO, dates, control, custom, related, other, email/newsletter

**Tests**:
- `internal/ghostapi/posts_extended_test.go` - Post extended fields tests

### Phase 3: Command Layer Preparation ✅

#### 3.1 Add --fields to RootFlags

**Modified Files**:
- `internal/cmd/root.go` - Add Fields field to RootFlags
  - `Fields string`: Field specification option
  - `-F` short option
  - `GHO_FIELDS` environment variable support

**Tests**:
- `internal/cmd/root_test.go` - Fields field tests (environment variables, flag priority, etc.)

### Phase 4: Output Layer Extension ✅

#### 4.1 Add Field Filtering to Formatter

**New Files**:
- `internal/outfmt/filter.go` - Field filtering functionality
  - `FilterFields(formatter, data, fields)`: Extract and output specified fields only
  - `filterMap()`: Extract specified fields from map
  - `filterStruct()`: Extract specified fields from struct
  - `StructToMap()`: Convert struct to map[string]interface{}

**Tests**:
- `internal/outfmt/filter_test.go` - Field filtering functionality tests

### Phase 5-6: posts list Command Implementation ✅

#### 5-6.1 --fields Support in posts list

**Modified Files**:
- `internal/cmd/posts.go` - Extend PostsListCmd.Run
  - JSON alone (without --fields): Display list of available fields
  - With field specification: Output only specified fields
  - Works with all JSON/Plain/Table formats

**Tests**:
- `internal/cmd/posts_test.go` - posts list fields support tests

## Usage Examples

```bash
# Display field list
gho posts list --json

# Get only specified fields (JSON)
gho posts list --json --fields id,title,status,excerpt

# Specify fields in Plain format (TSV)
gho posts list --plain --fields id,title,url

# Field specification also possible in table format
gho posts list --fields id,title,status,feature_image

# Specify with environment variable
export GHO_FIELDS="id,title,status"
gho posts list --json

# Short option
gho posts list --json -F id,title,url
```

## Quality Check Results

- ✅ All tests: PASS
- ✅ Type check: PASS
- ✅ Build: SUCCESS
- ✅ --fields option display: Verified

## Remaining Tasks

### Phase 7: Support for posts get Command

**Implementation**:
- Support `--fields` option in `posts get` command
- Display field list when JSON alone
- Filtering output when fields specified

**Modified Files**:
- `internal/cmd/posts.go` - Extend PostsInfoCmd.Run

### Phase 8: Horizontal Expansion to Other Resources

#### 8.1 Add Field Definitions

**New Files**:
- `internal/fields/pages.go` - Page field definitions (same as Post)
- `internal/fields/tags.go` - Tag field definitions
- `internal/fields/members.go` - Member field definitions
- `internal/fields/users.go` - User field definitions
- `internal/fields/newsletters.go` - Newsletter field definitions
- `internal/fields/tiers.go` - Tier field definitions
- `internal/fields/offers.go` - Offer field definitions
- `internal/fields/webhooks.go` - Webhook field definitions

#### 8.2 Extend Resource Structures

**Modified Files**:
- `internal/ghostapi/pages.go` - Page structure extension
- `internal/ghostapi/tags.go` - Tag structure extension
- `internal/ghostapi/members.go` - Member structure extension
- `internal/ghostapi/users.go` - User structure extension
- Similar for other resources

#### 8.3 Apply to Resource Commands

**Modified Files**:
- `internal/cmd/pages.go` - Support for pages list/get
- `internal/cmd/tags.go` - Support for tags list/get
- `internal/cmd/members.go` - Support for members list/get
- `internal/cmd/users.go` - Support for users list/get
- Similar for other resources

### Phase 9: API Layer Extension (Optional)

Use Ghost Admin API's `fields` parameter for server-side filtering.

**Changes**:
- Add `Fields []string` to `ListOptions`
- Add `fields` parameter to API requests

**Modified Files**:
- `internal/ghostapi/posts.go` - ListOptions extension
- `internal/ghostapi/pages.go` - Similar
- Similar for other resources

**Benefits**:
- Reduced network transfer
- Server-side filtering

**Notes**:
- Support may vary by Ghost Admin API version
- Low priority as existing client-side filtering works well

## Implementation Pattern (Reference for Expanding to Other Resources)

### 1. Create Field Definitions

```go
// internal/fields/tags.go
package fields

var TagFields = FieldSet{
    Default: []string{"id", "name", "slug"},
    Detail:  []string{"id", "name", "slug", "description", "visibility"},
    All:     []string{"id", "name", "slug", "description", "visibility", "created_at", "updated_at"},
}
```

### 2. Extend Commands

```go
// internal/cmd/tags.go
func (c *TagsListCmd) Run(ctx context.Context, root *RootFlags) error {
    // Display field list when JSON alone
    if root.JSON && root.Fields == "" {
        formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
        formatter.PrintMessage(fields.ListAvailable(fields.TagFields))
        return nil
    }

    // Parse field specification
    var selectedFields []string
    if root.Fields != "" {
        parsed, err := fields.Parse(root.Fields, fields.TagFields)
        if err != nil {
            return err
        }
        selectedFields = parsed
    }

    // Get data
    response, err := client.ListTags(...)
    if err != nil {
        return err
    }

    // Field filtering
    if len(selectedFields) > 0 {
        var tagsData []map[string]interface{}
        for _, tag := range response.Tags {
            tagMap, _ := outfmt.StructToMap(tag)
            tagsData = append(tagsData, tagMap)
        }
        return outfmt.FilterFields(formatter, tagsData, selectedFields)
    }

    // Default output
    return formatter.Print(response.Tags)
}
```

## Design Decisions

### Separating Field Definitions

**Decision**: Separate field definitions into dedicated package (`internal/fields/`)

**Reasons**:
- Separate field definitions from business logic
- Easier to test
- Referenceable from other packages

### Client-Side Filtering

**Decision**: Filter fields on client side

**Reasons**:
- Independent of server-side API version
- Easy to implement as it only processes already-fetched data
- Network transfer not a major issue (JSON compression works well)

**Future Extension**:
- Can add server-side filtering (`fields` parameter) as needed

### StructToMap Conversion

**Decision**: Implement generic conversion function based on reflection

**Reasons**:
- Reusable across all resources
- Gets accurate field names using JSON tags
- Supports `omitempty`

## Next Work Session Checklist

1. [ ] Review previous implementation (this document)
2. [ ] Phase 7: Support for posts get command
3. [ ] Phase 8: Horizontal expansion to other resources
   - [ ] pages (same fields as Post)
   - [ ] tags
   - [ ] members
   - [ ] users
   - [ ] newsletters
   - [ ] tiers
   - [ ] offers
   - [ ] webhooks
4. [ ] Strictly follow TDD cycle for each phase (RED → GREEN → REFACTOR)
5. [ ] Execute type check and tests for each phase
6. [ ] Update this document after completing implementation
