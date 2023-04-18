package util

import (
	"fmt"
	"strings"
)

func HumanSize(bytes int64) string {
	var abs = func(v int64) int64 {
		if v < 0 {
			return -v
		}
		return v
	}

	const unit = 1024
	if abs(bytes) < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; abs(n) >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func HumanSizeSign(bytes int64) string {
	str := HumanSize(bytes)
	if !strings.HasPrefix(str, "-") {
		return "+" + str
	}
	return str
}
