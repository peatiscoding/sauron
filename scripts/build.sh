#!/usr/bin/env bash
# Build Sauron: web → Go binary (single executable)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> Building web frontend…"
cd "$ROOT/web"
npm install
npm run build
echo "    → server/static/ populated"

echo "==> Building Go server…"
cd "$ROOT/server"
go mod tidy
go build -o "$ROOT/sauron" .
echo "    → $ROOT/sauron"

echo ""
echo "Done. Start with:"
echo "  $ROOT/sauron"
echo ""
echo "Then open: http://localhost:6905/focus"
