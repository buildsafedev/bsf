package cmd

import (
	"reflect"
	"testing"
)

func TestParseNixConfig(t *testing.T) {
	config := "key1=value1\nkey2=value2"
	expected := map[string]string{"key1": "value1", "key2": "value2"}
	result := parseNixConfig(config)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test case  failed: got %v, want %v", result, expected)
	}
}
