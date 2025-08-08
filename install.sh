#!/bin/bash
set -e

BIN_NAME="airdrop"
TARGET_DIR="/usr/local/sbin"
TARGET_PATH="$TARGET_DIR/$BIN_NAME"

echo "Installing $BIN_NAME to $TARGET_DIR ..."

if [[ $EUID -ne 0 ]]; then
  echo "This script requires root privileges. Please run with sudo."
  exit 1
fi

# mkdir -p "$TARGET_DIR"

cp -f "./$BIN_NAME" "$TARGET_PATH"

chmod 755 "$TARGET_PATH"

xattr -d com.apple.quarantine "$TARGET_PATH" 2>/dev/null || true

echo "Installed successfully to $TARGET_PATH"
