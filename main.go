package main

import (
	"fmt"
)

type Aktivitas struct {
	Nama     string
	Target   UkurTarget
	Progress int
	IsDone   bool
}

type UkurTarget struct {
	Sesi    int
	Periode int
}

func (act Aktivitas) disimpulkan() string {
	var status string
	if act.IsDone {
		status = "Selesai"
	} else {
		status = "Belum Selesai"
	}

	return fmt.Sprintf("%s dibuat dengan target %d/%d dengan status progres %d maka aktivitas ini %s",
		act.Nama, act.Target.Sesi, act.Target.Periode, act.Progress, status)
}

func main() {
	act1 := Aktivitas{Nama: "Olahraga", Target: UkurTarget{Sesi: 2, Periode: 7}, Progress: 2, IsDone: true}
	fmt.Println(act1.disimpulkan())
}
