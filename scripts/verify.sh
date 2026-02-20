#!/usr/bin/env bash
set -euo pipefail

module="pixelc"

check_violation() {
  local pkg="$1"
  local dep="$2"

  if [[ "$dep" != ${module}/* ]]; then
    return 1
  fi

  if [[ "$pkg" == ${module}/cmd/* ]]; then
    if [[ "$dep" == ${module}/core/* || "$dep" == ${module}/pkg/* ]]; then
      return 1
    fi
    echo "boundary violation: $pkg imports $dep (cmd can only import core/* or pkg/*)"
    return 0
  fi

  if [[ "$pkg" == ${module}/core/* ]]; then
    if [[ "$dep" == ${module}/pkg/* || "$dep" == ${module}/internal/* || "$dep" == ${module}/core/* ]]; then
      return 1
    fi
    echo "boundary violation: $pkg imports $dep (core can only import core/*, pkg/*, internal/*)"
    return 0
  fi

  if [[ "$pkg" == ${module}/pkg/* ]]; then
    if [[ "$dep" == ${module}/core/* || "$dep" == ${module}/cmd/* ]]; then
      echo "boundary violation: $pkg imports $dep (pkg cannot import core/* or cmd/*)"
      return 0
    fi
    return 1
  fi

  if [[ "$pkg" == ${module}/internal/* ]]; then
    if [[ "$dep" == ${module}/cmd/* ]]; then
      echo "boundary violation: $pkg imports $dep (internal cannot import cmd/*)"
      return 0
    fi
    return 1
  fi

  return 1
}

violations=0
while IFS='|' read -r pkg imports; do
  [[ -z "$pkg" ]] && continue
  [[ "$pkg" != ${module}/* ]] && continue

  for dep in $imports; do
    if msg=$(check_violation "$pkg" "$dep"); then
      echo "$msg"
      violations=$((violations + 1))
    fi
  done
done < <(go list -f '{{.ImportPath}}|{{join .Imports " "}}' ./...)

if [[ "$violations" -gt 0 ]]; then
  echo "boundary checks failed: $violations violation(s)"
  exit 1
fi

echo "boundary checks passed"
