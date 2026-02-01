/**
 * tags.go
 * Tag用フィールド定義
 *
 * Ghost Admin APIのTagリソースの全フィールドを定義します。
 */

package fields

// TagFields はTag用のフィールドセットです
var TagFields = FieldSet{
	// Default はlist用デフォルトフィールド（テーブル表示で使用）
	Default: []string{
		"id",
		"name",
		"slug",
		"visibility",
		"created_at",
	},

	// Detail はget用デフォルトフィールド（詳細表示で使用）
	Detail: []string{
		"id",
		"name",
		"slug",
		"description",
		"visibility",
		"created_at",
		"updated_at",
	},

	// All は全フィールド（Ghost Admin APIのTagリソースが持つ全フィールド）
	All: []string{
		"id",
		"name",
		"slug",
		"description",
		"visibility",
		"created_at",
		"updated_at",
	},
}
