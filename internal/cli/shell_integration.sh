cdbm() {
  case "$1" in
  add | list | delete | edit | help | "")
    command cdbm "$@"
    ;;
  *)
    eval "$(command cdbm "$1")"
    ;;
  esac
}
