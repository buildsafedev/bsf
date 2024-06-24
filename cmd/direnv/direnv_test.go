package direnv

import (
	"testing"
)

func TestValidateEnvVars(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "Valid env vars",
			args:    "KEY1=value1,KEY2=value2",
			wantErr: false,
		},
		{
			name:    "Invalid format (no value)",
			args:    "KEY1=value1,KEY2",
			wantErr: true,
		},
		{
			name:    "Invalid format (no key)",
			args:    "=value1,KEY2=value2",
			wantErr: true,
		},
		{
			name:    "Invalid format (empty)",
			args:    "",
			wantErr: true,
		},
		{
			name:    "Invalid format (no equals sign)",
			args:    "KEY1value1,KEY2=value2",
			wantErr: true,
		},
		{
			name:    "Invalid characters in key",
			args:    "KEY1 =value1,KEY2=value2",
			wantErr: true,
		},
		{
			name:    "Invalid characters in value",
			args:    "KEY1=value1,KEY2=value\x00",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEnvVars(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("validateEnvVars() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
