package main

import (
	"fmt"
)

// generate auto-completion as per https://cli.urfave.org/v2/#enabling
func generateAutocomplete(shell string) error {
	if shell == "zsh" {

		fmt.Println(`_cli_zsh_autocomplete() {
local -a opts
local cur
cur=${words[-1]}
if [[ "$cur" == "-"* ]]; then
	opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
else
	opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
fi

if [[ "${opts[1]}" != "" ]]; then
	_describe 'values' opts
else
	_files
fi
}

compdef _cli_zsh_autocomplete uses`)
	}
	if shell == "bash" {
		fmt.Println(`export PROG=uses
: ${PROG:=$(basename ${BASH_SOURCE})}

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete $PROG
unset PROG`)
	}
	return nil
}
