package vib

const (
	ProfileKind Kind = "Profile"
)

type ProfileSpec struct {
	// SetRefs is a list of reference to vib.Set.
	SetRefs []string
}

func (p *ProfileSpec) Render() (string, error) {
	// TODO implement me
	panic("not implemented yet")
}

func NewProfile(name string, spec ProfileSpec) *ResourceDefinition {
	return &ResourceDefinition{
		APIVersion: V1Alpha1,
		Kind:       ProfileKind,
		Metadata:   NewMetadata(name),
		Spec:       spec,
	}
}
