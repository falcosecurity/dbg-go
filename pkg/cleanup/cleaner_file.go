package cleanup

import (
	"github.com/falcosecurity/dbg-go/pkg/root"
	"os"
)

type fileCleaner struct {
	root.Looper
}

func NewFileCleaner() Cleaner {
	return &fileCleaner{Looper: root.NewFsLooper(root.BuildConfigPath)}
}

func (f *fileCleaner) Info() string {
	return "cleaning up local config files"
}

func (f *fileCleaner) Cleanup(opts Options) error {
	return f.LoopFiltered(opts.Options, "removing file", "config", func(driverVersion, configPath string) error {
		return os.Remove(configPath)
	})
}
