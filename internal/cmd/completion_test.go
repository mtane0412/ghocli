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

// TestCompletionScript_Bash verifies that bash script is generated
func TestCompletionScript_Bash(t *testing.T) {
	script, err := completionScript("bash")
	require.NoError(t, err, "bash script generation should succeed")
	assert.Contains(t, script, "_gho_complete", "bash script should contain _gho_complete function")
	assert.Contains(t, script, "gho __complete", "bash script should contain gho __complete command")
}

// TestCompletionScript_Zsh verifies that zsh script is generated
func TestCompletionScript_Zsh(t *testing.T) {
	script, err := completionScript("zsh")
	require.NoError(t, err, "zsh script generation should succeed")
	assert.Contains(t, script, "#compdef gho", "zsh script should contain #compdef gho")
	assert.Contains(t, script, "_gho_complete", "zsh script should contain _gho_complete function")
}

// TestCompletionScript_Fish verifies that fish script is generated
func TestCompletionScript_Fish(t *testing.T) {
	script, err := completionScript("fish")
	require.NoError(t, err, "fish script generation should succeed")
	assert.Contains(t, script, "__gho_complete", "fish script should contain __gho_complete function")
	assert.Contains(t, script, "gho __complete", "fish script should contain gho __complete command")
}

// TestCompletionScript_PowerShell verifies that powershell script is generated
func TestCompletionScript_PowerShell(t *testing.T) {
	script, err := completionScript("powershell")
	require.NoError(t, err, "powershell script generation should succeed")
	assert.Contains(t, script, "Register-ArgumentCompleter", "powershell script should contain Register-ArgumentCompleter")
	assert.Contains(t, script, "gho __complete", "powershell script should contain gho __complete command")
}

// TestCompletionScript_UnsupportedShell verifies that an error is returned for unsupported shells
func TestCompletionScript_UnsupportedShell(t *testing.T) {
	_, err := completionScript("unsupported")
	require.Error(t, err, "should return error for unsupported shell")
	assert.Contains(t, err.Error(), "unsupported shell", "error message should contain 'unsupported shell'")
}

// TestCompleteWords_Commands verifies that command completion works
func TestCompleteWords_Commands(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		cword    int
		wantOne  string
		wantMany []string
	}{
		{
			name:     "complete command starting with posts",
			words:    []string{"gho", "po"},
			cword:    1,
			wantOne:  "posts",
			wantMany: nil,
		},
		{
			name:     "p alias is also included in completion candidates",
			words:    []string{"gho", "p"},
			cword:    1,
			wantOne:  "posts",
			wantMany: []string{"posts", "pages"},
		},
		{
			name:     "complete command starting with tags",
			words:    []string{"gho", "ta"},
			cword:    1,
			wantOne:  "tags",
			wantMany: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := completeWords(tc.cword, tc.words)
			require.NoError(t, err, "completion candidates retrieval should succeed")

			if tc.wantOne != "" {
				assert.Contains(t, items, tc.wantOne, "completion candidates should contain %s", tc.wantOne)
			}
			if tc.wantMany != nil {
				for _, want := range tc.wantMany {
					assert.Contains(t, items, want, "completion candidates should contain %s", want)
				}
			}
		})
	}
}

// TestCompleteWords_Flags verifies that flag completion works
func TestCompleteWords_Flags(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		cword    int
		wantFlag string
	}{
		{
			name:     "complete --site flag",
			words:    []string{"gho", "--si"},
			cword:    1,
			wantFlag: "--site",
		},
		{
			name:     "complete --json flag",
			words:    []string{"gho", "--js"},
			cword:    1,
			wantFlag: "--json",
		},
		{
			name:     "complete -s short flag",
			words:    []string{"gho", "-s"},
			cword:    1,
			wantFlag: "-s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := completeWords(tc.cword, tc.words)
			require.NoError(t, err, "completion candidates retrieval should succeed")
			assert.Contains(t, items, tc.wantFlag, "completion candidates should contain %s", tc.wantFlag)
		})
	}
}

// TestCompleteWords_Subcommands verifies that subcommand completion works
func TestCompleteWords_Subcommands(t *testing.T) {
	testCases := []struct {
		name        string
		words       []string
		cword       int
		wantCommand string
	}{
		{
			name:        "complete posts list subcommand",
			words:       []string{"gho", "posts", "li"},
			cword:       2,
			wantCommand: "list",
		},
		{
			name:        "complete posts get subcommand",
			words:       []string{"gho", "posts", "ge"},
			cword:       2,
			wantCommand: "get",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := completeWords(tc.cword, tc.words)
			require.NoError(t, err, "completion candidates retrieval should succeed")
			assert.Contains(t, items, tc.wantCommand, "completion candidates should contain %s", tc.wantCommand)
		})
	}
}

// TestIsProgramName verifies program name detection
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

// TestBuildCompletionNode verifies that completion node is built from Kong model
func TestBuildCompletionNode(t *testing.T) {
	// Create parser
	parser, _, err := newParser()
	require.NoError(t, err, "parser creation should succeed")

	// Build completion node
	root := buildCompletionNode(parser.Model.Node)

	// Verify that top-level commands are included
	assert.Contains(t, root.children, "posts", "should contain posts command")
	assert.Contains(t, root.children, "post", "should contain post alias")
	assert.Contains(t, root.children, "p", "should contain p alias")
	assert.Contains(t, root.children, "tags", "should contain tags command")
	assert.Contains(t, root.children, "auth", "should contain auth command")

	// Verify that flags are included
	assert.Contains(t, root.flags, "--site", "should contain --site flag")
	assert.Contains(t, root.flags, "-s", "should contain -s short flag")
	assert.Contains(t, root.flags, "--json", "should contain --json flag")
}

// TestMatchingCommands verifies command matching
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
				assert.Contains(t, got, want, "matchingCommands(%q) should contain %s", tc.prefix, want)
			}
		})
	}
}

// TestMatchingFlags verifies flag matching
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
				assert.Contains(t, got, want, "matchingFlags(%q) should contain %s", tc.prefix, want)
			}
		})
	}
}

// TestSplitFlagToken verifies flag token splitting
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

// TestCompletionInternalCmd_Run verifies that __complete command works
func TestCompletionInternalCmd_Run(t *testing.T) {
	// Since __complete is a hidden command, test the completeWords function directly
	words := []string{"gho", "po"}
	cword := 1

	items, err := completeWords(cword, words)
	require.NoError(t, err, "completion candidates retrieval should succeed")
	assert.Contains(t, items, "posts", "completion candidates should contain posts")
}

// TestCompletionCmd_Run verifies that completion command works
func TestCompletionCmd_Run(t *testing.T) {
	// This test only verifies basic parsing since actual execution is difficult
	var cli CLI
	parser, err := NewParserForTest(&cli)
	require.NoError(t, err)

	// Verify that completion command can be parsed
	testCases := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range testCases {
		t.Run(shell, func(t *testing.T) {
			_, err = parser.Parse([]string{"completion", shell})
			require.NoError(t, err, "completion %s command should be parseable", shell)
		})
	}
}

// NewParserForTest creates a parser for testing
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
