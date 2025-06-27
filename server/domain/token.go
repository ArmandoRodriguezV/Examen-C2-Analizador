// domain/token.go
package domain

type Token struct {
	Token string `json:"token"`
	Tipo  string `json:"tipo"`
	Linea int    `json:"linea"`
}

type AnalyzeResponse struct {
	Tokens   []Token        `json:"tokens"`
	Conteo   map[string]int `json:"conteo"`
	Syntax   []string       `json:"syntax"`
	Semantic []string       `json:"semantic"`
}
