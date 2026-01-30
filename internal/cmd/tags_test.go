/**
 * tags_test.go
 * タグ管理コマンドのテストコード
 */

package cmd

import (
	"testing"
)

// TestTagsInfoCmd_構造体が存在すること
func TestTagsInfoCmd_構造体が存在すること(t *testing.T) {
	// TagsInfoCmdが定義されていることを確認
	_ = &TagsInfoCmd{}
}
