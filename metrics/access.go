package metrics

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const FILE = "/var/log/apache2/dev.eljur.access.log"

type Access struct {
	total   int
	api     int
	storage int
	other   int
}

func (a *Access) setTotal() {
	a.total = 3
}

func (a *Access) setApi() {
	a.api = 4
}

func (a *Access) setStorage() {
	a.storage = 5
}

func (a *Access) setOther() {
	a.other = 6
}

func (a *Access) Get(name string) map[string]int {
	return map[string]int{
		"pache.requests.total":    555,
		"apache.requests.api":   234,
		"apache.requests.storage": 435,
		"apache.requests.site":    123,
	}
}

func (a *Access) Get1(name string) (interface{}, error) {
	file, err := os.Open(FILE)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var scannerBytes []byte

	var countApi = 0
	var countStorage = 0
	var countOther = 0
	var countLine = 0

	charsApi := []byte("api")
	charsStorage := []byte("storage")

	w := bufio.NewWriterSize(os.Stdout, 1000000)
	for fileScanner.Scan() {
		scannerBytes = fileScanner.Bytes()

		countLine++

		if !bytes.Contains(scannerBytes, charsApi) && !bytes.Contains(scannerBytes, charsStorage) {
			countOther++
		}

		if bytes.Contains(scannerBytes, charsApi) {
			countApi++
		}

		if bytes.Contains(scannerBytes, charsStorage) {
			countStorage++
		}
	}
	_ = w.Flush()

	if err := fileScanner.Err(); err != nil {
		return 0, err
	}
	return 2, nil
}

func (a *Access) Set(name string, value interface{}) {
	fmt.Println(222)
}
