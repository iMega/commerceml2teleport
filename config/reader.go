package config

import (
	"io/ioutil"
	"os"
	"strings"
)

// GetConfigValue читает содержимое переменной окружения
func GetConfigValue(key string) (string, error) {
	filename := os.Getenv(key + "_FILE")
	if filename == "" {
		value := os.Getenv(key)
		return value, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(data), "\n"), nil
}
