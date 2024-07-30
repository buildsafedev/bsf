package df

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDF(t *testing.T) {
	tests := []struct {
		name        string
		dfType      string
		isHermetic  bool
		expected    string
	}{
		{
			name:       "generate Go Dockerfile with hermetic build",
			dfType:     "go",
			isHermetic: true,
			expected:   "RUN CGO_ENABLED=0 go build -mod=vendor -o /bin/server .",
		},
		{
			name:       "generate Go Dockerfile with non-hermetic build",
			dfType:     "go",
			isHermetic: false,
			expected:   "RUN CGO_ENABLED=0 go build -o /bin/server .",
		},
		{
			name:       "generate Python Dockerfile with hermetic build",
			dfType:     "python",
			isHermetic: true,
			expected:   "ENV PYTHONPATH=vendor",
		},
		{
			name:       "generate Python Dockerfile with non-hermetic build",
			dfType:     "python",
			isHermetic: false,
			expected:   "python -m pip install -r requirements.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := GenerateDF(&buf, tt.dfType, tt.isHermetic)
			assert.Nil(t, err)
			actualStr := strings.Join(strings.Fields(buf.String()), " ")
			expected := strings.Join(strings.Fields(tt.expected), " ")
			assert.Contains(t, actualStr, expected)
		})
	}
}


func TestGetDfTmpl(t *testing.T) {
	tests := []struct {
		name     string
		dfType   string
		expected string
	}{
		{
			name:     "get Go template",
			dfType:   "go",
			expected: goDfTmpl,
		},
		{
			name:     "get Python template",
			dfType:   "python",
			expected: pythonDfTmpl,
		},
		{
			name:     "get Rust template",
			dfType:   "rust",
			expected: rustDfTmpl,
		},
		{
			name:     "get empty template for unknown type",
			dfType:   "unknown",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDfTmpl(tt.dfType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
