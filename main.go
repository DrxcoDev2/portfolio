package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"drxco/utils" // Asegúrate de usar el path correcto si lo estás organizando en carpetas
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

	// Cachear template
	tmplPath := filepath.Join("templates", "index.html")
	tmpl = template.Must(template.ParseFiles(tmplPath))

	// Archivos estáticos
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	// Página principal
	http.HandleFunc("/", serveHome)

	// SSE para recarga automática
	http.HandleFunc("/events", handleSSE)

	log.Println("Servidor en http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error cargando template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error al renderizar template", http.StatusInternalServerError)
	}
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

	for {
		time.Sleep(1 * time.Second)
		modTime := utils.GetLatestModTime("templates", "styles")
		if modTime.After(lastModTime) {
			lastModTime = modTime
			_, _ = w.Write([]byte("data: reload\n\n"))
			flusher.Flush()
		}
	}
}
