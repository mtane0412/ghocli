/**
 * jwt_test.go
 * Ghost Admin API用のJWT生成機能のテストコード
 */

package ghostapi

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestGenerateJWT_正しいフォーマットのトークン生成
func TestGenerateJWT_正しいフォーマットのトークン生成(t *testing.T) {
	// テスト用のAPIキー情報
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"

	// JWTを生成
	token, err := GenerateJWT(keyID, secret)
	if err != nil {
		t.Fatalf("JWTの生成に失敗: %v", err)
	}

	// トークンが空でないことを確認
	if token == "" {
		t.Error("生成されたトークンが空です")
	}

	// トークンが3つのパートから構成されていることを確認（header.payload.signature）
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("トークンのパート数 = %d; want 3", len(parts))
	}
}

// TestGenerateJWT_トークンの検証
func TestGenerateJWT_トークンの検証(t *testing.T) {
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"

	// JWTを生成
	tokenString, err := GenerateJWT(keyID, secret)
	if err != nil {
		t.Fatalf("JWTの生成に失敗: %v", err)
	}

	// トークンをパースして検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムがHS256であることを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("予期しない署名アルゴリズム: %v", token.Header["alg"])
		}
		// シークレットを16進数からバイナリにデコード
		secretBytes, err := hex.DecodeString(secret)
		if err != nil {
			t.Fatalf("シークレットのデコードに失敗: %v", err)
		}
		return secretBytes, nil
	})

	if err != nil {
		t.Fatalf("トークンのパースに失敗: %v", err)
	}

	if !token.Valid {
		t.Error("トークンが無効です")
	}

	// クレームの検証
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("クレームの取得に失敗")
	}

	// audクレームがGhost Admin APIのパスであることを確認
	aud, ok := claims["aud"].(string)
	if !ok || aud != "/admin/" {
		t.Errorf("aud = %q; want %q", aud, "/admin/")
	}

	// iatクレームが現在時刻付近であることを確認
	iat, ok := claims["iat"].(float64)
	if !ok {
		t.Fatal("iatクレームの取得に失敗")
	}
	iatTime := time.Unix(int64(iat), 0)
	if time.Since(iatTime) > 10*time.Second {
		t.Errorf("iatが古すぎます: %v", iatTime)
	}

	// expクレームがiat + 5分であることを確認
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("expクレームの取得に失敗")
	}
	expectedExp := int64(iat) + 5*60
	if int64(exp) != expectedExp {
		t.Errorf("exp = %d; want %d", int64(exp), expectedExp)
	}
}

// TestGenerateJWT_ヘッダーにkidが含まれる
func TestGenerateJWT_ヘッダーにkidが含まれる(t *testing.T) {
	keyID := "64fac5417c4c6b0001234567"
	secret := "89abcdef01234567890123456789abcd01234567890123456789abcdef0123"

	// JWTを生成
	tokenString, err := GenerateJWT(keyID, secret)
	if err != nil {
		t.Fatalf("JWTの生成に失敗: %v", err)
	}

	// トークンをパース（検証なし）
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("トークンのパースに失敗: %v", err)
	}

	// ヘッダーのkidを確認
	kid, ok := token.Header["kid"].(string)
	if !ok {
		t.Fatal("ヘッダーにkidが含まれていない")
	}
	if kid != keyID {
		t.Errorf("kid = %q; want %q", kid, keyID)
	}
}

// TestGenerateJWT_空のキーIDでエラー
func TestGenerateJWT_空のキーIDでエラー(t *testing.T) {
	_, err := GenerateJWT("", "secret")
	if err == nil {
		t.Error("空のキーIDでエラーが返されなかった")
	}
}

// TestGenerateJWT_空のシークレットでエラー
func TestGenerateJWT_空のシークレットでエラー(t *testing.T) {
	_, err := GenerateJWT("keyid", "")
	if err == nil {
		t.Error("空のシークレットでエラーが返されなかった")
	}
}
