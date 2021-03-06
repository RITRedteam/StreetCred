package files

import (
	"bufio"
	"os"
)

func ReadList(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)

	if err != nil {
		os.Stderr.WriteString("ERROR: Failed to open specified file path: [ " + path + "]\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if scanner.Text() != "" {
			lines = append(lines, scanner.Text())
		}
	}

	return lines, nil
}

func ReadString(path string) (string, error) {
	var text string = ""
	file, err := os.Open(path)

	if err != nil {
		os.Stderr.WriteString("ERROR: Failed to open specified file path: [ " + path + "]\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}
	text = text[:len(text)-1]

	return text, nil
}
