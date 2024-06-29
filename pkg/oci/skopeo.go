package oci

import (
	"os"
	"os/exec"
)

// LoadDocker loads the image to the docker daemon
func LoadDocker(daemon, dir, imageName string) error {
	cmd := exec.Command("nix", "run", "nixpkgs#skopeo", "--", "copy", "--insecure-policy", "--dest-daemon-host="+daemon, "dir:"+dir, "docker-daemon:"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoadPodman loads the image to the poadman
func LoadPodman(dir, imageName string) error {
	cmd := exec.Command("nix", "run", "nixpkgs#skopeo", "--", "copy", "--insecure-policy", "dir:"+dir, "containers-storage:"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Push image to registry
func Push(dir, imageName string) error {
	if _, err := os.Stat("/etc/skopeo/config.json"); err != nil {
		cmd := exec.Command("nix", "run", "nixpkgs#skopeo", "--", "login", "--authfile", "~/.skopeo/config.json", "docker.io")

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("nix", "run", "nixpkgs#skopeo", "--", "copy", "--insecure-policy", "dir:"+dir, "docker://"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
