package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseTambahArgsDefaultTarget(t *testing.T) {
	judul, target, err := parseTambahArgs([]string{"Tilawah", "satu", "juz"})
	if err != nil {
		t.Fatal(err)
	}
	if judul != "Tilawah satu juz" {
		t.Fatalf("judul = %q", judul)
	}
	if target != (UkurTarget{Sesi: 3, Periode: 7}) {
		t.Fatalf("target = %#v", target)
	}
}

func TestParseTambahArgsCustomTarget(t *testing.T) {
	judul, target, err := parseTambahArgs([]string{"Tilawah satu juz", "--sesi", "5", "--periode=10"})
	if err != nil {
		t.Fatal(err)
	}
	if judul != "Tilawah satu juz" {
		t.Fatalf("judul = %q", judul)
	}
	if target != (UkurTarget{Sesi: 5, Periode: 10}) {
		t.Fatalf("target = %#v", target)
	}
}

func TestParseTambahArgsRejectsInvalidTargets(t *testing.T) {
	cases := [][]string{
		{"Tilawah", "--sesi", "0"},
		{"Tilawah", "--periode", "-2"},
		{"Tilawah", "--sesi"},
		{"Tilawah", "--unknown", "2"},
		{"--sesi", "2"},
	}
	for _, args := range cases {
		if _, _, err := parseTambahArgs(args); err == nil {
			t.Errorf("parseTambahArgs(%q) expected an error", args)
		}
	}
}

func TestParseListFlags(t *testing.T) {
	all, done, err := parseListFlags([]string{"--all"})
	if err != nil || !all || done {
		t.Fatalf("unexpected --all parse: all=%v done=%v err=%v", all, done, err)
	}
	all, done, err = parseListFlags([]string{"--done"})
	if err != nil || all || !done {
		t.Fatalf("unexpected --done parse: all=%v done=%v err=%v", all, done, err)
	}
	if _, _, err := parseListFlags([]string{"--unknown"}); err == nil {
		t.Fatal("unknown list flag should return an error")
	}
}

func TestUkurTargetJSONUsesLowercaseKeysAndReadsLegacyKeys(t *testing.T) {
	list := ListMutabaah{DaftarAktivitas: []Aktivitas{{
		Id: 1, Judul: "Tilawah", Target: UkurTarget{Sesi: 4, Periode: 9},
	}}}
	data, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\"sesi\":4") || !strings.Contains(string(data), "\"periode\":9") {
		t.Fatalf("expected lowercase target keys, got: %s", data)
	}

	path := filepath.Join(t.TempDir(), "legacy.json")
	legacy := "{\"daftarAktivitas\":[{\"id\":1,\"judul\":\"Old\",\"target\":{\"Sesi\":4,\"Periode\":9},\"proses\":1,\"ketuntasan\":false,\"dibuatSaat\":\"2026-10-20T00:00:00Z\"}]}"
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadAktivitas(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := loaded.DaftarAktivitas[0].Target; got != (UkurTarget{Sesi: 4, Periode: 9}) {
		t.Fatalf("legacy target = %#v", got)
	}
}

func TestLoadNormalizesLegacyCompletedProgress(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.json")
	legacy := "{\"daftarAktivitas\":[{\"id\":1,\"judul\":\"Old\",\"target\":{\"Sesi\":3,\"Periode\":7},\"proses\":0,\"ketuntasan\":true,\"dibuatSaat\":\"2026-10-01T00:00:00Z\"}]}"
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	list, err := LoadAktivitas(path)
	if err != nil {
		t.Fatal(err)
	}
	got := list.DaftarAktivitas[0]
	if !got.IsDone || got.Progress != got.Target.Sesi {
		t.Fatalf("legacy completed activity not normalized: %#v", got)
	}
}

func TestCatatIncrementsProgressAndCompletesAtTarget(t *testing.T) {
	var list ListMutabaah
	created, err := list.TambahDenganTarget("Tilawah", UkurTarget{Sesi: 2, Periode: 7})
	if err != nil {
		t.Fatal(err)
	}
	if created.Progress != 0 || created.IsDone {
		t.Fatalf("new activity = progress %d, done %v", created.Progress, created.IsDone)
	}

	first, err := list.Catat(created.Id)
	if err != nil {
		t.Fatal(err)
	}
	if first.Progress != 1 || first.IsDone {
		t.Fatalf("after first session = progress %d, done %v", first.Progress, first.IsDone)
	}

	second, err := list.Catat(created.Id)
	if err != nil {
		t.Fatal(err)
	}
	if second.Progress != 2 || !second.IsDone {
		t.Fatalf("after target = progress %d, done %v", second.Progress, second.IsDone)
	}
	if _, err := list.Catat(created.Id); err == nil {
		t.Fatal("recording a session for a completed activity should fail")
	}
}

func TestSudahCompletesProgressTarget(t *testing.T) {
	var list ListMutabaah
	activity, err := list.TambahDenganTarget("Sholat Dhuha", UkurTarget{Sesi: 3, Periode: 7})
	if err != nil {
		t.Fatal(err)
	}
	if err := list.Sudah(activity.Id); err != nil {
		t.Fatal(err)
	}
	got := list.DaftarAktivitas[0]
	if !got.IsDone || got.Progress != got.Target.Sesi {
		t.Fatalf("manual completion inconsistent: %#v", got)
	}
}

func TestTerlambatUsesTargetPeriodFromCreation(t *testing.T) {
	createdAt := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	activity := Aktivitas{DibuatSaat: createdAt, Target: UkurTarget{Sesi: 3, Periode: 7}}
	if activity.Terlambat(createdAt.AddDate(0, 0, 6)) {
		t.Fatal("activity should not be late before target period passes")
	}
	if !activity.Terlambat(createdAt.AddDate(0, 0, 8)) {
		t.Fatal("activity should be late after target period passes")
	}
	activity.IsDone = true
	if activity.Terlambat(createdAt.AddDate(0, 0, 8)) {
		t.Fatal("completed activity should not be overdue")
	}
}


func TestHelpDoesNotReadCorruptDataFile(t *testing.T) {
	if os.Getenv("MUTABAAH_HELP_TEST_CHILD") == "1" {
		os.Args = []string{"mutabaah-cli", "help"}
		main()
		return
	}

	home := t.TempDir()
	dataPath := filepath.Join(home, appDir, dataFile)
	if err := os.MkdirAll(filepath.Dir(dataPath), 0o700); err != nil {
		t.Fatal(err)
	}
	const corruptJSON = "this is not valid JSON"
	if err := os.WriteFile(dataPath, []byte(corruptJSON), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestHelpDoesNotReadCorruptDataFile$")
	cmd.Env = append(os.Environ(), "MUTABAAH_HELP_TEST_CHILD=1", "HOME="+home)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help should succeed despite corrupt persisted data: %v\n%s", err, output)
	}
	after, err := os.ReadFile(dataPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != corruptJSON {
		t.Fatalf("help changed the data file: %q", after)
	}
}
