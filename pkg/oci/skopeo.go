package oci

import (
    "context"
    "fmt"
    "os"
    "github.com/BalaadityaPatanjali/image/v5/copy"
    "github.com/BalaadityaPatanjali/image/v5/signature"
    "github.com/BalaadityaPatanjali/image/v5/transports/alltransports"
)

// PerformImageOperation performs the image operation based on provided source and destination references
func PerformImageOperation(srcRef, destRef string) error {
    ctx := context.Background()

    src, err := alltransports.ParseImageName(srcRef)
    if err != nil {
        return fmt.Errorf("parsing source image reference %s: %w", srcRef, err)
    }
    dest, err := alltransports.ParseImageName(destRef)
    if err != nil {
        return fmt.Errorf("parsing destination image reference %s: %w", destRef, err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, dest, src, &copy.Options{
        ReportWriter: os.Stdout,
    })
    if err != nil {
        return fmt.Errorf("copying image from %s to %s: %w", srcRef, destRef, err)
    }

    return nil
}

// LoadDocker loads the image to Docker using the specified daemon.
func LoadDocker(daemon, dir, imageName string) error {
    srcRef := "dir:" + dir
    destRef := "docker-daemon://" + daemon + "/" + imageName
    return PerformImageOperation(srcRef, destRef)
}

// LoadPodman loads the image to Podman
func LoadPodman(dir, imageName string) error {
    srcRef := "dir:" + dir
    destRef := "containers-storage:" + imageName
    return PerformImageOperation(srcRef, destRef)
}

// Push pushes the image to a registry
func Push(dir, imageName string) error {
    srcRef := "dir:" + dir
    destRef := "docker://" + imageName
    return PerformImageOperation(srcRef, destRef)
}
