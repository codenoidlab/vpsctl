// Package osadapter is the only place in the whole app that reads
// system files directly (/proc/meminfo etc). Everything else goes
// through core/runner. This makes the code easier to understand.
package osadapter

import (
	"os"
	"strconv"
	"strings"
)

// MemInfo holds RAM usage numbers in megabytes.
type MemInfo struct {
	TotalMB     int
	AvailableMB int
	UsedMB      int
}

// ReadMemInfo reads /proc/meminfo and returns RAM usage.
// Linux always has /proc/meminfo so this is reliable.
func ReadMemInfo() MemInfo {
	// Read the raw file
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemInfo{}
	}

	info := MemInfo{}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		// Each line looks like: "MemTotal:       8192000 kB"
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		// Parse the number (it's in kB, convert to MB)
		kb, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		mb := kb / 1024

		switch parts[0] {
		case "MemTotal:":
			info.TotalMB = mb
		case "MemAvailable:":
			info.AvailableMB = mb
		}
	}

	// Used = Total - Available
	info.UsedMB = info.TotalMB - info.AvailableMB
	return info
}

// ReadUptime reads /proc/uptime and returns seconds the server has been running.
func ReadUptime() float64 {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	// File looks like: "12345.67 23456.78"
	// First number is total uptime in seconds
	parts := strings.Fields(string(data))
	if len(parts) == 0 {
		return 0
	}
	seconds, _ := strconv.ParseFloat(parts[0], 64)
	return seconds
}

// FormatUptime turns seconds into a human-readable string like "2d 4h 30m".
func FormatUptime(seconds float64) string {
	total := int(seconds)
	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60

	if days > 0 {
		return strings.TrimSpace(
			strconv.Itoa(days) + "d " +
				strconv.Itoa(hours) + "h " +
				strconv.Itoa(minutes) + "m",
		)
	}
	if hours > 0 {
		return strconv.Itoa(hours) + "h " + strconv.Itoa(minutes) + "m"
	}
	return strconv.Itoa(minutes) + "m"
}
