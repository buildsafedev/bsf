package provenance

import (
	slsav1 "github.com/buildsafedev/bsf/pkg/slsa/v1"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

// Statement is a struct to hold the provenance statement
type Statement struct {
	intoto.StatementHeader
	Predicate slsav1.Provenance
}

// NewStatement creates a new provenance statement
func NewStatement() *Statement {
	st := Statement{}
	st.Type = "https://in-toto.io/Statement/v1"
	st.PredicateType = "https://slsa.dev/provenance/v1"
	return &st
}
