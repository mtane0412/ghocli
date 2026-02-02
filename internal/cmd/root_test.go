/**
 * root_test.go
 * Integration tests for environment variables in root.go
 */

package cmd

import (
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRootFlags_PlainEnvVar verifies that GHO_PLAIN environment variable is correctly loaded
func TestRootFlags_PlainEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_PLAIN", "1")
	defer os.Unsetenv("GHO_PLAIN")

	// Parse CLI (using posts list command)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// Since GHO_PLAIN is set, Plain should be true
	assert.True(t, cli.Plain, "When GHO_PLAIN environment variable is set, Plain should be true")
}

// TestRootFlags_VerboseEnvVar verifies that GHO_VERBOSE environment variable is correctly loaded
func TestRootFlags_VerboseEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_VERBOSE", "1")
	defer os.Unsetenv("GHO_VERBOSE")

	// Parse CLI (using posts list command)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// Since GHO_VERBOSE is set, Verbose should be true
	assert.True(t, cli.Verbose, "When GHO_VERBOSE environment variable is set, Verbose should be true")
}

// TestRootFlags_PlainFlagOverridesEnv verifies that flag takes precedence over environment variable
func TestRootFlags_PlainFlagOverridesEnv(t *testing.T) {
	// Set environment variable (value interpreted as false)
	os.Setenv("GHO_PLAIN", "0")
	defer os.Unsetenv("GHO_PLAIN")

	// Parse CLI (with --plain flag specified)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--plain", "posts", "list"})
	require.NoError(t, err)

	// Since flag takes precedence over environment variable, Plain should be true
	assert.True(t, cli.Plain, "--plain flag should take precedence over environment variable")
}

// TestRootFlags_VerboseFlagOverridesEnv verifies that flag takes precedence over environment variable
func TestRootFlags_VerboseFlagOverridesEnv(t *testing.T) {
	// Set environment variable (value interpreted as false)
	os.Setenv("GHO_VERBOSE", "0")
	defer os.Unsetenv("GHO_VERBOSE")

	// Parse CLI (with -v flag specified)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"-v", "posts", "list"})
	require.NoError(t, err)

	// Since flag takes precedence over environment variable, Verbose should be true
	assert.True(t, cli.Verbose, "-v flag should take precedence over environment variable")
}

// TestRootFlags_ColorEnvVar verifies that GHO_COLOR environment variable is correctly loaded
func TestRootFlags_ColorEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_COLOR", "never")
	defer os.Unsetenv("GHO_COLOR")

	// Parse CLI (using posts list command)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// Since GHO_COLOR is set, Color should be "never"
	assert.Equal(t, "never", cli.Color, "When GHO_COLOR environment variable is set, Color should be set correctly")
}

// TestRootFlags_ColorFlagOverridesEnv verifies that flag takes precedence over environment variable
func TestRootFlags_ColorFlagOverridesEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_COLOR", "never")
	defer os.Unsetenv("GHO_COLOR")

	// Parse CLI (with --color=always flag specified)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--color=always", "posts", "list"})
	require.NoError(t, err)

	// Since flag takes precedence over environment variable, Color should be "always"
	assert.Equal(t, "always", cli.Color, "--color flag should take precedence over environment variable")
}

// TestRootFlags_ColorDefault verifies that default value is "auto"
func TestRootFlags_ColorDefault(t *testing.T) {
	// Unset environment variable
	os.Unsetenv("GHO_COLOR")

	// Parse CLI
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// Default value should be "auto"
	assert.Equal(t, "auto", cli.Color, "Default value should be 'auto'")
}

// TestExecute_HelpMessage verifies help message is displayed
// Kong calls os.Exit(0) with --help, so this test is skipped
func TestExecute_HelpMessage(t *testing.T) {
	t.Skip("Cannot test because Kong calls os.Exit(0) with --help")
}

// TestExecute_VersionDisplay verifies version is displayed
// Kong calls os.Exit(0) with --version, so this test is skipped
func TestExecute_VersionDisplay(t *testing.T) {
	t.Skip("Cannot test because Kong calls os.Exit(0) with --version")
}

// TestExecute_InvalidCommand verifies behavior with invalid command
func TestExecute_InvalidCommand(t *testing.T) {
	// Pass non-existent command
	args := []string{"gho", "invalid-command"}

	// Should return error
	err := Execute(args)
	if err == nil {
		t.Error("Execute(invalid-command) returned nil, want error")
	}

	// Exit code should be non-zero
	code := ExitCode(err)
	if code == 0 {
		t.Error("Execute(invalid-command) exit code = 0, want non-zero")
	}
}

// TestExecute_SimilarCommandSuggestion verifies that commands with Levenshtein distance of 2 or less are suggested
func TestExecute_SimilarCommandSuggestion(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		wantErrorString string
	}{
		{
			name:            "suggest themes for typo of themes",
			args:            []string{"gho", "themse"},
			wantErrorString: "did you mean \"themes\"?",
		},
		{
			name:            "suggest posts for typo of posts",
			args:            []string{"gho", "postss"},
			wantErrorString: "did you mean \"posts\"?",
		},
		{
			name:            "suggest auth for typo of auth",
			args:            []string{"gho", "auht"},
			wantErrorString: "did you mean \"auth\"?",
		},
		{
			name:            "no suggestion for completely unknown command",
			args:            []string{"gho", "unknowncommand"},
			wantErrorString: "unexpected argument",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute command
			err := Execute(tc.args)

			// Error should be returned
			require.Error(t, err, "Execute(%v) should return error", tc.args)

			// Error message should contain expected string
			assert.Contains(t, err.Error(), tc.wantErrorString,
				"Execute(%v) error message should contain %q", tc.args, tc.wantErrorString)
		})
	}
}

