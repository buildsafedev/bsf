package cmd

import (
	"reflect"
	"testing"
)

func TestParseAppDetails(t *testing.T) {
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
				Digest:  "1vng6wj07s51jsgj338m24m0c0mw2i3k",
				Name:    "app",
				Version: "0.1.0",
			},
			wantErr: false,
		},
		{
			name: "Test Case 2",
			path: "/nix/store/da66gxmm6wy8shkw93x5m6c1x8gfj63r-caddy-2.7.6",
			wantApp: &App{
				Digest:  "da66gxmm6wy8shkw93x5m6c1x8gfj63r",
				Name:    "caddy",
				Version: "2.7.6",
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
			gotApp, err := parseAppDetails(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAppDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotApp, tt.wantApp) {
				t.Errorf("parseAppDetails() = %v, want %v", gotApp, tt.wantApp)
			}
		})
	}
}
