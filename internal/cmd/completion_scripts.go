/**
 * completion_scripts.go
 * シェル補完スクリプト
 *
 * 各シェル（bash/zsh/fish/powershell）用の補完スクリプトを提供します。
 */

package cmd

import "fmt"

// completionScript は指定されたシェル用の補完スクリプトを返します
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

// bashCompletionScript はbash用の補完スクリプトを返します
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

// zshCompletionScript はzsh用の補完スクリプトを返します
func zshCompletionScript() string {
	return `#compdef gho

autoload -Uz bashcompinit
bashcompinit
` + bashCompletionScript()
}

// fishCompletionScript はfish用の補完スクリプトを返します
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

// powerShellCompletionScript はpowershell用の補完スクリプトを返します
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
