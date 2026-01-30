/**
 * users_test.go
 * ユーザー管理コマンドのテストコード
 */

package cmd

import (
	"testing"
)

// TestUsersInfoCmd_構造体が存在すること
func TestUsersInfoCmd_構造体が存在すること(t *testing.T) {
	// UsersInfoCmdが定義されていることを確認
	_ = &UsersInfoCmd{}
}
