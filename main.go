package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"drxco/utils"
)

var (
	lastModTime time.Time
	tmpl        *template.Template
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmplPath := filepath.Join("templates", "index.html")
	var err error
	tmpl, err = template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error cargando template: %v", err)
	}

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/events", handleSSE)

	log.Println("Servidor en http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, nil)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}
	buf.WriteTo(w)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming no soportado", http.StatusInternalServerError)
		return
	}

	notify := r.Context().Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			modTime := utils.GetLatestModTime("templates", "styles")
			if modTime.After(lastModTime) {
				lastModTime = modTime
				_, err := w.Write([]byte("data: reload\n\n"))
				if err != nil {
					return
				}
			} else {
				_, err := w.Write([]byte(": ping\n\n"))
				if err != nil {
					return
				}
			}
			flusher.Flush()
		}
	}
}
