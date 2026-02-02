/**
 * posts_test.go
 * Post用フィールド定義のテスト
 *
 * Post用のフィールドセット定義のテストを提供します。
 */

package fields

import (
	"testing"
)

// TestPostFields_全フィールド数 はPostフィールドis definedします
func TestPostFields_全フィールド数(t *testing.T) {
	// PostFieldsis defined
	if PostFields.All == nil {
		t.Fatal("PostFields.Allが定義されていません")
	}

	// 全フィールド数が40個以上あることを確認（計画通り）
	if len(PostFields.All) < 40 {
		t.Errorf("PostFieldsの全フィールド数が不足: got=%d, want>=40", len(PostFields.All))
	}
}

// TestPostFields_基本フィールド は基本フィールドが含まれることを確認します
func TestPostFields_基本フィールド(t *testing.T) {
	// 基本フィールド
	basicFields := []string{"id", "uuid", "title", "slug", "status", "url"}

	// 各基本フィールドがAllに含まれることを確認
	for _, field := range basicFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.Allに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_コンテンツフィールド はコンテンツ系フィールドが含まれることを確認します
func TestPostFields_コンテンツフィールド(t *testing.T) {
	// コンテンツフィールド
	contentFields := []string{"html", "lexical", "excerpt", "custom_excerpt"}

	// 各コンテンツフィールドがAllに含まれることを確認
	for _, field := range contentFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.Allに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_画像フィールド は画像系フィールドが含まれることを確認します
func TestPostFields_画像フィールド(t *testing.T) {
	// 画像フィールド
	imageFields := []string{"feature_image", "feature_image_alt", "feature_image_caption", "og_image", "twitter_image"}

	// 各画像フィールドがAllに含まれることを確認
	for _, field := range imageFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.Allに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_SEOフィールド はSEO系フィールドが含まれることを確認します
func TestPostFields_SEOフィールド(t *testing.T) {
	// SEOフィールド
	seoFields := []string{
		"meta_title", "meta_description",
		"og_title", "og_description",
		"twitter_title", "twitter_description",
		"canonical_url",
	}

	// 各SEOフィールドがAllに含まれることを確認
	for _, field := range seoFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.Allに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_日時フィールド は日時系フィールドが含まれることを確認します
func TestPostFields_日時フィールド(t *testing.T) {
	// 日時フィールド
	dateFields := []string{"created_at", "updated_at", "published_at"}

	// 各日時フィールドがAllに含まれることを確認
	for _, field := range dateFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.Allに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_制御フィールド は制御系フィールドが含まれることを確認します
func TestPostFields_制御フィールド(t *testing.T) {
	// 制御フィールド
	controlFields := []string{"visibility", "featured", "email_only"}

	// 各制御フィールドがAllに含まれることを確認
	for _, field := range controlFields {
		if !contains(PostFields.All, field) {
			t.Errorf("PostFields.Allに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_デフォルトフィールド はデフォルトフィールドが適切に設定されていることを確認します
func TestPostFields_デフォルトフィールド(t *testing.T) {
	// Defaultフィールドが設定されていることを確認
	if PostFields.Default == nil {
		t.Fatal("PostFields.Defaultが定義されていません")
	}

	// デフォルトフィールド数が適切（5-10個程度）
	if len(PostFields.Default) < 3 || len(PostFields.Default) > 10 {
		t.Errorf("PostFields.Defaultのフィールド数が不適切: got=%d", len(PostFields.Default))
	}

	// 基本フィールドがDefaultに含まれることを確認
	requiredDefaults := []string{"id", "title", "status"}
	for _, field := range requiredDefaults {
		if !contains(PostFields.Default, field) {
			t.Errorf("PostFields.Defaultに'%s'が含まれていません", field)
		}
	}
}

// TestPostFields_詳細フィールド は詳細フィールドが適切に設定されていることを確認します
func TestPostFields_詳細フィールド(t *testing.T) {
	// Detailフィールドが設定されていることを確認
	if PostFields.Detail == nil {
		t.Fatal("PostFields.Detailが定義されていません")
	}

	// 詳細フィールド数がDefaultより多いことを確認
	if len(PostFields.Detail) <= len(PostFields.Default) {
		t.Errorf("PostFields.DetailがDefaultより多くありません: Detail=%d, Default=%d",
			len(PostFields.Detail), len(PostFields.Default))
	}
}

// ヘルパー関数: スライスに要素が含まれるかチェック
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
