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

func (f *fileCleaner) Cleanup(opts Options, driverVersion string) error {
	opts.DriverVersion = []string{driverVersion} // locally overwrite driverVersions to only match current driverVersion
	return root.LoopConfigsFiltered(opts.Options, "removing file", func(driverVersion, configPath string) error {
		return os.Remove(configPath)
	})
}
