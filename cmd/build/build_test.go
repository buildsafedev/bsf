package build

import "testing"

func TestIsNoFileError(t *testing.T) {
	tests := []struct {
		name string
		err  string
		want bool
	}{
		{
			name: "No such file or directory",
			err:  "open /tmp/foo: No such file or directory",
			want: true,
		},
		{
			name: "does not contain a 'bsf/flake.nix' file",
			err:  "error: source tree referenced by 'git+file:///root/caddy?dir=bsf&ref=refs/heads/master&rev=0dd0487eba1592b19cb2260fe247b7d375301fe5' does not contain a 'bsf/flake.nix' file",
			want: true,
		},
		{
			name: "other error",
			err:  "some other error",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNoFileError(tt.err); got != tt.want {
				t.Errorf("isNoFileError() = %v, want %v", got, tt.want)
			}
		})
	}
}
