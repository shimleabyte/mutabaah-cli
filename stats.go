package main

import (
	"fmt"
	"math"
)

// Tanggung jawabnya adalah persentase

func CompletionPercentage(activity Activity) int {
		var percentage int = (activity.Progress * 100) / activity.Target
		return percentage
		/* lebih praktisnya:
		return ((activity.Progress * 100) / activity.Target) */
} 


