package reporter

import (
	"os"
	"spacemon/internal/config"
	"spacemon/internal/scanner"
)

type ReportInterface interface {
	Update(result scanner.ScanResult)
	Save()
}

type BaseReport struct {
	lastReportOutput string
}

func (r *BaseReport) Save() {
	// save to file
	os.WriteFile(config.GetAppDir()+"/last-report.txt", []byte(r.lastReportOutput+"\n"), 0666)
}
