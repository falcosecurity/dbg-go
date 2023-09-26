package build

import (
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/driverkit/cmd"
)

type Options struct {
	root.Options
	SkipExisting   bool
	Publish        bool
	IgnoreErrors   bool
	RedirectErrors string
}

type publishVal struct {
	driverVersion string
	out           cmd.OutputOptions
}
