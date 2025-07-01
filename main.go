package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"drxco/utils" // Ajusta el path si es necesario
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

	// Cachear template una vez al iniciar el servidor
	tmplPath := filepath.Join("templates", "index.html")
	var err error
	tmpl, err = template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatalf("Error cargando template: %v", err)
	}

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
	// No reparsear el template cada vez, usar el cacheado en la variable global tmpl

	err := tmpl.Execute(w, nil)
	if err != nil {
		// Solo escribir error una vez y no llamar WriteHeader dos veces
		http.Error(w, "Error al renderizar template", http.StatusInternalServerError)
		return
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
