package skopeo

type Policy struct{}

type PolicyContext struct{}

func NewPolicyContext(policy *Policy) (*PolicyContext, error) {
	return &PolicyContext{}, nil
}

func (pc *PolicyContext) Destroy() {}
