package generate

import (
	"fmt"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
)

type UnsupportedTargetErr struct {
	target builder.Type
}

func (err *UnsupportedTargetErr) Error() string {
	return fmt.Sprintf("target %s is unsupported by driverkit", err.target.String())
}
