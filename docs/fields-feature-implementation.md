# フィールド選択機能の実装

## 概要

Ghost Admin APIの全フィールドに対応し、gh CLI風の`--fields`オプションによるフィールド選択機能を実装。

## 実装完了状況（2026-02-01）

### Phase 1: 基盤整備 ✅

#### 1.1 フィールド定義パッケージの作成

**新規ファイル**:
- `internal/fields/fields.go` - フィールド定義基盤
  - `FieldSet`: フィールドセット構造体（Default, Detail, All）
  - `Parse(input, fieldSet)`: カンマ区切りフィールド指定のパース
  - `Validate(fields, available)`: フィールドのバリデーション
  - `ListAvailable(fieldSet)`: 利用可能なフィールド一覧表示

- `internal/fields/posts.go` - Post用フィールド定義（40フィールド）
  - Default: list用デフォルトフィールド（id, title, status, created_at, published_at）
  - Detail: get用デフォルトフィールド（13フィールド）
  - All: 全フィールド（40フィールド）
    - 基本情報: id, uuid, title, slug, status, url
    - コンテンツ: html, lexical, excerpt, custom_excerpt
    - 画像: feature_image, feature_image_alt, feature_image_caption, og_image, twitter_image
    - SEO: meta_title, meta_description, og_*, twitter_*, canonical_url
    - 日時: created_at, updated_at, published_at
    - 制御: visibility, featured, email_only
    - カスタム: codeinjection_head, codeinjection_foot, custom_template
    - 関連: tags, authors, primary_author, primary_tag
    - その他: comment_id, reading_time
    - メール・ニュースレター: email_segment, newsletter_id, send_email_when_published

**テスト**:
- `internal/fields/fields_test.go` - 基盤機能のテスト
- `internal/fields/posts_test.go` - Postフィールド定義のテスト

### Phase 2: 構造体の拡張 ✅

#### 2.1 共通型定義

**新規ファイル**:
- `internal/ghostapi/types.go` - 共通型定義
  - `Author`: 著者情報（ID, Name, Slug, Email, Bio, Location, Website）
  - Tag構造体は既存のtags.goで定義済み

**テスト**:
- `internal/ghostapi/types_test.go` - Author/TagのJSON変換テスト

#### 2.2 Post構造体の全フィールド対応

**変更ファイル**:
- `internal/ghostapi/posts.go` - Post構造体を40フィールド以上に拡張
  - 基本情報、コンテンツ、画像、SEO、日時、制御、カスタム、関連、その他、メール・ニュースレター

**テスト**:
- `internal/ghostapi/posts_extended_test.go` - Post拡張フィールドのテスト

### Phase 3: コマンド層の準備 ✅

#### 3.1 RootFlagsへの--fields追加

**変更ファイル**:
- `internal/cmd/root.go` - RootFlagsにFieldsフィールドを追加
  - `Fields string`: フィールド指定オプション
  - `-F` ショートオプション
  - `GHO_FIELDS` 環境変数サポート

**テスト**:
- `internal/cmd/root_test.go` - Fieldsフィールドのテスト（環境変数、フラグ優先度等）

### Phase 4: 出力層の拡張 ✅

#### 4.1 Formatterへのフィールドフィルタリング機能追加

**新規ファイル**:
- `internal/outfmt/filter.go` - フィールドフィルタリング機能
  - `FilterFields(formatter, data, fields)`: 指定フィールドのみを抽出して出力
  - `filterMap()`: マップから指定フィールドを抽出
  - `filterStruct()`: 構造体から指定フィールドを抽出
  - `StructToMap()`: 構造体をmap[string]interface{}に変換

**テスト**:
- `internal/outfmt/filter_test.go` - フィールドフィルタリング機能のテスト

### Phase 5-6: posts listコマンド実装 ✅

#### 5-6.1 posts listでの--fields対応

**変更ファイル**:
- `internal/cmd/posts.go` - PostsListCmd.Runを拡張
  - JSON単独（--fieldsなし）時：利用可能なフィールド一覧を表示
  - フィールド指定時：指定フィールドのみを出力
  - JSON/Plain/Table形式すべてで動作

**テスト**:
- `internal/cmd/posts_test.go` - posts listのfields対応テスト

## 使用例

```bash
# フィールド一覧を表示
gho posts list --json

# 指定フィールドのみ取得（JSON）
gho posts list --json --fields id,title,status,excerpt

# Plain形式（TSV）でフィールド指定
gho posts list --plain --fields id,title,url

# テーブル形式でもフィールド指定可能
gho posts list --fields id,title,status,feature_image

# 環境変数で指定
export GHO_FIELDS="id,title,status"
gho posts list --json

# ショートオプション
gho posts list --json -F id,title,url
```

## 品質チェック結果

- ✅ 全テスト: PASS
- ✅ 型チェック: PASS
- ✅ ビルド: 成功
- ✅ --fieldsオプション表示: 確認済み

## 残りのタスク

### Phase 7: posts getコマンドへの対応

**実装内容**:
- `posts get` コマンドで`--fields`オプションをサポート
- JSON単独時にフィールド一覧表示
- フィールド指定時にフィルタリング出力

