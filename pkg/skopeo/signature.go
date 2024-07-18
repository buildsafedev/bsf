package skopeo

import "fmt"

type Policy struct{}

type PolicyContext struct{}

func NewPolicyContext(policy *Policy) (*PolicyContext, error) {
	if policy == nil {
		return nil, fmt.Errorf("policy cannot be nil")
	}
	return &PolicyContext{}, nil
}

func (pc *PolicyContext) Destroy() {
	if pc == nil {
		fmt.Println("Warning: Destroy called on a nil PolicyContext")
		return
	}
}
