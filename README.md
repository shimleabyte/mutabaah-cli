<div align="center">

# mutabaah

**A compact CLI for daily activities and short notes.**

Go · Local JSON storage · Standard library only

<img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
<img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />

</div>

---

Mutabaah organizes its commands by domain: **acts** tracks activities and their session targets, while **notes** stores small text notes that can be reviewed quickly from the terminal.

## Features

- Add activities with configurable session targets and deadlines.
- Track activity progress with a terminal progress bar.
- Manage short notes containing one to five logical lines.
- List notes in a compact table with wrapped text, newest first.
- Keep note storage separate from activity storage.
- Preserve notes IDs when an entry is deleted.
- Use local JSON storage without external dependencies.

## Requirements

- [Go](https://go.dev/dl/) **1.26.5 or newer**, as declared in go.mod.
- No external dependencies.

## Getting started

Clone the repository:

```bash
git clone https://github.com/shimleabyte/mutabaah-cli.git
cd mutabaah-cli
```

Build the shorter executable:

```bash
go build -o mutabaah .
./mutabaah help
```

Or run from the source directory:

```bash
go run . help
```

## Command structure

### Activities

| Command | Purpose |
|---|---|
| `mutabaah acts list` | List unfinished activities |
| `mutabaah acts list --done` | List completed activities |
| `mutabaah acts list --all` | List all activities |
| `mutabaah acts tambah "Activity title" --sesi 5 --periode 10` | Add an activity with a target |
| `mutabaah acts catat <id>` | Record one session |
| `mutabaah acts sudah <id>` | Mark an activity complete manually |
| `mutabaah acts ubah <id> "New title"` | Rename an activity |
| `mutabaah acts hapus <id>` | Delete an activity |

### Notes

| Command | Purpose |
|---|---|
| `mutabaah notes add "line 1" ["line 2" ...]` | Add a note with one to five lines |
| `mutabaah notes list` | List notes in a compact table |
| `mutabaah notes show <id>` | Display a note in full |
| `mutabaah notes edit <id> "line 1" ["line 2" ...]` | Replace a note's content |
| `mutabaah notes delete <id>` | Delete a note |
| `mutabaah notes help` | Show notes help |

Each quoted text argument becomes one logical line. Notes require one to five non-empty lines; blank lines are rejected.

## Notes workflow

Add a short, one-line note:

```bash
./mutabaah notes add "Review error handling before release"
```

Add a note with multiple lines. Each quoted argument becomes one line in the stored note:

```bash
./mutabaah notes add \
  "Migration must be idempotent" \
  "Do not delete the source file automatically" \
  "Add a regression test"
```

List notes:

```bash
./mutabaah notes list
```

Example output:

```text
ID  DATE   NOTE
12  23 Oct Migration must be idempotent. Do not delete
          the source file before the new data is verified.
11  22 Oct Progress should remain consistent with session
          targets. A completed item must not contradict
          its own target.
```

The NOTE column wraps at a fixed width of 48 Unicode code points. Wrapping is for display only: it does not add lines to the saved note or change its content. Notes are ordered newest first.

Show or edit one note:

```bash
./mutabaah notes show 12
./mutabaah notes edit 12 \
  "Migration must be idempotent" \
  "Validate the destination before cleanup"
```

Editing replaces all existing note lines. It updates updatedAt while preserving createdAt.

## Notes data model

Each note contains an ID, one to five logical text lines, a creation timestamp, and an update timestamp. The collection also persists nextId so removing the highest-numbered note does not cause that ID to be reused.

Notes are stored independently at:

```text
~/.mtcli/notes.json
```

Activity data remains in its own file:

```text
~/.mtcli/daftaraktivitas.json
```

If notes.json is missing, the notes list starts empty. Invalid JSON or invalid note records return an error; the CLI does not silently replace a corrupt notes file with an empty collection.

Example JSON:

```json
{
  "nextId": 13,
  "notes": [
    {
      "id": 12,
      "lines": [
        "Migration must be idempotent.",
        "Do not delete the source file automatically."
      ],
      "createdAt": "2026-10-23T08:30:00Z",
      "updatedAt": "2026-10-23T08:30:00Z"
    }
  ]
}
```

## Project structure

| File | Responsibility |
|---|---|
| main.go | Root and domain command routing, activity table, activity commands |
| mutabaah.go | Activity model, session progress, completion, and deadlines |
| inventory.go | Activity JSON persistence and legacy migration |
| notes.go | Note model, one-to-five-line validation, CRUD operations, ID sequencing |
| notes_store.go | Separate notes storage and JSON validation |
| notes_cli.go | Notes commands, detail view, wrapping, and list table |
| notes_test.go | Notes validation, storage, CRUD, wrapping, and routing tests |
| table_test.go | Activity progress table and root routing tests |

## Current boundaries

- Notes are plain text, without tags, categories, search, attachments, or rich-text formatting.
- A note has a maximum of five logical lines; visual wrapping does not count toward that limit.
- Notes and activities are separate JSON files. They are local to the current user account and do not sync across devices.
- Activity periods remain one-time deadlines; recurring activity cycles are not implemented.

## License

Released under the [MIT License](LICENSE).
