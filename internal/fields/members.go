/**
 * members.go
 * Member用フィールド定義
 *
 * Ghost Admin APIのMemberリソースの全フィールドを定義します。
 */

package fields

// MemberFields はMember用のフィールドセットです
var MemberFields = FieldSet{
	// Default はlist用デフォルトフィールド（テーブル表示で使用）
	Default: []string{
		"id",
		"email",
		"name",
		"status",
		"created_at",
	},

	// Detail はget用デフォルトフィールド（詳細表示で使用）
	Detail: []string{
		"id",
		"uuid",
		"email",
		"name",
		"note",
		"status",
		"labels",
		"created_at",
		"updated_at",
	},

	// All は全フィールド（Ghost Admin APIのMemberリソースが持つ全フィールド）
	All: []string{
		"id",
		"uuid",
		"email",
		"name",
		"note",
		"status",
		"labels",
		"created_at",
		"updated_at",
	},
}
