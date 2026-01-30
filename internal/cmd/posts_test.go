/**
 * posts_test.go
 * 投稿管理コマンドのテストコード
 *
 * Phase 1〜4で追加される新規コマンドのテストを含みます。
 */

package cmd

import (
	"testing"
)

// TestPostsDraftsCmd_構造体が存在すること
func TestPostsDraftsCmd_構造体が存在すること(t *testing.T) {
	// PostsDraftsCmdが定義されていることを確認
	_ = &PostsDraftsCmd{}
}

// TestPostsPublishedCmd_構造体が存在すること
func TestPostsPublishedCmd_構造体が存在すること(t *testing.T) {
	// PostsPublishedCmdが定義されていることを確認
	_ = &PostsPublishedCmd{}
}

// TestPostsScheduledCmd_構造体が存在すること
func TestPostsScheduledCmd_構造体が存在すること(t *testing.T) {
	// PostsScheduledCmdが定義されていることを確認
	_ = &PostsScheduledCmd{}
}

// TestPostsURLCmd_構造体が存在すること
func TestPostsURLCmd_構造体が存在すること(t *testing.T) {
	// PostsURLCmdが定義されていることを確認
	_ = &PostsURLCmd{}
}

// TestPostsUnpublishCmd_構造体が存在すること
func TestPostsUnpublishCmd_構造体が存在すること(t *testing.T) {
	// PostsUnpublishCmdが定義されていることを確認
	_ = &PostsUnpublishCmd{}
}

// TestPostsScheduleCmd_構造体が存在すること
func TestPostsScheduleCmd_構造体が存在すること(t *testing.T) {
	// PostsScheduleCmdが定義されていることを確認
	_ = &PostsScheduleCmd{}
}

// TestPostsBatchPublishCmd_構造体が存在すること
func TestPostsBatchPublishCmd_構造体が存在すること(t *testing.T) {
	// PostsBatchPublishCmdが定義されていることを確認
	_ = &PostsBatchPublishCmd{}
}

// TestPostsBatchDeleteCmd_構造体が存在すること
func TestPostsBatchDeleteCmd_構造体が存在すること(t *testing.T) {
	// PostsBatchDeleteCmdが定義されていることを確認
	_ = &PostsBatchDeleteCmd{}
}

// TestPostsSearchCmd_構造体が存在すること
func TestPostsSearchCmd_構造体が存在すること(t *testing.T) {
	// PostsSearchCmdが定義されていることを確認
	_ = &PostsSearchCmd{}
}
