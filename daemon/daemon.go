package main

import (
	"encoding/base64"
	"fmt"
	"github.com/creack/pty"
	"github.com/fatih/color"
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
	"strings"
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

	var result scanner.ScanResult
	for result = range scanResultsChan {
		report.Update(result)
		//html := report.Render()
		html := ansihtml.ConvertToHTML([]byte(report.Render()))
		_, err := ws.Write([]byte(html))
		if err != nil {
			log.Println("Socket Write error 484:", err)
		}
	}

	// todo: if not dry run
	storage.SaveResult(result)
	report.Save()
	storage.Cleanup(cfg.MaxHistorySize)
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
	fileServer = basicAuthMiddleware(fileServer)

	// Root route handler
	http.Handle("/", http.StripPrefix("/", fileServer))
	cfg := config.LoadConfig()
	log.Printf("Starting server on http://localhost:%d ...\n", cfg.DaemonPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.DaemonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func init() {
	color.NoColor = false
}

func basicAuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" {
			authParts := strings.SplitN(auth, " ", 2)
			if len(authParts) == 2 && authParts[0] == "Basic" {
				userPass, err := base64.StdEncoding.DecodeString(authParts[1])
				if err == nil {
					parts := strings.SplitN(string(userPass), ":", 2)
					cfg := config.LoadConfig()
					if len(parts) == 2 && parts[0] == cfg.DaemonBasicUsername && parts[1] == cfg.DaemonBasicPassword {
						handler.ServeHTTP(w, r)
						return
					}
				}
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
	})
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
