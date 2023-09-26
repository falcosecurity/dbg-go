package stats

import (
	"github.com/falcosecurity/dbg-go/pkg/root"
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
	Info() string
	GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error)
}
