package main

import (
	"fmt"
	"strings"
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
	Sesi    int `json:"sesi"`
	Periode int `json:"periode"`
}

type ListMutabaah struct {
	DaftarAktivitas []Aktivitas `json:"daftarAktivitas"`
}

func (lm *ListMutabaah) Tambah(judul string) Aktivitas {
	aktivitas, _ := lm.TambahDenganTarget(judul, UkurTarget{
		Sesi:    defaultTargetSesi,
		Periode: defaultTargetPeriode,
	})
	return aktivitas
}

func (lm *ListMutabaah) TambahDenganTarget(judul string, target UkurTarget) (Aktivitas, error) {
	judul = strings.TrimSpace(judul)
	if judul == "" {
		return Aktivitas{}, fmt.Errorf("judul aktivitas tidak boleh kosong")
	}
	if target.Sesi <= 0 || target.Periode <= 0 {
		return Aktivitas{}, fmt.Errorf("target sesi dan periode harus berupa angka positif")
	}

	idBaru := 1
	if len(lm.DaftarAktivitas) > 0 {
		idBaru = lm.DaftarAktivitas[len(lm.DaftarAktivitas)-1].Id + 1
	}

	aktivitas := Aktivitas{
		Id:         idBaru,
		Judul:      judul,
		Target:     target,
		Progress:   0,
		IsDone:     false,
		DibuatSaat: time.Now(),
	}
	lm.DaftarAktivitas = append(lm.DaftarAktivitas, aktivitas)
	return aktivitas, nil
}

// Catat menambah satu sesi progress. Aktivitas selesai saat jumlah sesi
// mencapai target. Periode digunakan sebagai batas waktu sejak aktivitas dibuat.
func (lm *ListMutabaah) Catat(id int) (Aktivitas, error) {
	for i := range lm.DaftarAktivitas {
		aktivitas := &lm.DaftarAktivitas[i]
		if aktivitas.Id != id {
			continue
		}
		if aktivitas.IsDone {
			return Aktivitas{}, fmt.Errorf("aktivitas dengan ID %d sudah selesai", id)
		}
		if aktivitas.Target.Sesi <= 0 {
			return Aktivitas{}, fmt.Errorf("target sesi aktivitas %d tidak valid", id)
		}

		aktivitas.Progress++
		if aktivitas.Progress >= aktivitas.Target.Sesi {
			aktivitas.Progress = aktivitas.Target.Sesi
			aktivitas.IsDone = true
		}
		return *aktivitas, nil
	}
	return Aktivitas{}, fmt.Errorf("aktivitas dengan ID %d tidak ditemukan; periksa dengan perintah list", id)
}

// Sudah merupakan shortcut manual untuk menuntaskan target aktivitas.
func (lm *ListMutabaah) Sudah(id int) error {
	for i := range lm.DaftarAktivitas {
		aktivitas := &lm.DaftarAktivitas[i]
		if aktivitas.Id != id {
			continue
		}
		aktivitas.IsDone = true
		if aktivitas.Target.Sesi > 0 {
			aktivitas.Progress = aktivitas.Target.Sesi
		}
		return nil
	}
	return fmt.Errorf("aktivitas dengan ID %d tidak ditemukan; periksa dengan perintah list", id)
}

func (a Aktivitas) Terlambat(now time.Time) bool {
	if a.IsDone || a.Target.Periode <= 0 || a.DibuatSaat.IsZero() {
		return false
	}
	deadline := a.DibuatSaat.AddDate(0, 0, a.Target.Periode)
	return now.After(deadline)
}

func (lm *ListMutabaah) Hapus(id int) error {
	for i, m := range lm.DaftarAktivitas {
		if m.Id == id {
			lm.DaftarAktivitas = append(lm.DaftarAktivitas[:i], lm.DaftarAktivitas[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Aktivitas yang anda (id: %d) cari dengan perintah 'list' apakah sudah dihapus atau memang tidak dibuat sebelumnya", id)
}

func (lm *ListMutabaah) Ubah(id int, judulBaru string) error {
	for i := range lm.DaftarAktivitas {
		if lm.DaftarAktivitas[i].Id == id {
			lm.DaftarAktivitas[i].Judul = judulBaru
			return nil
		}
	}

	return fmt.Errorf("Aktivitas dengan ID %d tidak ditemukan, cek lagi daftarnya dengan perintah list", id)
}
