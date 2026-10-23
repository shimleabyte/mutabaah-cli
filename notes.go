package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const maxNoteLines = 5

// Note stores a short note as logical lines. Terminal wrapping is presentation
// only and never changes Lines.
type Note struct {
	ID        int       `json:"id"`
	Lines     []string  `json:"lines"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NotesCollection owns notes and the next ID. NextID is persisted so deleting
// the newest note does not cause its ID to be reused.
type NotesCollection struct {
	NextID int    `json:"nextId"`
	Notes  []Note `json:"notes"`
}

func EmptyNotesCollection() NotesCollection {
	return NotesCollection{
		NextID: 1,
		Notes:  []Note{},
	}
}

// ValidateNoteLines trims surrounding whitespace and enforces one to five
// non-empty logical lines.
func ValidateNoteLines(lines []string) ([]string, error) {
	if len(lines) < 1 || len(lines) > maxNoteLines {
		return nil, fmt.Errorf("catatan harus berisi 1–%d baris; diterima %d", maxNoteLines, len(lines))
	}

	validated := make([]string, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			return nil, fmt.Errorf("baris %d kosong; setiap baris harus berisi teks", i+1)
		}
		validated[i] = line
	}
	return validated, nil
}

func (collection *NotesCollection) Add(lines []string, now time.Time) (Note, error) {
	validated, err := ValidateNoteLines(lines)
	if err != nil {
		return Note{}, err
	}
	collection.normalizeNextID()

	note := Note{
		ID:        collection.NextID,
		Lines:     validated,
		CreatedAt: now,
		UpdatedAt: now,
	}
	collection.NextID++
	collection.Notes = append(collection.Notes, note)
	return note, nil
}

func (collection *NotesCollection) Get(id int) (Note, error) {
	for _, note := range collection.Notes {
		if note.ID == id {
			return note, nil
		}
	}
	return Note{}, fmt.Errorf("catatan dengan ID %d tidak ditemukan", id)
}

func (collection *NotesCollection) Edit(id int, lines []string, now time.Time) (Note, error) {
	validated, err := ValidateNoteLines(lines)
	if err != nil {
		return Note{}, err
	}
	for i := range collection.Notes {
		if collection.Notes[i].ID == id {
			collection.Notes[i].Lines = validated
			collection.Notes[i].UpdatedAt = now
			return collection.Notes[i], nil
		}
	}
	return Note{}, fmt.Errorf("catatan dengan ID %d tidak ditemukan", id)
}

func (collection *NotesCollection) Delete(id int) error {
	for i := range collection.Notes {
		if collection.Notes[i].ID == id {
			collection.Notes = append(collection.Notes[:i], collection.Notes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("catatan dengan ID %d tidak ditemukan", id)
}

func (collection NotesCollection) NewestFirst() []Note {
	notes := append([]Note(nil), collection.Notes...)
	sort.Slice(notes, func(i, j int) bool {
		if notes[i].CreatedAt.Equal(notes[j].CreatedAt) {
			return notes[i].ID > notes[j].ID
		}
		return notes[i].CreatedAt.After(notes[j].CreatedAt)
	})
	return notes
}

func (collection NotesCollection) nextAvailableID() int {
	maxID := 0
	for _, note := range collection.Notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	return maxID + 1
}

func (collection *NotesCollection) normalizeNextID() {
	minimum := collection.nextAvailableID()
	if collection.NextID < minimum {
		collection.NextID = minimum
	}
	if collection.NextID < 1 {
		collection.NextID = 1
	}
	if collection.Notes == nil {
		collection.Notes = []Note{}
	}
}
