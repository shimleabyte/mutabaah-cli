package main

import "fmt"

// Tanggung jawabnya adalah persentase

func (activity Activity) CompletionPercentage() int {
	if activity.Target == 0 {
		return 0
	}
	return (activity.Progress * 100) / activity.Target
}

func (activity Activity) ActivitySummary() string {
	return fmt.Sprintf("%s: %d/%d (%d%%)", activity.Name, activity.Progress, activity.Target, activity.CompletionPercentage())
}
