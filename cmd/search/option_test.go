package search

import "testing"

func TestDeriveAV(t *testing.T) {
	tests := []struct {
		name   string
		vector string
		want   string
	}{
		{
			name:   "Test Case 1",
			vector: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Network",
		},
		{
			name:   "Test Case 2",
			vector: "CVSS:3.1/AV:A/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Adjacent Network",
		},

		{
			name:   "Test Case 3",
			vector: "CVSS:3.1/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "",
		},
		{
			name:   "Test Case 4",
			vector: "CVSS:3.1/AV:L/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Local",
		},
		{
			name:   "Test Case 5",
			vector: "CVSS:3.1/AV:P/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Physical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeriveAV(tt.vector); got != tt.want {
				t.Errorf("deriveAV() = %v, want %v", got, tt.want)
			}
		})
	}
}
