package precheck

import (
	"reflect"
	"testing"
)

func TestCheckVersionConstraint(t *testing.T) {
	tests := []struct {
		currentVer   string
		requiredVer  string
		expectResult bool
	}{
		{"v1.2.3", "v1.0.0", true},
		{"v1.2.3", "v1.2.3", true},
		{"v1.2.3", "v1.2.4", false},
		{"v1.2.3", "v1.3.0", false},
	}

	for _, test := range tests {
		result := checkVersionGreater(test.currentVer, test.requiredVer)
		if result != test.expectResult {
			t.Errorf("checkVersionConstraint(%s, %s) = %t, want %t", test.currentVer, test.requiredVer, result, test.expectResult)
		}
	}
}

func TestIsSnapshotterEnabled(t *testing.T) {

	tests := []struct {
		name     string
		data     string
		expected bool
	}{
		{
			name:     "ContainerdSnapshotter is true",
			data:     " '[[driver-type io.containerd.snapshotter.v1]]' ",
			expected: true,
		},
		{
			name:     "ContainerSnapshotter is false",
			data:     "[[Backing Filesystem extfs] [Supports d_type true] [Using metacopy false] [Native Overlay Diff true] [userxattr false]]",
			expected: false,
		},
	}

	for _, testCases := range tests {
		t.Run(testCases.name, func(t *testing.T) {
			resp := isSnapshotterEnabled(testCases.data)
			if !reflect.DeepEqual(testCases.expected, resp) {
				t.Fail()
			}
		})
	}
}
