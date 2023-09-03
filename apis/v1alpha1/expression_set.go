package v1alpha1

type ExpressionSetSpec struct {
	// ArbitraryKeys is used for special resolvers, such as plain that does not require associated values
	ArbitraryKeys []string `json:"arbitraryKeys" yaml:"arbitraryKeys"`

	// KeyValues uses a list of map to avoid reordered key-values
	KeyValues []map[string]string `json:"keyValues"   yaml:"keyValues"`

	// ResolverRef
	ResolverRef string `json:"resolverRef"   yaml:"resolverRef"`
}
