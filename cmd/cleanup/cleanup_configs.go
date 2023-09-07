package cleanup

import (
	"github.com/fededp/dbg-go/pkg/cleanup"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
)

func NewCleanupConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup dbg configs",
		RunE:  executeConfigs,
	}
	return cmd
}

func executeConfigs(c *cobra.Command, args []string) error {
	return cleanup.Run(cleanup.Options{Options: root.LoadRootOptions()}, cleanup.NewFileCleaner())
}
