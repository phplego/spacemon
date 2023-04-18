package comparer

import (
	"fmt"
	"spacemon/internal/scanner"
	"time"
)

type ComparisonResult struct {
	ScanResult       *scanner.ScanResult
	DirectoryResults []DirectoryComparisonResult
	FreeSpaceDiff    int64
	PrevScanTime     time.Time
	CurrentScanTime  time.Time
}

type DirectoryComparisonResult struct {
	DirectoryPath   string
	FileCountDiff   int64
	FolderCountDiff int64
	SizeDiff        int64
}

func CompareResults(prevResult, result *scanner.ScanResult) ComparisonResult {
	// Create a map for fast lookup of previous scanning results
	prevResultMap := make(map[string]scanner.DirectoryResult)
	for _, r := range prevResult.DirectoryResults {
		prevResultMap[r.DirectoryPath] = r
	}

	// Create a slice to store comparison results
	directoryCompareResults := make([]DirectoryComparisonResult, 0, len(result.DirectoryResults))

	// Iterate through the current scanning results
	for _, dirResult := range result.DirectoryResults {
		prevDirResult, ok := prevResultMap[dirResult.DirectoryPath]
		if !ok {
			fmt.Printf("Warning: No previous result found for directory: '%s'\n", dirResult.DirectoryPath)
			continue
		}

		// Calculate the difference between the current and previous result
		dirCompareResult := DirectoryComparisonResult{
			DirectoryPath:   dirResult.DirectoryPath,
			FileCountDiff:   dirResult.FileCount - prevDirResult.FileCount,
			FolderCountDiff: dirResult.FolderCount - prevDirResult.FolderCount,
			SizeDiff:        dirResult.TotalSize - prevDirResult.TotalSize,
		}

		// Add the comparison result to the slice
		directoryCompareResults = append(directoryCompareResults, dirCompareResult)
	}

	// Calculate the difference in free disk space
	freeSpaceDiff := result.FreeSpace - prevResult.FreeSpace

	// Return the ComparisonResult
	return ComparisonResult{
		ScanResult:       result,
		DirectoryResults: directoryCompareResults,
		FreeSpaceDiff:    freeSpaceDiff,
		PrevScanTime:     prevResult.StartTime,
		CurrentScanTime:  result.StartTime,
	}
}
