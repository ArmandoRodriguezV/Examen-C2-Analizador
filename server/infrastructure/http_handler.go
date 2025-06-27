package infrastructure

import (
	"Analizador/application"
	"Analizador/domain"
	"encoding/json"
	"net/http"
	"strings"
)

type AnalyzeRequest struct {
	Query string `json:"query"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Validar que hay contenido para analizar
	if strings.TrimSpace(req.Query) == "" {
		response := domain.AnalyzeResponse{
			Tokens:   []domain.Token{},
			Conteo:   map[string]int{"PR": 0, "Global": 0, "Símbolo": 0, "ID": 0, "Número": 0, "Cadenas": 0, "Comentario": 0, "Error": 0},
			Syntax:   []string{},
			Semantic: []string{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	tokens, conteo := application.AnalyzeCode(req.Query)

	lines := strings.Split(req.Query, "\n")

	syntaxErrors := application.SyntaxCheck(lines)
	semanticErrors := application.SemanticCheck(lines)

	response := domain.AnalyzeResponse{
		Tokens:   tokens,
		Conteo:   conteo,
		Syntax:   syntaxErrors,
		Semantic: semanticErrors,
	}

	// Enviar respuesta
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
		return
	}
}
