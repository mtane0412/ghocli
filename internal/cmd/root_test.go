/**
 * root_test.go
 * root.goの環境変数統合テスト
 */

package cmd

import (
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRootFlags_PlainEnvVar はGHO_PLAIN環境変数が正しく読み込まれることをテストします
func TestRootFlags_PlainEnvVar(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GHO_PLAIN", "1")
	defer os.Unsetenv("GHO_PLAIN")

	// CLIをパース（posts listコマンドを使用）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// GHO_PLAINが設定されているので、Plainがtrueになるはず
	assert.True(t, cli.Plain, "GHO_PLAIN環境変数が設定されている場合、Plainがtrueになるべき")
}

// TestRootFlags_VerboseEnvVar はGHO_VERBOSE環境変数が正しく読み込まれることをテストします
func TestRootFlags_VerboseEnvVar(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GHO_VERBOSE", "1")
	defer os.Unsetenv("GHO_VERBOSE")

	// CLIをパース（posts listコマンドを使用）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// GHO_VERBOSEが設定されているので、Verboseがtrueになるはず
	assert.True(t, cli.Verbose, "GHO_VERBOSE環境変数が設定されている場合、Verboseがtrueになるべき")
}

// TestRootFlags_PlainFlagOverridesEnv はフラグが環境変数より優先されることをテストします
func TestRootFlags_PlainFlagOverridesEnv(t *testing.T) {
	// 環境変数を設定（falseとして解釈される値）
	os.Setenv("GHO_PLAIN", "0")
	defer os.Unsetenv("GHO_PLAIN")

	// CLIをパース（--plainフラグを指定）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--plain", "posts", "list"})
	require.NoError(t, err)

	// フラグが環境変数より優先されるので、Plainがtrueになるはず
	assert.True(t, cli.Plain, "--plainフラグは環境変数より優先されるべき")
}

// TestRootFlags_VerboseFlagOverridesEnv はフラグが環境変数より優先されることをテストします
func TestRootFlags_VerboseFlagOverridesEnv(t *testing.T) {
	// 環境変数を設定（falseとして解釈される値）
	os.Setenv("GHO_VERBOSE", "0")
	defer os.Unsetenv("GHO_VERBOSE")

	// CLIをパース（-vフラグを指定）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"-v", "posts", "list"})
	require.NoError(t, err)

	// フラグが環境変数より優先されるので、Verboseがtrueになるはず
	assert.True(t, cli.Verbose, "-vフラグは環境変数より優先されるべき")
}

// TestRootFlags_ColorEnvVar はGHO_COLOR環境変数が正しく読み込まれることをテストします
func TestRootFlags_ColorEnvVar(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GHO_COLOR", "never")
	defer os.Unsetenv("GHO_COLOR")

	// CLIをパース（posts listコマンドを使用）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// GHO_COLORが設定されているので、Colorが"never"になるはず
	assert.Equal(t, "never", cli.Color, "GHO_COLOR環境変数が設定されている場合、Colorが正しく設定されるべき")
}

// TestRootFlags_ColorFlagOverridesEnv はフラグが環境変数より優先されることをテストします
func TestRootFlags_ColorFlagOverridesEnv(t *testing.T) {
	// 環境変数を設定
	os.Setenv("GHO_COLOR", "never")
	defer os.Unsetenv("GHO_COLOR")

	// CLIをパース（--color=alwaysフラグを指定）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--color=always", "posts", "list"})
	require.NoError(t, err)

	// フラグが環境変数より優先されるので、Colorが"always"になるはず
	assert.Equal(t, "always", cli.Color, "--colorフラグは環境変数より優先されるべき")
}

// TestRootFlags_ColorDefault はデフォルト値が"auto"であることをテストします
func TestRootFlags_ColorDefault(t *testing.T) {
	// 環境変数を削除
	os.Unsetenv("GHO_COLOR")

	// CLIをパース
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// デフォルト値は"auto"になるはず
	assert.Equal(t, "auto", cli.Color, "デフォルト値は'auto'であるべき")
}
