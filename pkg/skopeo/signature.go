package skopeo

import (
	"errors"
	"fmt"
)

type Policy struct{}

type PolicyContext struct {
	policy *Policy
}

func NewPolicyContext(policy *Policy) (*PolicyContext, error) {
	if policy == nil {
		return nil, errors.New("policy cannot be nil")
	}
	return &PolicyContext{policy: policy}, nil
}

func (pc *PolicyContext) Destroy() {
	if pc == nil {
		fmt.Println("Warning: Destroy called on a nil PolicyContext")
		return
	}
	pc.policy = nil
}
