// Package library ...
package library

import (
	"bytes"
	"encoding/json"
	"log"
)

// PrettyPrintJSON ... returns JSON string
func PrettyPrintJSON(data interface{}) string {
	var out bytes.Buffer

	jsonBytes, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	err = json.Indent(&out, jsonBytes, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}
