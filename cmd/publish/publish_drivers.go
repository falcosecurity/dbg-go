package publish

import (
	"github.com/fededp/dbg-go/pkg/publish"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
)

func NewPublishDriversCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish local drivers to remote bucket",
		RunE:  executeDrivers,
	}
	return cmd
}

func executeDrivers(_ *cobra.Command, _ []string) error {
	options := publish.Options{
		Options: root.LoadRootOptions(),
	}
	return publish.Run(options)
}
