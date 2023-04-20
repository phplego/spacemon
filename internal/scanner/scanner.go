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
	ScanSetup        ScanSetup
	DirectoryResults *SafeMap[string, DirectoryResult]
	Errors           []ScanError
	FreeSpace        int64
	StartTime        time.Time
	Completed        bool
}

type DirectoryResult struct {
	DirectoryPath string
	FileCount     int64
	FolderCount   int64
	TotalSize     int64
	ScanDuration  time.Duration
	Completed     bool
}

type ScanSetup struct {
	Directories []string
	Title       string
}

func ScanDirectories(setup ScanSetup, resultsChan chan<- ScanResult) {
	result := ScanResult{
		ScanSetup:        setup,
		DirectoryResults: NewSafeMap[string, DirectoryResult](),
		StartTime:        time.Now(),
		FreeSpace:        0, // Update FreeSpace later
	}
	result.FreeSpace, _ = GetFreeSpace()

	resultsChan <- result

	for _, directory := range setup.Directories {
		ch := make(chan DirectoryResult)
		go ScanDirectory(directory, ch)
		for res := range ch {
			result.DirectoryResults.Set(directory, res)
			resultsChan <- result
		}
	}

	result.Completed = true
	resultsChan <- result

	close(resultsChan)
}

func ScanDirectory(directory string, dirResultChan chan<- DirectoryResult) DirectoryResult {
	directory = AbsPath(directory)
	var directoryResult = DirectoryResult{
		DirectoryPath: directory,
	}
	startTime := time.Now()
	tmpTime := time.Now()
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
			directoryResult.ScanDuration = time.Now().Sub(startTime).Round(time.Millisecond)

			if time.Now().After(tmpTime.Add(50 * time.Millisecond)) {
				tmpTime = time.Now()
				dirResultChan <- directoryResult
			}

		}
		return nil
	})
	directoryResult.Completed = true
	dirResultChan <- directoryResult
	close(dirResultChan)
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
