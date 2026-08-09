package main

import "fmt"

/* import (
	"fmt"
	"math"
) */

// Tanggung jawabnya adalah persentase

func CompletionPercentage(activity Activity) int {
	var ActivityPercentage int = (activity.Progress * 100) / activity.Target
	return ActivityPercentage
	/* lebih praktisnya:
	return ((activity.Progress * 100) / activity.Target) */
}

func ActivitySummary(Act1 Activity) string {
	return fmt.Sprintf("%s: %d/%d (%d%%)", Act1.Name, Act1.Progress, Act1.Target, CompletionPercentage(Act1))
}
