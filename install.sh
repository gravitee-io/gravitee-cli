#!/bin/sh
# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# gctl installer for macOS and Linux.
#
#   curl -fsSL https://raw.githubusercontent.com/gravitee-io/gravitee-cli/main/install.sh | sh
#
# Environment overrides:
#   GCTL_BIN  binary to install: "gctl" (default) or "gctl-ro"
#   GCTL_DIR  install directory (default: /usr/local/bin, falls back to ~/.local/bin)

set -eu

REPO="gravitee-io/gravitee-cli"
BIN="${GCTL_BIN:-gctl}"

case "$BIN" in
  gctl)    ASSET_PREFIX="gctl" ;;
  gctl-ro) ASSET_PREFIX="gctl_readonly" ;;
  *) echo "error: GCTL_BIN must be 'gctl' or 'gctl-ro', got '$BIN'" >&2; exit 1 ;;
esac

err() { echo "error: $*" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || err "curl is required"
command -v tar  >/dev/null 2>&1 || err "tar is required"

# Detect OS.
os=$(uname -s)
case "$os" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *) err "unsupported OS '$os' (on Windows, download a binary manually from the releases page)" ;;
esac

# Detect architecture.
arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) err "unsupported architecture '$arch'" ;;
esac

# Resolve the latest release tag.
api="https://api.github.com/repos/${REPO}/releases/latest"
TAG=$(curl -fsSL "$api" | grep '"tag_name"' | head -1 | cut -d'"' -f4)
[ -n "$TAG" ] || err "could not determine the latest release tag from $REPO"
VERSION="${TAG#v}"

ASSET="${ASSET_PREFIX}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

# Choose an install directory we can write to.
DIR="${GCTL_DIR:-/usr/local/bin}"
SUDO=""
if [ ! -d "$DIR" ] || [ ! -w "$DIR" ]; then
  if [ "$DIR" = "/usr/local/bin" ]; then
    if command -v sudo >/dev/null 2>&1; then
      SUDO="sudo"
    else
      DIR="${HOME}/.local/bin"
    fi
  fi
fi
mkdir -p "$DIR" 2>/dev/null || $SUDO mkdir -p "$DIR"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading ${BIN} ${TAG} (${OS}/${ARCH})..."
curl -fsSL -o "$tmp/$ASSET" "$URL" || err "download failed: $URL"

tar -xzf "$tmp/$ASSET" -C "$tmp"
[ -f "$tmp/$BIN" ] || err "binary '$BIN' not found in archive"
chmod +x "$tmp/$BIN"
$SUDO mv "$tmp/$BIN" "$DIR/$BIN"

echo "Installed $BIN to $DIR/$BIN"
case ":${PATH}:" in
  *":${DIR}:"*) ;;
  *) echo "note: $DIR is not on your PATH; add it to use '$BIN' directly" ;;
esac
"$DIR/$BIN" version || true
