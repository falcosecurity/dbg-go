package build

import "github.com/fededp/dbg-go/pkg/root"

type Options struct {
	root.Options
	SkipExisting bool
	Publish      bool
	IgnoreErrors bool
	AwsProfile   string
}
