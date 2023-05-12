package main

import (
	"context"
	"flag"
	"spacemon/internal/config"
	"spacemon/internal/reporter"
	"spacemon/internal/scanner"
	"spacemon/internal/storage"
	"spacemon/internal/util"
)

var argDryRun = flag.Bool("dry", false, "Dry run (don't save scan result)")
var argJson = flag.Bool("json", false, "Print the report as JSON")

func init() {
	flag.Parse()
}

func main() {
	cfg := config.LoadConfig()
	scanResultsChan := make(chan scanner.ScanResult)
	//ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*300)
	ctx := context.TODO()
	go scanner.ScanDirectories(ctx, scanner.ScanSetup{
		Directories: cfg.Directories,
		Title:       cfg.Title,
	}, scanResultsChan)

	prevResult, err := storage.LoadPreviousResults()
	if err != nil {
		// Handle error
	}

	var report reporter.ReportInterface

	if prevResult == nil {
		report = &reporter.SingleScanReport{}
	} else {
		report = reporter.NewComparisonReport(*prevResult)
	}

	var result scanner.ScanResult
	for result = range scanResultsChan {
		report.Update(result)
		if *argJson {
			util.ClearAndPrint(report.RenderJson())
		} else {
			util.ClearAndPrint(report.Render())
		}
	}

	if !*argDryRun {
		storage.SaveResult(result)
		report.Save()
		storage.Cleanup(cfg.MaxHistorySize)
	}

}
