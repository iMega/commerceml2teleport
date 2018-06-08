package parser

import (
	"os"
	"testing"
)

func Test_getFileType(t *testing.T) {
	xmlFile, err := os.Open("./fixture/import.xml")
	if err != nil {
		t.Fatalf("failed open fixture import.xml, %s", err)
	}
	defer xmlFile.Close()

	fileType := getFileType(xmlFile)
	if fileType != fileTypeClassifier {
		t.Error("fixture import.xml is not classifier", err)
	}
}
