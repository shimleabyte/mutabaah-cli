package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const notesTableTextWidth = 48

func runNotes(args []string, out io.Writer) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		printNotesHelp(out)
		return nil
	}

	command := args[0]
	var (
		id    int
		lines []string
		err   error
	)

	// Validate the command shape and note text before resolving or loading storage.
	switch command {
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("penggunaan: mutabaah notes list")
		}
	case "add":
		lines, err = ValidateNoteLines(args[1:])
		if err != nil {
			return err
		}
	case "show", "delete":
		if len(args) != 2 {
			return fmt.Errorf("penggunaan: mutabaah notes %s <id>", command)
		}
		id, err = parseNoteID(args[1])
		if err != nil {
			return err
		}
	case "edit":
		if len(args) < 3 {
			return fmt.Errorf("penggunaan: mutabaah notes edit <id> \"baris 1\" [\"baris 2\" ...]")
		}
		id, err = parseNoteID(args[1])
		if err != nil {
			return err
		}
		lines, err = ValidateNoteLines(args[2:])
		if err != nil {
			return err
		}
	default:
		printNotesHelp(out)
		return fmt.Errorf("subcommand notes %q tidak dikenal", command)
	}

	path, err := NotesPath()
	if err != nil {
		return err
	}
	collection, err := LoadNotes(path)
	if err != nil {
		return err
	}

	switch command {
	case "list":
		if len(collection.Notes) == 0 {
			fmt.Fprintln(out, "Belum ada catatan. Tambahkan dengan: mutabaah notes add \"isi catatan\"")
			return nil
		}
		printNotesTable(out, collection.NewestFirst(), notesTableTextWidth)
		return nil

	case "add":
		note, err := collection.Add(lines, time.Now())
		if err != nil {
			return err
		}
		if err := SaveNotes(path, collection); err != nil {
			return fmt.Errorf("gagal menyimpan catatan: %w", err)
		}
		fmt.Fprintf(out, "Catatan %d ditambahkan (%d baris).\n", note.ID, len(note.Lines))
		return nil

	case "show":
		note, err := collection.Get(id)
		if err != nil {
			return err
		}
		printNoteDetail(out, note)
		return nil

	case "edit":
		note, err := collection.Edit(id, lines, time.Now())
		if err != nil {
			return err
		}
		if err := SaveNotes(path, collection); err != nil {
			return fmt.Errorf("gagal menyimpan perubahan catatan: %w", err)
		}
		fmt.Fprintf(out, "Catatan %d diperbarui (%d baris).\n", note.ID, len(note.Lines))
		return nil

	case "delete":
		if err := collection.Delete(id); err != nil {
			return err
		}
		if err := SaveNotes(path, collection); err != nil {
			return fmt.Errorf("gagal menghapus catatan dari penyimpanan: %w", err)
		}
		fmt.Fprintf(out, "Catatan %d dihapus.\n", id)
		return nil
	}

	return nil
}

func parseNoteID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("ID catatan harus berupa angka positif")
	}
	return id, nil
}

func printNotesHelp(out io.Writer) {
	fmt.Fprint(out, "mutabaah notes — catatan pendek 1–5 baris\n\n"+
		"Penggunaan:\n"+
		"  mutabaah notes add \"baris 1\" [\"baris 2\" ...]\n"+
		"  mutabaah notes list\n"+
		"  mutabaah notes show <id>\n"+
		"  mutabaah notes edit <id> \"baris 1\" [\"baris 2\" ...]\n"+
		"  mutabaah notes delete <id>\n"+
		"  mutabaah notes help\n\n"+
		"Setiap argumen teks menjadi satu baris logis. Catatan harus memiliki 1–5 baris.\n"+
		"Tabel list membungkus teks panjang tanpa mengubah isi yang disimpan.\n"+
		"Data notes disimpan terpisah dari data acts.\n")
}

func printNoteDetail(out io.Writer, note Note) {
	fmt.Fprintf(out, "Note #%d\n", note.ID)
	fmt.Fprintf(out, "Created: %s\n", note.CreatedAt.Format("02 Jan 2006 15:04 MST"))
	fmt.Fprintf(out, "Updated: %s\n", note.UpdatedAt.Format("02 Jan 2006 15:04 MST"))
	fmt.Fprintln(out, strings.Repeat("-", 48))
	for _, line := range note.Lines {
		fmt.Fprintln(out, line)
	}
}

func printNotesTable(out io.Writer, notes []Note, width int) {
	if width < 1 {
		width = notesTableTextWidth
	}
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tDATE\tNOTE")

	for _, note := range notes {
		displayLines := make([]string, 0, len(note.Lines))
		for _, line := range note.Lines {
			displayLines = append(displayLines, wrapNoteLine(line, width)...)
		}
		if len(displayLines) == 0 {
			displayLines = []string{""}
		}
		fmt.Fprintf(writer, "%d\t%s\t%s\n", note.ID, note.CreatedAt.Format("02 Jan"), displayLines[0])
		for _, line := range displayLines[1:] {
			fmt.Fprintf(writer, "\t\t%s\n", line)
		}
	}
	_ = writer.Flush()
}

// wrapNoteLine wraps one logical note line for display only. It counts Unicode
// code points rather than UTF-8 bytes and breaks unusually long words safely.
func wrapNoteLine(line string, width int) []string {
	if width < 1 {
		width = notesTableTextWidth
	}
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	result := make([]string, 0, 1)
	current := make([]rune, 0, width)
	for _, word := range words {
		remainingWord := []rune(word)
		for len(remainingWord) > 0 {
			space := 0
			if len(current) > 0 {
				space = 1
			}
			available := width - len(current) - space
			if available <= 0 {
				result = append(result, string(current))
				current = make([]rune, 0, width)
				continue
			}
			if len(remainingWord) <= available {
				if space == 1 {
					current = append(current, ' ')
				}
				current = append(current, remainingWord...)
				remainingWord = nil
				continue
			}
			if len(current) > 0 {
				result = append(result, string(current))
				current = make([]rune, 0, width)
				continue
			}
			result = append(result, string(remainingWord[:width]))
			remainingWord = remainingWord[width:]
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}
	return result
}
