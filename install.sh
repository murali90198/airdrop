#!/bin/bash
set -e

BIN_NAME="airdrop"
TARGET_DIR="/usr/local/sbin"
TARGET_PATH="$TARGET_DIR/$BIN_NAME"

echo "Installing $BIN_NAME to $TARGET_DIR ..."

# check root rights or required sudo
if [[ $EUID -ne 0 ]]; then
  echo "This script requires root privileges. Please run with sudo."
  exit 1
fi

# create directory
mkdir -p "$TARGET_DIR"

# copy
cp -f "./$BIN_NAME" "$TARGET_PATH"

# make rights for running +x
chmod 755 "$TARGET_PATH"

echo "Installed successfully to $TARGET_PATH"
