package cleanup

import "github.com/fededp/dbg-go/pkg/root"

type Options struct {
	root.Options
}

type Cleaner interface {
	Remove(path string) error
	RemoveAll(path string) error
}
