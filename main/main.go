package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	helpMessage      string = "this command line tool adds/changes parameters in text files\n example: configWriter -file=\"config.txt\" -parameter=username -value=root"
	defaultParameter string = "username"
	defaultValue     string = "root"
)

var (
	help      *bool
	file      *string
	parameter *string
	value     *string
)

var (
	WrongParameterError     = errors.New("pramaters missing")
	WrongFileExtensionError = errors.New("stat %s: wrong file extension")
	FileIsDirError          = errors.New("stat %s: given file is directory")
	FileWriteError          = errors.New("stat %s: could not write file")
	FileCreateError         = errors.New("stat %s: could not create file")
)

const (
	InvalidArgumentError int = 128
	GeneralError         int = 0
)

func init() {
	help = flag.Bool("help", false, helpMessage)
	file = flag.String("file", "file.txt", "the name of the text file (extension \"*.txt\", which should be processed")
	parameter = flag.String("parameter", defaultParameter, "the parameter name which should created or changed")
	value = flag.String("value", defaultValue, "the content of the parameter")

	flag.Parse()

	if *help || len(os.Args) < 4 || *parameter == defaultParameter || *value == defaultValue {
		flag.Usage()
		errorExit(WrongParameterError, GeneralError)
	}
}

func main() {
	content := string(loadFile(*file))

	content = replaceValue(&content, *parameter, *value)

	err := os.WriteFile(*file, []byte(content), 0644)

	if err != nil {
		errorExit(fmt.Errorf(FileWriteError.Error(), *file), InvalidArgumentError)
	}

	fmt.Println("content written successful")
}

func replaceValue(content *string, name string, value string) string {
	lines := strings.Fields(*content)

	parameterFound := false

	for pos, line := range lines {
		if false == strings.HasPrefix(line, name+"=") {
			continue
		}

		parameterFound = true

		lines[pos] = name + "=" + value
	}

	if parameterFound == false {
		lines = append(lines, name+"="+value)
	}

	return strings.Join(lines, "\n")
}

func loadFile(file string) []byte {

	if strings.HasSuffix(file, ".txt") == false {
		errorExit(fmt.Errorf(WrongFileExtensionError.Error(), file), InvalidArgumentError)
	}

	if fileInfo, err := os.Stat(file); err != nil {

		if f, err := os.Create(file); err != nil {
			errorExit(fmt.Errorf(FileCreateError.Error(), file), GeneralError)
		} else {
			if err := f.Close(); err != nil {
				errorExit(fmt.Errorf(FileCreateError.Error(), file), GeneralError)
			}
		}

	} else if fileInfo.IsDir() == true {
		errorExit(fmt.Errorf(FileIsDirError.Error(), file), InvalidArgumentError)
	}

	content, err := os.ReadFile(file)

	if err != nil {
		errorExit(err, InvalidArgumentError)
	}

	return content
}

func errorExit(err error, exit int) {
	fmt.Println(err.Error())
	os.Exit(exit)
}
