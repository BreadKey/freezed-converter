package freezedconverter

import (
	"log"
	"strings"
	"unicode"
)

type ParameterToken struct {
	Name     string
	Type     string
	Nullable bool
}

type Freezed struct {
	Name       string
	Parameters []ParameterToken
	blockCount int
}

type convertForMap func(p *ParameterToken) string

func Parse(filePath string) []Freezed {
	result := make([]Freezed, 0, 2)
	tokenizer := NewTokenizer(filePath)

	var currentFreezed *Freezed

parse:
	for {
		token := tokenizer.Next()

		switch token {
		case "":
			break parse
		case "@freezed":
			currentFreezed = &Freezed{}
		case "{":
			if currentFreezed != nil {
				currentFreezed.blockCount++
			}
		case "}":
			if currentFreezed != nil {
				currentFreezed.blockCount--

				if currentFreezed.blockCount == 0 {
					result = append(result, *currentFreezed)
					currentFreezed = nil
				}
			}
		case "class":
			if currentFreezed != nil {
				name := tokenizer.Next()
				currentFreezed.Name = name
			}
		case "factory":
			if currentFreezed != nil {
				parseParameters(tokenizer, currentFreezed)
			}
		}
	}

	return result
}

func parseParameters(tokenizer *Tokenizer, freezed *Freezed) {
	bracketCount := 0

	var currentParameter *ParameterToken

	freezedName := tokenizer.Next()

	if freezedName != freezed.Name {
		if freezedName == freezed.Name+".fromJson" {
			return
		}

		log.Fatalf("Not implemented multiple freezed! %s", freezedName)
	}

	for {
		token := tokenizer.Next()

		if IsBracket(token) {
			if IsOpeningBracket(token) {
				bracketCount++
			} else {
				bracketCount--
			}

			if bracketCount == 0 {
				return
			}
			continue
		}

		if currentParameter == nil {
			currentParameter = &ParameterToken{}

			typeName, include := getParameterType(tokenizer, token)

			name := tokenizer.Next()

			currentParameter.Type = typeName
			currentParameter.Name = name

			if isNullable(typeName) {
				currentParameter.Nullable = true
				currentParameter.Type = typeName[:len(typeName)-1]
			}

			if include {
				freezed.Parameters = append(freezed.Parameters, *currentParameter)
			}
			currentParameter = nil
		}
	}
}

func getParameterType(tokenizer *Tokenizer, currentToken string) (string, bool) {
	include := true
	typeName := currentToken
	for {
		if typeName[0] == '@' {
			tokenizer.Next()
			bracketCount := 1

			for {
				next := tokenizer.Next()

				if IsOpeningBracket(next) {
					bracketCount++
				} else if IsClosingBracket(next) {
					bracketCount--

					if bracketCount == 0 {
						typeName = tokenizer.Next()
						break
					}
				} else if next == "includeToJson" {
					includeToJson := tokenizer.Next()
					include = includeToJson != "false"
				}
			}
		} else {
			break
		}
	}

	if typeName == "required" {
		typeName = tokenizer.Next()
	}

	return typeName, include
}

func isNullable(typeName string) bool {
	runes := []rune(typeName)

	return runes[len(typeName)-1] == '?'
}

func TranslateToGo(freezed *Freezed) string {
	sb := strings.Builder{}

	sb.WriteString("type " + freezed.Name + " struct {\n")

	goParameters := make([]ParameterToken, len(freezed.Parameters))

	for i, p := range freezed.Parameters {
		goParameters[i] = *translateToGoParameter(&p)
	}

	maxNameLength := maxStrLength(goParameters, func(p *ParameterToken) string { return p.Name })
	maxTypeLength := maxStrLength(goParameters, func(p *ParameterToken) string { return p.Type })

	for i, goP := range goParameters {
		firestoreName := freezed.Parameters[i].Name

		if goP.Nullable {
			firestoreName += ",omitempty"
		}

		line := "\t" + pad(goP.Name, maxNameLength) + " " + pad(goP.Type, maxTypeLength) + " `firestore:\"" + firestoreName + "\"`\n"
		sb.WriteString(line)
	}
	sb.WriteString("}")

	return sb.String()
}

func capitalize(str string) string {
	runes := []rune(str)

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func maxStrLength(parameters []ParameterToken, convert convertForMap) int {
	strs := make([]string, len(parameters))

	for i, p := range parameters {
		strs[i] = convert(&p)
	}

	length := 0
	for _, s := range strs {
		l := len(s)

		if l > length {
			length = l
		}
	}

	return length
}

func pad(name string, length int) string {
	lName := len(name)

	if lName == length {
		return name
	}

	sb := &strings.Builder{}
	sb.WriteString(name)

	for i := lName; i < length; i++ {
		sb.WriteRune(' ')
	}

	return sb.String()
}

func translateToGoParameter(p *ParameterToken) *ParameterToken {
	goName := capitalize(p.Name)
	var goType string

	switch p.Type {
	case "int":
		goType = "int"
	case "double":
		goType = "float64"
	default:
		goType = "string"
	}

	if p.Nullable {
		goType = "*" + goType
	}

	return &ParameterToken{
		Name:     goName,
		Type:     goType,
		Nullable: p.Nullable,
	}
}
