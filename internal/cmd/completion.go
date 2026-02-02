/**
 * completion.go
 * Shell completion functionality
 *
 * Generates completion scripts for bash/zsh/fish/powershell.
 */

package cmd

import (
	"context"
	"fmt"
	"os"
)

// CompletionCmd is the command to generate shell completion scripts
type CompletionCmd struct {
	Shell string `arg:"" name:"shell" help:"Shell (bash|zsh|fish|powershell)" enum:"bash,zsh,fish,powershell"`
}

// Run executes the completion command
func (c *CompletionCmd) Run(_ context.Context) error {
	script, err := completionScript(c.Shell)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, script)
	return err
}

// CompletionInternalCmd is a hidden command that provides completion candidates
type CompletionInternalCmd struct {
	Cword int      `name:"cword" help:"Index of the current word" default:"-1"`
	Words []string `arg:"" optional:"" name:"words" help:"Words to complete"`
}

// Run executes the __complete command
func (c *CompletionInternalCmd) Run(_ context.Context) error {
	items, err := completeWords(c.Cword, c.Words)
	if err != nil {
		return err
	}
	for _, item := range items {
		if _, err := fmt.Fprintln(os.Stdout, item); err != nil {
			return err
		}
	}
	return nil
}
