package strings

import (
	"reflect"
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

func TestPreferNewSliceElements(t *testing.T) {
	parseFunc := func(s string) string {
		return s
	}

	tests := []struct {
		name     string
		existing []string
		new      []string
		want     []string
	}{
		{
			name:     "Test 1",
			existing: []string{"a", "b", "c"},
			new:      []string{"d", "e", "f"},
			want:     []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			name:     "Test 2",
			existing: []string{"a", "b", "c"},
			new:      []string{"a", "b", "c"},
			want:     []string{"a", "b", "c"},
		},
		{
			name:     "Test 3",
			existing: []string{"a", "b", "c"},
			new:      []string{"a", "b", "d"},
			want:     []string{"c", "a", "b", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreferNewSliceElements(tt.existing, tt.new, parseFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreferNewSliceElements() = %v, want %v", got, tt.want)
			}
		})
	}
}
