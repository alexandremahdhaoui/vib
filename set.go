package vib

const (
	SetKind = "Set"
)

type SetSpec struct {
	// ExpressionRefs is a list of reference to vib.Expression
	ExpressionRefs []string
}

func (p *SetSpec) Render() (string, error) {
	// TODO implement me
	panic("not implemented yet")
}

func NewSet(name string, spec SetSpec) *ResourceDefinition {
	return &ResourceDefinition{
		APIVersion: V1Alpha1,
		Kind:       SetKind,
		Metadata:   NewMetadata(name),
		Spec:       spec,
	}
}
