/**
 * cmd_test.go
 * Test code for common command functionality
 *
 * Includes tests for flag aliases added in Phase 5.
 */

package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
)

// TestLimitAliases verifies that Limit flag aliases (--max, -n) work correctly
func TestLimitAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want int
	}{
		{"--limit flag", []string{"posts", "list", "--limit=10"}, 10},
		{"--max alias", []string{"posts", "list", "--max=10"}, 10},
		{"--n alias", []string{"posts", "list", "--n=10"}, 10},
		{"-l short flag", []string{"posts", "list", "-l", "10"}, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize CLI
			cli := &CLI{}

			// Create Kong parser
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // Don't exit during tests
			)
			if err != nil {
				t.Fatalf("failed to create Kong parser: %v", err)
			}

			// Parse command line
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("failed to parse command line: %v", err)
			}

			// Verify Limit field is set correctly
			if cli.Posts.List.Limit != tc.want {
				t.Errorf("Limit field not set correctly: got=%d, want=%d", cli.Posts.List.Limit, tc.want)
			}
		})
	}
}

// TestFilterAliases verifies that Filter flag aliases (--where, -w) work correctly
func TestFilterAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{"--filter flag", []string{"members", "list", "--filter=status:paid"}, "status:paid"},
		{"--where alias", []string{"members", "list", "--where=status:paid"}, "status:paid"},
		{"--w alias", []string{"members", "list", "--w=status:paid"}, "status:paid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize CLI
			cli := &CLI{}

			// Create Kong parser
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // Don't exit during tests
			)
			if err != nil {
				t.Fatalf("failed to create Kong parser: %v", err)
			}

			// Parse command line
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("failed to parse command line: %v", err)
			}

			// Verify Filter field is set correctly
			if cli.Members.List.Filter != tc.want {
				t.Errorf("Filter field not set correctly: got=%s, want=%s", cli.Members.List.Filter, tc.want)
			}
		})
	}
}

// TestPagesFlagAliases verifies that Limit aliases work correctly for Pages command
func TestPagesFlagAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want int
	}{
		{"--limit flag", []string{"pages", "list", "--limit=20"}, 20},
		{"--max alias", []string{"pages", "list", "--max=20"}, 20},
		{"--n alias", []string{"pages", "list", "--n=20"}, 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize CLI
			cli := &CLI{}

			// Create Kong parser
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // Don't exit during tests
			)
			if err != nil {
				t.Fatalf("failed to create Kong parser: %v", err)
			}

			// Parse command line
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("failed to parse command line: %v", err)
			}

			// Verify Limit field is set correctly
			if cli.Pages.List.Limit != tc.want {
				t.Errorf("Pages.List.Limit not set correctly: got=%d, want=%d", cli.Pages.List.Limit, tc.want)
			}
		})
	}
}

// TestTagsFlagAliases verifies that Limit aliases work correctly for Tags command
func TestTagsFlagAliases(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want int
	}{
		{"--limit flag", []string{"tags", "list", "--limit=30"}, 30},
		{"--max alias", []string{"tags", "list", "--max=30"}, 30},
		{"--n alias", []string{"tags", "list", "--n=30"}, 30},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize CLI
			cli := &CLI{}

			// Create Kong parser
			parser, err := kong.New(cli,
				kong.Name("gho"),
				kong.Exit(func(int) {}), // Don't exit during tests
			)
			if err != nil {
				t.Fatalf("failed to create Kong parser: %v", err)
			}

			// Parse command line
			_, err = parser.Parse(tc.args)
			if err != nil {
				t.Fatalf("failed to parse command line: %v", err)
			}

			// Verify Limit field is set correctly
			if cli.Tags.List.Limit != tc.want {
				t.Errorf("Tags.List.Limit not set correctly: got=%d, want=%d", cli.Tags.List.Limit, tc.want)
			}
		})
	}
}
