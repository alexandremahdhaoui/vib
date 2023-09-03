package v1alpha1

type SetSpec struct {
	// ExpressionRefs is a list of reference to vib.Expression or vib.Exp
	ExpressionRefs []string `json:"expressionRefs" yaml:"expressionRefs"`
}
