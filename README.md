# windows-shell-cracker

Bubble Tea TUI for quickly opening Windows Explorer special folders via `shell:` URIs.

Commands come from the Windows shell cheat sheet (Shell Command column only), sorted alphabetically.

## Requirements

- Go 1.22+
- Windows (to actually launch Explorer views)

## Run

```bash
go run .
```

Build for Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o windows-shell-cracker.exe .
```

## Controls

| Input | Action |
| --- | --- |
| Mouse wheel | Scroll the tumbler one item |
| ↑ / ↓ | Move selection (stops at first/last item; no wrap) |
| Type in filter | Substring filter on command names |
| Enter | Run highlighted command (`start explorer <shell:…>`) |
| Esc | Clear filter and highlight the first command |
| Ctrl+C | Exit (saves last highlight) |

The last highlighted command is stored under your user config directory in `windows-shell-cracker/state.json`.

## Layout

13 lines: header (app name + filter), 11-row tumbler with the selection on line 6, and a Bubble Tea-style footer.
