package stats

import (
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/dbg-go/pkg/stats"
	"github.com/spf13/cobra"
)

func NewStatsDriversCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Fetch stats about remote drivers",
		RunE:  executeDrivers,
	}
	return cmd
}

func executeDrivers(_ *cobra.Command, _ []string) error {
	statter, err := stats.NewS3Statter()
	if err != nil {
		return err
	}
	return stats.Run(stats.Options{Options: root.LoadRootOptions()}, statter)
}