// TestRootFlags_FieldsField verifies that Fields field exists in RootFlags
func TestRootFlags_FieldsField(t *testing.T) {
	// Create RootFlags instance
	flags := &RootFlags{
		Fields: "id,title,status",
	}

	// Verify that Fields field is set
	assert.Equal(t, "id,title,status", flags.Fields, "Fields field should be set correctly")
}

// TestRootFlags_FieldsDefault verifies the default value of Fields
func TestRootFlags_FieldsDefault(t *testing.T) {
	// Create RootFlags instance (with default value)
	flags := &RootFlags{}

	// Verify that default value of Fields is empty string
	assert.Equal(t, "", flags.Fields, "Default value of Fields should be empty string")
}

// TestRootFlags_FieldsEnvVar verifies that GHO_FIELDS environment variable is correctly loaded
func TestRootFlags_FieldsEnvVar(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_FIELDS", "id,title,url")
	defer os.Unsetenv("GHO_FIELDS")

	// Parse CLI (using posts list command)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"posts", "list"})
	require.NoError(t, err)

	// Since GHO_FIELDS is set, Fields should be "id,title,url"
	assert.Equal(t, "id,title,url", cli.Fields, "When GHO_FIELDS environment variable is set, Fields should be set correctly")
}

// TestRootFlags_FieldsFlagOverridesEnv verifies that flag takes precedence over environment variable
func TestRootFlags_FieldsFlagOverridesEnv(t *testing.T) {
	// Set environment variable
	os.Setenv("GHO_FIELDS", "id,title")
	defer os.Unsetenv("GHO_FIELDS")

	// Parse CLI (with --fields flag specified)
	var cli CLI
	parser, err := kong.New(&cli)
	require.NoError(t, err)

	_, err = parser.Parse([]string{"--fields=id,title,status,url", "posts", "list"})
	require.NoError(t, err)

	// Since flag takes precedence over environment variable, Fields should be "id,title,status,url"
	assert.Equal(t, "id,title,status,url", cli.Fields, "--fields flag should take precedence over environment variable")
}

// TestCommandAliases_Posts verifies that posts command aliases work
func TestCommandAliases_Posts(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"posts command", "posts"},
		{"post alias", "post"},
		{"p alias", "p"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Tags verifies that tags command aliases work
func TestCommandAliases_Tags(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"tags command", "tags"},
		{"tag alias", "tag"},
		{"t alias", "t"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Pages verifies that pages command aliases work
func TestCommandAliases_Pages(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"pages command", "pages"},
		{"page alias", "page"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Members verifies that members command aliases work
func TestCommandAliases_Members(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"members command", "members"},
		{"member alias", "member"},
		{"m alias", "m"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Users verifies that users command aliases work
func TestCommandAliases_Users(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"users command", "users"},
		{"user alias", "user"},
		{"u alias", "u"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Newsletters verifies that newsletters command aliases work
func TestCommandAliases_Newsletters(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"newsletters command", "newsletters"},
		{"newsletter alias", "newsletter"},
		{"nl alias", "nl"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Tiers verifies that tiers command aliases work
func TestCommandAliases_Tiers(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"tiers command", "tiers"},
		{"tier alias", "tier"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Offers verifies that offers command aliases work
func TestCommandAliases_Offers(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"offers command", "offers"},
		{"offer alias", "offer"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Webhooks verifies that webhooks command aliases work
func TestCommandAliases_Webhooks(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"webhooks command", "webhooks"},
		{"webhook alias", "webhook"},
		{"wh alias", "wh"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			// webhooks does not have list subcommand, so verify with create
			// It will fail due to missing required arguments, but the command itself should be recognized
			_, err = parser.Parse([]string{tc.command, "create"})
			// Error occurs, but verify it's not "unexpected argument"
			if err != nil {
				assert.NotContains(t, err.Error(), "unexpected argument", "%s command should be recognized", tc.command)
			}
		})
	}
}

// TestCommandAliases_Settings verifies that settings command aliases work
func TestCommandAliases_Settings(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"settings command", "settings"},
		{"setting alias", "setting"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			_, err = parser.Parse([]string{tc.command, "list"})
			require.NoError(t, err, "%s command should be parsed correctly", tc.command)
		})
	}
}

// TestCommandAliases_Images verifies that images command aliases work
func TestCommandAliases_Images(t *testing.T) {
	testCases := []struct {
		name    string
		command string
	}{
		{"images command", "images"},
		{"image alias", "image"},
		{"img alias", "img"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cli CLI
			parser, err := kong.New(&cli)
			require.NoError(t, err)

			// images has upload subcommand
			// It will fail due to missing required arguments, but the command itself should be recognized
			_, err = parser.Parse([]string{tc.command, "upload"})
			// Error occurs, but verify it's not "unexpected argument"
			if err != nil {
				assert.NotContains(t, err.Error(), "unexpected argument", "%s command should be recognized", tc.command)
			}
		})
	}
}
