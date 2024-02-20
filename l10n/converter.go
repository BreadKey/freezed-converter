package l10nconverter

import (
	fc "breadkey/freezed/converter"
	"strings"
)

type Select struct {
	Name        string
	PlaceHolder string
	Pairs       []SelectPair
	Prefix      string
	Suffix      string
}

type SelectPair struct {
	Key   string
	Value string
}

func ParseSelect(filePath string) []Select {
	result := make([]Select, 0, 2)
	tokenizer := fc.NewTokenizer(filePath)

	tokenizer.Next()
parse:
	for {
		token := tokenizer.Next()

		switch token {
		case "", "}":
			break parse
		default:
			key := token
			if key[1] == '@' && key[2] != '@' {
				tokenizer.Next()
				bracketCount := 1
				for {
					nextToken := tokenizer.Next()
					switch {
					case fc.IsOpeningBracket(nextToken):
						bracketCount++
					case fc.IsClosingBracket(nextToken):
						bracketCount--
						if bracketCount == 0 {
							continue parse
						}
					}
				}
			}
			value := tokenizer.Next()
			s := parseSelect(key, value)
			if s != nil {
				result = append(result, *s)
			}
		}
	}

	return result
}

func parseSelect(key string, value string) *Select {
	last := ""
	current := ""

	s := &Select{
		Name: key[1 : len(key)-1],
	}
	pairs := make([]SelectPair, 0, 2)
	meetBracket := false
	isSelect := false

	i := 0

checkSelect:
	for {
		if i >= len(value) {
			break
		}

		r := value[i]
		if i > 0 {
			switch r {
			case ' ':
				if !meetBracket {
					current += string(r)
				}
			case '{':
				meetBracket = true
				s.Prefix = current
				current = ""
			case ',':
				if current == "select" {
					isSelect = true
					s.PlaceHolder = last
					if value[i+1] == ' ' {
						i = i + 2
					} else {
						i = i + 1
					}
					break checkSelect
				} else {
					last = current
					current = ""
				}
			default:
				current += string(r)
			}

		}
		i++
	}

	if isSelect {
		parts := make([]string, 0, 2)

		currentPart := ""
		selectValue := value[i : len(value)-1]

		bracketCount := 0

		for _, r := range selectValue {
			switch r {
			case ' ':
				if bracketCount >= 0 && len(currentPart) == 0 {
					continue
				} else {
					if bracketCount >= 0 {
						currentPart += string(r)
					} else {
						s.Suffix += string(r)
					}
				}
			case '}':
				bracketCount--
				if bracketCount >= 0 {
					currentPart += string(r)
					parts = append(parts, currentPart)
					currentPart = ""
				}
			default:
				if r == '{' {
					bracketCount++
				}

				if bracketCount >= 0 {
					currentPart += string(r)
				} else {
					s.Suffix += string(r)
				}
			}
		}

		for _, part := range parts {
			key := ""
			keyEnd := false
			value := ""

			for _, r := range part {
				if !keyEnd {
					if r == '{' {
						keyEnd = true
					} else {
						key += string(r)
					}
				} else {
					if r != '}' {
						value += string(r)
					}
				}
			}

			pairs = append(pairs, SelectPair{
				Key:   key,
				Value: value,
			})
		}
		s.Pairs = pairs

		return s
	} else {
		return nil
	}
}

func TranslateToGo(selectSyntax *Select) string {
	sb := strings.Builder{}

	placeHolderName := fc.ToGoName(selectSyntax.PlaceHolder)

	sb.WriteString("type " + placeHolderName + " int\n")
	sb.WriteString("const (\n")

	var otherValue string

	assignIota := false
	for _, pair := range selectSyntax.Pairs {
		if pair.Key == "other" {
			otherValue = pair.Value
			continue
		}
		if !assignIota {
			sb.WriteString("\t" + fc.ToGoName(pair.Key) + " " + placeHolderName + " = iota\n")
			assignIota = true
		} else {
			sb.WriteString("\t" + fc.ToGoName(pair.Key) + "\n")
		}
	}
	sb.WriteString(")\n")
	sb.WriteString("func " + placeHolderName + "ToL10n(" + selectSyntax.PlaceHolder + " " + placeHolderName + ") string {\n")
	sb.WriteString("\tswitch {\n")
	for _, pair := range selectSyntax.Pairs {
		if pair.Key == "other" {
			continue
		}

		keyName := fc.ToGoName(pair.Key)
		sb.WriteString("\tcase " + keyName + ":\n")
		sb.WriteString("\t\treturn \"" + pair.Key + "\"\n")
	}
	sb.WriteString("\tdefault:\n\t\treturn \"other\"\n")
	sb.WriteString("\t}\n}\n\n")

	sb.WriteString("func " + fc.ToGoName(selectSyntax.Name) + "(" + selectSyntax.PlaceHolder + " " + placeHolderName + ") string {\n")
	sb.WriteString("\tswitch " + selectSyntax.PlaceHolder + " {\n")
	for _, pair := range selectSyntax.Pairs {
		if pair.Key == "other" {
			continue
		}

		keyName := fc.ToGoName(pair.Key)
		sb.WriteString("\tcase " + keyName + ":\n")
		sb.WriteString("\t\treturn \"" + selectSyntax.Prefix + pair.Value + selectSyntax.Suffix + "\"\n")
	}
	if len(otherValue) > 0 {
		sb.WriteString("\tdefault:\n\t\treturn \"" + otherValue + "\"\n")
	}
	sb.WriteString("\t}\n}")
	return sb.String()
}
