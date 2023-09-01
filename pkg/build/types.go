package build

import "github.com/fededp/dbg-go/pkg/root"

type Options struct {
	root.Options
	DriverName   string
	SkipExisting bool
	Publish      bool
	AwsProfile   string
}
