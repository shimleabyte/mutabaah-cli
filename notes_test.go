package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestValidateNoteLines(t *testing.T) {
	valid := [][]string{
		{"one"},
		{"one", "two"},
		{"one", "two", "three", "four", "five"},
	}
	for _, lines := range valid {
		got, err := ValidateNoteLines(lines)
		if err != nil {
			t.Errorf("ValidateNoteLines(%q): %v", lines, err)
		}
		if len(got) != len(lines) {
			t.Errorf("got %d lines for input %d", len(got), len(lines))
		}
	}

	invalid := [][]string{
		nil,
		{},
		{"   "},
		{"one", ""},
		{"embedded\nnewline"},
		{"1", "2", "3", "4", "5", "6"},
	}
	for _, lines := range invalid {
		if _, err := ValidateNoteLines(lines); err == nil {
			t.Errorf("ValidateNoteLines(%q) expected an error", lines)
		}
	}
}

func TestNotesCRUDAndIDDoesNotMoveBackAfterDelete(t *testing.T) {
	collection := EmptyNotesCollection()
	firstTime := time.Date(2026, 10, 20, 10, 0, 0, 0, time.UTC)
	secondTime := firstTime.Add(time.Hour)
	updatedTime := secondTime.Add(time.Hour)

	first, err := collection.Add([]string{"first"}, firstTime)
	if err != nil {
		t.Fatal(err)
	}
	second, err := collection.Add([]string{"second", "second line"}, secondTime)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != 1 || second.ID != 2 {
		t.Fatalf("unexpected IDs: first=%d second=%d", first.ID, second.ID)
	}

	edited, err := collection.Edit(first.ID, []string{"updated", "second line"}, updatedTime)
	if err != nil {
		t.Fatal(err)
	}
	if !edited.UpdatedAt.Equal(updatedTime) || !edited.CreatedAt.Equal(firstTime) {
		t.Fatalf("timestamps after edit are wrong: %#v", edited)
	}

	if err := collection.Delete(second.ID); err != nil {
		t.Fatal(err)
	}
	third, err := collection.Add([]string{"third"}, updatedTime.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if third.ID != 3 {
		t.Fatalf("ID after deleting highest ID = %d, want 3", third.ID)
	}
	if _, err := collection.Get(second.ID); err == nil {
		t.Fatal("deleted note should not be found")
	}
}

func TestNotesNewestFirstDoesNotMutateStoredOrder(t *testing.T) {
	collection := EmptyNotesCollection()
	now := time.Now()
	older, err := collection.Add([]string{"older"}, now)
	if err != nil {
		t.Fatal(err)
	}
	newer, err := collection.Add([]string{"newer"}, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	sorted := collection.NewestFirst()
	if sorted[0].ID != newer.ID || sorted[1].ID != older.ID {
		t.Fatalf("unexpected order: %#v", sorted)
	}
	if collection.Notes[0].ID != older.ID {
		t.Fatal("NewestFirst must not mutate the stored collection order")
	}
}

func TestNotesStoreRoundTripPreservesNextID(t *testing.T) {
	path := filepath.Join(t.TempDir(), appDir, notesDataFile)
	collection := EmptyNotesCollection()
	created, err := collection.Add([]string{"keep this ID"}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := collection.Delete(created.ID); err != nil {
		t.Fatal(err)
	}
	if collection.NextID != 2 {
		t.Fatalf("NextID before save = %d, want 2", collection.NextID)
	}
	if err := SaveNotes(path, collection); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadNotes(path)
	if err != nil {
		t.Fatal(err)
	}
	next, err := loaded.Add([]string{"new entry"}, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if next.ID != 2 {
		t.Fatalf("ID after storage round trip = %d, want 2", next.ID)
	}
}

func TestNotesStoreRejectsCorruptJSONWithoutOverwritingIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), notesDataFile)
	const corrupt = "{invalid JSON"
	if err := os.WriteFile(path, []byte(corrupt), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadNotes(path); err == nil {
		t.Fatal("LoadNotes should reject corrupt JSON")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != corrupt {
		t.Fatalf("corrupt file was modified: %q", got)
	}
}

func TestWrapNoteLineKeepsRunesAndRespectsWidth(t *testing.T) {
	line := "A long note about 日本語 and a surprisinglylongwordthatmustbewrapped"
	wrapped := wrapNoteLine(line, 12)
	if len(wrapped) < 2 {
		t.Fatalf("expected multiple wrapped lines, got %#v", wrapped)
	}
	for _, part := range wrapped {
		if got := len([]rune(part)); got > 12 {
			t.Errorf("wrapped line has %d runes, over width: %q", got, part)
		}
	}
	joined := strings.Join(wrapped, " ")
	if !strings.Contains(joined, "日本語") {
		t.Fatalf("wrapping dropped Unicode text: %q", joined)
	}
	longWord := "surprisinglylongwordthatmustbewrapped"
	longWordWrapped := wrapNoteLine(longWord, 12)
	if got := strings.Join(longWordWrapped, ""); got != longWord {
		t.Fatalf("wrapping changed a long word: %q", got)
	}
}

func TestPrintNotesListFramesCardsAndShowsNewestFirst(t *testing.T) {
	now := time.Date(2026, 10, 23, 15, 0, 0, 0, time.UTC)
	notes := []Note{
		{ID: 1, Lines: []string{"older note"}, CreatedAt: now.Add(-time.Hour)},
		{ID: 2, Lines: []string{"a short note that should wrap to another visual line in the compact table"}, CreatedAt: now},
	}

	var output bytes.Buffer
	printNotesList(&output, []Note{notes[1], notes[0]}, 18)
	got := output.String()
	if !strings.Contains(got, "NOTES  /  2 entries") {
		t.Fatalf("missing list heading: %s", got)
	}
	if !strings.Contains(got, "┌─ #2") || !strings.Contains(got, "└") {
		t.Fatalf("notes should be rendered as framed cards: %s", got)
	}
	if strings.Index(got, "2") > strings.Index(got, "1") {
		t.Fatalf("notes are not displayed newest first: %s", got)
	}
	if !strings.Contains(got, "another") || !strings.Contains(got, "visual") {
		t.Fatalf("wrapped table lost note text: %s", got)
	}
}

func withNotesHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	oldHome, hadHome := os.LookupEnv("HOME")
	oldProfile, hadProfile := os.LookupEnv("USERPROFILE")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("USERPROFILE", home); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadHome {
			_ = os.Setenv("HOME", oldHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
		if hadProfile {
			_ = os.Setenv("USERPROFILE", oldProfile)
		} else {
			_ = os.Unsetenv("USERPROFILE")
		}
	})
	return home
}

func TestRunNotesCRUDAndIndependentStore(t *testing.T) {
	home := withNotesHome(t)
	var output bytes.Buffer

	if err := runNotes([]string{"add", "first line", "second line"}, &output); err != nil {
		t.Fatal(err)
	}
	if err := runNotes([]string{"add", "another note"}, &output); err != nil {
		t.Fatal(err)
	}

	notesPath := filepath.Join(home, appDir, notesDataFile)
	actsPath, err := DataPath()
	if err != nil {
		t.Fatal(err)
	}
	if notesPath == actsPath {
		t.Fatal("notes and acts must not share a storage file")
	}

	output.Reset()
	if err := runNotes([]string{"list"}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "first line") || !strings.Contains(output.String(), "another note") {
		t.Fatalf("list output missing notes: %s", output.String())
	}

	output.Reset()
	if err := runNotes([]string{"show", "1"}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "first line") || !strings.Contains(output.String(), "second line") {
		t.Fatalf("show output missing logical note lines: %s", output.String())
	}

	output.Reset()
	if err := runNotes([]string{"edit", "1", "edited line"}, &output); err != nil {
		t.Fatal(err)
	}
	note, err := LoadNotes(notesPath)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := note.Get(1); err != nil || got.Lines[0] != "edited line" {
		t.Fatalf("edit did not persist: note=%#v err=%v", got, err)
	}

	if err := runNotes([]string{"delete", "1"}, &output); err != nil {
		t.Fatal(err)
	}
	afterDelete, err := LoadNotes(notesPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := afterDelete.Get(1); err == nil {
		t.Fatal("deleted note should not remain in persisted store")
	}
	if _, err := note.Get(1); err != nil {
		t.Fatalf("expected earlier in-memory fixture copy to remain unchanged: %v", err)
	}
}

func TestNotesHelpDoesNotTouchCorruptStore(t *testing.T) {
	home := withNotesHome(t)
	path := filepath.Join(home, appDir, notesDataFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	const corrupt = "do not change me"
	if err := os.WriteFile(path, []byte(corrupt), 0o600); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := runNotes([]string{"help"}, &output); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != corrupt {
		t.Fatalf("notes help modified storage: %q", after)
	}
}
