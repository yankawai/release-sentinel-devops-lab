#!/usr/bin/env sh
set -eu

find deploy observability .github \( -name '*.yml' -o -name '*.yaml' \) | while read -r file; do
  case "$file" in
    *node_modules*) continue ;;
    deploy/helm/*/templates/*) continue ;;
  esac
  ruby -e 'require "yaml"; YAML.load_stream(File.read(ARGV[0])); puts "valid yaml: #{ARGV[0]}"' "$file"
done
