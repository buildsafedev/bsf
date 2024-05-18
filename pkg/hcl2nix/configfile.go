package hcl2nix

import "strings"

// ConfigFiles holds config files to export
type ConfigFiles struct {
	Name string `hcl:"name,label"`
	// Name of files to copy from root of project
	Files []string `hcl:"files"`
	// DestinationDir is the directory to copy config files to in the container
	// This directory will be created in the root of the image
	DestinationDir string `hcl:"destinationDir,optional"`
}

// Validate validates ConfigFiles
func (c *ConfigFiles) Validate() *string {
	if c.Name == "" {
		return pointerTo("Name of config file cannot be empty")
	}

	if len(c.Files) == 0 {
		return pointerTo("Files to copy cannot be empty")
	}

	if c.DestinationDir == "" {
		c.DestinationDir = "/"
	}

	blockedChars := []string{";", "||", "&&", "&", "|", ">", "<", "$"}

	// check if the DestinationDir contains any blocked characters to prevent potential command injection misuse
	for _, char := range blockedChars {
		if strings.Contains(c.DestinationDir, char) {
			return pointerTo("Destination directory cannot contain special characters")
		}
	}

	return nil
}
