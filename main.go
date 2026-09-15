package main

import (
	"fmt"
	"time"
)

func main() {

	aktivitas := Aktivitas{
		ID:         1,
		Judul:      "Test",
		Target:     UkurTarget{Sesi: 1, Periode: 2},
		Progress:   0,
		IsDone:     false,
		DibuatSaat: time.Now(),
	}

	fmt.Printf("Aktivitas \n %s ditambahkan dengan target %d/%d \n Dibuat %s", aktivitas.Judul, aktivitas.Target.Sesi, aktivitas.Target.Periode, aktivitas.DibuatSaat)

}
