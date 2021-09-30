package files

import (
	"bufio"
	"os"
)

func ReadList(path string) ([]string, error) {
	var lines []string
	//b, err := os.ReadFile(path)
	// if err != nil {
	// 	return nil, err
	// }
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
	//return strings.Split(string(b), "\r\n"), nil
}
