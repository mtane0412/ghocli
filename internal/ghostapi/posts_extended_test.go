/**
 * posts_extended_test.go
 * Post構造体の拡張フィールドのテスト
 *
 * Ghost Admin APIの全フィールドに対応したPost構造体のテストを提供します。
 */

package ghostapi

import (
	"encoding/json"
	"testing"
)

// TestPost_全フィールドJSON変換 はPost構造体の全フィールドがJSON変換できることを確認します
func TestPost_全フィールドJSON変換(t *testing.T) {
	// テストデータ：全フィールドを含むPost
	post := Post{
		// 基本情報
		ID:     "post123",
		UUID:   "uuid123",
		Title:  "テスト記事",
		Slug:   "test-post",
		Status: "published",
		URL:    "https://example.com/test-post",

		// コンテンツ
		HTML:          "<p>HTMLコンテンツ</p>",
		Lexical:       "{}",
		Excerpt:       "抜粋",
		CustomExcerpt: "カスタム抜粋",

		// 画像
		FeatureImage:        "https://example.com/image.jpg",
		FeatureImageAlt:     "画像の説明",
		FeatureImageCaption: "画像のキャプション",
		OGImage:             "https://example.com/og.jpg",
		TwitterImage:        "https://example.com/twitter.jpg",

		// SEO
		MetaTitle:          "メタタイトル",
		MetaDescription:    "メタ説明",
		OGTitle:            "OGタイトル",
		OGDescription:      "OG説明",
		TwitterTitle:       "Twitterタイトル",
		TwitterDescription: "Twitter説明",
		CanonicalURL:       "https://example.com/canonical",

		// 制御
		Visibility: "public",
		Featured:   true,
		EmailOnly:  false,

		// カスタム
		CodeinjectionHead: "<script>head</script>",
		CodeinjectionFoot: "<script>foot</script>",
		CustomTemplate:    "custom-template",

		// その他
		CommentID:   "comment123",
		ReadingTime: 5,

		// メール・ニュースレター
		EmailSegment:             "all",
		NewsletterID:             "newsletter123",
		SendEmailWhenPublished:   true,
	}

	// JSONに変換
	jsonData, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("JSONマーシャルに失敗: %v", err)
	}

	// JSONから復元
	var restored Post
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("JSONアンマーシャルに失敗: %v", err)
	}

	// 主要フィールドが一致することを確認
	if restored.ID != post.ID {
		t.Errorf("IDが一致しません: got=%s, want=%s", restored.ID, post.ID)
	}
	if restored.Title != post.Title {
		t.Errorf("Titleが一致しません: got=%s, want=%s", restored.Title, post.Title)
	}
	if restored.FeatureImage != post.FeatureImage {
		t.Errorf("FeatureImageが一致しません: got=%s, want=%s", restored.FeatureImage, post.FeatureImage)
	}
	if restored.MetaTitle != post.MetaTitle {
		t.Errorf("MetaTitleが一致しません: got=%s, want=%s", restored.MetaTitle, post.MetaTitle)
	}
	if restored.Visibility != post.Visibility {
		t.Errorf("Visibilityが一致しません: got=%s, want=%s", restored.Visibility, post.Visibility)
	}
	if restored.Featured != post.Featured {
		t.Errorf("Featuredが一致しません: got=%v, want=%v", restored.Featured, post.Featured)
	}
}

// TestPost_関連フィールドJSON変換 はPost構造体の関連フィールド（tags, authors）がJSON変換できることを確認します
func TestPost_関連フィールドJSON変換(t *testing.T) {
	// テストデータ：関連フィールドを含むPost
	post := Post{
		ID:    "post123",
		Title: "テスト記事",
		Tags: []Tag{
			{ID: "tag1", Name: "技術", Slug: "tech"},
			{ID: "tag2", Name: "Go言語", Slug: "golang"},
		},
		Authors: []Author{
			{ID: "author1", Name: "山田太郎", Slug: "yamada"},
		},
		PrimaryAuthor: &Author{
			ID:   "author1",
			Name: "山田太郎",
			Slug: "yamada",
		},
		PrimaryTag: &Tag{
			ID:   "tag1",
			Name: "技術",
			Slug: "tech",
		},
	}

	// JSONに変換
	jsonData, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("JSONマーシャルに失敗: %v", err)
	}

	// JSONから復元
	var restored Post
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("JSONアンマーシャルに失敗: %v", err)
	}

	// Tagsが一致することを確認
	if len(restored.Tags) != len(post.Tags) {
		t.Errorf("Tagsの長さが一致しません: got=%d, want=%d", len(restored.Tags), len(post.Tags))
	}
	if len(restored.Tags) > 0 && restored.Tags[0].Name != post.Tags[0].Name {
		t.Errorf("Tags[0].Nameが一致しません: got=%s, want=%s", restored.Tags[0].Name, post.Tags[0].Name)
	}

	// Authorsが一致することを確認
	if len(restored.Authors) != len(post.Authors) {
		t.Errorf("Authorsの長さが一致しません: got=%d, want=%d", len(restored.Authors), len(post.Authors))
	}
	if len(restored.Authors) > 0 && restored.Authors[0].Name != post.Authors[0].Name {
		t.Errorf("Authors[0].Nameが一致しません: got=%s, want=%s", restored.Authors[0].Name, post.Authors[0].Name)
	}

	// PrimaryAuthorが一致することを確認
	if restored.PrimaryAuthor == nil || restored.PrimaryAuthor.Name != post.PrimaryAuthor.Name {
		t.Errorf("PrimaryAuthor.Nameが一致しません")
	}

	// PrimaryTagが一致することを確認
	if restored.PrimaryTag == nil || restored.PrimaryTag.Name != post.PrimaryTag.Name {
		t.Errorf("PrimaryTag.Nameが一致しません")
	}
}
