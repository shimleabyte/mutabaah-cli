package main

import (
	"time"
)

type Aktivitas struct {
	ID         int        `json:"id"`
	Judul      string     `json:"judul"`
	Target     UkurTarget `json:"target"`
	Progress   int        `json:"proses"`
	IsDone     bool       `json:"ketuntasan"`
	DibuatSaat time.Time  `json:"dibuatSaat"`
}

type UkurTarget struct {
	Sesi    int
	Periode int
}

type ListMutabaah struct {
	DaftarAktivitas []Aktivitas `json:"daftarAktivitas"`
}

func (lm *ListMutabaah) Tambah(judul string) Aktivitas {

	idBaru := 1
	target := UkurTarget{3, 7}

	if len(lm.DaftarAktivitas) > 0 {
		idBaru = lm.DaftarAktivitas[len(lm.DaftarAktivitas)-1].ID + 1
	}

	aktivitas := Aktivitas{
		ID:         idBaru,
		Judul:      judul,
		Target:     target,
		Progress:   0,
		IsDone:     false,
		DibuatSaat: time.Now(),
	}

	lm.DaftarAktivitas = append(lm.DaftarAktivitas, aktivitas)
	return aktivitas
}
