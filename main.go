package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Para desarrollo local
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "�Hola desde Go en Render!")
	})

	log.Println("Servidor corriendo en el puerto " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
