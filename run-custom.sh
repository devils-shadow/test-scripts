#!/usr/bin/env bash
set -euo pipefail

echo "🚀 LiveKit Load Test Launcher"
echo "============================="

read -p "📦 Number of rooms: " NUM_ROOMS
read -p "👥 Bots per room: " BOTS_PER_ROOM
read -p "⏱️ Duration (e.g., 60s or 10m): " DURATION
read -p "🐛 Enable debug logging? (y/n): " DEBUG_LOG
read -p "💾 Save logs to JSON file? (y/n): " SAVE_LOG

for f in sample720p.ivf sample.ogg; do
  if [[ ! -f $f ]]; then
    echo "Missing media file: $f" >&2
    exit 1
  fi
done

ARGS="-rooms $NUM_ROOMS -bots $BOTS_PER_ROOM -d $DURATION \
-url wss://meet.cst.ro -token-url https://meet.cst.ro/token \
-video sample720p.ivf -audio sample.ogg"

if [[ "$DEBUG_LOG" =~ ^[yY]$ ]]; then
  ARGS="$ARGS -debug=true"
  echo "🐛 Debug logging enabled"
else
  echo "🐛 Debug logging disabled"
fi

if [[ "$SAVE_LOG" =~ ^[yY]$ ]]; then
  ARGS="$ARGS -log=true"
  echo "🗒️ Logging enabled → output: last_run.json"
else
  echo "🗒️ Logging disabled"
fi

echo ""
echo "🎯 Launching $((NUM_ROOMS * BOTS_PER_ROOM)) bots across $NUM_ROOMS room(s)..."
echo "🕒 Each bot will stay connected for $DURATION"
echo ""

go run loadbot.go $ARGS
