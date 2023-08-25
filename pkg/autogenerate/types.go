package autogenerate

import (
	"github.com/fededp/dbg-go/pkg/root"
)

type Options struct {
	root.Options
	DriverName string
	Target     string
}
