package requestUtils

import (
	"bufio"
	"bytes"
	"os"
)

func SubStringCommand(substring string) (int, error) {
	file, err := os.Open("/var/log/apache2/dev.eljur.access.log")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	var scannerBytes []byte
	var count = 0
	chars := []byte(substring)

	w := bufio.NewWriterSize(os.Stdout, 1000000)
	for fileScanner.Scan() {
		scannerBytes = fileScanner.Bytes()

		if bytes.Contains(scannerBytes, chars) {
			count++
		}
	}
	_ = w.Flush()

	if err := fileScanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
