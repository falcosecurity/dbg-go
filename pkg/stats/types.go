package stats

import "github.com/fededp/dbg-go/pkg/root"

type Options struct {
	root.Options
}

type driverStats struct {
	NumProbes            int64
	NumModules           int64
	NumHeaders           int64
	NumKernelConfigDatas int64
}
