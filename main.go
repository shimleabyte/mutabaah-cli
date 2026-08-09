package main

import "fmt"

// Tanggung jawab main.go saat ini adalah membuat Activity

func main() {
	activities := []Activity{
		{Name: "Shalat", Target: 5, Progress: 3},
		{Name: "Puasa", Target: 30, Progress: 15},
		{Name: "Sedekah", Target: 10, Progress: 7},
	}

	for _, activity := range activities {
		fmt.Println(ActivitySummary(activity))
	}
}
