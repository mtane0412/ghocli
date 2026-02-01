/**
 * posts.go
 * Post用フィールド定義
 *
 * Ghost Admin APIのPost/Pageリソースの全フィールドを定義します。
 */

package fields

// PostFields はPost/Page用のフィールドセットです
var PostFields = FieldSet{
	// Default はlist用デフォルトフィールド（テーブル表示で使用）
	Default: []string{
		"id",
		"title",
		"status",
		"created_at",
		"published_at",
	},

	// Detail はget用デフォルトフィールド（詳細表示で使用）
	Detail: []string{
		"id",
		"uuid",
		"title",
		"slug",
		"status",
		"url",
		"excerpt",
		"feature_image",
		"created_at",
		"updated_at",
		"published_at",
		"visibility",
		"featured",
	},

	// All は全フィールド（Ghost Admin APIのPostリソースが持つ全フィールド）
	All: []string{
		// 基本情報
		"id",
		"uuid",
		"title",
		"slug",
		"status",
		"url",

		// コンテンツ
		"html",
		"lexical",
		"excerpt",
		"custom_excerpt",

		// 画像
		"feature_image",
		"feature_image_alt",
		"feature_image_caption",
		"og_image",
		"twitter_image",

		// SEO
		"meta_title",
		"meta_description",
		"og_title",
		"og_description",
		"twitter_title",
		"twitter_description",
		"canonical_url",

		// 日時
		"created_at",
		"updated_at",
		"published_at",

		// 制御
		"visibility",
		"featured",
		"email_only",

		// カスタム
		"codeinjection_head",
		"codeinjection_foot",
		"custom_template",

		// 関連
		"tags",
		"authors",
		"primary_author",
		"primary_tag",

		// その他
		"comment_id",
		"reading_time",

		// メール・ニュースレター
		"email_segment",
		"newsletter_id",
		"send_email_when_published",
	},
}
