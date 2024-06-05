package misc

import (
	"bufio"
	"os"
)

func LoadBankMetadata(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	res := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res, scanner.Err()
}
