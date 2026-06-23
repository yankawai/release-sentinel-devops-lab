#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"

check() {
  path="$1"
  expected="$2"
  status="$(curl -sS -o /tmp/release-sentinel-smoke.json -w '%{http_code}' "$BASE_URL$path")"
  if [ "$status" != "$expected" ]; then
    echo "unexpected status for $path: got $status, expected $expected"
    cat /tmp/release-sentinel-smoke.json
    exit 1
  fi
  echo "ok $path $status"
}

check /healthz 200
check /readyz 200
check /version 200
check /work 200
check /metrics 200
