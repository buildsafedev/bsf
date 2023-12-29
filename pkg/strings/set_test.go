package strings

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSliceToSet(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"package1", "package2", "package3"},
			expected: []string{"package1", "package2", "package3"},
		},
		{
			name:     "with duplicates",
			input:    []string{"package1", "package2", "package2", "package3"},
			expected: []string{"package1", "package2", "package3"},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceToSet(tt.input)
			less := func(a, b string) bool { return a < b }

			equalIgnoreOrder := cmp.Diff(got, tt.expected, cmpopts.SortSlices(less)) == ""
			if !equalIgnoreOrder {
				t.Errorf("got %v, want %v", got, tt.expected)
				t.Fail()
			}
		})
	}
}
