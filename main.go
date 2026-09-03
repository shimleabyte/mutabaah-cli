package main

import (
	"bufio"
	"fmt"
	"go-project/mutabaah-cli/utils"
	"os"
	"strconv"
	"strings"
)

// ... (Masukkan definisi struct Activity dan Enum dari poin 1 di sini) ...

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== TAMBAH AKTIVITAS BARU ===")

	// 1. Input Nama
	fmt.Print("Nama Aktivitas: ")
	scanner.Scan()
	nama := strings.TrimSpace(scanner.Text())

	// 2. Opsi Pilihan Jenis Aktivitas
	fmt.Println("\nPilih Jenis Aktivitas:")
	fmt.Println("1. Checklist (Tinggal centang selesai/belum)")
	fmt.Println("2. Counter (Input angka target harian/mingguan)")
	fmt.Println("3. Note (Berupa catatan teks/keterangan)")
	fmt.Print("Masukkan pilihan (1/2/3): ")
	scanner.Scan()
	pilihan := strings.TrimSpace(scanner.Text())

	var jenis utils.JenisAktivitas
	var target int

	// 3. Logika penentuan Target berdasarkan Jenis
	switch pilihan {
	case "1":
		jenis = utils.JenisChecklist
		target = 1 // Checklist goalsnya cuma 1 (selesai)
	case "2":
		jenis = utils.JenisCounter
		fmt.Print("Masukkan angka target (misal: 5 untuk 5 gelas air): ")
		scanner.Scan()
		target, _ = strconv.Atoi(strings.TrimSpace(scanner.Text()))
	case "3":
		jenis = utils.JenisNote
		target = 0 // Note tidak pakai target angka
	default:
		fmt.Println("Pilihan tidak valid!")
		return
	}

	// 4. Bikin objek baru (Progres otomatis 0 karena tidak kita isi)
	newActivity := utils.Activity{
		Nama:   nama,
		Jenis:  jenis,
		Target: target,
	}

	// 5. Tinggal append ke slice storage kamu seperti sebelumnya
	var todoList []utils.Activity
	todoList = append(todoList, newActivity)

	fmt.Printf("\nSukses membuat aktivitas '%s' dengan tipe [%s]!\n", newActivity.Nama, newActivity.Jenis)
}
