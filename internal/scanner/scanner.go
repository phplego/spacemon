package scanner

import (
	"github.com/shirou/gopsutil/disk"
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
	StartTime        time.Time
	IsCompleted      bool
}

type DirectoryResult struct {
	DirectoryPath string
	FileCount     int64
	FolderCount   int64
	TotalSize     int64
	ScanDuration  time.Duration
}

func ScanDirectories(directories []string, resultsChan chan<- ScanResult) {
	result := ScanResult{}
	result.DirectoryResults = make([]DirectoryResult, 0)
	result.StartTime = time.Now()
	result.FreeSpace, _ = GetFreeSpace()

	resultsChan <- result

	for _, directory := range directories {
		directoryResult := ScanDirectory(directory)
		result.DirectoryResults = append(result.DirectoryResults, directoryResult)
		resultsChan <- result
	}

	result.IsCompleted = true
	resultsChan <- result

	close(resultsChan)
}

func ScanDirectory(directory string) DirectoryResult {
	directory = AbsPath(directory)
	var directoryResult = DirectoryResult{
		DirectoryPath: directory,
	}
	startTime := time.Now()
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
	directoryResult.ScanDuration = time.Now().Sub(startTime).Round(time.Millisecond)
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

func ShorifyPath(absPath string) string {
	usr, err := user.Current()
	if err != nil {
		return absPath
	}
	if strings.HasPrefix(absPath, usr.HomeDir) {
		return strings.Replace(absPath, usr.HomeDir, "~", 1)
	}
	return absPath
}

func GetFreeSpace() (int64, error) {
	di, err := disk.Usage(".")
	if err != nil {
		return 0, err
	}
	return int64(di.Free), nil
}
