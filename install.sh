#!/bin/sh
set -e

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64)  ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

URL="https://github.com/jhgundersen/videogen/releases/latest/download/videogen-${OS}-${ARCH}"
DEST="${HOME}/.local/bin/videogen"

mkdir -p "$(dirname "$DEST")"
echo "Downloading videogen for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o "$DEST"
chmod +x "$DEST"
echo "Installed to $DEST"
