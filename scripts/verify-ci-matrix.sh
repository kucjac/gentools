#!/usr/bin/env bash
set -euo pipefail

workflow=${1:-.github/workflows/ci.yml}

require() {
  if ! grep -Fq -- "$1" "$workflow"; then
    printf 'missing required workflow content: %s\n' "$1" >&2
    exit 1
  fi
}

require 'actions/checkout@v6'
require 'actions/setup-go@v6'
require 'golangci/golangci-lint-action@v9'
require '"1.27.1"'
require 'ubuntu-latest'
require 'macos-latest'
require 'windows-latest'
require 'Root module dependency consistency'
require 'Root module formatting'
require 'Root module vet'
require 'Root module tests'
require 'Nested integration module tests'
require 'Go 1.27.1 / ubuntu-latest / lint'
