/**
 * completion.go
 * シェル補完機能
 *
 * bash/zsh/fish/powershell用の補完スクリプトを生成します。
 */

package cmd

import (
	"context"
	"fmt"
	"os"
)

// CompletionCmd はシェル補完スクリプトを生成するコマンドです
type CompletionCmd struct {
	Shell string `arg:"" name:"shell" help:"Shell (bash|zsh|fish|powershell)" enum:"bash,zsh,fish,powershell"`
}

// Run はcompletionコマンドを実行します
func (c *CompletionCmd) Run(_ context.Context) error {
	script, err := completionScript(c.Shell)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(os.Stdout, script)
	return err
}

// CompletionInternalCmd は補完候補を提供する隠しコマンドです
type CompletionInternalCmd struct {
	Cword int      `name:"cword" help:"Index of the current word" default:"-1"`
	Words []string `arg:"" optional:"" name:"words" help:"Words to complete"`
}

// Run は__completeコマンドを実行します
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
