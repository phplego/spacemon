package reporter

import (
	"space-monitor2/internal/scanner"
)

type Report interface {
	Update(result scanner.ScanResult)
	Save()
}
