package main

import (
	"fmt"
	"time"
)

type Aktivitas struct {
	Id         int        `json:"id"`
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

	IdBaru := 1
	target := UkurTarget{3, 7}

	if len(lm.DaftarAktivitas) > 0 {
		IdBaru = lm.DaftarAktivitas[len(lm.DaftarAktivitas)-1].Id + 1
	}

	aktivitas := Aktivitas{
		Id:         IdBaru,
		Judul:      judul,
		Target:     target,
		Progress:   0,
		IsDone:     false,
		DibuatSaat: time.Now(),
	}

	lm.DaftarAktivitas = append(lm.DaftarAktivitas, aktivitas)
	return aktivitas
}

func (lm *ListMutabaah) Sudah(id int) error {
	done := true

	for i := range lm.DaftarAktivitas {
		if lm.DaftarAktivitas[i].Id == id {
			lm.DaftarAktivitas[i].IsDone = done
			return nil
		}
	}
	return fmt.Errorf("Lupakah dengan agenda anda? \n id %d tidak ada dalam list mutabaah yang sudah anda buat", id)
}

func (lm *ListMutabaah) Hapus(id int) error {
	for i, m := range lm.DaftarAktivitas {
		if m.Id == id {
			lm.DaftarAktivitas = append(lm.DaftarAktivitas[:i], lm.DaftarAktivitas[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Aktivitas yang anda (id: %d) cari sudah dihapus atau memang tidak dibuat sebelumnya", id)
}
