package precheck

import "testing"

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
