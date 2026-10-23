package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestProgressBarShowsSessionRatio(t *testing.T) {
	if got := progressBar(2, 5, 10); got != "[####------]" {
		t.Fatalf("progressBar(2, 5, 10) = %q", got)
	}
	if got := progressBar(5, 5, 10); got != "[##########]" {
		t.Fatalf("completed progress bar = %q", got)
	}
	if got := progressBar(-1, 5, 10); got != "[----------]" {
		t.Fatalf("negative progress should clamp to zero, got %q", got)
	}
	if got := progressBar(8, 5, 10); got != "[##########]" {
		t.Fatalf("progress above target should clamp to full, got %q", got)
	}
}

func TestDeadlineLabelAndTableStates(t *testing.T) {
	now := time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC)
	activities := []Aktivitas{
		{
			Id:         1,
			Judul:      "Tilawah satu juz",
			Target:     UkurTarget{Sesi: 5, Periode: 7},
			Progress:   2,
			DibuatSaat: now.AddDate(0, 0, -2),
		},
		{
			Id:         2,
			Judul:      "Murajaah",
			Target:     UkurTarget{Sesi: 5, Periode: 7},
			Progress:   1,
			DibuatSaat: now.AddDate(0, 0, -10),
		},
		{
			Id:         3,
			Judul:      "Sholat Dhuha",
			Target:     UkurTarget{Sesi: 3, Periode: 7},
			Progress:   3,
			IsDone:     true,
			DibuatSaat: now.AddDate(0, 0, -2),
		},
	}

	var output bytes.Buffer
	shown := printActsTable(&output, activities, true, false, now)
	if shown != 3 {
		t.Fatalf("shown = %d, want 3", shown)
	}
	for _, want := range []string{
		"ID", "STATE", "ACTIVITY", "PROGRESS", "DEADLINE",
		"OPEN", "LATE", "DONE",
		"Tilawah satu juz", "[####------] 2/5", "5d left",
		"Murajaah", "late 3d",
		"Sholat Dhuha", "[##########] 3/3", "complete",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("table output does not contain %q:\n%s", want, output.String())
		}
	}
}

func TestPrintActsTableFiltersCompletedAndOpenActivities(t *testing.T) {
	now := time.Date(2026, 10, 23, 12, 0, 0, 0, time.UTC)
	activities := []Aktivitas{
		{Id: 1, Judul: "Open", Target: UkurTarget{Sesi: 2, Periode: 7}, DibuatSaat: now},
		{Id: 2, Judul: "Done", Target: UkurTarget{Sesi: 2, Periode: 7}, Progress: 2, IsDone: true, DibuatSaat: now},
	}

	var openOutput bytes.Buffer
	if got := printActsTable(&openOutput, activities, false, false, now); got != 1 {
		t.Fatalf("default list shows %d activities, want 1", got)
	}
	if strings.Contains(openOutput.String(), "Done") || !strings.Contains(openOutput.String(), "Open") {
		t.Fatalf("default list filtering is wrong:\n%s", openOutput.String())
	}

	var doneOutput bytes.Buffer
	if got := printActsTable(&doneOutput, activities, false, true, now); got != 1 {
		t.Fatalf("--done shows %d activities, want 1", got)
	}
	if strings.Contains(doneOutput.String(), "Open") || !strings.Contains(doneOutput.String(), "Done") {
		t.Fatalf("--done filtering is wrong:\n%s", doneOutput.String())
	}
}

func TestHelpRoutingDoesNotRequireDataStore(t *testing.T) {
	for _, args := range [][]string{
		{"help"},
		{"acts", "help"},
		{"notes"},
		{"notes", "help"},
	} {
		var output bytes.Buffer
		if err := run(args, &output); err != nil {
			t.Errorf("run(%q) returned error: %v", args, err)
		}
		if output.Len() == 0 {
			t.Errorf("run(%q) printed no help text", args)
		}
	}
}
