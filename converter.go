package freezedconverter

import (
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
}

type stringMapper func(p *ParameterToken) string

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
		case "class":
			if currentFreezed != nil {
				name := tokenizer.Next()
				currentFreezed.Name = name
			}
		case "factory":
			if currentFreezed != nil {
				factoryName := tokenizer.Next()

				containsDot, at := contains([]rune(factoryName), '.')
				if containsDot {
					if factoryName[at:] == ".fromJson" {
						currentFreezed = &Freezed{
							Name: currentFreezed.Name,
						}
						continue
					}

					currentFreezed.Name = factoryName
				}

				parseParameters(tokenizer, currentFreezed)
			}
		case "=":
			if currentFreezed != nil {
				name := []rune(currentFreezed.Name)

				if containsDot, at := contains(name, '.'); containsDot {
					constructorName := tokenizer.Next()

					if constructorName[0] != '_' {
						currentFreezed.Name = constructorName
					} else {
						name[at+1] = unicode.ToUpper(name[at+1])
						currentFreezed.Name = string(name[:at]) + string(name[at+1:])
					}
				}

				if len(currentFreezed.Parameters) > 0 {
					result = append(result, *currentFreezed)
				}
				currentFreezed = &Freezed{}
			}
		}
	}

	return result
}

func parseParameters(tokenizer *Tokenizer, freezed *Freezed) {
	bracketCount := 0

	var currentParameter *ParameterToken

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
			if name == "<" {
				typeName = parseTypedParameter(tokenizer, typeName)
				name = tokenizer.Next()
				if name == "?" {
					typeName += "?"
					name = tokenizer.Next()
				}
			}

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

func parseTypedParameter(tokenizer *Tokenizer, typeName string) string {
	typeName += "<"
	bracketCount := 1

	for {
		token := tokenizer.Next()
		typeName += token

		if IsOpeningBracket(token) {
			bracketCount++
		} else if IsClosingBracket(token) {
			bracketCount--
			if bracketCount == 0 {
				return typeName
			}
		}
	}
}

func isNullable(typeName string) bool {
	runes := []rune(typeName)

	return runes[len(typeName)-1] == '?'
}

func TranslateToGo(freezed *Freezed, format string) string {
	sb := strings.Builder{}

	sb.WriteString("type " + freezed.Name + " struct {\n")

	goParameters := make([]ParameterToken, len(freezed.Parameters))

	for i, p := range freezed.Parameters {
		goParameters[i] = *translateToGoParameter(&p)
	}

	maxNameLength := maxStrLength(goParameters, func(p *ParameterToken) string { return p.Name })
	maxTypeLength := maxStrLength(goParameters, func(p *ParameterToken) string { return p.Type })

	for i, goP := range goParameters {
		formattedName := freezed.Parameters[i].Name

		if goP.Nullable {
			formattedName += ",omitempty"
		}

		line := "\t" + Pad(goP.Name, maxNameLength) + " " + Pad(goP.Type, maxTypeLength) + " `" + format + ":\"" + formattedName + "\"`\n"
		sb.WriteString(line)
	}
	sb.WriteString("}")

	return sb.String()
}

func ToGoName(str string) string {
	runes := []rune(str)

	runes[0] = unicode.ToUpper(runes[0])
	length := len(str)

	for i := 0; i < length-1; i++ {
		r := runes[i]
		isStartOfId := r == 'I' || (i == 0 && r == 'i')

		if isStartOfId {
			if runes[i+1] == 'd' {
				if i+2 < length {
					rAfterId := runes[i+2]
					if unicode.IsLower(rAfterId) {
						continue
					}
				}
				runes[i+1] = 'D'
			}
		}
	}

	return string(runes)
}

func maxStrLength(parameters []ParameterToken, mapper stringMapper) int {
	strs := make([]string, len(parameters))

	for i, p := range parameters {
		strs[i] = mapper(&p)
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

func Pad(name string, length int) string {
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
	goName := ToGoName(p.Name)
	var goType string

	switch {
	case p.Type == "int":
		goType = "int"
	case p.Type == "bool":
		goType = "bool"
	case p.Type == "double":
		goType = "float64"
	case p.Type[:4] == "List", p.Type[:3] == "Map":
		goType = "?"
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
