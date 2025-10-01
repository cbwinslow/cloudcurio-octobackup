#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
go build -o octobackup ./cmd/octobackup
echo "Built ./octobackup"
