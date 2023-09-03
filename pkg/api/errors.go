package api

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

func ErrKind(kind Kind) error {
	return logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("kind %q is not supported", kind))
}

func ErrApiVersion(apiVersion APIVersion, kind Kind) error {
	return logger.NewErrAndLog(logger.ErrType, fmt.Sprintf("APIVersion %q for Kind %q is not supported", apiVersion, kind))
}

func ErrReference(reference string, kind Kind) error {
	return logger.NewErrAndLog(logger.ErrReference, fmt.Sprintf("cannot resolve reference %q to Kind %q", reference, kind))
}
