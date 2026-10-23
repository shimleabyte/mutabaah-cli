<div align="center">

# Mutabaah CLI

**A local-first Go CLI for tracking activities and keeping short notes.**

Go · JSON persistence · No third-party runtime dependencies

<img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
<img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />

</div>

---

Mutabaah brings two independent workflows into one terminal application:

- **Activities (`acts`)** — track session targets, record progress, and manage completion and one-time deadlines.
- **Notes (`notes`)** — keep short text entries of one to five logical lines, displayed as timestamped journal cards.

The app uses Go's standard library and local JSON files. Activity data and notes are stored separately under the user's home directory; the app does not sync data to a cloud service.

## Requirements

- Go **1.26.5**.
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

During development, run commands as `go run . <domain> <command>`; after building, use the compiled executable instead.

## Commands

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

A new activity defaults to **3 sessions in 7 days**. The period is a one-time deadline calculated from the creation date; recurring cycles are not implemented. Activity states are `OPEN`, `LATE`, and `DONE`.

### Notes

| Command | Purpose |
|---|---|
| `mutabaah notes add "baris 1" ["baris 2" ...]` | Create a note with one to five lines |
| `mutabaah notes list` | Show notes newest first |
| `mutabaah notes show <id>` | View a note and its timestamps |
| `mutabaah notes edit <id> "baris 1" ["baris 2" ...]` | Replace a note's content |
| `mutabaah notes delete <id>` | Delete a note |
| `mutabaah notes help` | Show note help |

Each quoted text argument becomes one logical line. Lines must be non-empty; embedded newlines in a single argument are rejected. Editing replaces the complete note and updates its modification timestamp.

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

Long lines wrap for display only; the stored note text remains unchanged.

## Local data and privacy

Mutabaah stores its data under `~/.mtcli/`:

| File | Contents |
|---|---|
| `daftaraktivitas.json` | Activity records and progress |
| `notes.json` | Notes, timestamps, and the next note ID |

The app creates the data directory with restricted permissions and writes data files with owner-only permissions on supported systems. The two data stores are independent.

If `notes.json` is missing, the notes collection starts empty. Invalid JSON or invalid note records produce an error rather than silently replacing the file. The activity command can migrate the legacy `./daftaraktivitas.json` file to the user data directory when the new activity file does not exist.

**Back up `~/.mtcli/` yourself.** Mutabaah is local-only and does not provide device synchronization. Local data files are excluded from Git tracking.

## Development

Run the checks individually:

```bash
go test -race ./...
go vet ./...
go build -o mutabaah .
```

Or use the included Makefile:

```bash
make check
make build
make clean
```

`make check` validates Go formatting, runs race-enabled tests, runs `go vet`, and builds the CLI. GitHub Actions runs these checks for pushes to `main` and pull requests targeting `main`.

## Project structure

The code is split by responsibility: root command routing, activity domain/CLI/storage, and note domain/CLI/storage. This keeps argument parsing and terminal output separate from the JSON persistence layer.

```text
.
├── .editorconfig
├── .github/workflows/go.yml
├── .gitignore
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── main.go              # Root command routing
├── acts.go              # Activity models and domain operations
├── acts_cli.go          # Activity commands and terminal output
├── acts_store.go        # Activity JSON storage and legacy migration
├── acts_cli_test.go     # Activity CLI and storage tests
├── acts_table_test.go   # Activity table and root routing tests
├── acts_test.go         # Activity domain tests
├── notes.go             # Note models, validation, and CRUD operations
├── notes_cli.go         # Note commands and journal-card output
├── notes_store.go       # Note persistence and data validation
└── notes_test.go        # Note validation, storage, CRUD, and CLI tests
```

## Scope

Notes are plain text and do not currently support tags, search, attachments, or rich formatting. Activity deadlines do not repeat automatically.

## License

Released under the [MIT License](LICENSE).
