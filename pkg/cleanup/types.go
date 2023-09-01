package cleanup

import "github.com/fededp/dbg-go/pkg/root"

type Options struct {
	root.Options
}

type Cleaner interface {
	Info() string
	Cleanup(opts Options, driverVersion string) error
}
