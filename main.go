package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// run routes the root command before touching the data store. Help, unknown
// domains, and the not-yet-implemented notes domain do not load activity data.
func run(args []string, out io.Writer) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		printRootHelp(out)
		return nil
	}

	switch args[0] {
	case "acts":
		return runActs(args[1:], out)
	case "notes":
		return runNotes(args[1:], out)
	default:
		printRootHelp(out)
		return fmt.Errorf("domain %q tidak dikenal", args[0])
	}
}

func runNotes(args []string, out io.Writer) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		printNotesHelp(out)
		return nil
	}
	return fmt.Errorf("modul notes belum diimplementasikan; command ini belum membaca atau mengubah data")
}

func runActs(args []string, out io.Writer) error {
	if len(args) == 0 || isHelpArg(args[0]) {
		printActsHelp(out)
		return nil
	}

	command := args[0]
	var (
		judul   string
		target  UkurTarget
		id      int
		showAll bool
		showDone bool
		err     error
	)

	// Validate the subcommand and its arguments before resolving or loading data.
	switch command {
	case "list":
		showAll, showDone, err = parseListFlags(args[1:])
		if err != nil {
			return err
		}
	case "tambah":
		judul, target, err = parseTambahArgs(args[1:])
		if err != nil {
			return err
		}
	case "catat", "sudah", "hapus":
		id, err = parseID(args[1:], command)
		if err != nil {
			return err
		}
	case "ubah":
		if len(args) < 3 {
			return fmt.Errorf("penggunaan: mutabaah acts ubah <id> \"judul baru\"")
		}
		id, err = parseID(args[1:2], command)
		if err != nil {
			return err
		}
		judul = strings.TrimSpace(strings.Join(args[2:], " "))
		if judul == "" {
			return fmt.Errorf("judul baru tidak boleh kosong")
		}
	default:
		printActsHelp(out)
		return fmt.Errorf("subcommand acts %q tidak dikenal", command)
	}

	path, err := DataPath()
	if err != nil {
		return err
	}
	if err := MigrateLegacy(path); err != nil {
		return fmt.Errorf("gagal migrasi data lama: %w", err)
	}

	list, err := LoadAktivitas(path)
	if err != nil {
		return fmt.Errorf("gagal memuat data: %w", err)
	}

	changed := false
	switch command {
	case "list":
		shown := printActsTable(out, list.DaftarAktivitas, showAll, showDone, time.Now())
		if shown == 0 {
			switch {
			case len(list.DaftarAktivitas) == 0:
				fmt.Fprintln(out, "Belum ada aktivitas. Tambahkan dengan: mutabaah acts tambah \"judul aktivitas\"")
			case showDone:
				fmt.Fprintln(out, "Belum ada aktivitas yang selesai. Catat sesi dengan: mutabaah acts catat <id>")
			case !showAll:
				fmt.Fprintln(out, "Tidak ada aktivitas yang belum selesai. Gunakan: mutabaah acts list --all")
			}
		}
	case "tambah":
		aktivitas, err := list.TambahDenganTarget(judul, target)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Aktivitas %q ditambahkan (ID %d, target %d sesi / %d hari)\n",
			aktivitas.Judul, aktivitas.Id, aktivitas.Target.Sesi, aktivitas.Target.Periode)
		changed = true
	case "catat":
		aktivitas, err := list.Catat(id)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Sesi dicatat untuk aktivitas %d: %d/%d\n", id, aktivitas.Progress, aktivitas.Target.Sesi)
		if aktivitas.IsDone {
			fmt.Fprintln(out, "Target sesi tercapai; aktivitas ditandai selesai.")
		} else if aktivitas.Terlambat(time.Now()) {
			fmt.Fprintln(out, "Perhatian: target periode sudah lewat.")
		}
		changed = true
	case "sudah":
		if err := list.Sudah(id); err != nil {
			return err
		}
		fmt.Fprintf(out, "Aktivitas %d ditandai selesai secara manual.\n", id)
		changed = true
	case "hapus":
		if err := list.Hapus(id); err != nil {
			return err
		}
		fmt.Fprintf(out, "Aktivitas %d dihapus.\n", id)
		changed = true
	case "ubah":
		if err := list.Ubah(id, judul); err != nil {
			return err
		}
		fmt.Fprintf(out, "Judul aktivitas %d diperbarui menjadi %q.\n", id, judul)
		changed = true
	}

	if changed {
		if err := SimpanAktivitas(path, list); err != nil {
			return fmt.Errorf("gagal menyimpan data: %w", err)
		}
	}
	return nil
}

