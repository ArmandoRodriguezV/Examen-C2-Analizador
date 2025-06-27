package main

import (
	"Analizador/infrastructure"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/tokens", infrastructure.TokenHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Println("Endpoints disponibles:")
	log.Println("  POST /tokens - Analizar código")
	log.Println("  GET /health - Estado del servidor")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
