
_envdo() {
  COMPREPLY=()
  local word="${COMP_WORDS[COMP_CWORD]}"

  if [ "$COMP_CWORD" -eq 1 ]; then
    COMPREPLY=( $(compgen -W "$(envdo ls --plain) help add ls --gpg-recipient --version --directory" -- "$word") )
  elif [[ $COMP_CWORD -eq 2 ]]; then
    case "${COMP_WORDS[1]}" in
      add)
        if [[ $COMP_CWORD -le 2 ]]; then
          COMPREPLY=($(compgen -W "$(envdo ls --plain)" -- ${word}));
          return
        fi
        ;;
    esac
    COMPREPLY=($(compgen -A function -abck -- "${word}"));
  fi
}

complete -o default -F _envdo envdo
