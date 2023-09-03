package v1alpha1

type ExpressionSpec struct {
	Key         string `json:"key"         yaml:"key"`
	Value       string `json:"value"       yaml:"value"`
	ResolverRef string `json:"resolverRef" yaml:"resolverRef"`
}
