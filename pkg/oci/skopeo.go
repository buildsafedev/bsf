package oci

import (
    "context"
    "fmt"
    "io"
    "github.com/containers/image/v5/copy"
    "github.com/containers/image/v5/signature"
    "github.com/buildsafedev/bsf/pkg/skopeo"
)

// LoadDocker loads the image to the docker daemon
func LoadDocker(daemon, dir, imageName string, out io.Writer) error {
    ctx := context.Background()

    // Parsed the source and destination image references
    srcRef, err := skopeo.ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := skopeo.ParseImageName("docker-daemon://" + daemon + "/" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    // Created a policy context
    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    // Copied the image
    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: out,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}

// LoadPodman loads the image to the podman
func LoadPodman(dir, imageName string, out io.Writer) error {
    ctx := context.Background()

    srcRef, err := skopeo.ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := skopeo.ParseImageName("containers-storage:" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: out,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}

// Push image to registry
func Push(dir, imageName string, out io.Writer) error {
    ctx := context.Background()

    srcRef, err := skopeo.ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := skopeo.ParseImageName("docker://" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: out,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}
