package comparer

import (
	"fmt"
	"spacemon/internal/scanner"
	"time"
)

type ComparisonResult struct {
	ScanResult      *scanner.ScanResult
	DirectoryDiffs  map[string]DirectoryDiff
	FreeSpaceDiff   int64
	PrevScanTime    time.Time
	CurrentScanTime time.Time
}

type DirectoryDiff struct {
	DirectoryPath   string
	FileCountDiff   int64
	FolderCountDiff int64
	SizeDiff        int64
}

func CompareResults(prevResult, result *scanner.ScanResult) ComparisonResult {

	// Create a map to store comparison results
	directoryCompareResults := make(map[string]DirectoryDiff, 0)

	// Iterate through the current scanning results
	for dir, dirResult := range result.DirectoryResults {
		prevDirResult, ok := prevResult.DirectoryResults[dir]
		if !ok {
			fmt.Printf("Warning: No previous result found for directory: '%s'\n", dirResult.DirectoryPath)
			continue
		}

		// Calculate the difference between the current and previous result
		dirCompareResult := DirectoryDiff{
			DirectoryPath:   dirResult.DirectoryPath,
			FileCountDiff:   dirResult.FileCount - prevDirResult.FileCount,
			FolderCountDiff: dirResult.FolderCount - prevDirResult.FolderCount,
			SizeDiff:        dirResult.TotalSize - prevDirResult.TotalSize,
		}

		// Add the comparison result to the slice
		directoryCompareResults[dir] = dirCompareResult
	}

	// Calculate the difference in free disk space
	freeSpaceDiff := result.FreeSpace - prevResult.FreeSpace

	// Return the ComparisonResult
	return ComparisonResult{
		ScanResult:      result,
		DirectoryDiffs:  directoryCompareResults,
		FreeSpaceDiff:   freeSpaceDiff,
		PrevScanTime:    prevResult.StartTime,
		CurrentScanTime: result.StartTime,
	}
}
