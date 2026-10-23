package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	appDir     = ".mtcli"
	dataFile   = "daftaraktivitas.json"
	legacyData = "daftaraktivitas.json"
)

func DataPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Direktori raib: %w", err)
	}
	return filepath.Join(home, appDir, dataFile), nil
}

func LoadAktivitas(path string) (ListMutabaah, error) {
	var list ListMutabaah

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return list, nil
		}
		return list, err
	}

	if err := json.Unmarshal(data, &list); err != nil {
		return list, err
	}

	// Normalize older records that may lack target metadata. encoding/json matches
	// keys case-insensitively, so legacy "Sesi"/"Periode" fields remain readable.
	for i := range list.DaftarAktivitas {
		aktivitas := &list.DaftarAktivitas[i]
		if aktivitas.Target.Sesi <= 0 {
			aktivitas.Target.Sesi = defaultTargetSesi
		}
		if aktivitas.Target.Periode <= 0 {
			aktivitas.Target.Periode = defaultTargetPeriode
		}
		if aktivitas.Progress < 0 {
			aktivitas.Progress = 0
		}
		if aktivitas.Progress > aktivitas.Target.Sesi {
			aktivitas.Progress = aktivitas.Target.Sesi
		}
		// Older versions stored completion without maintaining Progress.
		if aktivitas.IsDone && aktivitas.Progress < aktivitas.Target.Sesi {
			aktivitas.Progress = aktivitas.Target.Sesi
		}
		if aktivitas.Progress >= aktivitas.Target.Sesi {
			aktivitas.Progress = aktivitas.Target.Sesi
			aktivitas.IsDone = true
		}
	}
	return list, nil

}

func SimpanAktivitas(path string, list ListMutabaah) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func MigrateLegacy(newPath string) error {
	if _, err := os.Stat(newPath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	data, err := os.ReadFile(legacyData)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if err := os.MkdirAll(filepath.Dir(newPath), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(newPath, data, 0o600); err != nil {
		return err
	}
	fmt.Printf("Data lama dari ./%s disalin ke %s\n", legacyData, newPath)
	return nil
}
