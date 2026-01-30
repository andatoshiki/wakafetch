#!/bin/sh
# Install wakafetch from GitHub releases.
# Usage: curl -fsSL https://raw.githubusercontent.com/andatoshiki/wakafetch/master/scripts/install.sh | sh
#   Install latest:  curl ... | sh
#   Install version: VERSION=v2.1.1 curl ... | sh   (or INSTALL_DIR=~/bin VERSION=v1.0.0 curl ... | sh)
set -e

REPO="andatoshiki/wakafetch"
INSTALL_DIR="${INSTALL_DIR:-}"
VERSION="${VERSION:-latest}"

# Detect OS
os=$(uname -s)
case "$os" in
  Darwin)  platform="darwin" ;;
  Linux)   platform="linux" ;;
  FreeBSD) platform="freebsd" ;;
  OpenBSD) platform="openbsd" ;;
  NetBSD)  platform="netbsd" ;;
  *)
    echo "error: unsupported OS: $os"
    exit 1
    ;;
esac

# Detect arch
arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  arch="x86-64" ;;
  aarch64|arm64) arch="aarch64" ;;
  i386|i686)     arch="i386" ;;
  armv7l)        arch="armv7" ;;
  armv6l)        arch="armv6" ;;
  *)
    echo "error: unsupported arch: $arch"
    exit 1
    ;;
esac

# Choose install directory
if [ -n "$INSTALL_DIR" ]; then
  bindir="$INSTALL_DIR"
else
  if [ -w "/usr/local/bin" ] 2>/dev/null; then
    bindir="/usr/local/bin"
  else
    bindir="${HOME}/.local/bin"
  fi
fi

mkdir -p "$bindir"

# Resolve version (tag)
if [ "$VERSION" = "latest" ]; then
  tag=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
  [ -z "$tag" ] && { echo "error: could not get latest release tag"; exit 1; }
else
  tag="$VERSION"
  case "$tag" in v*) ;; *) tag="v${tag}" ;; esac
fi

# Asset name: wakafetch-{platform}-{arch}-{tag}.tar.gz (zip on Windows)
suffix="tar.gz"
asset_name="wakafetch-${platform}-${arch}-${tag}.${suffix}"
download_url="https://github.com/${REPO}/releases/download/${tag}/${asset_name}"

echo "Installing wakafetch ${tag} (${platform}-${arch}) to ${bindir}"

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

if ! curl -sSfL -o "${tmpdir}/archive.${suffix}" "$download_url"; then
  echo "error: download failed for $download_url"
  echo "Check https://github.com/${REPO}/releases for available builds."
  exit 1
fi

cd "$tmpdir"
tar -xzf "archive.${suffix}"
chmod +x wakafetch
mv wakafetch "${bindir}/wakafetch"

echo "Installed: ${bindir}/wakafetch"
if ! command -v wakafetch >/dev/null 2>&1; then
  echo "Add ${bindir} to your PATH if it is not already."
fi
