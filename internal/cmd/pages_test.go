/**
 * pages_test.go
 * ページ管理コマンドのテストコード
 *
 * Phase 1, 2で追加される新規コマンドのテストを含みます。
 */

package cmd

import (
	"testing"
)

// TestPagesInfoCmd_構造体が存在すること
func TestPagesInfoCmd_構造体が存在すること(t *testing.T) {
	// PagesInfoCmdが定義されていることを確認
	_ = &PagesInfoCmd{}
}

// TestPagesURLCmd_構造体が存在すること
func TestPagesURLCmd_構造体が存在すること(t *testing.T) {
	// PagesURLCmdが定義されていることを確認
	_ = &PagesURLCmd{}
}

// TestPagesPublishCmd_構造体が存在すること
func TestPagesPublishCmd_構造体が存在すること(t *testing.T) {
	// PagesPublishCmdが定義されていることを確認
	_ = &PagesPublishCmd{}
}

// TestPagesUnpublishCmd_構造体が存在すること
func TestPagesUnpublishCmd_構造体が存在すること(t *testing.T) {
	// PagesUnpublishCmdが定義されていることを確認
	_ = &PagesUnpublishCmd{}
}

// TestPagesCatCmd_構造体が存在すること
func TestPagesCatCmd_構造体が存在すること(t *testing.T) {
	// PagesCatCmdが定義されていることを確認
	_ = &PagesCatCmd{}
}
