package storage

import (
	"spacemon/internal/scanner"
	"time"
)
import "github.com/gorepos/storage"

func LoadPreviousResults() (*scanner.ScanResult, error) {
	keys := storage.Keys("scans/")
	if len(keys) == 0 {
		// todo: consider to return err
		return nil, nil
	}
	prevkey := keys[len(keys)-1]
	var prevResult scanner.ScanResult
	err := storage.Get(prevkey, &prevResult)
	if err != nil {
		return nil, err
	}
	return &prevResult, nil

}

func buildKey(result scanner.ScanResult) string {
	return "scans/" + result.ScanTime.Format(time.RFC3339Nano)
}

func SaveResult(result scanner.ScanResult) {
	err := storage.Put(buildKey(result), result)
	if err != nil {
		panic(err)
	}
}
