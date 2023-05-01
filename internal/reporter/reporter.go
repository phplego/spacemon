package reporter

import (
	"fmt"
	"github.com/robert-nix/ansihtml"
	"os"
	"spacemon/internal/config"
	"spacemon/internal/scanner"
)

type ReportInterface interface {
	Update(result scanner.ScanResult)
	Render() string
	Save()
	RenderJson() string
}

type BaseReport struct {
	lastReportOutput string
}

func (r *BaseReport) Save() {
	// save to file
	os.WriteFile(config.GetAppDir()+"/last-report.txt", []byte(r.lastReportOutput+"\n"), 0666)
	html := ansihtml.ConvertToHTML([]byte(r.lastReportOutput))
	html = []byte(fmt.Sprintf("<pre style='background-color: black; color: white'>%s</pre>", html))
	os.WriteFile(config.GetAppDir()+"/last-report.html", html, 0666)

}
