/**
 * cmd_test.go
 * コマンド共通機能のテストコード
 *
 * Phase 5で追加されるフラグエイリアスのテストを含みます。
 */

package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
)

// TestLimitAliases はLimitフラグのエイリアス（--max, -n）が正しく機能することを確認します
func TestLimitAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want int
	}{
		{"--limit フラグ", []string{"posts", "list", "--limit=10"}, 10},
		{"--max エイリアス", []string{"posts", "list", "--max=10"}, 10},
		{"--n エイリアス", []string{"posts", "list", "--n=10"}, 10},
		{"-l ショートフラグ", []string{"posts", "list", "-l", "10"}, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// CLIを初期化
			cli := &CLI{}

			// Kongパーサーを作成
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // テスト時は終了しない
			)
			if err != nil {
				t.Fatalf("Kongパーサーの作成に失敗: %v", err)
			}

			// コマンドラインをパース
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("コマンドラインのパースに失敗: %v", err)
			}

			// Limitフィールドが正しく設定されているか確認
			if cli.Posts.List.Limit != tc.want {
				t.Errorf("Limitフィールドが正しく設定されていません: got=%d, want=%d", cli.Posts.List.Limit, tc.want)
			}
		})
	}
}

// TestFilterAliases はFilterフラグのエイリアス（--where, -w）が正しく機能することを確認します
func TestFilterAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{"--filter フラグ", []string{"members", "list", "--filter=status:paid"}, "status:paid"},
		{"--where エイリアス", []string{"members", "list", "--where=status:paid"}, "status:paid"},
		{"--w エイリアス", []string{"members", "list", "--w=status:paid"}, "status:paid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// CLIを初期化
			cli := &CLI{}

			// Kongパーサーを作成
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // テスト時は終了しない
			)
			if err != nil {
				t.Fatalf("Kongパーサーの作成に失敗: %v", err)
			}

			// コマンドラインをパース
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("コマンドラインのパースに失敗: %v", err)
			}

			// Filterフィールドが正しく設定されているか確認
			if cli.Members.List.Filter != tc.want {
				t.Errorf("Filterフィールドが正しく設定されていません: got=%s, want=%s", cli.Members.List.Filter, tc.want)
			}
		})
	}
}

// TestPagesFlagAliases はPagesコマンドのLimitエイリアスが正しく機能することを確認します
func TestPagesFlagAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want int
	}{
		{"--limit フラグ", []string{"pages", "list", "--limit=20"}, 20},
		{"--max エイリアス", []string{"pages", "list", "--max=20"}, 20},
		{"--n エイリアス", []string{"pages", "list", "--n=20"}, 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// CLIを初期化
			cli := &CLI{}

			// Kongパーサーを作成
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // テスト時は終了しない
			)
			if err != nil {
				t.Fatalf("Kongパーサーの作成に失敗: %v", err)
			}

			// コマンドラインをパース
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("コマンドラインのパースに失敗: %v", err)
			}

			// Limitフィールドが正しく設定されているか確認
			if cli.Pages.List.Limit != tc.want {
				t.Errorf("Pages.List.Limitが正しく設定されていません: got=%d, want=%d", cli.Pages.List.Limit, tc.want)
			}
		})
	}
}

// TestTagsFlagAliases はTagsコマンドのLimitエイリアスが正しく機能することを確認します
func TestTagsFlagAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want int
	}{
		{"--limit フラグ", []string{"tags", "list", "--limit=30"}, 30},
		{"--max エイリアス", []string{"tags", "list", "--max=30"}, 30},
		{"--n エイリアス", []string{"tags", "list", "--n=30"}, 30},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// CLIを初期化
			cli := &CLI{}

			// Kongパーサーを作成
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // テスト時は終了しない
			)
			if err != nil {
				t.Fatalf("Kongパーサーの作成に失敗: %v", err)
			}

			// コマンドラインをパース
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("コマンドラインのパースに失敗: %v", err)
			}

			// Limitフィールドが正しく設定されているか確認
			if cli.Tags.List.Limit != tc.want {
				t.Errorf("Tags.List.Limitが正しく設定されていません: got=%d, want=%d", cli.Tags.List.Limit, tc.want)
			}
		})
	}
}
