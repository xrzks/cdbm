cdbm() {
  if [[ "$1" == "add" ]] || [[ "$1" == "list" ]]; then
    command cdbm "$@"
  else
    eval "$(command cdbm "$1")"
  fi
}
