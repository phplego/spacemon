package reporter

import "spacemon/internal/scanner"

type SingleScanReport struct{}

func (s *SingleScanReport) Update(result scanner.ScanResult) {
	// todo: implement report generation
}

func (s *SingleScanReport) Save() {
	// save to file
}
