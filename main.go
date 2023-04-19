package main

import (
	"flag"
	"log"
	"spacemon/internal/config"
	"spacemon/internal/reporter"
	"spacemon/internal/scanner"
	"spacemon/internal/storage"
)

var argDryRun = flag.Bool("dry", false, "Dry run (don't save scan result)")

func init() {
	flag.Parse()
}

func main() {
	cfg := config.LoadConfig()
	scanResultsChan := make(chan scanner.ScanResult)
	go scanner.ScanDirectories(scanner.ScanSetup{
		Directories: cfg.Directories,
		Title:       cfg.Title,
	}, scanResultsChan)

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

	if !*argDryRun {
		log.Println("save")
		storage.SaveResult(lastResult)
		report.Save()
	}
}
