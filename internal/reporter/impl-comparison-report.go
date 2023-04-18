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
	prevResult scanner.ScanResult
}

func NewComparisonReport(prevResult scanner.ScanResult) *ComparisonReport {
	return &ComparisonReport{
		prevResult: prevResult,
	}
}

var lastPrint = ""

func (r *ComparisonReport) Update(result scanner.ScanResult) {
	comparisonResult := comparer.CompareResults(&r.prevResult, &result)

	ClearScreen(true)
	PrintTable(comparisonResult)
	// todo: implement report generation
}

func (r *ComparisonReport) Save() {
	// save to file
}

// PrintTable prints summarized table
func PrintTable(comparisonResult comparer.ComparisonResult) {
	title := "test"
	if title == "" {
		title, _ = os.Hostname()
		title = strings.ToUpper(title)
	}
	tableWriter := table.NewWriter()
	tableWriter.SetTitle("%s - %d directories", color.New(color.Bold, color.FgHiYellow).Sprintf(title), len(comparisonResult.DirectoryResults))
	tableWriter.SetStyle(table.StyleRounded)
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{"path", "size", "dirs", "files", "scan duration"})

	for i, dirResult := range comparisonResult.ScanResult.DirectoryResults {
		dirCompResult := comparisonResult.DirectoryResults[i]
		tableWriter.AppendRow([]interface{}{
			color.HiBlueString(scanner.ShorifyPath(dirResult.DirectoryPath)),
			HumanSize(dirResult.TotalSize) + " " + color.HiMagentaString(HumanSizeSign(dirCompResult.SizeDiff)),
			fmt.Sprintf("%d (%+d)", dirResult.FolderCount, dirCompResult.FolderCountDiff),
			fmt.Sprintf("%d (%+d)", dirResult.FileCount, dirCompResult.FileCountDiff),
			dirResult.ScanDuration,
		})
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
		TimeAgo(comparisonResult.CurrentScanTime),
	})
	tableWriter.AppendSeparator()

	deltaFreeSpace := HumanSizeSign(comparisonResult.FreeSpaceDiff)

	tableWriter.AppendRow(table.Row{
		"FREE SPACE",
		color.HiGreenString(HumanSize(comparisonResult.ScanResult.FreeSpace)) + " " + color.HiMagentaString(deltaFreeSpace),
		"", "",
		time.Since(comparisonResult.ScanResult.StartTime).Round(time.Millisecond),
	})
	tableWriter.Render()
}

func ColorHeader(str string, a ...interface{}) string {
	return color.New(color.FgBlue).Add(color.Bold).Sprintf(str, a...)
}

func ColorHeaderHi(str string, a ...interface{}) string {
	return color.New(color.FgHiBlue).Add(color.Bold).Sprintf(str, a...)
}

func ColorPale(str string, a ...interface{}) string {
	return color.New(color.FgHiYellow). /*.Add(color.Bold)*/ Sprintf(str, a...)
}
