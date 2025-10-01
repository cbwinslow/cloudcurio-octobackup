#!/usr/bin/env bash
set -euo pipefail
log(){ printf "[cloudcurio] %s\n" "$*"; }
OS="$(uname -s)"; log "OS: $OS"
if command -v apt >/dev/null 2>&1; then
  sudo apt update -y && sudo apt install -y curl git ca-certificates
elif command -v dnf >/dev/null 2>&1; then
  sudo dnf install -y curl git ca-certificates
elif command -v pacman >/dev/null 2>&1; then
  sudo pacman -Syu --noconfirm curl git ca-certificates
fi
DOTS="${DOTS:-https://github.com/youruser/dotfiles.git}"
log "Cloning dotfiles from $DOTS"
if [ ! -d "$HOME/.dotfiles" ]; then
  git clone --depth=1 "$DOTS" "$HOME/.dotfiles"
  (cd "$HOME/.dotfiles" && ./install.sh || true)
else
  (cd "$HOME/.dotfiles" && git pull)
fi
log "Bootstrap complete ✔"
