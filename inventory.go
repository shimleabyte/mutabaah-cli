package main

import (
	"encoding/json"
	"os"
)

const dataFile = "daftaraktivitas.json"

func LoadAktivitas(path string) (ListMutabaah, error) {
	var list ListMutabaah

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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
	return os.WriteFile(path, data, 0644)
}
