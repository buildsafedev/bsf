package jsonl

import (
	"bufio"
	"bytes"
	"encoding/json"
)

// ValidateIsJSONL validates if the file is JSONL
func ValidateIsJSONL(file []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		var data interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
