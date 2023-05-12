package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/robert-nix/ansihtml"
	"golang.org/x/net/websocket"
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

func cmdScan(ch chan string, dryRun bool) {
	cfg := config.LoadConfig()
	scanContextCancel()
	scanContext, scanContextCancel = context.WithCancel(context.Background())
	scanResultsChan := make(chan scanner.ScanResult)
	go scanner.ScanDirectories(scanContext, scanner.ScanSetup{
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
		html := ansihtml.ConvertToHTML([]byte(report.Render()))
		bytes, _ := json.Marshal(map[string]string{
			"output": string(html),
			"title":  cfg.Title,
		})
		ch <- string(bytes)
	}

	if !dryRun {
		storage.SaveResult(result)
		report.Save()
		storage.Cleanup(cfg.MaxHistorySize)
	}
	close(ch)
}

func wsHandler(ws *websocket.Conn) {
	defer ws.Close()
	var dry = ws.Request().URL.Query().Get("dry")
	println("dry=", dry)
	var ch = make(chan string)
	go cmdScan(ch, dry != "")
	for msg := range ch {
		_, err := ws.Write([]byte(msg))
		if err != nil {
			log.Println("Socket Write error 484:", err)
			return
		}
	}
}

func RunWebserver() {
	// Websocket handler /ws
	http.Handle("/ws", websocket.Handler(wsHandler))

	// File server for HTML and JS files
	fileServer := http.FileServer(http.Dir("static"))
	fileServer = basicAuthMiddleware(fileServer, config.LoadConfig().DaemonBasicUsername, config.LoadConfig().DaemonBasicPassword)

	// Root route handler
	http.Handle("/", http.StripPrefix("/", fileServer))
	cfg := config.LoadConfig()
	log.Printf("Starting server on http://localhost:%d ...\n", cfg.DaemonPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.DaemonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	RunWebserver()
}
