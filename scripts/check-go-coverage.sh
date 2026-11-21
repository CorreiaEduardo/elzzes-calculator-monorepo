#!/usr/bin/env bash
# Enforces the per-package coverage floor from spec 004. The profile's directory
# must be the Go module root.
set -euo pipefail

profile=${1:-coverage.out}
threshold=${2:-80}

if [[ ! -f "$profile" ]]; then
  echo "coverage: profile not found: $profile" >&2
  exit 1
fi

cd "$(dirname "$profile")"
profile=$(basename "$profile")

missing=0
while IFS= read -r pkg; do
  if ! awk -v prefix="$pkg/" 'index($1, prefix) == 1 { found = 1; exit } END { exit !found }' "$profile"; then
    echo "coverage: no statements recorded for $pkg" >&2
    missing=1
  fi
done < <(go list ./internal/...)

awk -v threshold="$threshold" '
  NR == 1 { next }
  {
    file = $1
    sub(/:[0-9].*$/, "", file)
    pkg = file
    sub(/\/[^\/]+$/, "", pkg)
    statements[pkg] += $2
    if ($3 + 0 > 0) covered[pkg] += $2
  }
  END {
    failed = 0
    for (pkg in statements) {
      if (pkg !~ /(^|\/)internal(\/|$)/) continue
      if (statements[pkg] == 0) continue
      pct = 100 * covered[pkg] / statements[pkg]
      printf "%-58s %6.1f%%\n", pkg, pct
      if (pct < threshold) failed = 1
    }
    if (failed) {
      printf "coverage: one or more packages are below %.0f%%\n", threshold > "/dev/stderr"
      exit 1
    }
    printf "coverage: every internal package meets %.0f%%\n", threshold
  }
' "$profile"

if [[ "$missing" -ne 0 ]]; then
  echo "coverage: add tests for the packages listed above" >&2
  exit 1
fi
