package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "help" || command == "--help" || command == "-h" {
		printHelp()
		return
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

	changed := true

	switch command {
	case "tambah":
		judul, target, err := parseTambahArgs(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		aktivitas, err := list.TambahDenganTarget(judul, target)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Aktivitas %q dengan ID %d berhasil ditambahkan (target %d sesi dalam %d hari)\n",
			aktivitas.Judul, aktivitas.Id, aktivitas.Target.Sesi, aktivitas.Target.Periode)

	case "list":
		changed = false
		showAll := false
		showDone := false

		for _, arg := range os.Args[2:] {
			switch arg {
			case "--all", "-all":
				showAll = true
			case "--done", "-done":
				showDone = true
			default:
				fmt.Printf("flag %q raib, coba: --all atau --done\n", arg)
				os.Exit(1)
			}
		}

		if showAll && showDone {
			fmt.Println("--all dan --done tidak bisa dipakai bersamaan")
			os.Exit(1)
		}

		if len(list.DaftarAktivitas) == 0 {
			fmt.Println("Anda belum menambahkan agenda apapun, gunakan perintah: mutabaah-cli tambah \"judul aktivitas\"")
			return
		}

		shown := 0

		for _, a := range list.DaftarAktivitas {
			if a.IsDone && !showAll && !showDone {
				continue
			}
			if !a.IsDone && showDone {
				continue
			}
			status := "|Ncan|"
			if a.IsDone {
				status = "|Ntos|"
			}
			detail := fmt.Sprintf("[%d/%d sesi; target %d hari]", a.Progress, a.Target.Sesi, a.Target.Periode)
			if a.Terlambat(time.Now()) {
				detail += " [TERLAMBAT]"
			}
			fmt.Printf("%s %d. %s %s\n", status, a.Id, a.Judul, detail)
			shown++
		}
		if shown == 0 {
			if showDone {
				fmt.Println("Belum ada aktivitas yang selesai. Gunakan: mutabaah-cli sudah <id> atau catat sesi dengan mutabaah-cli catat <id>")
			} else {
				fmt.Println("Semua aktivitas sudah selesai. Gunakan mutabaah-cli list --all untuk melihat semuanya.")
			}
		}

	case "catat":
		id := parseID(os.Args, "catat")
		aktivitas, err := list.Catat(id)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Printf("Sesi dicatat untuk aktivitas %d. Progress: %d/%d sesi\n", id, aktivitas.Progress, aktivitas.Target.Sesi)
		if aktivitas.IsDone {
			fmt.Println("Target sesi tercapai; aktivitas ditandai selesai.")
		} else if aktivitas.Terlambat(time.Now()) {
			fmt.Println("Catatan tersimpan, tetapi target periodenya sudah lewat.")
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

func parseTambahArgs(args []string) (string, UkurTarget, error) {
	target := UkurTarget{Sesi: defaultTargetSesi, Periode: defaultTargetPeriode}
	judulParts := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--sesi" || arg == "--periode":
			if i+1 >= len(args) {
				return "", UkurTarget{}, fmt.Errorf("Flag %s membutuhkan angka positif", arg)
			}
			value, err := parsePositiveInt(arg, args[i+1])
			if err != nil {
				return "", UkurTarget{}, err
			}
			if arg == "--sesi" {
				target.Sesi = value
			} else {
				target.Periode = value
			}
			i++
		case strings.HasPrefix(arg, "--sesi=") || strings.HasPrefix(arg, "--periode="):
			flag, rawValue, _ := strings.Cut(arg, "=")
			value, err := parsePositiveInt(flag, rawValue)
			if err != nil {
				return "", UkurTarget{}, err
			}
			if flag == "--sesi" {
				target.Sesi = value
			} else {
				target.Periode = value
			}
		case strings.HasPrefix(arg, "-"):
			return "", UkurTarget{}, fmt.Errorf("Flag %q tidak dikenal. Gunakan --sesi dan --periode", arg)
		default:
			judulParts = append(judulParts, arg)
		}
	}

	judul := strings.TrimSpace(strings.Join(judulParts, " "))
	if judul == "" {
		return "", UkurTarget{}, fmt.Errorf("Penggunaan: mutabaah-cli tambah \"judul aktivitas\" [--sesi N] [--periode HARI]")
	}
	return judul, target, nil
}

func parsePositiveInt(flag, rawValue string) (int, error) {
	value, err := strconv.Atoi(rawValue)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("Nilai %s harus berupa angka positif", flag)
	}
	return value, nil
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
	fmt.Println(`mutabaah-cli — tracker ibadah dan aktivitas berbasis CLI

Penggunaan:
  mutabaah-cli tambah "judul aktivitas" [--sesi N] [--periode HARI]
  mutabaah-cli list [--all | --done]
  mutabaah-cli catat <id>                 Catat satu sesi aktivitas
  mutabaah-cli sudah <id>                 Tandai target aktivitas selesai manual
  mutabaah-cli ubah <id> "judul baru"     Ubah judul aktivitas
  mutabaah-cli hapus <id>                 Hapus aktivitas
  mutabaah-cli help                       Tampilkan bantuan

Perintah list:
  (tanpa flag)   Tampilkan aktivitas yang belum selesai
  --done         Tampilkan aktivitas yang sudah selesai
  --all          Tampilkan semua aktivitas

Target default aktivitas baru: 3 sesi dalam 7 hari.
Periode dihitung sejak aktivitas dibuat; aktivitas yang melewati target waktu
akan ditandai TERLAMBAT selama belum selesai. Gunakan --sesi dan --periode
untuk mengatur target masing-masing aktivitas.`)
}
