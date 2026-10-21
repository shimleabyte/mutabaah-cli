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
