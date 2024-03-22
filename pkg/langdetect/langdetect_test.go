package langdetect

import (
	"testing"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func TestBinaryFromModule(t *testing.T) {
	tests := []struct {
		name     string
		mod      *modfile.File
		wantPath string
	}{
		{
			name: "Test Case 1",
			mod: &modfile.File{
				Module: &modfile.Module{
					Mod: module.Version{
						Path: "github.com/user/project",
					},
				},
			},
			wantPath: "project",
		},
		{
			name: "Test Case 2",
			mod: &modfile.File{
				Module: &modfile.Module{
					Mod: module.Version{
						Path: "github.com/user/anotherproject/v2",
					},
				},
			},
			wantPath: "anotherproject",
		},
		{
			name: "Test Case 3",
			mod: &modfile.File{
				Module: &modfile.Module{
					Mod: module.Version{
						Path: "newproject",
					},
				},
			},
			wantPath: "newproject",
		},
		{
			name: "Test Case 4",
			mod: &modfile.File{
				Module: &modfile.Module{
					Mod: module.Version{
						Path: "git.bsf.com/longproject/subporject/subproject/newproject",
					},
				},
			},
			wantPath: "newproject",
		},

		{
			name: "Test Case 5",
			mod: &modfile.File{
				Module: &modfile.Module{
					Mod: module.Version{
						Path: "git.bsf.com/longproject/subporject/subproject/newproject/v25",
					},
				},
			},
			wantPath: "newproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPath := binaryFromModule(tt.mod); gotPath != tt.wantPath {
				t.Errorf("binaryFromModule() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
