#!/usr/bin/env bash
set -euo pipefail
src=${1:-/opt/netupdown}
dst=${2:-/backup/netupdown}
stamp=$(date +%Y%m%d-%H%M)
mkdir -p "$dst"
sqlite3 "$src/data/netupdown.db" ".backup '$dst/db-$stamp.db'"
tar -czf "$dst/data-$stamp.tar.gz" -C "$src" config.yaml data/secret data/uploads data/themes
find "$dst" -type f -mtime +7 -delete

