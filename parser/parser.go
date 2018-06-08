package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	_ = iota
	fileTypeClassifier
	fileTypeOffer
)

// Parse commerce ML file
func Parse(login string) error {
	//xmlFile, err := os.Open(fmt.Sprintf("/data/source/%s/unzipped/", login))
	return nil
}

func getFileType(f *os.File) int64 {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.Contains(t, "Классификатор") {
			return fileTypeClassifier
		}
		if strings.Contains(t, "ПакетПредложений") {
			return fileTypeOffer
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("reading file: %v", err)
		os.Exit(1)
	}
	return 0
}
