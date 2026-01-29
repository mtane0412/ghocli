/**
 * jwt.go
 * Ghost Admin API用のJWT生成
 *
 * Ghost Admin APIはHS256アルゴリズムで署名されたJWTを要求します。
 * トークンの有効期限は5分です。
 */

package ghostapi

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT はGhost Admin API用のJWTトークンを生成します。
// keyID: Admin APIキーのID部分
// secret: Admin APIキーのシークレット部分
func GenerateJWT(keyID, secret string) (string, error) {
	if keyID == "" {
		return "", errors.New("キーIDが空です")
	}
	if secret == "" {
		return "", errors.New("シークレットが空です")
	}

	// 現在時刻（秒単位）
	now := time.Now().Unix()

	// JWTクレームを設定
	claims := jwt.MapClaims{
		"iat": now,           // 発行時刻
		"exp": now + 5*60,    // 有効期限（5分後）
		"aud": "/admin/",     // Ghost Admin APIのパス
	}

	// トークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// ヘッダーにキーIDを設定
	token.Header["kid"] = keyID

	// シークレットで署名
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