**変更ファイル**:
- `internal/cmd/posts.go` - PostsInfoCmd.Runを拡張

### Phase 8: 他リソースへの横展開

#### 8.1 フィールド定義の追加

**新規ファイル**:
- `internal/fields/pages.go` - Page用フィールド定義（Postと同一）
- `internal/fields/tags.go` - Tag用フィールド定義
- `internal/fields/members.go` - Member用フィールド定義
- `internal/fields/users.go` - User用フィールド定義
- `internal/fields/newsletters.go` - Newsletter用フィールド定義
- `internal/fields/tiers.go` - Tier用フィールド定義
- `internal/fields/offers.go` - Offer用フィールド定義
- `internal/fields/webhooks.go` - Webhook用フィールド定義

#### 8.2 各リソース構造体の拡張

**変更ファイル**:
- `internal/ghostapi/pages.go` - Page構造体の拡張
- `internal/ghostapi/tags.go` - Tag構造体の拡張
- `internal/ghostapi/members.go` - Member構造体の拡張
- `internal/ghostapi/users.go` - User構造体の拡張
- 他のリソースも同様

#### 8.3 各リソースコマンドへの適用

**変更ファイル**:
- `internal/cmd/pages.go` - pages list/getへの対応
- `internal/cmd/tags.go` - tags list/getへの対応
- `internal/cmd/members.go` - members list/getへの対応
- `internal/cmd/users.go` - users list/getへの対応
- 他のリソースも同様

### Phase 9: API層の拡張（オプション）

Ghost Admin APIの`fields`パラメータを使用してサーバーサイドでフィルタリング。

**変更内容**:
- `ListOptions`に`Fields []string`を追加
- APIリクエスト時に`fields`パラメータを付与

**変更ファイル**:
- `internal/ghostapi/posts.go` - ListOptions拡張
- `internal/ghostapi/pages.go` - 同様
- 他のリソースも同様

**メリット**:
- ネットワーク転送量の削減
- サーバー側でのフィルタリング

**注意点**:
- Ghost Admin APIのバージョンによってサポート状況が異なる可能性
- 既存のクライアント側フィルタリングで十分に機能するため、優先度は低い

## 実装パターン（他リソース展開時の参考）

### 1. フィールド定義の作成

```go
// internal/fields/tags.go
package fields

var TagFields = FieldSet{
    Default: []string{"id", "name", "slug"},
    Detail:  []string{"id", "name", "slug", "description", "visibility"},
    All:     []string{"id", "name", "slug", "description", "visibility", "created_at", "updated_at"},
}
```

### 2. コマンドの拡張

```go
// internal/cmd/tags.go
func (c *TagsListCmd) Run(ctx context.Context, root *RootFlags) error {
    // JSON単独時はフィールド一覧表示
    if root.JSON && root.Fields == "" {
        formatter := outfmt.NewFormatter(os.Stdout, root.GetOutputMode())
        formatter.PrintMessage(fields.ListAvailable(fields.TagFields))
        return nil
    }

    // フィールド指定をパース
    var selectedFields []string
    if root.Fields != "" {
        parsed, err := fields.Parse(root.Fields, fields.TagFields)
        if err != nil {
            return err
        }
        selectedFields = parsed
    }

    // データ取得
    response, err := client.ListTags(...)
    if err != nil {
        return err
    }

    // フィールドフィルタリング
    if len(selectedFields) > 0 {
        var tagsData []map[string]interface{}
        for _, tag := range response.Tags {
            tagMap, _ := outfmt.StructToMap(tag)
            tagsData = append(tagsData, tagMap)
        }
        return outfmt.FilterFields(formatter, tagsData, selectedFields)
    }

    // デフォルト出力
    return formatter.Print(response.Tags)
}
```

## 設計判断

### フィールド定義の分離

**判断**: フィールド定義を専用パッケージ（`internal/fields/`）に分離

**理由**:
- フィールド定義とビジネスロジックを分離
- テストが容易
- 他のパッケージからも参照可能

### クライアント側フィルタリング

**判断**: クライアント側でフィールドをフィルタリング

**理由**:
- サーバー側のAPIバージョンに依存しない
- すでに取得したデータを加工するだけなので実装が容易
- ネットワーク転送量は大きな問題ではない（JSON圧縮が効く）

**将来の拡張**:
- 必要に応じてサーバー側フィルタリング（`fields`パラメータ）も追加可能

### StructToMap変換

**判断**: reflectionベースの汎用的な変換関数を実装

**理由**:
- すべてのリソースで再利用可能
- JSONタグを使用して正確なフィールド名を取得
- `omitempty`対応

## 次回作業時のチェックリスト

1. [ ] 前回の実装内容を確認（このドキュメント）
2. [ ] Phase 7: posts getコマンドへの対応
3. [ ] Phase 8: 他リソースへの横展開
   - [ ] pages（Postと同じフィールド）
   - [ ] tags
   - [ ] members
   - [ ] users
   - [ ] newsletters
   - [ ] tiers
   - [ ] offers
   - [ ] webhooks
4. [ ] 各フェーズでTDDサイクルを厳守（RED → GREEN → REFACTOR）
5. [ ] 各フェーズで型チェックとテストを実行
6. [ ] 実装完了後、このドキュメントを更新
