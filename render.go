package vib

import "fmt"

type ExpressionRenderer interface {
	Render(key, value string) (string, error)
}

type Renderer interface {
	Render() (string, error)
}

func Render(resource *ResourceDefinition) (string, error) {
	var renderer Renderer
	var ok bool

	switch resource.Kind {
	case ProfileKind:
		renderer, ok = resource.Spec.(*ProfileSpec)
	case SetKind:
		renderer, ok = resource.Spec.(*ProfileSpec)
	case ExpressionKind:
		renderer, ok = resource.Spec.(*ExpressionSpec)
	default:
		ok = false
	}

	if !ok {
		return "", NewErrAndLog(ErrType, fmt.Sprintf("Kind %q does not support Render", resource.Kind))
	}

	return renderer.Render() //nolint:wrapcheck
}
