package storage

import (
	"errors"
	"spacemon/internal/config"
	"spacemon/internal/scanner"
	"time"
)
import "github.com/gorepos/storage"

func init() {
	// Set the storage directory for the application.
	storage.SetOptions(storage.Options{
		Dir: config.GetAppDir() + "/storage",
	})
}

// LoadPreviousResults loads the most recent scan result from the storage.
// It returns a pointer to the ScanResult and an error if any occurred.
func LoadPreviousResults() (*scanner.ScanResult, error) {
	// Retrieve the list of available keys in the "scans/" namespace.
	keys := storage.Keys("scans/")
	if len(keys) == 0 {
		// If there are no keys, return nil
		return nil, errors.New("there is no previous result")
	}
	// Get the latest scan result using the last key in the sorted list.
	prevkey := keys[len(keys)-1]
	var prevResult scanner.ScanResult
	err := storage.Get(prevkey, &prevResult)
	if err != nil {
		return nil, err
	}
	return &prevResult, nil
}

// buildKey constructs a key for the given scan result using its start time.
func buildKey(result scanner.ScanResult) string {
	return "scans/" + result.StartTime.Format(time.RFC3339Nano)
}

// SaveResult saves the given scan result to the storage using the generated key.
// It panics if an error occurs while saving the result.
func SaveResult(result scanner.ScanResult) {
	err := storage.Put(buildKey(result), result)
	if err != nil {
		panic(err)
	}
}
