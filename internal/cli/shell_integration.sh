cdbm() {
  case "$1" in
  add | list | remove | update | help | "")
    command cdbm "$@"
    ;;
  *)
    eval "$(command cdbm "$1")"
    ;;
  esac
}
