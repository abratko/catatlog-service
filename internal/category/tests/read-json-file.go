package tests

import (
	"bytes"
	"encoding/json"
	"os"
)

func ReadJsonFile(path string) string {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = json.Compact(&buffer, jsonData)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
