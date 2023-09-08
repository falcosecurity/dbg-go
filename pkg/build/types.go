package build

import (
	"github.com/falcosecurity/driverkit/cmd"
	"github.com/fededp/dbg-go/pkg/root"
)

type Options struct {
	root.Options
	SkipExisting   bool
	Publish        bool
	IgnoreErrors   bool
	RedirectErrors string
	AwsProfile     string
}

type publishVal struct {
	driverVersion string
	out           cmd.OutputOptions
}
