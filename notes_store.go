package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const notesDataFile = "notes.json"

// NotesPath intentionally returns a different file from DataPath so invalid or
// missing notes data cannot interfere with the acts store.
func NotesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("tidak dapat menentukan direktori home: %w", err)
	}
	return filepath.Join(home, appDir, notesDataFile), nil
}

func LoadNotes(path string) (NotesCollection, error) {
	collection := EmptyNotesCollection()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return collection, nil
		}
		return collection, err
	}

	if err := json.Unmarshal(data, &collection); err != nil {
		return EmptyNotesCollection(), fmt.Errorf("format JSON notes tidak valid: %w", err)
	}

	seen := make(map[int]struct{}, len(collection.Notes))
	for _, note := range collection.Notes {
		if note.ID < 1 {
			return EmptyNotesCollection(), fmt.Errorf("data notes rusak: ID catatan harus positif")
		}
		if _, exists := seen[note.ID]; exists {
			return EmptyNotesCollection(), fmt.Errorf("data notes rusak: ID %d duplikat", note.ID)
		}
		seen[note.ID] = struct{}{}
		if _, err := ValidateNoteLines(note.Lines); err != nil {
			return EmptyNotesCollection(), fmt.Errorf("data notes rusak pada ID %d: %w", note.ID, err)
		}
		if note.CreatedAt.IsZero() || note.UpdatedAt.IsZero() {
			return EmptyNotesCollection(), fmt.Errorf("data notes rusak pada ID %d: timestamp tidak lengkap", note.ID)
		}
	}

	collection.normalizeNextID()
	return collection, nil
}

func SaveNotes(path string, collection NotesCollection) error {
	collection.normalizeNextID()
	seen := make(map[int]struct{}, len(collection.Notes))
	for _, note := range collection.Notes {
		if note.ID < 1 {
			return fmt.Errorf("ID catatan harus positif")
		}
		if _, exists := seen[note.ID]; exists {
			return fmt.Errorf("ID catatan %d duplikat", note.ID)
		}
		seen[note.ID] = struct{}{}
		if _, err := ValidateNoteLines(note.Lines); err != nil {
			return fmt.Errorf("catatan ID %d tidak valid: %w", note.ID, err)
		}
		if note.CreatedAt.IsZero() || note.UpdatedAt.IsZero() {
			return fmt.Errorf("timestamp catatan %d tidak boleh kosong", note.ID)
		}
	}

	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	// os.WriteFile does not change the mode of an existing file.
	if err := os.Chmod(path, 0o600); err != nil {
		return err
	}
	return nil
}
