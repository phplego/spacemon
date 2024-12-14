// Webserver for spacemon
// It serves the web interface and provides a websocket API for the frontend

package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorepos/asciigraph"
	"github.com/robert-nix/ansihtml"
	"golang.org/x/net/websocket"
	"io/fs"
	"log"
	"net/http"
	"spacemon/internal/config"
	"spacemon/internal/reporter"
	"spacemon/internal/scanner"
	"spacemon/internal/storage"
)

var scanContext, scanContextCancel = context.WithCancel(context.Background())

func init() {
	color.NoColor = false
}

func actionScan(ctx context.Context, ch chan<- string, dryRun bool) {
	cfg := config.LoadConfig()
	scanResultsChan := make(chan scanner.ScanResult)
	go scanner.ScanDirectories(ctx, scanner.ScanSetup{
		Directories: cfg.Directories,
		Title:       cfg.Title,
	}, scanResultsChan)

	prevResult, err := storage.LoadPreviousResults()
	if err != nil {
		// Handle error
	}

	var report reporter.ReportInterface

	if prevResult == nil {
		report = &reporter.SingleScanReport{}
	} else {
		report = reporter.NewComparisonReport(*prevResult)
	}

	var result scanner.ScanResult
	for result = range scanResultsChan {
		report.Update(result)
		html := ansihtml.ConvertToHTMLWithClasses([]byte(report.Render()), "", true)
		bytes, _ := json.Marshal(map[string]string{
			"output": string(html),
			"title":  cfg.Title,
		})
		ch <- string(bytes)
	}

	if !dryRun && result.Error == "" {
		storage.SaveResult(result)
		report.Save()
		storage.Cleanup(cfg.MaxHistorySize)
	}
	close(ch)
}

func actionLastReport(ch chan<- string) {

	prevResult, err := storage.LoadPreviousResults()
	if err != nil {
		ch <- "Empty result. No previous scan"
		close(ch)
		return
	}
	prevPrevResult, err := storage.LoadPreviousResultsN(2)

	var report reporter.ReportInterface

	if prevPrevResult == nil {
		report = &reporter.SingleScanReport{}
	} else {
		report = reporter.NewComparisonReport(*prevPrevResult)
	}

	report.Update(*prevResult)
	html := ansihtml.ConvertToHTMLWithClasses([]byte(report.Render()), "", true)
	bytes, _ := json.Marshal(map[string]string{
		"output": string(html),
		"title":  prevResult.ScanSetup.Title,
	})
	ch <- string(bytes)
	close(ch)
}

type GetterFunc func(result scanner.ScanResult) float64

func chart(allResults []scanner.ScanResult, fun GetterFunc, caption, unit string) string {
	var data []float64
	for i := len(allResults) - 1; i >= 0; i-- {
		res := allResults[i]
		val := fun(res)
		if val < 0 {
			continue
		}
		data = append(data, val)
	}

	graph := asciigraph.Plot(data,
		asciigraph.Height(5),
		asciigraph.UnitPostfix(" "+unit),
		asciigraph.SeriesColors(asciigraph.Green),
		asciigraph.Caption(caption),
	)
	return string(ansihtml.ConvertToHTMLWithClasses([]byte(graph), "", true))
}

func actionGraph(ch chan<- string) {
	results := []scanner.ScanResult{}
	n := 1
	for {
		r, _ := storage.LoadPreviousResultsN(n)
		if r == nil {
			break
		}
		results = append(results, *r)
		n++
	}
	output := ""
	output += chart(results, func(result scanner.ScanResult) float64 {
		return float64(result.FreeSpace / 1024 / 1024)
	}, "free space", "MB")
	output += "\n\n"

	for _, dir := range results[0].ScanSetup.Directories {
		output += chart(results, func(result scanner.ScanResult) float64 {
			if dr, ok := result.DirectoryResults.Get(dir); ok {
				return float64(dr.TotalSize) / 1024 / 1024
			}
			return -1
		}, dir, "MB")

		output += "\n\n"
	}

	bytes, _ := json.Marshal(map[string]string{
		"output": output,
		"title":  "Graph",
	})
	ch <- string(bytes)
	close(ch)
}

func wsHandler(ws *websocket.Conn) {
	defer ws.Close()
	var dry = ws.Request().URL.Query().Get("dry")
	var act = ws.Request().URL.Query().Get("action")
	var ch = make(chan string)

	scanContextCancel() // cancel previous scan. Only one scan at a time
	scanContext, scanContextCancel = context.WithCancel(context.Background())

	switch act {
	case "last":
		go actionLastReport(ch)
	case "graph":
		go actionGraph(ch)
	default:
		go actionScan(scanContext, ch, dry != "")
	}

	for msg := range ch {
		_, err := ws.Write([]byte(msg))
		if err != nil {
			log.Println("Socket Write error 484:", err)
			return
		}
	}
}

// embed static files
//
//go:embed static
var embeddedFiles embed.FS

func RunWebserver() {
	// Websocket handler /ws
	http.Handle("/ws", authMiddleware(websocket.Handler(wsHandler), config.LoadConfig().DaemonUsername, config.LoadConfig().DaemonPassword))

	// File server for HTML and JS files
	//fileServer := http.FileServer(http.Dir("static"))
	staticFiles, _ := fs.Sub(embeddedFiles, "static")
	fileServer := http.FileServer(http.FS(staticFiles))
	fileServer = authMiddleware(fileServer, config.LoadConfig().DaemonUsername, config.LoadConfig().DaemonPassword)

	// Root route handler
	http.Handle("/", http.StripPrefix("/", fileServer))
	cfg := config.LoadConfig()
	log.Printf("Starting server on http://%s:%d ...\n", cfg.DaemonBindAddr, cfg.DaemonPort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.DaemonBindAddr, cfg.DaemonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	RunWebserver()
}
