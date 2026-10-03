package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	list, err := LoadTodos(dataFile)
	if err != nil {
		fmt.Println("Gagal load data:", err)
		os.Exit(1)
	}

}
