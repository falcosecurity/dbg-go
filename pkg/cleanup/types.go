package cleanup

import "github.com/falcosecurity/dbg-go/pkg/root"

type Options struct {
	root.Options
}

type Cleaner interface {
	Info() string
	Cleanup(opts Options) error
}
