package reporter

import (
	"spacemon/internal/scanner"
)

type Report interface {
	Update(result scanner.ScanResult)
	Save()
}
