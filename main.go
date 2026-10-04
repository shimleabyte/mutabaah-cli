package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	list, err := LoadAktivitas(dataFile)
	if err != nil {
		fmt.Println("Gagal load data:", err)
		os.Exit(1)
	}

	command := os.Args[1]
	changed := true

	switch command {
	case "tambah":
		if len(os.Args) < 3 {
			fmt.Println("Penggunaan: mutabaah-cli tambah \"Judul Aktivitas\"")
			os.Exit(1)
		}
		judul := strings.Join(os.Args[2:], " ")
		aktivitas := list.Tambah(judul)
		fmt.Printf("Aktivitas %s dengan ID %d\n", aktivitas.Judul, aktivitas.Id)

	case "list":
		changed = false
		if len(list.DaftarAktivitas) == 0 {
			fmt.Println("Anda belum menambahkah agenda apapun, gunakan perintah: tambah \"judul aktivitas\"")
			return
		}
		for _, a := range list.DaftarAktivitas {
			status := "|Ncan|"
			if a.IsDone {
				status = "|Ntos|"
			}
			fmt.Printf("%s %d. %s\n", status, a.Id, a.Judul)
		}

	case "sudah":
		id := parseID(os.Args, "sudah")
		if err := list.Sudah(id); err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		fmt.Printf("Bagoes! Aktivitas %d sudah ditandai terlaksana", id)

	case "hapus":
		id := parseID(os.Args, "hapus")
		if err := list.Hapus(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Aktivitas %d dihapus\n", id)

	case "ubah":
	case "help":
		changed = false
		printHelp()
	default:
		changed = false
		fmt.Println("Command naon ai sia:", command)
		printHelp()
		os.Exit(1)

	}

	if changed {
		if err := SimpanAktivitas(dataFile, list); err != nil {
			fmt.Println("Gagal disimpan:", err)
			os.Exit(1)
		}
	}
}

func parseID(args []string, cmd string) int {
	if len(args) < 3 {
		fmt.Printf("Penggunaan: mutabaah-cli %s <id>\n", cmd)
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("ID harus berupa angka")
		os.Exit(1)
	}
	return id
}

func printHelp() {
	fmt.Println(`mutabaah-cli tracker ibadah sederhana
	Penggunaan:
	mutabaah-cli tambah "judul aktivitas" Tambah aktivitas baru ke list mutabaah
	`)
}
