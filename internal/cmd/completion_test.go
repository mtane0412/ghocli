/**
 * completion_test.go
 * Tests for shell completion functionality
 */

package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompletionScript_Bash はbashスクリプトが生成されることをテストします
func TestCompletionScript_Bash(t *testing.T) {
	script, err := completionScript("bash")
	require.NoError(t, err, "bashスクリプトの生成に成功するべき")
	assert.Contains(t, script, "_gho_complete", "bashスクリプトに_gho_complete関数が含まれるべき")
	assert.Contains(t, script, "gho __complete", "bashスクリプトにgho __completeコマンドが含まれるべき")
}

// TestCompletionScript_Zsh はzshスクリプトが生成されることをテストします
func TestCompletionScript_Zsh(t *testing.T) {
	script, err := completionScript("zsh")
	require.NoError(t, err, "zshスクリプトの生成に成功するべき")
	assert.Contains(t, script, "#compdef gho", "zshスクリプトに#compdef ghoが含まれるべき")
	assert.Contains(t, script, "_gho_complete", "zshスクリプトに_gho_complete関数が含まれるべき")
}

// TestCompletionScript_Fish はfishスクリプトが生成されることをテストします
func TestCompletionScript_Fish(t *testing.T) {
	script, err := completionScript("fish")
	require.NoError(t, err, "fishスクリプトの生成に成功するべき")
	assert.Contains(t, script, "__gho_complete", "fishスクリプトに__gho_complete関数が含まれるべき")
	assert.Contains(t, script, "gho __complete", "fishスクリプトにgho __completeコマンドが含まれるべき")
}

// TestCompletionScript_PowerShell はpowershellスクリプトが生成されることをテストします
func TestCompletionScript_PowerShell(t *testing.T) {
	script, err := completionScript("powershell")
	require.NoError(t, err, "powershellスクリプトの生成に成功するべき")
	assert.Contains(t, script, "Register-ArgumentCompleter", "powershellスクリプトにRegister-ArgumentCompleterが含まれるべき")
	assert.Contains(t, script, "gho __complete", "powershellスクリプトにgho __completeコマンドが含まれるべき")
}

// TestCompletionScript_UnsupportedShell は未サポートのシェルでエラーを返すことをテストします
func TestCompletionScript_UnsupportedShell(t *testing.T) {
	_, err := completionScript("unsupported")
	require.Error(t, err, "未サポートのシェルでエラーを返すべき")
	assert.Contains(t, err.Error(), "unsupported shell", "エラーメッセージに'unsupported shell'が含まれるべき")
}

// TestCompleteWords_Commands はコマンド補完が動作することをテストします
func TestCompleteWords_Commands(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		cword    int
		wantOne  string
		wantMany []string
	}{
		{
			name:     "postsで始まるコマンドを補完",
			words:    []string{"gho", "po"},
			cword:    1,
			wantOne:  "posts",
			wantMany: nil,
		},
		{
			name:     "pエイリアスも補完候補に含まれる",
			words:    []string{"gho", "p"},
			cword:    1,
			wantOne:  "posts",
			wantMany: []string{"posts", "pages"},
		},
		{
			name:     "tagsで始まるコマンドを補完",
			words:    []string{"gho", "ta"},
			cword:    1,
			wantOne:  "tags",
			wantMany: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := completeWords(tc.cword, tc.words)
			require.NoError(t, err, "補完候補の取得に成功するべき")

			if tc.wantOne != "" {
				assert.Contains(t, items, tc.wantOne, "補完候補に%sが含まれるべき", tc.wantOne)
			}
			if tc.wantMany != nil {
				for _, want := range tc.wantMany {
					assert.Contains(t, items, want, "補完候補に%sが含まれるべき", want)
				}
			}
		})
	}
}

// TestCompleteWords_Flags はフラグ補完が動作することをテストします
func TestCompleteWords_Flags(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		cword    int
		wantFlag string
	}{
		{
			name:     "--siteフラグを補完",
			words:    []string{"gho", "--si"},
			cword:    1,
			wantFlag: "--site",
		},
		{
			name:     "--jsonフラグを補完",
			words:    []string{"gho", "--js"},
			cword:    1,
			wantFlag: "--json",
		},
		{
			name:     "-sショートフラグを補完",
			words:    []string{"gho", "-s"},
			cword:    1,
			wantFlag: "-s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := completeWords(tc.cword, tc.words)
			require.NoError(t, err, "補完候補の取得に成功するべき")
			assert.Contains(t, items, tc.wantFlag, "補完候補に%sが含まれるべき", tc.wantFlag)
		})
	}
}

// TestCompleteWords_Subcommands はサブコマンド補完が動作することをテストします
func TestCompleteWords_Subcommands(t *testing.T) {
	testCases := []struct {
		name        string
		words       []string
		cword       int
		wantCommand string
	}{
		{
			name:        "posts listサブコマンドを補完",
			words:       []string{"gho", "posts", "li"},
			cword:       2,
			wantCommand: "list",
		},
		{
			name:        "posts getサブコマンドを補完",
			words:       []string{"gho", "posts", "ge"},
			cword:       2,
			wantCommand: "get",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := completeWords(tc.cword, tc.words)
			require.NoError(t, err, "補完候補の取得に成功するべき")
			assert.Contains(t, items, tc.wantCommand, "補完候補に%sが含まれるべき", tc.wantCommand)
		})
	}
}

