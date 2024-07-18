package oci

import (
	"context"
	"fmt"
	"io"
	"github.com/buildsafedev/bsf/pkg/skopeo"
)

func handleImageOperation(srcType, destType, dir, imageName string, out io.Writer) error {
	ctx := context.Background()

	srcRef, err := skopeo.ParseImageName(srcType + ":" + dir)
	if err != nil {
		return fmt.Errorf("parsing source image reference: %w", err)
	}
	destRef, err := skopeo.ParseImageName(destType + ":" + imageName)
	if err != nil {
		return fmt.Errorf("parsing destination image reference: %w", err)
	}

	policyContext, err := skopeo.NewPolicyContext(&skopeo.Policy{})
	if err != nil {
		return fmt.Errorf("creating policy context: %w", err)
	}
	defer policyContext.Destroy()

	_, err = skopeo.CopyImage(ctx, policyContext, destRef, srcRef, &skopeo.Options{
		ReportWriter: out,
	})
	if err != nil {
		return fmt.Errorf("copying image: %w", err)
	}

	return nil
}

func LoadDocker(daemon, dir, imageName string, out io.Writer) error {
	return handleImageOperation("dir", "docker-daemon://"+daemon+"/"+imageName, dir, imageName, out)
}

func LoadPodman(dir, imageName string, out io.Writer) error {
	return handleImageOperation("dir", "containers-storage", dir, imageName, out)
}

func Push(dir, imageName string, out io.Writer) error {
	return handleImageOperation("dir", "docker", dir, imageName, out)
}
