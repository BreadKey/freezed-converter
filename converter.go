package freezedconverter

import "log"

type ParameterToken struct {
	Name string
	Type string
}

type Freezed struct {
	Name       string
	Parameters []ParameterToken
	BlockCount int
}

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
				currentFreezed.BlockCount++
			}
		case "}":
			if currentFreezed != nil {
				currentFreezed.BlockCount--

				if currentFreezed.BlockCount == 0 {
					result = append(result, *currentFreezed)
					currentFreezed = nil
				}
			}
		case "class":
			name := tokenizer.Next()
			currentFreezed.Name = name
		case "factory":
			parseParameters(tokenizer, currentFreezed)
		}
	}

	return result
}

func parseParameters(tokenizer *Tokenizer, freezed *Freezed) {
	bracketCount := 0

	var currentParameter *ParameterToken

	freezedName := tokenizer.Next()

	if freezedName != freezed.Name {
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
			typeName := tokenizer.Next()
			if typeName == "required" {
				typeName = tokenizer.Next()
			}

			name := tokenizer.Next()
			currentParameter.Type = typeName
			currentParameter.Name = name

			freezed.Parameters = append(freezed.Parameters, *currentParameter)
			currentParameter = nil
		}
	}
}
