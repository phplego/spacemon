package main

import (
	"spacemon/internal/config"
	"spacemon/internal/reporter"
	"spacemon/internal/scanner"
	"spacemon/internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	scanResultsChan := make(chan scanner.ScanResult)
	go scanner.ScanDirectories(cfg.Directories, scanResultsChan)

	prevResult, err := storage.LoadPreviousResults()
	if err != nil {
		// Handle error
	}

	var report reporter.Report

	if prevResult == nil {
		report = &reporter.SingleScanReport{}
	} else {
		report = reporter.NewComparisonReport(*prevResult)
	}

	var lastResult scanner.ScanResult
	for result := range scanResultsChan {
		report.Update(result)
		lastResult = result
	}

	storage.SaveResult(lastResult)
	report.Save()
}
