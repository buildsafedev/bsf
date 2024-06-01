package cmd

import (
	"testing"
)

func TestParseNixStorePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantApp *App
		wantErr bool
	}{
		{
			name: "Test Case 1",
			path: "/nix/store/1vng6wj07s51jsgj338m24m0c0mw2i3k-python3.11-app-0.1.0",
			wantApp: &App{
				ResultDigest: "1vng6wj07s51jsgj338m24m0c0mw2i3k",
				Name:         "python3.11-app",
				Version:      "0.1.0",
			},
			wantErr: false,
		},
		{
			name: "Test Case 2",
			path: "/nix/store/da66gxmm6wy8shkw93x5m6c1x8gfj63r-caddy-2.7.6",
			wantApp: &App{
				ResultDigest: "da66gxmm6wy8shkw93x5m6c1x8gfj63r",
				Name:         "caddy",
				Version:      "2.7.6",
			},
		},
		{
			name:    "Test Case 3",
			path:    "/nix/store/1vng6wj07s51jsgj338m24m0c0mw2i3k",
			wantApp: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digest, version, name, err := parseNixStorePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAppDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantApp == nil {
				return
			}
			if tt.wantApp.ResultDigest != digest {
				t.Errorf("parseAppDetails() digest = %v, want %v", digest, tt.wantApp.ResultDigest)
			}

			if tt.wantApp.Name != name {
				t.Errorf("parseAppDetails() name = %v, want %v", name, tt.wantApp.Name)
			}

			if tt.wantApp.Version != version {
				t.Errorf("parseAppDetails() version = %v, want %v", version, tt.wantApp.Version)
			}

		})
	}
}
