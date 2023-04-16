package scanner

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type ScanError struct {
	DirectoryPath string
	Err           error
}

type ScanResult struct {
	DirectoryResults []DirectoryResult
	Errors           []ScanError
	FreeSpace        int64
	ScanTime         time.Time
}

type DirectoryResult struct {
	DirectoryPath string
	FileCount     int
	FolderCount   int
	TotalSize     int64
	ScanDuration  time.Duration
}

func ScanDirectories(directories []string, resultsChan chan<- ScanResult) {
	result := ScanResult{}
	result.DirectoryResults = make([]DirectoryResult, 0)
	result.ScanTime = time.Now()
	resultsChan <- result

	for _, directory := range directories {
		directoryResult := ScanDirectory(directory)
		result.DirectoryResults = append(result.DirectoryResults, directoryResult)
		resultsChan <- result
	}

	close(resultsChan)
}

func ScanDirectory(directory string) DirectoryResult {
	directory = AbsPath(directory)
	var directoryResult = DirectoryResult{
		DirectoryPath: directory,
	}
	_ = filepath.Walk(directory, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			//log.Println(err)
			//return err // return error if you want to break walking
		} else {

			if fileInfo.IsDir() {
				directoryResult.FolderCount++
			} else {
				directoryResult.FileCount++
				directoryResult.TotalSize += fileInfo.Size()
			}
		}
		return nil
	})

	return directoryResult
}

func AbsPath(path string) string {
	usr, _ := user.Current()

	if path == "~" {
		return usr.HomeDir
	} else if strings.HasPrefix(path, "~/") {
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}
