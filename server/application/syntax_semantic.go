package application

import (
	"fmt"
	"regexp"
	"strings"
)

type VarInfo struct {
	Type     string
	Declared bool
	Line     int
}

func SyntaxCheck(lines []string) []string {
	errors := []string{}
	stack := []rune{}
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Verificar indentación después de ':'
		if i > 0 && strings.HasSuffix(lines[i-1], ":") {
			currentIndent := len(line) - len(trimmed)
			previousIndent := len(lines[i-1]) - len(strings.TrimSpace(lines[i-1]))
			if currentIndent <= previousIndent {
				errors = append(errors, fmt.Sprintf("Línea %d: Se esperaba indentación después de ':'", i+1))
			}
		}

		// Verificar balance de paréntesis y corchetes
		for _, c := range trimmed {
			switch c {
			case '(', '[', '{':
				stack = append(stack, c)
			case ')':
				if len(stack) == 0 || stack[len(stack)-1] != '(' {
					errors = append(errors, fmt.Sprintf("Línea %d: Paréntesis de cierre sin apertura", i+1))
				} else {
					stack = stack[:len(stack)-1]
				}
			case ']':
				if len(stack) == 0 || stack[len(stack)-1] != '[' {
					errors = append(errors, fmt.Sprintf("Línea %d: Corchete de cierre sin apertura", i+1))
				} else {
					stack = stack[:len(stack)-1]
				}
			case '}':
				if len(stack) == 0 || stack[len(stack)-1] != '{' {
					errors = append(errors, fmt.Sprintf("Línea %d: Llave de cierre sin apertura", i+1))
				} else {
					stack = stack[:len(stack)-1]
				}
			}
		}
	}

	for _, c := range stack {
		errors = append(errors, fmt.Sprintf("Símbolo '%c' sin cierre correspondiente", c))
	}

	return errors
}
func SemanticCheck(lines []string) []string {
	errors := []string{}
	declaredVars := map[string]VarInfo{}
	usedVars := map[string]bool{}

	reAssign := regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.+)`)
	reFuncDef := regexp.MustCompile(`^def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)`)
	reInt := regexp.MustCompile(`^-?\d+$`)
	reFloat := regexp.MustCompile(`^-?\d+\.\d+$`)
	reStr := regexp.MustCompile(`^".*"$|^'.*'$`)
	reBool := regexp.MustCompile(`^(True|False)$`)
	reLiterals := regexp.MustCompile(`"([^"\\]|\\.)*"|'([^'\\]|\\.)*'`)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Detectar definición de función y registrar parámetros
		if match := reFuncDef.FindStringSubmatch(trimmed); len(match) == 3 {
			fnName := match[1]
			paramList := match[2]
			declaredVars[fnName] = VarInfo{Type: "function", Declared: true, Line: i + 1}

			params := strings.Split(paramList, ",")
			for _, p := range params {
				name := strings.TrimSpace(strings.Split(p, "=")[0])
				if name != "" {
					declaredVars[name] = VarInfo{Type: "param", Declared: true, Line: i + 1}
				}
			}
			continue
		}

		// Detectar asignaciones
		if matches := reAssign.FindStringSubmatch(trimmed); len(matches) == 3 {
			varName := matches[1]
			value := strings.TrimSpace(matches[2])
			usedVars[varName] = true

			var tipo string
			switch {
			case reInt.MatchString(value):
				tipo = "int"
			case reFloat.MatchString(value):
				tipo = "float"
			case reStr.MatchString(value):
				tipo = "str"
			case reBool.MatchString(value):
				tipo = "bool"
			default:
				tipo = "unknown"
			}

			if prev, exists := declaredVars[varName]; exists {
				if prev.Type != tipo && tipo != "unknown" {
					errors = append(errors,
						fmt.Sprintf("Línea %d: Variable '%s' fue declarada como '%s', pero se le asigna un valor de tipo '%s'",
							i+1, varName, prev.Type, tipo))
				}
			} else {
				declaredVars[varName] = VarInfo{Type: tipo, Declared: true, Line: i + 1}
			}
			continue
		}

		// Limpiar literales antes de detectar identificadores
		cleaned := reLiterals.ReplaceAllString(trimmed, "")

		identifierRegex := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\b`)
		matches := identifierRegex.FindAllStringSubmatch(cleaned, -1)

		for _, m := range matches {
			v := m[1]
			keywords := map[string]bool{
				"if": true, "elif": true, "else": true, "def": true, "return": true,
				"for": true, "while": true, "in": true, "print": true, "and": true,
				"or": true, "not": true, "True": true, "False": true, "None": true,
			}
			if keywords[v] {
				continue
			}
			if _, declared := declaredVars[v]; !declared {
				errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' usada sin ser asignada previamente", i+1, v))
			} else {
				usedVars[v] = true
			}
		}
	}

	// Detectar variables nunca usadas (opcional)
	for name, info := range declaredVars {
		if !usedVars[name] && info.Type != "function" {
			errors = append(errors, fmt.Sprintf("Línea %d: Variable '%s' declarada pero nunca usada", info.Line, name))
		}
	}

	return errors
}
