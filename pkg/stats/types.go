package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
)

type Options struct {
	root.Options
}

type driverStats struct {
	NumProbes  int64
	NumModules int64
}

type driverStatsByDriverVersion map[string]driverStats

type Statter interface {
	GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error)
}
