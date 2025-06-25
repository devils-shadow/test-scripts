#!/usr/bin/env bash
set -euo pipefail

BASE_URL="https://raw.githubusercontent.com/livekit/client-sdk-js/main/examples/media"

for file in sample720p.ivf sample.ogg; do
  if [[ -f $file ]]; then
    echo "$file already exists, skipping"
    continue
  fi
  echo "Downloading $file..."
  curl -fL "$BASE_URL/$file" -o "$file"

done
