package cleanup

import (
	"github.com/falcosecurity/dbg-go/pkg/cleanup"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/spf13/cobra"
)

func NewCleanupDriversCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup desired remote drivers",
		RunE:  executeDrivers,
	}
	return cmd
}

func executeDrivers(_ *cobra.Command, _ []string) error {
	cleaner, err := cleanup.NewS3Cleaner()
	if err != nil {
		return err
	}
	return cleanup.Run(cleanup.Options{Options: root.LoadRootOptions()}, cleaner)
}
