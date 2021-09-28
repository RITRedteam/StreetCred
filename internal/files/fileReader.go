package files

import (
	"os"
	"strings"
)

func ReadList(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\r\n"), nil
}

