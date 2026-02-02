/**
 * completion_scripts.go
 * Shell completion scripts
 *
 * Provides completion scripts for each shell (bash/zsh/fish/powershell).
 */

package cmd

import "fmt"

// completionScript returns the completion script for the specified shell
func completionScript(shell string) (string, error) {
	switch shell {
	case "bash":
		return bashCompletionScript(), nil
	case "zsh":
		return zshCompletionScript(), nil
	case "fish":
		return fishCompletionScript(), nil
	case "powershell":
		return powerShellCompletionScript(), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// bashCompletionScript returns the completion script for bash
func bashCompletionScript() string {
	return `#!/usr/bin/env bash

_gho_complete() {
  local IFS=$'\n'
  local completions
  completions=$(gho __complete --cword "$COMP_CWORD" -- "${COMP_WORDS[@]}")
  COMPREPLY=()
  if [[ -n "$completions" ]]; then
    COMPREPLY=( $completions )
  fi
}

complete -F _gho_complete gho
`
}

// zshCompletionScript returns the completion script for zsh
func zshCompletionScript() string {
	return `#compdef gho

autoload -Uz bashcompinit
bashcompinit
` + bashCompletionScript()
}

// fishCompletionScript returns the completion script for fish
func fishCompletionScript() string {
	return `function __gho_complete
  set -l words (commandline -opc)
  set -l cur (commandline -ct)
  set -l cword (count $words)
  if test -n "$cur"
    set cword (math $cword - 1)
  end
  gho __complete --cword $cword -- $words
end

complete -c gho -f -a "(__gho_complete)"
`
}

// powerShellCompletionScript returns the completion script for powershell
func powerShellCompletionScript() string {
	return `Register-ArgumentCompleter -CommandName gho -ScriptBlock {
  param($commandName, $wordToComplete, $cursorPosition, $commandAst, $fakeBoundParameter)
  $elements = $commandAst.CommandElements | ForEach-Object { $_.ToString() }
  $cword = $elements.Count - 1
  $completions = gho __complete --cword $cword -- $elements
  foreach ($completion in $completions) {
    [System.Management.Automation.CompletionResult]::new($completion, $completion, 'ParameterValue', $completion)
  }
}
`
}
