package parser

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_getFileType(t *testing.T) {
	xmlFile, err := os.Open("./fixture/import.xml")
	if err != nil {
		t.Errorf("failed open fixture import.xml, %s", err)
	}
	defer xmlFile.Close()

	fileType, err := getFileType(xmlFile)
	if err != nil {
		t.Errorf("test failed, %s", err)
	}
	if fileType != fileTypeStore {
		t.Error("fixture import.xml is not store", err)
	}
}

func Test_findXMLFiles(t *testing.T) {
	actual, err := findXMLFiles("./fixture")
	if err != nil {
		t.Errorf("failed to find fixture, %s", err)
	}
	expected := []string{
		"fixture/import.xml",
		"fixture/offer.xml",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("not equal")
	}
}

func Test_Parse(t *testing.T) {
	err := Parse("./fixture")
	if err != nil {
		fmt.Printf("failed to parse file, %s\n", err)
	}
}
