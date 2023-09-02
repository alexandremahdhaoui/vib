package vib

const (
	ExpressionKind = "Expression"
)

type ExpressionSpec struct {
	Key         string `json:"key" yaml:"key"`
	Value       string `json:"value" yaml:"value"`
	ResolverRef string `json:"resolverRef" yaml:"resolverRef"`
}

func (p *ExpressionSpec) Render() (string, error) {
	// TODO implement me
	panic("not implemented yet")
}

func NewExpression(name string, spec ExpressionSpec) *ResourceDefinition {
	return &ResourceDefinition{
		APIVersion: V1Alpha1,
		Kind:       ExpressionKind,
		Metadata:   NewMetadata(name),
		Spec:       spec,
	}
}

func NewExpressionWithDefaultName(spec ExpressionSpec) *ResourceDefinition {
	return NewExpression(spec.Key, spec)
}
