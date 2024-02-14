package main

import (
	"breadkey/freezed/converter"
	"log"
	"os"
	"strings"
)

const (
	fullCommand          = "--"
	outputCommand        = "output"
	languageCommand      = "language"
	shortCommand         = "-"
	shortLanguageCommand = "l"
	shortOutputCommand   = "o"
	defaultLanguage      = "go"

	enterCorrectFileName = "Please enter dart file name like **/*.dart!"
	enterCorrectCommand  = "Please enter correct command!"
	unsupportedLanguage  = "Unsupported language!"
)

var supportedLanguages = [...]string{defaultLanguage}

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalln(enterCorrectFileName)
	}

	fileName := args[1]

	if fileName[0] == '-' {
		log.Fatalln(enterCorrectFileName)
	}

	pointer := 2

	language := defaultLanguage
	fileNameParts := strings.Split(fileName, string(os.PathSeparator))

	dartFileName := fileNameParts[len(fileNameParts)-1]
	dartFileParts := strings.Split(dartFileName, ".")

	if dartFileParts[len(dartFileParts)-1] != "dart" {
		log.Fatalln(enterCorrectFileName)
	}

	argsLength := len(args)

	var outputFileName string

	for {
		if pointer >= argsLength {
			break
		}

		if !isCommand(args[pointer]) || len(args[pointer]) < 2 {
			log.Fatalln(enterCorrectCommand)
		} else {
			isShort := args[pointer][1] != '-'

			var command string

			if isShort {
				command = args[pointer][1:]
			} else {
				command = args[pointer][2:]
			}

			command = translateCommand(command, isShort)

			if pointer+1 >= argsLength {
				log.Fatalln(enterCorrectCommand)
			}

			arg := args[pointer+1]

			if isCommand(arg) {
				log.Fatalln(enterCorrectCommand)
			}

			switch command {
			case outputCommand:
				outputFileName = arg
			case languageCommand:
				language = arg
			}

			pointer += 2
		}
	}

	if !supported(language) {
		log.Fatalln("")
	}

	if outputFileName == "" {
		outputFileName = dartFileParts[0] + "." + language
	}

	useDefaultDirectory := true

	for _, char := range outputFileName {
		if os.IsPathSeparator(uint8(char)) {
			useDefaultDirectory = false
			break
		}
	}

	if useDefaultDirectory {
		os.Mkdir("outputs", os.ModePerm)
		outputFileName = "outputs" + string(os.PathSeparator) + outputFileName
	}

	freezeds := freezedconverter.Parse(fileName)

	if language == defaultLanguage {
		translates := make([]string, len(freezeds))

		for i, f := range freezeds {
			translates[i] = freezedconverter.TranslateToGo(&f)
		}

		result := strings.Join(translates, "\n\n")

		err := os.WriteFile(outputFileName, []byte(result), os.ModePerm)
	
		if err != nil {
			log.Fatalf("Write file error! %v", err)
		}
	}
}

func supported(language string) bool {
	for _, l := range supportedLanguages {
		if l == language {
			return true
		}
	}

	return false
}

func isCommand(str string) bool {
	return str[0] == '-'
}

func translateCommand(command string, isShort bool) string {
	if isShort {
		switch command {
		case shortOutputCommand:
			return outputCommand
		case shortLanguageCommand:
			return languageCommand
		}
	} else {
		switch command {
		case outputCommand, languageCommand:
			return command
		}
	}

	log.Fatalln(enterCorrectCommand)
	return ""
}
