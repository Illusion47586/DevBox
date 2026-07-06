#!/usr/bin/env sh
set -eu

limit=400
tmp="${TMPDIR:-/tmp}/devbox-loc-files.$$"
trap 'rm -f "$tmp"' EXIT

find . -type f \
  \( -name '*.go' -o -name '*.sh' -o -name 'Dockerfile*' -o -name 'Makefile' \) \
  -not -path './.git/*' \
  -not -path './bin/*' \
  -not -path './dist/*' > "$tmp"

while IFS= read -r file; do
  lines=$(wc -l < "$file" | tr -d ' ')
  if [ "$lines" -gt "$limit" ]; then
    printf '%s has %s lines; limit is %s\n' "$file" "$lines" "$limit"
    exit 1
  fi
done < "$tmp"

exit 0
