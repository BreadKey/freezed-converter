package freezedconverter

import (
	"log"
	"os"
)

var openBracketRunes = []rune{'(', '{', '['}
var closeBracketRunes = []rune{')', '}', ']'}
var bracketRunes = append(openBracketRunes, closeBracketRunes...)
var endRunes = append(bracketRunes, '\n', '\t', ' ', ',', ';')

func contains(runes []rune, r rune) bool {
	for _, c := range runes {
		if c == r {
			return true
		}
	}

	return false
}

func IsBracket(str string) bool {
	for _, r := range str {
		return contains(bracketRunes, r)
	}
	return false
}

func IsOpeningBracket(str string) bool {
	for _, r := range str {
		return contains(openBracketRunes, r)
	}
	return false
}

func IsClosingBracket(str string) bool {
	for _, r := range str {
		return contains(closeBracketRunes, r)
	}
	return false
}

type Tokenizer struct {
	FilePath     string
	Pointer      int
	file         *os.File
	currentToken string
}

func NewTokenizer(filePath string) *Tokenizer {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("File open error! %v", err)
	}

	return &Tokenizer{
		FilePath:     filePath,
		file:         file,
		currentToken: "",
	}
}

func (it *Tokenizer) Next() string {
	bytes := make([]byte, 5)

	it.file.Seek(int64(it.Pointer), 0)
	result, err := it.file.Read(bytes)

	if result == 0 {
		it.file.Close()

		if len(it.currentToken) == 0 {
			return ""
		} else {
			return it.returnToken(it.currentToken)
		}
	}

	if err != nil {
		log.Fatalf("Read error! %v", err)
	}

	str := string(bytes)

	startPointer := it.Pointer

	for pos, char := range str {
		it.Pointer = startPointer + pos + 1

		switch {
		case contains(endRunes, char):
			if len(it.currentToken) != 0 {
				it.Pointer = startPointer + pos

				return it.returnToken(it.currentToken)
			} else {
				if contains(bracketRunes, char) {
					return it.returnToken(string(char))
				}
			}
		default:
			it.currentToken += string(char)
		}
	}

	return it.Next()
}

func (it *Tokenizer) returnToken(token string) string {
	it.currentToken = ""

	return token
}
