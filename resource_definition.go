package vib

const (
	V1Alpha1 = "vib.alexandre.mahdhaoui.com/v1alpha1"
)

type ResourceDefinition struct {
	APIVersion APIVersion  `json:"apiVersion" yaml:"apiVersion"`
	Kind       Kind        `json:"kind"       yaml:"kind"`
	Metadata   Metadata    `json:"metadata"   yaml:"metadata"`
	Spec       interface{} `json:"spec"       yaml:"spec"`
}

type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"      yaml:"labels,omitempty"`
	Name        string            `json:"name"                  yaml:"name"`
}

func NewMetadata(name string) Metadata {
	return Metadata{Name: name} //nolint:exhaustruct,exhaustivestruct
}

func NewResourceDefinition(apiVersion APIVersion, kind Kind, name string, spec interface{}) *ResourceDefinition {
	return &ResourceDefinition{
		APIVersion: apiVersion,
		Kind:       kind,
		Metadata:   NewMetadata(name),
		Spec:       spec,
	}
}
