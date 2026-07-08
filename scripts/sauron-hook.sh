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

# Detect git stage status for the file.
GIT_STATUS=""
if command -v git &>/dev/null; then
  GIT_RAW="$(git -C "$(dirname "$FILE")" status --porcelain -- "$FILE" 2>/dev/null | head -1)"
  if [ -n "$GIT_RAW" ]; then
    XY="${GIT_RAW:0:2}"
    case "$XY" in
      "??")               GIT_STATUS="untracked" ;;
      " M"|" D")          GIT_STATUS="modified"  ;;
      "MM"|"AM"|"RM")     GIT_STATUS="modified"  ;;
      *)                  GIT_STATUS="staged"    ;;
    esac
  fi
fi

PAYLOAD="$(jq -n --arg c "$CONTENT" --arg t "$FILETYPE" --arg f "$FILENAME" --arg g "$GIT_STATUS" --arg fp "$FILE" \
  '{ content: $c, filetype: $t, filename: $f, gitStatus: $g, filepath: $fp }')"

curl -s -X POST "$SAURON_URL/api/focus" \
  -H "Content-Type: application/json" \
  --data-raw "$PAYLOAD" \
  --max-time 2 \
  --connect-timeout 1 \
  || true   # never block Claude if server is down
