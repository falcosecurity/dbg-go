package cleanup

import (
	"github.com/fededp/dbg-go/pkg/root"
	"os"
)

type fileCleaner struct{}

func NewFileCleaner() Cleaner {
	return &fileCleaner{}
}

func (f *fileCleaner) Info() string {
	return "cleaning up local config files"
}

func (f *fileCleaner) Cleanup(opts Options) error {
	return root.LoopPathFiltered(opts.Options, root.BuildConfigPath, "removing file", "config", func(driverVersion, configPath string) error {
		return os.Remove(configPath)
	})
}