// TestIsProgramName はプログラム名の判定をテストします
func TestIsProgramName(t *testing.T) {
	testCases := []struct {
		name string
		word string
		want bool
	}{
		{"gho", "gho", true},
		{"/usr/local/bin/gho", "/usr/local/bin/gho", true},
		{"gho.exe", "gho.exe", true},
		{"GHO", "GHO", true},
		{"posts", "posts", false},
		{"other", "other", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isProgramName(tc.word)
			assert.Equal(t, tc.want, got, "isProgramName(%q) = %v, want %v", tc.word, got, tc.want)
		})
	}
}

// TestBuildCompletionNode はKongモデルから補完ノードを構築することをテストします
func TestBuildCompletionNode(t *testing.T) {
	// パーサーを作成
	parser, _, err := newParser()
	require.NoError(t, err, "パーサーの作成に成功するべき")

	// 補完ノードを構築
	root := buildCompletionNode(parser.Model.Node)

	// トップレベルコマンドが含まれることを確認
	assert.Contains(t, root.children, "posts", "postsコマンドが含まれるべき")
	assert.Contains(t, root.children, "post", "postエイリアスが含まれるべき")
	assert.Contains(t, root.children, "p", "pエイリアスが含まれるべき")
	assert.Contains(t, root.children, "tags", "tagsコマンドが含まれるべき")
	assert.Contains(t, root.children, "auth", "authコマンドが含まれるべき")

	// フラグが含まれることを確認
	assert.Contains(t, root.flags, "--site", "--siteフラグが含まれるべき")
	assert.Contains(t, root.flags, "-s", "-sショートフラグが含まれるべき")
	assert.Contains(t, root.flags, "--json", "--jsonフラグが含まれるべき")
}

// TestMatchingCommands はコマンドマッチングをテストします
func TestMatchingCommands(t *testing.T) {
	node := &completionNode{
		children: map[string]*completionNode{
			"posts":  {},
			"post":   {},
			"p":      {},
			"pages":  {},
			"page":   {},
			"auth":   {},
			"config": {},
		},
	}

	testCases := []struct {
		name   string
		prefix string
		want   []string
	}{
		{"po prefix", "po", []string{"posts", "post"}},
		{"p prefix", "p", []string{"posts", "post", "p", "pages", "page"}},
		{"pa prefix", "pa", []string{"pages", "page"}},
		{"au prefix", "au", []string{"auth"}},
		{"empty prefix", "", []string{"posts", "post", "p", "pages", "page", "auth", "config"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchingCommands(node, tc.prefix)
			for _, want := range tc.want {
				assert.Contains(t, got, want, "matchingCommands(%q)に%sが含まれるべき", tc.prefix, want)
			}
		})
	}
}

// TestMatchingFlags はフラグマッチングをテストします
func TestMatchingFlags(t *testing.T) {
	node := &completionNode{
		flags: map[string]completionFlag{
			"--site":   {takesValue: true},
			"-s":       {takesValue: true},
			"--json":   {takesValue: false},
			"--plain":  {takesValue: false},
			"--fields": {takesValue: true},
			"-F":       {takesValue: true},
		},
	}

	testCases := []struct {
		name   string
		prefix string
		want   []string
	}{
		{"--s prefix", "--s", []string{"--site"}},
		{"--j prefix", "--j", []string{"--json"}},
		{"-s prefix", "-s", []string{"-s"}},
		{"-- prefix", "--", []string{"--site", "--json", "--plain", "--fields"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchingFlags(node, tc.prefix)
			for _, want := range tc.want {
				assert.Contains(t, got, want, "matchingFlags(%q)に%sが含まれるべき", tc.prefix, want)
			}
		})
	}
}

// TestSplitFlagToken はフラグトークンの分割をテストします
func TestSplitFlagToken(t *testing.T) {
	testCases := []struct {
		name      string
		word      string
		wantToken string
		wantValue bool
	}{
		{"--site=myblog", "--site=myblog", "--site", true},
		{"--json", "--json", "--json", false},
		{"-s=myblog", "-s=myblog", "-s", true},
		{"-s", "-s", "-s", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotToken, gotValue := splitFlagToken(tc.word)
			assert.Equal(t, tc.wantToken, gotToken, "splitFlagToken(%q) token", tc.word)
			assert.Equal(t, tc.wantValue, gotValue, "splitFlagToken(%q) hasValue", tc.word)
		})
	}
}

// TestCompletionInternalCmd_Run は__completeコマンドが動作することをテストします
func TestCompletionInternalCmd_Run(t *testing.T) {
	// __completeは隠しコマンドのため、completeWords関数を直接テストする
	words := []string{"gho", "po"}
	cword := 1

	items, err := completeWords(cword, words)
	require.NoError(t, err, "補完候補の取得に成功するべき")
	assert.Contains(t, items, "posts", "補完候補にpostsが含まれるべき")
}

// TestCompletionCmd_Run はcompletionコマンドが動作することをテストします
func TestCompletionCmd_Run(t *testing.T) {
	// このテストはactualの実行は難しいため、基本的なパース確認のみ
	var cli CLI
	parser, err := NewParserForTest(&cli)
	require.NoError(t, err)

	// completionコマンドがパースできることを確認
	testCases := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range testCases {
		t.Run(shell, func(t *testing.T) {
			_, err = parser.Parse([]string{"completion", shell})
			require.NoError(t, err, "completion %sコマンドがパースできるべき", shell)
		})
	}
}

// NewParserForTest はテスト用のパーサーを作成します
func NewParserForTest(cli *CLI) (*kong.Kong, error) {
	parser, err := kong.New(cli,
		kong.Name("gho"),
		kong.Description("Ghost Admin API CLI"),
	)
	if err != nil {
		return nil, err
	}
	return parser, nil
}
