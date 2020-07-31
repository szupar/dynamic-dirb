package file_mng

import (
	"bufio"
	"dynamic-dirb/internal/string_mng"
	"fmt"
	"os"
)

// function fileExists verifies if a file exists in a given path (fileName)
func FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

// function readLines return the content of a file
func readLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()
	return lines, scanner.Err()
}

// function readWordlist safely return the content of a file
func ReadFileByLine(fileName string) []string {
	if FileExists(fileName) {
		lines, _ := readLines(fileName)
		return lines
	}
	return nil
}

func Append(fileName string, data string) bool {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	if _, err := f.WriteString(data + "\n"); err != nil {
		f.Close()
		return false
	}
	if err := f.Close(); err != nil {
		return false
	}

	return true
}

// override a file with a empty one
func emptyFile(fileName string) bool {

	if FileExists(fileName) {
		if RequireConfirmation("Output file already exists, do you want to override it? y/n") {
			f, _ := os.Create(fileName)
			f.Close()
			return true
		}
		return false
	}

	f, _ := os.Create(fileName)
	f.Close()
	return true

}

// override a file with a empty one without asking
func EmptyFileOverride(fileName string) bool {
	f, _ := os.Create(fileName)
	f.Close()
	return true

}

func RequireConfirmation(question string) bool {
	string_mng.PrintNotice(question)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	if text == "y" {
		return true
	}
	return false
}

func WriteGraph(fileName string, root string, node string) bool {
	if !FileExists(fileName) {
		EmptyFileOverride(fileName)
	}

	if IsEmpty(fileName) == 0 {
		if !Append(fileName, "digraph {") {
			return false
		}
	}

	data := "\t\"" + root + "\"" + " -> " + "\"" + node + "\";"

	return Append(fileName, data)
}

func IsEmpty(fileName string) int {
	f, _ := os.Stat(fileName)
	return int(f.Size())
}
