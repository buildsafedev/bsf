package skopeo

import (
	"errors"
	"strings"
)

type ImageReference interface {
	String() string
}

type imageRef struct {
	ref string
}

func (i *imageRef) String() string {
	return i.ref
}

func ParseImageName(ref string) (ImageReference, error) {
	if ref == "" || !strings.Contains(ref, ":") {
		return nil, errors.New("invalid image reference")
	}
	return &imageRef{ref: ref}, nil
}
