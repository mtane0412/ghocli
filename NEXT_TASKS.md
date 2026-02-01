# gho UX改善 - 残りタスク（Phase 4 & 5）

## 概要

gogcliのUX改善機能をghoに適用するプロジェクトの続きです。Phase 1-3は完了しており、Phase 4-5が残っています。

## 完了済み

- ✅ **Phase 1**: コマンドエイリアス（`gho p list`、`gho t list`など）
- ✅ **Phase 2**: シェル補完（bash/zsh/fish/powershell対応）
- ✅ **Phase 3**: カスタムヘルプ（色付け、端末幅調整）

## 残りタスク

### Phase 4: エラーメッセージの改善

**目的**: ユーザーフレンドリーなエラーメッセージを提供し、解決方法を提示する

**実装内容**:

1. **認証エラーの改善**
   - 現状: `No API key configured`
   - 改善後:
     ```
     No API key configured for site "myblog".

     Add credentials:
       gho auth add myblog https://myblog.com
     ```

2. **サイト未指定エラーの改善**
   - 現状: `No site specified`
   - 改善後:
     ```
     No site specified.

     Specify with --site flag or set default:
       gho config set default_site myblog
     ```

3. **不明なフラグエラーの改善**
   - 現状: `unknown flag --foo`
   - 改善後:
     ```
     unknown flag --foo
     Run with --help to see available flags
     ```

**変更ファイル**: `internal/errfmt/errfmt.go`

**実装手順**:

1. **TDD原則に従う**: まずテストを書いてから実装
2. `internal/errfmt/errfmt.go`を編集
3. 以下の関数を追加/修正:
   - `FormatAuthError(site string) string` - 認証エラーのフォーマット
   - `FormatSiteError() string` - サイト未指定エラーのフォーマット
   - `FormatFlagError(flag string) string` - 不明なフラグエラーのフォーマット

**参考ファイル**:
- `/Users/mtane0412/dev/gogcli/internal/errfmt/errfmt.go`

**テストケース例**:
```go
func TestFormatAuthError(t *testing.T) {
    msg := FormatAuthError("myblog")
    assert.Contains(t, msg, "No API key configured")
    assert.Contains(t, msg, "gho auth add")
}
```

---

### Phase 5: フラグエイリアスの追加

**目的**: よく使うフラグに短縮形を提供する

**実装内容**:

主要なフラグにエイリアスを追加します：

| フラグ | エイリアス | 対象コマンド |
|--------|-----------|-------------|
| `--limit` | `--max`, `-n` | list系コマンド |
| `--filter` | `--where`, `-w` | list系コマンド |
| `--output` | `--format`, `-o` | 全コマンド（将来用） |

**変更ファイル**:
- `internal/cmd/posts.go`
- `internal/cmd/pages.go`
- `internal/cmd/tags.go`
- `internal/cmd/members.go`
- `internal/cmd/users.go`
- `internal/cmd/newsletters.go`
- `internal/cmd/tiers.go`
- `internal/cmd/offers.go`

**実装手順**:

1. **TDD原則に従う**: まずテストを書いてから実装
2. 各コマンドの`ListCmd`構造体を編集
3. 例（`internal/cmd/posts.go`）:
   ```go
   // Before
   Limit  int    `name:"limit" help:"Maximum number of posts to return"`
   Filter string `name:"filter" help:"Filter posts"`

   // After
   Limit  int    `name:"limit" aliases:"max,n" help:"Maximum number of posts to return"`
   Filter string `name:"filter" aliases:"where,w" help:"Filter posts"`
   ```

**テストケース例**:
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

## 実装ガイドライン

### TDD原則（厳格に適用）

1. **RED**: 失敗するテストを先に書く
2. **GREEN**: テストを通す最小限のコードを書く
3. **REFACTOR**: コードを整理する

### 品質チェック

実装後、必ず以下を実行してください：

```bash
# テスト実行
make test

# 型チェック
make type-check

# Lint
make lint

# ビルド
make build
```

### コミットメッセージ

各Phaseごとにコミットを作成してください：

```bash
# Phase 4
git add <変更ファイル>
git commit -m "feat: エラーメッセージの改善（Phase 4）

<変更内容の説明>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

# Phase 5
git add <変更ファイル>
git commit -m "feat: フラグエイリアスの追加（Phase 5）

<変更内容の説明>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

### ブランチ

作業は既存のfeatureブランチで続けてください：

```bash
# 現在のブランチを確認
git branch --show-current
# => feature/gogcli-ux-improvements

# このブランチで作業を継続
```

---

## 参考情報

### 既存の実装（gogcli）

参考にするべきgogcliのファイル：
- `/Users/mtane0412/dev/gogcli/internal/errfmt/errfmt.go` - エラーメッセージフォーマット
- `/Users/mtane0412/dev/gogcli/internal/cmd/root.go` - フラグエイリアスの例

### ghoの現状

- コマンドエイリアス: ✅ 実装済み（`posts` → `post`, `p`）
- シェル補完: ✅ 実装済み（bash/zsh/fish/powershell）
- カスタムヘルプ: ✅ 実装済み（色付け、端末幅調整）
- エラーメッセージ改善: ❌ 未実装
- フラグエイリアス: ❌ 未実装

### 優先度

1. **Phase 5（フラグエイリアス）** - より簡単で影響範囲が小さい
2. **Phase 4（エラーメッセージ）** - エラーハンドリングの設計が必要

Phase 5から始めることを推奨します。

---

## 成功基準

- [ ] Phase 4のテストがすべて成功する
- [ ] Phase 5のテストがすべて成功する
- [ ] 既存のテストがすべて成功する（回帰テストなし）
- [ ] `make lint`がエラーなく完了する
- [ ] `make type-check`がエラーなく完了する
- [ ] 実際にビルドして動作確認ができる
- [ ] コミットメッセージが明確で、変更内容が理解できる

---

## 質問・不明点

不明な点があれば、以下を参照してください：

1. **TDD**: `@rules/tdd.md`
2. **品質チェック**: `@rules/quality-checks.md`
3. **Gitワークフロー**: `@rules/git-workflow.md`
4. **テストコード**: `@rules/testing.md`

または、ユーザーに質問してください。
