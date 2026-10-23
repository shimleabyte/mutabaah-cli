package main

import (
	"testing"
	"time"
)

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
	activity := Aktivitas{
		DibuatSaat: createdAt,
		Target:     UkurTarget{Sesi: 3, Periode: 7},
	}
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
