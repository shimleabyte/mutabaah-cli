package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// run routes root and domain commands before touching activity storage. Notes
// use an independent store and help commands return before either store is read.
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

func isHelpArg(arg string) bool {
	return arg == "help" || arg == "-h" || arg == "--help"
}

func printRootHelp(out io.Writer) {
	fmt.Fprint(out, "mutabaah — pencatatan aktivitas dan catatan pribadi\n\n"+
		"Penggunaan:\n"+
		"  mutabaah acts <command>     Kelola aktivitas\n"+
		"  mutabaah notes <command>    Kelola catatan pendek (1–5 baris)\n"+
		"  mutabaah help               Tampilkan bantuan\n\n"+
		"Contoh:\n"+
		"  mutabaah acts list\n"+
		"  mutabaah acts tambah \"Tilawah satu juz\" --sesi 3 --periode 7\n"+
		"  mutabaah acts catat 1\n"+
		"  mutabaah notes add \"Review error handling\"\n"+
		"  mutabaah notes list\n")
}
