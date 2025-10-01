# CloudCurio OctoBackup 🐙⚡

Neon Octopus Edition — a Bubble Tea TUI to stream Linux backups directly to your homelab over SSH.
No local temp files. Multiple strategies: dd, rsync, Borg, ZFS send, Btrfs send.

## Build
```bash
go build -o octobackup ./cmd/octobackup
./octobackup
```
