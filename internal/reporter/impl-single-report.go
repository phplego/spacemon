package reporter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"spacemon/internal/scanner"
	. "spacemon/internal/util"
	"strings"
	"time"
)

type SingleScanReport struct{}

func (s *SingleScanReport) Update(result scanner.ScanResult) {
	ClearScreen(true)
	out := RenderSingleScanTable(result)
	fmt.Println(out)
}

// RenderSingleScanTable prints summarized table
func RenderSingleScanTable(result scanner.ScanResult) string {
	title := result.ScanSetup.Title
	if title == "" {
		title, _ = os.Hostname()
	}
	title = strings.ToUpper(title)
	tableWriter := table.NewWriter()
	tableWriter.SetTitle("%s - %d directories", C("title", title), len(result.ScanSetup.Directories))
	tableWriter.SetStyle(table.StyleRounded)
	tableWriter.AppendHeader(table.Row{"path", "size", "dirs", "files", "scan duration"})

	for i, dir := range result.ScanSetup.Directories {
		if i < len(result.DirectoryResults) {
			dirResult := result.DirectoryResults[i]
			tableWriter.AppendRow([]interface{}{
				C("dirs", scanner.ShorifyPath(dir)),
				HumanSize(dirResult.TotalSize),
				fmt.Sprintf("%d", dirResult.FolderCount),
				fmt.Sprintf("%d", dirResult.FileCount),
				dirResult.ScanDuration,
			})
		} else { // not scanned yet
			tableWriter.AppendRow([]interface{}{
				C("dirs", scanner.ShorifyPath(dir)),
				"…",
				"…",
				"…",
				0,
			})
		}
	}

	tableWriter.AppendSeparator()
	tableWriter.AppendRow(table.Row{ // print previous start time
		C("header", "prev stime (t₀)"),
		C("pale", "-"),
	})
	tableWriter.AppendRow(table.Row{
		C("headerHi", "curr stime (t₁)"),
		result.StartTime.Format("02 Jan 15:04"),
	})
	tableWriter.AppendSeparator()

	tableWriter.AppendRow(table.Row{
		"FREE SPACE",
		C("free", HumanSize(result.FreeSpace)),
		"", "",
		time.Since(result.StartTime).Round(time.Millisecond),
	})
	return tableWriter.Render()
}

func (s *SingleScanReport) Save() {
	// save to file
}
