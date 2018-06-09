package parser

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const (
	_ = iota
	fileTypeStore
	fileTypeOffer
)

type entry struct {
	Store *os.File
	Offer *os.File
}

// Parse commerce ML file
func Parse(path string) error {
	ent := &entry{}

	files, err := findXMLFiles(path)
	if err != nil {
		return fmt.Errorf("failed to find files, %s", err)
	}
	for _, v := range files {
		xmlFile, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("failed to open file %s, %s", v, err)
		}

		t, err := getFileType(xmlFile)
		if err != nil {
			return err
		}

		switch t {
		case fileTypeStore:
			ent.Store = xmlFile
		case fileTypeOffer:
			ent.Offer = xmlFile
		}
	}

	readXML(ent.Store)

	return nil
}

func getFileType(f *os.File) (int64, error) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.Contains(t, "Классификатор") {
			return fileTypeStore, nil
		}
		if strings.Contains(t, "ПакетПредложений") {
			return fileTypeOffer, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to reading file, %s", err)
	}
	return 0, fmt.Errorf("failed to identify file")
}

func findXMLFiles(path string) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	var files []string
	err := filepath.Walk(path, func(p string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString("xml", f.Name())
			if err == nil && r {
				files = append(files, p)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func readXML(f *os.File) {
	f.Seek(0, 0)
	decoder := xml.NewDecoder(f)
	for {
		t, err := decoder.Token()
		fmt.Printf("failed to decode, %s\n", err)
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			entityType, err := CommerceMLType(se.Name.Local)
			if err == nil {
				entity := reflect.New(entityType.Elem()).Interface().(CommerceMLInterface)
				decoder.DecodeElement(&entity, &se)
				fmt.Printf("@@@ %#v\n", entity)
				fmt.Printf("$$$ %s\n", entity.String())
			}
		}
	}
}
