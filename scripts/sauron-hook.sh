#!/usr/bin/env bash
# Claude Code PostToolUse hook — push edited .md/.yaml files to Sauron.
#
# Add to ~/.claude/settings.json:
#
#   "hooks": {
#     "PostToolUse": [{
#       "matcher": "Edit|Write|MultiEdit",
#       "hooks": [{ "type": "command", "command": "/path/to/depict/scripts/sauron-hook.sh" }]
#     }]
#   }
#
# Claude sets CLAUDE_TOOL_INPUT (JSON) and CLAUDE_TOOL_NAME env vars.

set -euo pipefail

SAURON_URL="${SAURON_URL:-http://localhost:6905}"

# Extract file path from tool input (Edit/Write both use file_path)
FILE=$(echo "${CLAUDE_TOOL_INPUT:-}" | jq -r '.file_path // empty' 2>/dev/null)

[ -z "$FILE" ] && exit 0
[ -f "$FILE" ] || exit 0

EXT="${FILE##*.}"
case "$EXT" in
  md)           FILETYPE="markdown" ;;
  yaml|yml)     FILETYPE="yaml"     ;;
  json)         FILETYPE="json"     ;;
  *)            exit 0              ;;
esac

FILENAME="$(basename "$FILE")"
CONTENT="$(cat "$FILE")"
PAYLOAD="$(jq -n --arg c "$CONTENT" --arg t "$FILETYPE" --arg f "$FILENAME" \
  '{ content: $c, filetype: $t, filename: $f }')"

curl -s -X POST "$SAURON_URL/api/focus" \
  -H "Content-Type: application/json" \
  --data-raw "$PAYLOAD" \
  --max-time 2 \
  --connect-timeout 1 \
  || true   # never block Claude if server is down
