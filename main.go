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

	path, err := DataPath()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if err := MigrateLegacy(path); err != nil {
		fmt.Println("Gagal migrasi data lama:", err)
		os.Exit(1)
	}

	list, err := LoadAktivitas(path)
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
		fmt.Printf("Aktivitas %s dengan ID %d berhasil ditambahkan\n", aktivitas.Judul, aktivitas.Id)

	case "list":
		changed = false
		showAll := false

		for _, arg := range os.Args[2:] {
			if arg != "--all" && arg != "-all" {
				fmt.Printf("flag yang anda ketik %q raib, coba flag ini: --all atau -all\n", arg)
				os.Exit(1)
			}
			showAll = true
		}

		if len(list.DaftarAktivitas) == 0 {
			fmt.Println("Anda belum menambahkan agenda apapun, gunakan perintah: mutabaah-cli tambah \"judul aktivitas\"")
			return
		}

		shown := 0

		for _, a := range list.DaftarAktivitas {
			if a.IsDone && !showAll {
				continue
			}
			status := "|Ncan|"
			if a.IsDone {
				status = "|Ntos|"
			}
			fmt.Printf("%s %d. %s\n", status, a.Id, a.Judul)
			shown++
		}
		if shown == 0 {
			fmt.Println("Aiight bro that's it you're complete your goals, lihat semua goal dengan: mutabaah-cli list --all")
		}

	case "sudah":
		id := parseID(os.Args, "sudah")
		if err := list.Sudah(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Bagoes! Aktivitas %d sudah ditandai terlaksana\n", id)

	case "hapus":
		id := parseID(os.Args, "hapus")
		if err := list.Hapus(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Aktivitas %d dihapus\n", id)

	case "ubah":
		if len(os.Args) < 4 {
			fmt.Println("Penggunaan: mutabaah-cli ubah <id> \"Judul Baru\"")
			os.Exit(1)
		}

		id := parseID(os.Args, "ubah")
		judulBaru := strings.Join(os.Args[3:], " ")

		if err := list.Ubah(id, judulBaru); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Printf("Aktivitas %d berhasil diubah menjadi \"%s\"\n", id, judulBaru)

	case "help":
		changed = false
		printHelp()
	default:
		changed = false
		fmt.Printf("Perintah %q tidak dikenal. Berikut daftar perintah yang tersedia:\n\n", command)
		printHelp()
		os.Exit(1)

	}

	if changed {
		if err := SimpanAktivitas(path, list); err != nil {
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
  mutabaah-cli tambah "judul aktivitas"    Tambah aktivitas baru
  mutabaah-cli list                        Tampilkan semua aktivitas
  mutabaah-cli sudah <id>                  Tandai aktivitas selesai
  mutabaah-cli hapus <id>                  Hapus aktivitas
  mutabaah-cli ubah <id> "judul baru"      Ubah judul aktivitas
	`)
}