func isHelpArg(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func printRootHelp(out io.Writer) {
	fmt.Fprint(out, "mutabaah — pencatatan aktivitas dan catatan pribadi\n\n"+
		"Penggunaan:\n"+
		"  mutabaah acts <command>     Kelola aktivitas\n"+
		"  mutabaah notes              Catatan pribadi (belum tersedia)\n"+
		"  mutabaah help               Tampilkan bantuan\n\n"+
		"Contoh:\n"+
		"  mutabaah acts list\n"+
		"  mutabaah acts tambah \"Tilawah satu juz\" --sesi 3 --periode 7\n"+
		"  mutabaah acts catat 1\n"+
		"  mutabaah notes\n")
}

func printActsHelp(out io.Writer) {
	fmt.Fprint(out, "mutabaah acts — kelola aktivitas\n\n"+
		"Penggunaan:\n"+
		"  mutabaah acts tambah \"judul aktivitas\" [--sesi N] [--periode HARI]\n"+
		"  mutabaah acts list [--all | --done]\n"+
		"  mutabaah acts catat <id>             Catat satu sesi\n"+
		"  mutabaah acts sudah <id>             Tandai selesai secara manual\n"+
		"  mutabaah acts ubah <id> \"judul baru\"\n"+
		"  mutabaah acts hapus <id>\n"+
		"  mutabaah acts help\n\n"+
		"list tanpa flag menampilkan aktivitas yang belum selesai.\n"+
		"--done menampilkan aktivitas yang selesai; --all menampilkan semuanya.\n"+
		"Target default: 3 sesi dalam 7 hari sejak aktivitas dibuat.\n")
}

func printNotesHelp(out io.Writer) {
	fmt.Fprint(out, "mutabaah notes — catatan pribadi\n\n"+
		"Modul notes belum diimplementasikan. Command ini belum membaca atau mengubah data.\n"+
		"Rencana command akan dipisahkan dari acts, misalnya:\n"+
		"  mutabaah notes list\n"+
		"  mutabaah notes add \"judul catatan\"\n"+
		"  mutabaah notes show <id>\n")
}

func parseTambahArgs(args []string) (string, UkurTarget, error) {
	target := UkurTarget{Sesi: defaultTargetSesi, Periode: defaultTargetPeriode}
	judulParts := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--sesi" || arg == "--periode":
			if i+1 >= len(args) {
				return "", UkurTarget{}, fmt.Errorf("flag %s membutuhkan angka positif", arg)
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
			return "", UkurTarget{}, fmt.Errorf("flag %q tidak dikenal; gunakan --sesi dan --periode", arg)
		default:
			judulParts = append(judulParts, arg)
		}
	}

	judul := strings.TrimSpace(strings.Join(judulParts, " "))
	if judul == "" {
		return "", UkurTarget{}, fmt.Errorf("penggunaan: mutabaah acts tambah \"judul aktivitas\" [--sesi N] [--periode HARI]")
	}
	return judul, target, nil
}

func parsePositiveInt(flag, rawValue string) (int, error) {
	value, err := strconv.Atoi(rawValue)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("nilai %s harus berupa angka positif", flag)
	}
	return value, nil
}

func parseListFlags(args []string) (showAll, showDone bool, err error) {
	for _, arg := range args {
		switch arg {
		case "--all", "-all":
			showAll = true
		case "--done", "-done":
			showDone = true
		default:
			return false, false, fmt.Errorf("flag %q tidak dikenal; gunakan --all atau --done", arg)
		}
	}
	if showAll && showDone {
		return false, false, fmt.Errorf("--all dan --done tidak bisa dipakai bersamaan")
	}
	return showAll, showDone, nil
}

func parseID(args []string, command string) (int, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("penggunaan: mutabaah acts %s <id>", command)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("ID harus berupa angka positif")
	}
	return id, nil
}

// printActsTable renders a terminal-friendly table. Progress is displayed as
// sessions completed / target sessions, with a fixed-width ASCII progress bar.
func printActsTable(out io.Writer, activities []Aktivitas, showAll, showDone bool, now time.Time) int {
	selected := make([]Aktivitas, 0, len(activities))
	for _, activity := range activities {
		if activity.IsDone && !showAll && !showDone {
			continue
		}
		if !activity.IsDone && showDone {
			continue
		}
		selected = append(selected, activity)
	}
	if len(selected) == 0 {
		return 0
	}

	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tSTATE\tACTIVITY\tPROGRESS\tDEADLINE")
	for _, activity := range selected {
		state := "OPEN"
		if activity.IsDone {
			state = "DONE"
		} else if activity.Terlambat(now) {
			state = "LATE"
		}
		bar := progressBar(activity.Progress, activity.Target.Sesi, 10)
		progress := fmt.Sprintf("%s %d/%d", bar, activity.Progress, activity.Target.Sesi)
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\n",
			activity.Id, state, activity.Judul, progress, deadlineLabel(activity, now))
	}
	_ = writer.Flush()
	return len(selected)
}

func progressBar(progress, target, width int) string {
	if width < 1 {
		width = 10
	}
	if target <= 0 {
		return "[" + strings.Repeat("-", width) + "]"
	}
	if progress < 0 {
		progress = 0
	}
	if progress > target {
		progress = target
	}
	filled := progress * width / target
	if progress == target {
		filled = width
	}
	return "[" + strings.Repeat("#", filled) + strings.Repeat("-", width-filled) + "]"
}

func deadlineLabel(activity Aktivitas, now time.Time) string {
	if activity.IsDone {
		return "complete"
	}
	if activity.Target.Periode <= 0 || activity.DibuatSaat.IsZero() {
		return "unknown"
	}

	deadline := activity.DibuatSaat.AddDate(0, 0, activity.Target.Periode)
	if now.After(deadline) {
		daysLate := int(math.Ceil(now.Sub(deadline).Hours() / 24))
		if daysLate < 1 {
			daysLate = 1
		}
		return fmt.Sprintf("late %dd", daysLate)
	}
	hoursLeft := deadline.Sub(now).Hours()
	if hoursLeft < 24 {
		return "due soon"
	}
	daysLeft := int(math.Ceil(hoursLeft / 24))
	return fmt.Sprintf("%dd left", daysLeft)
}
