package reporter

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"spacemon/internal/comparer"
	"spacemon/internal/scanner"
	. "spacemon/internal/util"
	. "spacemon/internal/util/timeago"
	"strings"
	"time"
)

type ComparisonReport struct {
	BaseReport
	prevResult scanner.ScanResult
}

func NewComparisonReport(prevResult scanner.ScanResult) *ComparisonReport {
	return &ComparisonReport{
		prevResult: prevResult,
	}
}

func (r *ComparisonReport) Update(result scanner.ScanResult) {
	comparisonResult := comparer.CompareResults(&r.prevResult, &result)

	r.lastReportOutput = RenderTable(comparisonResult)
	ClearScreen(true)
	fmt.Println(r.lastReportOutput)
}

// RenderTable prints summarized table
func RenderTable(comparisonResult comparer.ComparisonResult) string {
	title := comparisonResult.ScanResult.ScanSetup.Title
	if title == "" {
		title, _ = os.Hostname()
	}
	title = strings.ToUpper(title)
	tableWriter := table.NewWriter()
	tableWriter.SetTitle("%s - %d directories", C("title", title), len(comparisonResult.ScanResult.ScanSetup.Directories))
	tableWriter.SetStyle(table.StyleRounded)
	tableWriter.AppendHeader(table.Row{"path", "size", "dirs", "files", "scan duration"})

	for _, dir := range comparisonResult.ScanResult.ScanSetup.Directories {
		dirResult, ok := comparisonResult.ScanResult.DirectoryResults.Get(dir)
		if !ok {
			// not started yet
			tableWriter.AppendRow([]interface{}{
				C("dirs", scanner.ShorifyPath(dir)),
				"…",
				"…",
				"…",
				0,
			})
			continue
		}

		if diff, ok := comparisonResult.DirectoryDiffs[dir]; ok && dirResult.Completed {
			tableWriter.AppendRow([]interface{}{
				C("dirs", scanner.ShorifyPath(dir)),
				HumanSize(dirResult.TotalSize) + " " + C("diff", HumanSizeSign(diff.SizeDiff)),
				fmt.Sprintf("%d %s", dirResult.FolderCount, C("diff", "%+d", diff.FolderCountDiff)),
				fmt.Sprintf("%d %s", dirResult.FileCount, C("diff", "%+d", diff.FileCountDiff)),
				dirResult.ScanDuration,
			})
		} else {
			// in progress
			tableWriter.AppendRow([]interface{}{
				C("dirs", scanner.ShorifyPath(dir)),
				HumanSize(dirResult.TotalSize),
				fmt.Sprintf("%d", dirResult.FolderCount),
				fmt.Sprintf("%d", dirResult.FileCount),
				dirResult.ScanDuration,
			})
		}
	}

	tableWriter.AppendSeparator()
	tableWriter.AppendRow(table.Row{ // print previous start time
		C("header", "prev stime (t₀)"),
		C("pale", comparisonResult.PrevScanTime.Format("02 Jan 15:04")),
		C("pale", TimeAgo(comparisonResult.PrevScanTime)),
	})
	tableWriter.AppendRow(table.Row{
		C("headerHi", "curr stime (t₁)"),
		comparisonResult.CurrentScanTime.Format("02 Jan 15:04"),
		//TimeAgo(comparisonResult.CurrentScanTime),
	})
	tableWriter.AppendSeparator()

	deltaFreeSpace := HumanSizeSign(comparisonResult.FreeSpaceDiff)

	tableWriter.AppendRow(table.Row{
		"FREE SPACE",
		C("free", HumanSize(comparisonResult.ScanResult.FreeSpace)) + " " + color.HiMagentaString(deltaFreeSpace),
		"", "",
		time.Since(comparisonResult.ScanResult.StartTime).Round(time.Millisecond),
	})
	return tableWriter.Render()
}
