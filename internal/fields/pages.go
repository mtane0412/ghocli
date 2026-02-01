/**
 * pages.go
 * Page用フィールド定義
 *
 * Ghost Admin APIのPageリソースはPostと同じスキーマを持つため、
 * PostFieldsのエイリアスとして定義します。
 */

package fields

// PageFields はPage用のフィールドセットです（Postと同一）
var PageFields = PostFields
