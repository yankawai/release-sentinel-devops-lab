#!/usr/bin/env sh
set -eu

find deploy observability .github -name '*.yml' -o -name '*.yaml' | while read -r file; do
  case "$file" in
    *node_modules*) continue ;;
  esac
  python3 - "$file" <<'PY'
import pathlib
import sys
try:
    import yaml
except Exception as exc:
    raise SystemExit(f"PyYAML is required for YAML validation: {exc}")

path = pathlib.Path(sys.argv[1])
with path.open("r", encoding="utf-8") as handle:
    list(yaml.safe_load_all(handle))
print(f"valid yaml: {path}")
PY
done
