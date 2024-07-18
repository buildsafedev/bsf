package skopeo

import (
	"context"
	"errors"
	"io"
)

type Options struct {
	ReportWriter io.Writer
}

type Image struct{}

func CopyImage(ctx context.Context, policyContext *PolicyContext, destRef, srcRef ImageReference, options *Options) (*Image, error) {
	if srcRef == nil || destRef == nil {
		return nil, errors.New("invalid image reference")
	}

	if options.ReportWriter != nil {
		io.WriteString(options.ReportWriter, "Copying image...\n")
	}

	return &Image{}, nil
}
