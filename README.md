<div align="center">

<h1>mutabaah</h1>
<p><strong>CLI untuk mencatat aktivitas harian dengan target sesi dan tenggat.</strong></p>
<p>Go · Penyimpanan JSON lokal · Tanpa dependency eksternal</p>
<p>
  <img src="https://img.shields.io/badge/Go-1.26.5-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.26.5" />
  <img src="https://img.shields.io/badge/License-MIT-0B3D2E?style=flat-square" alt="MIT License" />
</p>

</div>

---

*Mutabaah* is a small command-line tool for tracking daily activities. The CLI is organized by domain: `acts` for activity tracking and `notes` as a planned, separate area for personal notes.

This project is also a way to learn Go through practical implementation: command routing, structs and slices, JSON persistence, input validation, progress tracking, and storage migration.

## Features

- Add activities with configurable session targets and deadlines.
- Record each completed session and show progress as a terminal progress bar.
- Automatically mark an activity complete when it reaches its session target.
- List open activities, completed activities, or all activities.
- Show overdue activities in the table.
- Update activity titles and delete activities.
- Store data locally in JSON and migrate the legacy data file when needed.
- Keep root and domain help commands independent of the activity data store.

## Requirements

- [Go](https://go.dev/dl/) **1.26.5 or newer**, as declared in `go.mod`.
- Standard library only; no external dependencies.

## Getting started

Clone the repository:

```bash
git clone https://github.com/shimleabyte/mutabaah-cli.git
cd mutabaah-cli
```

Run without building a binary:

```bash
go run . help
```

Or build the shorter `mutabaah` executable:

```bash
go build -o mutabaah .
./mutabaah help
```

On Windows:

```powershell
go build -o mutabaah.exe .
.\\mutabaah.exe help
```

## Command structure

The root command separates features by domain.

| Command | Purpose |
|---|---|
| `mutabaah acts list` | List activities |
| `mutabaah acts tambah "Activity title"` | Add an activity |
| `mutabaah acts catat <id>` | Record one session |
| `mutabaah acts sudah <id>` | Mark an activity complete manually |
| `mutabaah acts ubah <id> "New title"` | Rename an activity |
| `mutabaah acts hapus <id>` | Delete an activity |
| `mutabaah acts help` | Show activity command help |
| `mutabaah notes` | Notes module placeholder; not implemented yet |
| `mutabaah help` | Show root help |

`mutabaah notes` is reserved for a separate notes domain. Notes commands do not currently create, read, update, or delete notes.

## Activity targets

New activities use a default target of **3 sessions within 7 days**. Set a per-activity target with `--sesi` and `--periode`:

```bash
./mutabaah acts tambah "Tilawah satu juz" --sesi 5 --periode 10
```

- `--sesi` is the number of sessions needed to complete the activity.
- `--periode` is a one-time deadline in days, measured from when the activity was created.
- `catat <id>` increments progress by one session.
- Once progress reaches the target, the activity is marked complete.
- `sudah <id>` remains available as a manual completion shortcut.
- The deadline does not reset automatically. Recurring periods are not implemented yet.

## Listing activities

The default list shows only activities that are not complete.

```bash
./mutabaah acts list
./mutabaah acts list --done
./mutabaah acts list --all
```

Example output:

```text
ID  STATE  ACTIVITY           PROGRESS          DEADLINE
1   OPEN   Tilawah satu juz   [####------] 2/5  5d left
2   LATE   Murajaah           [##--------] 1/5  late 3d
3   DONE   Sholat Dhuha       [##########] 3/3  complete
```

| State | Meaning |
|---|---|
| `OPEN` | Activity is not complete and is within its deadline |
| `LATE` | Activity is not complete and its deadline has passed |
| `DONE` | Session target was reached or completion was marked manually |

Progress bars use `#` for recorded sessions and `-` for remaining sessions. `--done` lists only completed activities; `--all` includes every activity.

## Data storage

Activity data is stored in the user's home directory:

```text
~/.mtcli/daftaraktivitas.json
```

The CLI creates the directory and file when a change needs to be saved. Commands such as `mutabaah help` and `mutabaah acts help` return before locating, migrating, or loading stored activity data, so help remains available when the JSON file is invalid.

### Legacy migration

If the new data file does not exist and `daftaraktivitas.json` is found in the current working directory, the legacy file is copied to the new location. The original file is not removed automatically.

### JSON format

```json
{
  "daftarAktivitas": [
    {
      "id": 1,
      "judul": "Tilawah satu juz",
      "target": {
        "sesi": 5,
        "periode": 10
      },
      "proses": 2,
      "ketuntasan": false,
      "dibuatSaat": "2026-10-23T07:30:00Z"
    }
  ]
}
```

The target fields use lowercase JSON keys (`sesi` and `periode`). Existing data with the older `Sesi` and `Periode` keys remains readable.

## Project structure

| File | Responsibility |
|---|---|
| `main.go` | Root/domain routing, argument parsing, table output, and activity commands |
| `mutabaah.go` | Activity model, session progress, completion, and deadline logic |
| `inventory.go` | Data path, JSON persistence, legacy-data normalization, and migration |
| `main_test.go` | Argument parsing, JSON compatibility, and help-routing tests |
| `mutabaah_test.go` | Session progress, completion, and deadline tests |
| `table_test.go` | Progress-bar, table-rendering, and command-routing tests |

## Current design boundaries

- `acts` is implemented; `notes` is only a placeholder.
- A period is a one-time deadline, not a recurring cycle.
- IDs are derived from the last remaining activity ID plus one, so deleting the last activity can cause an ID to be reused.
- `sudah <id>` can mark an activity complete without recording each individual session.

## License

Released under the [MIT License](LICENSE).
