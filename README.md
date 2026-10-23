<div align="center">

# Mutabaah CLI

**A local-first CLI for tracking activities and keeping short notes.**

Go · JSON files · No third-party runtime dependencies

<img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
<img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />

</div>

---

Mutabaah keeps two small workflows in one terminal app:

- **Activities (`acts`)** — track session goals, progress, completion, and one-time deadlines.
- **Notes (`notes`)** — save short text entries of one to five logical lines.

Data stays on the local machine. The project uses Go's standard library and JSON files; it does not sync data to a cloud service.

## Requirements

- Go **1.26.5 or newer**.
- `make` is optional; direct Go commands work without it.

## Quick start

```bash
git clone https://github.com/shimleabyte/mutabaah-cli.git
cd mutabaah-cli

go run . help
go build -o mutabaah .
./mutabaah acts list
```

On Windows, build with `go build -o mutabaah.exe .` and run `mutabaah.exe`.

## Commands

Use `go run .` from the repository root during development, or replace it with the compiled `./mutabaah` executable.

### Activities

| Command | Purpose |
|---|---|
| `mutabaah acts list` | List unfinished activities |
| `mutabaah acts list --done` | List completed activities |
| `mutabaah acts list --all` | List all activities |
| `mutabaah acts tambah "Tilawah satu juz" --sesi 3 --periode 7` | Add an activity with session and day targets |
| `mutabaah acts catat <id>` | Record one session |
| `mutabaah acts sudah <id>` | Mark an activity complete manually |
| `mutabaah acts ubah <id> "Judul baru"` | Rename an activity |
| `mutabaah acts hapus <id>` | Delete an activity |
| `mutabaah acts help` | Show activity help |

A new activity defaults to **3 sessions in 7 days**. The period is a one-time deadline from the creation date; recurring cycles are not implemented. Activity states are `OPEN`, `LATE`, and `DONE`.

### Notes

| Command | Purpose |
|---|---|
| `mutabaah notes add "baris 1" ["baris 2" ...]` | Create a note with one to five lines |
| `mutabaah notes list` | Show notes newest first |
| `mutabaah notes show <id>` | View a note and its timestamps |
| `mutabaah notes edit <id> "baris 1" ["baris 2" ...]` | Replace a note's content |
| `mutabaah notes delete <id>` | Delete a note |
| `mutabaah notes help` | Show note help |

Each quoted text argument becomes one logical line. Lines must be non-empty; embedded newlines in a single argument are rejected. Editing replaces the whole note and updates its modification timestamp.

Example `notes list` output:

```text
NOTES  /  2 entries

┌─ #12 · 23 Oct 2026 08:30 ─────────────────────────────────────┐
│  Migration must be idempotent                                  │
│  Validate the destination before cleanup                      │
└────────────────────────────────────────────────────────────────┘

┌─ #11 · 22 Oct 2026 19:10 · edited ─────────────────────────────┐
│  Review error handling before release                          │
└────────────────────────────────────────────────────────────────┘
```

Long lines wrap for display only. Wrapping does not change the text stored in the note.

## Local data and privacy

Mutabaah stores data under `~/.mtcli/`:

| File | Contents |
|---|---|
| `daftaraktivitas.json` | Activity records and progress |
| `notes.json` | Notes, timestamps, and the next note ID |

The app creates the directory with restricted permissions and writes data files with owner-only permissions on supported systems. Activity and note data are independent, and neither is included in the repository.

If `notes.json` is missing, the notes collection starts empty. Invalid JSON or invalid note records produce an error rather than silently replacing the file. The activity command can migrate the legacy `./daftaraktivitas.json` file into the user data directory if the new file does not exist.

**Back up `~/.mtcli/` yourself.** Mutabaah is local-only and does not provide device synchronization.

## Development

Run the checks:

```bash
go test ./...
go vet ./...
go build -o mutabaah .
```

Or use the included Makefile:

```bash
make check
make build
make clean
```

GitHub Actions checks Go formatting, runs race-enabled tests and `go vet`, and builds the CLI for pull requests targeting `main` and pushes to `main`.

## Project layout

```text
.
├── main.go             # Entry point and root command routing
├── acts.go             # Activity models and domain operations
├── acts_cli.go         # Activity commands, argument parsing, and terminal output
├── acts_store.go       # Activity JSON storage and legacy migration
├── notes.go            # Note models, validation, and CRUD operations
├── notes_cli.go        # Note commands and journal-style terminal output
├── notes_store.go      # Note persistence and data validation
├── acts_test.go        # Activity model tests
├── acts_cli_test.go    # Activity CLI and storage tests
├── acts_table_test.go  # Activity table and root routing tests
├── notes_test.go       # Note validation, storage, CRUD, and CLI tests
├── Makefile            # Common development commands
└── .github/workflows/  # CI checks
```

## Scope

Notes are plain text and do not currently support tags, search, attachments, or rich formatting. Activity deadlines do not repeat automatically.

## License

Released under the [MIT License](LICENSE).
