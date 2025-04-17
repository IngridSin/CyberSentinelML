package capture

import (
	"math"
	sort "sort"
	"time"
)

func detectBulkFeatures(timestamps []time.Time, lengths []int) (avgBytesBulk float64, avgPacketsBulk float64, bulkRate float64) {
	var bulks [][]int
	var current []int

	for i := 1; i < len(timestamps); i++ {
		delta := timestamps[i].Sub(timestamps[i-1])
		if delta <= time.Millisecond {
			if len(current) == 0 {
				current = append(current, i-1)
			}
			current = append(current, i)
		} else {
			if len(current) >= 4 {
				bulks = append(bulks, current)
			}
			current = nil
		}
	}
	if len(current) >= 4 {
		bulks = append(bulks, current)
	}

	var totalPackets, totalBytes int
	var totalDuration time.Duration

	for _, bulk := range bulks {
		totalPackets += len(bulk)
		for _, idx := range bulk {
			totalBytes += lengths[idx]
		}
		duration := timestamps[bulk[len(bulk)-1]].Sub(timestamps[bulk[0]])
		totalDuration += duration
	}

	count := float64(len(bulks))
	if count == 0 {
		return 0, 0, 0
	}

	return float64(totalBytes) / count, float64(totalPackets) / count, float64(totalBytes) / totalDuration.Seconds()
}

func detectActiveIdleFeatures(timestamps []time.Time) (activeMean, activeStd, activeMax, activeMin float64,
	idleMean, idleStd, idleMax, idleMin float64) {

	if len(timestamps) < 2 {
		return
	}

	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].Before(timestamps[j])
	})

	var activeDurations []float64
	var idleDurations []float64

	activeStart := timestamps[0]
	for i := 1; i < len(timestamps); i++ {
		delta := timestamps[i].Sub(timestamps[i-1]).Seconds()
		if delta > 1.0 {
			activeDur := timestamps[i-1].Sub(activeStart).Seconds()
			if activeDur > 0 {
				activeDurations = append(activeDurations, activeDur)
			}
			idleDurations = append(idleDurations, delta)
			activeStart = timestamps[i]
		}
	}
	if last := timestamps[len(timestamps)-1].Sub(activeStart).Seconds(); last > 0 {
		activeDurations = append(activeDurations, last)
	}

	activeMean, activeStd, activeMax, activeMin = stats(activeDurations)
	idleMean, idleStd, idleMax, idleMin = stats(idleDurations)
	return
}

func stats(values []float64) (mean, std, max, min float64) {
	if len(values) == 0 {
		return
	}
	sum := 0.0
	max = values[0]
	min = values[0]
	for _, v := range values {
		sum += v
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	mean = sum / float64(len(values))

	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	std = math.Sqrt(variance / float64(len(values)))
	return
}

func mergeAndSort(a, b []time.Time) []time.Time {
	merged := append([]time.Time{}, a...)
	merged = append(merged, b...)
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Before(merged[j])
	})
	return merged
}
