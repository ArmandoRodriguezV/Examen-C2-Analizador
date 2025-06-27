package application

import (
	"Analizador/domain"
	"regexp"
	"strings"
)

// Palabras reservadas de Python
var reservedWords = map[string]bool{
	"False": true, "None": true, "True": true, "and": true, "as": true,
	"assert": true, "async": true, "await": true, "break": true, "class": true,
	"continue": true, "def": true, "del": true, "elif": true, "else": true,
	"except": true, "finally": true, "for": true, "from": true, "global": true,
	"if": true, "import": true, "in": true, "is": true, "lambda": true,
	"nonlocal": true, "not": true, "or": true, "pass": true, "raise": true,
	"return": true, "try": true, "while": true, "with": true, "yield": true,
}

// Símbolos válidos de Python
var symbolRegex = regexp.MustCompile(`[()

\[\]

:.,+\-*/%=<>!|]`)

func AnalyzeCode(code string) ([]domain.Token, map[string]int) {
	var tokens []domain.Token
	conteo := map[string]int{
		"PR": 0, "Símbolo": 0, "ID": 0,
		"Número": 0, "Cadenas": 0, "Comentario": 0, "Error": 0,
	}

	lines := strings.Split(code, "\n")

	for lineNum, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			tokens = append(tokens, domain.Token{Token: originalLine, Tipo: "Comentario", Linea: lineNum + 1})
			conteo["Comentario"]++
			continue
		}

		if line != "" {
			tokens = append(tokens, processLine(line, lineNum+1, conteo)...)
		}
	}

	return tokens, conteo
}

func processLine(line string, lineNum int, conteo map[string]int) []domain.Token {
	var tokens []domain.Token
	i := 0
	for i < len(line) {
		if line[i] == ' ' || line[i] == '\t' {
			i++
			continue
		}
		if line[i] == '"' || line[i] == '\'' {
			quote := line[i]
			start := i
			i++
			for i < len(line) && line[i] != quote {
				if line[i] == '\\' && i+1 < len(line) {
					i += 2
				} else {
					i++
				}
			}
			if i < len(line) {
				i++
			} else {
				tokens = append(tokens, domain.Token{Token: line[start:], Tipo: "Error", Linea: lineNum})
				conteo["Error"]++
				break
			}
			tokens = append(tokens, domain.Token{Token: line[start:i], Tipo: "Cadenas", Linea: lineNum})
			conteo["Cadenas"]++
			continue
		}
		if isDigit(line[i]) || (line[i] == '.' && i+1 < len(line) && isDigit(line[i+1])) {
			start := i
			hasDecimal := false
			if line[i] == '.' {
				hasDecimal = true
				i++
			}
			for i < len(line) && (isDigit(line[i]) || (line[i] == '.' && !hasDecimal)) {
				if line[i] == '.' {
					hasDecimal = true
				}
				i++
			}
			token := line[start:i]
			tokens = append(tokens, domain.Token{Token: token, Tipo: "Número", Linea: lineNum})
			conteo["Número"]++
			continue
		}
		if i+1 < len(line) {
			twoChar := line[i : i+2]
			if isTwoCharOperator(twoChar) {
				tokens = append(tokens, domain.Token{Token: twoChar, Tipo: "Símbolo", Linea: lineNum})
				conteo["Símbolo"]++
				i += 2
				continue
			}
		}
		if symbolRegex.MatchString(string(line[i])) {
			tokens = append(tokens, domain.Token{Token: string(line[i]), Tipo: "Símbolo", Linea: lineNum})
			conteo["Símbolo"]++
			i++
			continue
		}
		if isLetter(line[i]) || line[i] == '_' {
			start := i
			for i < len(line) && (isLetter(line[i]) || isDigit(line[i]) || line[i] == '_') {
				i++
			}
			token := line[start:i]
			tipo := "ID"
			if reservedWords[token] {
				tipo = "PR"
				conteo["PR"]++
			} else {
				conteo["ID"]++
			}
			tokens = append(tokens, domain.Token{Token: token, Tipo: tipo, Linea: lineNum})
			continue
		}
		tokens = append(tokens, domain.Token{Token: string(line[i]), Tipo: "Error", Linea: lineNum})
		conteo["Error"]++
		i++
	}
	return tokens
}

func isTwoCharOperator(op string) bool {
	return map[string]bool{
		"==": true, "!=": true, "<=": true, ">=": true,
		"**": true, "//": true, ":=": true,
		"+=": true, "-=": true, "*=": true, "/=": true,
		"%=": true,
	}[op]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
