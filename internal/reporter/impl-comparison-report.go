package reporter

import (
	"fmt"
	"spacemon/internal/comparer"
	"spacemon/internal/scanner"
)

type ComparisonReport struct {
	prevResult scanner.ScanResult
}

func NewComparisonReport(prevResult scanner.ScanResult) *ComparisonReport {
	return &ComparisonReport{
		prevResult: prevResult,
	}
}

func (r *ComparisonReport) Update(result scanner.ScanResult) {
	comparisonResult := comparer.CompareResults(&r.prevResult, &result)
	fmt.Printf("comparisonResult: %+v\n", comparisonResult)
	// todo: implement report generation
}

func (r *ComparisonReport) Save() {
	// save to file
}
