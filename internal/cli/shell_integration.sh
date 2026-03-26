cdbm() {
  case "$1" in
  add | list | delete | edit | help | init | "" | -*)
    command cdbm "$@"
    ;;
  *)
    output="$(command cdbm "$1")"
    if [[ "$output" =~ ^cd\ .+$ ]]; then
      eval "$output"
    else
      echo "Invalid output from cdbm" >&2
      return 1
    fi
    ;;
  esac
}
