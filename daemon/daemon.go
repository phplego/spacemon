package main

import (
	"fmt"
	"github.com/creack/pty"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"spacemon/internal/scanner"
	"sync"
	"time"
)

func wsHandler(ws *websocket.Conn) {
	defer ws.Close()
	var _ = scanner.ScanResult{}
	for {
		message := fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC1123))
		_, err := ws.Write([]byte(message))
		if err != nil {
			log.Println("Socket Write error 484:", err)
			break
		}
		time.Sleep(1 * time.Second)
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
