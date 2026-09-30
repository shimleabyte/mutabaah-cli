package main

import (
	"fmt"
	"time"
)

func main() {

	aktivitas := Aktivitas{
		Id:         1,
		Judul:      "Test",
		Target:     UkurTarget{Sesi: 1, Periode: 2},
		Progress:   0,
		IsDone:     false,
		DibuatSaat: time.Now(),
	}

	fmt.Printf(`Update Berhasil! Aktivitas ditambahkan dengan keterangan:
	Judul: %s
	Target(Sesi/Periode): %d/%d
	Direncanakan pada: %s`, aktivitas.Judul, aktivitas.Target.Sesi, aktivitas.Target.Periode, aktivitas.DibuatSaat)

}
