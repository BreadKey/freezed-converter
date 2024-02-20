package main

import (
	"breadkey/freezed/converter"
	l10nconverter "breadkey/l10n/converter"
	"log"
	"os"
	"strings"
)

const (
	fullCommand          = "--"
	outputCommand        = "output"
	languageCommand      = "language"
	formatCommand        = "format"
	shortCommand         = "-"
	shortLanguageCommand = "l"
	shortOutputCommand   = "o"
	shortFormatCommand   = "f"

	defaultLanguage = "go"

	defaultFormat = "firestore"

	enterCorrectFileName = "Please enter dart file name like **/*.dart! or **/*.arb!"
	enterCorrectCommand  = "Please enter correct command!"
	unsupportedLanguage  = "Unsupported language!"
)

var supportedFileExtensions = [...]string{"dart", "arb"}
var supportedLanguages = [...]string{defaultLanguage}
var supportedFormats = [...]string{defaultFormat, "json"}

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
	format := defaultFormat
	fileNameParts := strings.Split(fileName, string(os.PathSeparator))

	parsingFileName := fileNameParts[len(fileNameParts)-1]
	parsingFileNameParts := strings.Split(parsingFileName, ".")
	extension := parsingFileNameParts[len(parsingFileNameParts)-1]

	if !isSupportedFileExtension(extension) {
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
			case formatCommand:
				format = arg
			}

			pointer += 2
		}
	}

	if !isSupportedLanguage(language) {
		log.Fatalln("Unsupported Language! Please enter language in", supportedFormats)
	}

	if !isSupportedFormat(format) {
		log.Fatalln("Unsupported format! Please enter language in", supportedFormats)
	}

	if outputFileName == "" {
		outputFileName = parsingFileNameParts[0] + "." + language
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

	var translates []string

	switch extension {
	case "dart":
		freezeds := freezedconverter.Parse(fileName)

		translates = make([]string, len(freezeds))

		for i, f := range freezeds {
			translates[i] = freezedconverter.TranslateToGo(&f, format)
		}
	case "arb":
		selects := l10nconverter.ParseSelect(fileName)

		translates = make([]string, len(selects))

		for i, s := range selects {
			translates[i] = l10nconverter.TranslateToGo(&s)
		}
	}

	if language == defaultLanguage {
		result := strings.Join(translates, "\n\n")

		err := os.WriteFile(outputFileName, []byte(result), os.ModePerm)

		if err != nil {
			log.Fatalf("Write file error! %v", err)
		}
	}
}

func isSupportedFileExtension(extension string) bool {
	for _, e := range supportedFileExtensions {
		if e == extension {
			return true
		}
	}

	return false
}

func isSupportedLanguage(language string) bool {
	for _, l := range supportedLanguages {
		if l == language {
			return true
		}
	}

	return false
}

func isSupportedFormat(format string) bool {
	for _, l := range supportedFormats {
		if l == format {
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
		case shortFormatCommand:
			return formatCommand
		}
	} else {
		switch command {
		case outputCommand, languageCommand, formatCommand:
			return command
		}
	}

	log.Fatalln(enterCorrectCommand)
	return ""
}
