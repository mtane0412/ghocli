/**
 * users.go
 * User用フィールド定義
 *
 * Ghost Admin APIのUserリソースの全フィールドを定義します。
 */

package fields

// UserFields はUser用のフィールドセットです
var UserFields = FieldSet{
	// Default はlist用デフォルトフィールド（テーブル表示で使用）
	Default: []string{
		"id",
		"name",
		"slug",
		"email",
		"created_at",
	},

	// Detail はget用デフォルトフィールド（詳細表示で使用）
	Detail: []string{
		"id",
		"name",
		"slug",
		"email",
		"bio",
		"location",
		"website",
		"profile_image",
		"cover_image",
		"roles",
		"created_at",
		"updated_at",
	},

	// All は全フィールド（Ghost Admin APIのUserリソースが持つ全フィールド）
	All: []string{
		"id",
		"name",
		"slug",
		"email",
		"bio",
		"location",
		"website",
		"profile_image",
		"cover_image",
		"roles",
		"created_at",
		"updated_at",
	},
}
