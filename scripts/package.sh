#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
ZIP="octobackup-$(date +%Y%m%d).zip"
rm -f "$ZIP"
zip -r "$ZIP" . -x "*.git*"
echo "Created $ZIP"
