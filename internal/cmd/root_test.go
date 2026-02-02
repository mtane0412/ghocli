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
	// Set environment variable
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
	// Set environment variable
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
	// Set environment variable（falseとして解釈される値）
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
	// Set environment variable（falseとして解釈される値）
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
	// Set environment variable
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
	// Set environment variable
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

// TestExecute_ヘルプメッセージが表示される
// Kongは--helpでos.Exit(0)を呼び出すため、このテストはスキップする
func TestExecute_ヘルプメッセージが表示される(t *testing.T) {
	t.Skip("Kongは--helpでos.Exit(0)を呼び出すため、テストできない")
}

// TestExecute_バージョン表示
// Kongは--versionでos.Exit(0)を呼び出すため、このテストはスキップする
func TestExecute_バージョン表示(t *testing.T) {
	t.Skip("Kongは--versionでos.Exit(0)を呼び出すため、テストできない")
}

// TestExecute_不正なコマンド
func TestExecute_不正なコマンド(t *testing.T) {
	// 存在しないコマンドを渡す
	args := []string{"gho", "invalid-command"}

	// エラーを返すべき
	err := Execute(args)
	if err == nil {
		t.Error("Execute(invalid-command) returned nil, want error")
	}

	// 終了コードは0以外
	code := ExitCode(err)
	if code == 0 {
		t.Error("Execute(invalid-command) exit code = 0, want non-zero")
	}
}

// TestExecute_近似コマンド提案 はレーベンシュタイン距離が2以下のコマンドを提案することをテストします
func TestExecute_近似コマンド提案(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		wantErrorString string
	}{
		{
			name:            "themesのタイポをthemesと提案",
			args:            []string{"gho", "themse"},
			wantErrorString: "did you mean \"themes\"?",
		},
		{
			name:            "postsのタイポをpostsと提案",
			args:            []string{"gho", "postss"},
			wantErrorString: "did you mean \"posts\"?",
		},
		{
			name:            "authのタイポをauthと提案",
			args:            []string{"gho", "auht"},
			wantErrorString: "did you mean \"auth\"?",
		},
		{
			name:            "完全に不明なコマンドは提案なし",
			args:            []string{"gho", "unknowncommand"},
			wantErrorString: "unexpected argument",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// コマンドを実行
			err := Execute(tc.args)

			// エラーが返されるべき
			require.Error(t, err, "Execute(%v) should return error", tc.args)

			// エラーメッセージに期待する文字列が含まれるべき
			assert.Contains(t, err.Error(), tc.wantErrorString,
				"Execute(%v) error message should contain %q", tc.args, tc.wantErrorString)
		})
	}
}

// TestRootFlags_Fieldsフィールド はRootFlagsにFieldsフィールドが存在することを確認します
func TestRootFlags_Fieldsフィールド(t *testing.T) {
	// RootFlagsインスタンスを作成
	flags := &RootFlags{
		Fields: "id,title,status",
	}

	// Fieldsフィールドが設定されることを確認
	assert.Equal(t, "id,title,status", flags.Fields, "Fieldsフィールドが正しく設定されるべき")
}

// TestRootFlags_Fieldsデフォルト値 はFieldsのデフォルト値を確認します
func TestRootFlags_Fieldsデフォルト値(t *testing.T) {
	// RootFlagsインスタンスを作成（デフォルト値）
	flags := &RootFlags{}

	// Fieldsのデフォルト値が空文字列であることを確認
	assert.Equal(t, "", flags.Fields, "Fieldsのデフォルト値は空文字列であるべき")
}

// TestRootFlags_FieldsEnvVar はGHO_FIELDS環境変数が正しく読み込まれることをテストします
func TestRootFlags_FieldsEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_FIELDS", "id,title,url")
	defer os.Unsetenv("GHO_FIELDS")

	// CLIをパース（posts listコマンドを使用）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// GHO_FIELDSが設定されているので、Fieldsが"id,title,url"になるはず
	assert.Equal(t, "id,title,url", cli.Fields, "GHO_FIELDS環境変数が設定されている場合、Fieldsが正しく設定されるべき")
}

// TestRootFlags_FieldsFlagOverridesEnv はフラグが環境変数より優先されることをテストします
func TestRootFlags_FieldsFlagOverridesEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_FIELDS", "id,title")
	defer os.Unsetenv("GHO_FIELDS")

	// CLIをパース（--fieldsフラグを指定）
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--fields=id,title,status,url", "posts", "list"})
	require.NoError(t, err)

	// フラグが環境変数より優先されるので、Fieldsが"id,title,status,url"になるはず
	assert.Equal(t, "id,title,status,url", cli.Fields, "--fieldsフラグは環境変数より優先されるべき")
}

// TestCommandAliases_Posts はpostsコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Posts(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"postsコマンド", "posts"},
		{"postエイリアス", "post"},
		{"pエイリアス", "p"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Tags はtagsコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Tags(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"tagsコマンド", "tags"},
		{"tagエイリアス", "tag"},
		{"tエイリアス", "t"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Pages はpagesコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Pages(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"pagesコマンド", "pages"},
		{"pageエイリアス", "page"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Members はmembersコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Members(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"membersコマンド", "members"},
		{"memberエイリアス", "member"},
		{"mエイリアス", "m"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Users はusersコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Users(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"usersコマンド", "users"},
		{"userエイリアス", "user"},
		{"uエイリアス", "u"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Newsletters はnewslettersコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Newsletters(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"newslettersコマンド", "newsletters"},
		{"newsletterエイリアス", "newsletter"},
		{"nlエイリアス", "nl"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Tiers はtiersコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Tiers(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"tiersコマンド", "tiers"},
		{"tierエイリアス", "tier"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Offers はoffersコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Offers(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"offersコマンド", "offers"},
		{"offerエイリアス", "offer"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Webhooks はwebhooksコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Webhooks(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"webhooksコマンド", "webhooks"},
		{"webhookエイリアス", "webhook"},
		{"whエイリアス", "wh"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			// webhooksはlistサブコマンドを持たないため、createで確認
			// 必須引数がないため失敗するが、コマンド自体は認識される
			_, err = parser.Parse([]string{tc.command, "create"})
			// エラーは発生するが、"unexpected argument"ではないことを確認
			if err != nil {
				assert.NotContains(t, err.Error(), "unexpected argument", "%sコマンドは認識されるべき", tc.command)
			}
		})
	}
}

// TestCommandAliases_Settings はsettingsコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Settings(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"settingsコマンド", "settings"},
		{"settingエイリアス", "setting"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%sコマンドは正しくパースされるべき", tc.command)
		})
	}
}

// TestCommandAliases_Images はimagesコマンドのエイリアスが動作することをテストします
func TestCommandAliases_Images(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"imagesコマンド", "images"},
		{"imageエイリアス", "image"},
		{"imgエイリアス", "img"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			// imagesはuploadサブコマンドを持つ
			// 必須引数がないため失敗するが、コマンド自体は認識される
			_, err = parser.Parse([]string{tc.command, "upload"})
			// エラーは発生するが、"unexpected argument"ではないことを確認
			if err != nil {
				assert.NotContains(t, err.Error(), "unexpected argument", "%sコマンドは認識されるべき", tc.command)
			}
		})
	}
}
