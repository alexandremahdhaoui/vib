package v1alpha1

type ProfileSpec struct {
	// SetRefs is a list of reference to vib.Set.
	SetRefs []string `json:"setRefs" yaml:"setRefs"`
}
