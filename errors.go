package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/api"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

func ErrCannotTypeAssertKind(kind api.Kind) error {
	return logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("Kind %q does not support Render", kind))
}
