/**
 * types_test.go
 * 共通型定義のテスト
 *
 * Author、Tag等の共通型のテストを提供します。
 */

package ghostapi

import (
	"encoding/json"
	"testing"
)

// TestAuthor_JSON変換 はAuthor構造体がJSON変換できることを確認します
func TestAuthor_JSON変換(t *testing.T) {
	// テストデータ
	author := Author{
		ID:    "abc123",
		Name:  "山田太郎",
		Slug:  "yamada-taro",
		Email: "yamada@example.com",
	}

	// JSONに変換
	jsonData, err := json.Marshal(author)
	if err != nil {
		t.Fatalf("JSONマーシャルに失敗: %v", err)
	}

	// JSONから復元
	var restored Author
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("JSONアンマーシャルに失敗: %v", err)
	}

	// 値が一致することを確認
	if restored.ID != author.ID {
		t.Errorf("IDが一致しません: got=%s, want=%s", restored.ID, author.ID)
	}
	if restored.Name != author.Name {
		t.Errorf("Nameが一致しません: got=%s, want=%s", restored.Name, author.Name)
	}
	if restored.Slug != author.Slug {
		t.Errorf("Slugが一致しません: got=%s, want=%s", restored.Slug, author.Slug)
	}
	if restored.Email != author.Email {
		t.Errorf("Emailが一致しません: got=%s, want=%s", restored.Email, author.Email)
	}
}

// TestTag_JSON変換 はTag構造体がJSON変換できることを確認します
func TestTag_JSON変換(t *testing.T) {
	// テストデータ
	tag := Tag{
		ID:          "tag123",
		Name:        "技術",
		Slug:        "tech",
		Description: "技術関連の記事",
	}

	// JSONに変換
	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("JSONマーシャルに失敗: %v", err)
	}

	// JSONから復元
	var restored Tag
	if err := json.Unmarshal(jsonData, &restored); err != nil {
		t.Fatalf("JSONアンマーシャルに失敗: %v", err)
	}

	// 値が一致することを確認
	if restored.ID != tag.ID {
		t.Errorf("IDが一致しません: got=%s, want=%s", restored.ID, tag.ID)
	}
	if restored.Name != tag.Name {
		t.Errorf("Nameが一致しません: got=%s, want=%s", restored.Name, tag.Name)
	}
	if restored.Slug != tag.Slug {
		t.Errorf("Slugが一致しません: got=%s, want=%s", restored.Slug, tag.Slug)
	}
	if restored.Description != tag.Description {
		t.Errorf("Descriptionが一致しません: got=%s, want=%s", restored.Description, tag.Description)
	}
}
