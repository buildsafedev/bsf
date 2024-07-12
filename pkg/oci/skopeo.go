package oci

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

var registries = []string{"docker.io"}

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
	currentuser, err := user.Current()
	if err != nil {
		return err
	}
	authFile := currentuser.HomeDir + "/.skopeo/config.json"
	cmd := exec.Command("nix", "run", "nixpkgs#skopeo", "--", "copy", "--authfile", authFile, "--insecure-policy", "dir:"+dir, "docker://"+imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func Auth(registryName string) error {

	for _, reg := range registries {
		if reg == registryName {
			currentuser, err := user.Current()
			if err != nil {
				return err
			}
			fmt.Println("Logging in to registry")
			authFile := currentuser.HomeDir + "/.skopeo/config.json"
			cmd := exec.Command("nix", "run", "nixpkgs#skopeo", "--", "login", "--authfile", authFile, registryName)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				return err
			}
			return nil
		}

	}

	return fmt.Errorf("Registry not found")
}
