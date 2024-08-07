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
func Push(dir, imageName string, destcreds string, digestPath string) error {
	skoepoCmd := []string{"nix", "run", "nixpkgs#skopeo", "--", "copy", "--insecure-policy", "dir:" + dir}

	var cmd *exec.Cmd
	if digestPath != "" {
		skoepoCmd = append(skoepoCmd, "docker://"+imageName+"@@unknown-digest@@", "--digestfile", digestPath)
	} else {
		skoepoCmd = append(skoepoCmd, "docker://"+imageName)
	}
	if destcreds != "" {
		skoepoCmd = append(skoepoCmd, "--dest-creds", destcreds)
	}

	cmd = exec.Command(skoepoCmd[0], skoepoCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
