package jsonl

import (
	"encoding/json"
)

func ValidateIsJSONL(line []byte) error {
	var data interface{}
	if err := json.Unmarshal(line, &data); err != nil {
		return err
	}
	return nil
}
