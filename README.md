# 📌 tack

**Sticky context for your terminal.** Pin notes, TODOs, and warnings to directories. Your filesystem becomes a spatial notebook.

## Install

```bash
go install github.com/user/tack@latest
```

Or build from source:

```bash
git clone https://github.com/user/tack
cd tack
go build -o tack .
```

## Usage

```bash
# Pin a note to the current directory
tack pin "remember: API key rotates monthly"

# Add a TODO
tack todo "fix flaky test in auth_test.go"

# Pin a warning (shows bold + colored)
tack warn "DO NOT deploy from this branch"

# See everything pinned here
tack

# Mark a TODO as done
tack done a3f1

# Remove a tack
tack rm a3f1

# Bird's-eye view of all tacks everywhere
tack board

# Show only directories with stale TODOs (7+ days old)
tack board --stale

# Adjust stale threshold
tack board --stale --days 14

# JSON output (for piping)
tack --json
tack board --json
```

## Shell Hook (auto-display on `cd`)

Get a summary every time you enter a directory:

```bash
# Print the hook for your shell
tack hook zsh   # or bash, fish

# Add it to your config
tack hook zsh >> ~/.zshrc
```

After restarting your shell, you'll see something like this when entering a directory with tacks:

```
📋 📌 2 notes · ☐ 1 todo · ⚠️ 1 warning
```

## Storage

All data lives in `~/.tack/tacks.json` — a single human-readable JSON file. Back it up, sync it, or edit it by hand. No database, no cloud, no accounts.

## Philosophy

- **Spatial memory**: context lives where the work lives
- **Zero friction**: one command to add, bare `tack` to see
- **Power-user friendly**: `--json` on everything, short IDs, shell hooks
- **Playful**: emoji, colors, and a personality — but never in the way

## License

MIT
