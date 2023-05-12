package scanner

import (
	"github.com/shirou/gopsutil/disk"
	"golang.org/x/net/context"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type ScanResult struct {
	ScanSetup        ScanSetup
	DirectoryResults *SafeMap[string, DirectoryResult]
	FreeSpace        int64
	StartTime        time.Time
	Completed        bool
	Error            string
}

type DirectoryResult struct {
	DirectoryPath string
	FileCount     int64
	FolderCount   int64
	TotalSize     int64
	ScanDuration  time.Duration
	Completed     bool
	Canceled      bool
	Error         string
}

type ScanSetup struct {
	Directories []string
	Title       string
}

func ScanDirectories(ctx context.Context, setup ScanSetup, resultsChan chan<- ScanResult) {
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
		go ScanDirectory(ctx, directory, ch)
		for res := range ch {
			if res.Error != "" {
				result.Error = res.Error
			}
			result.DirectoryResults.Set(directory, res)
			resultsChan <- result
		}

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			break
		default:
		}
	}

	result.Completed = true
	resultsChan <- result

	close(resultsChan)
}

func ScanDirectory(ctx context.Context, directory string, dirResultChan chan<- DirectoryResult) {
	directory = AbsPath(directory)
	var directoryResult = DirectoryResult{
		DirectoryPath: directory,
	}
	startTime := time.Now()
	tmpTime := time.Now()
	err := filepath.Walk(directory, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			//log.Println(err)
			//return err // return error if you want to break walking
		} else {

			// Check if context is cancelled
			select {
			case <-ctx.Done():
				directoryResult.Canceled = true
				directoryResult.Error = ctx.Err().Error()
				return ctx.Err() // returning error means break walking
			default:
			}

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
	if err != nil {
		directoryResult.Error = err.Error()
	}
	directoryResult.Completed = true
	dirResultChan <- directoryResult
	close(dirResultChan)
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
