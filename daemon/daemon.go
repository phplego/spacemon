package main

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/robert-nix/ansihtml"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"spacemon/internal/config"
	"spacemon/internal/reporter"
	"spacemon/internal/scanner"
	"spacemon/internal/storage"
	"sync"
)

func wsHandler(ws *websocket.Conn) {
	defer ws.Close()
	cfg := config.LoadConfig()
	scanResultsChan := make(chan scanner.ScanResult)
	go scanner.ScanDirectories(scanner.ScanSetup{
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

	for result := range scanResultsChan {
		report.Update(result)
		//html := report.Render()
		html := ansihtml.ConvertToHTML([]byte(report.Render()))
		_, err := ws.Write([]byte(html))
		if err != nil {
			log.Println("Socket Write error 484:", err)
			break
		}
	}
}

// example how to proxy spacemon console output
func exec1() {
	cmd := exec.Command("./spacemon")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer func() { _ = ptmx.Close() }()
	_, _ = io.Copy(os.Stdout, ptmx)
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Command exited with error: %v\n", err)
		os.Exit(1)
	}
}

func RunWebserver() {
	// Websocket handler /ws
	http.Handle("/ws", websocket.Handler(wsHandler))

	// File server for HTML and JS files
	fileServer := http.FileServer(http.Dir("static"))

	// Root route handler
	http.Handle("/", http.StripPrefix("/", fileServer))

	log.Println("Starting server on http://localhost:8080 ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		RunWebserver()
	}()

	wg.Wait() // Блокирует, пока веб-сервер работает
}
