package crypto

import "testing"

func TestHexToBase64(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		want    string
		wantErr bool
	}{
		{
			name:    "empty string",
			hexStr:  "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "valid hex string",
			hexStr:  "0cebc68454e3fd3570cc92d56566fc7087cbd2bfffffd8baf3c39fe49c876e32",
			want:    "DOvGhFTj/TVwzJLVZWb8cIfL0r///9i688Of5JyHbjI=",
			wantErr: false,
		},
		{
			name:    "invalid hex string",
			hexStr:  "GG",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToBase64(tt.hexStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToBase64() = %v, want %v", got, tt.want)
			}
		})
	}
}
