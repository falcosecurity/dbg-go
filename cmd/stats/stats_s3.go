package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/stats"
	"github.com/spf13/cobra"
)

func NewStatsS3Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Fetch stats about remote drivers",
		RunE:  executeS3,
	}
	return cmd
}

func executeS3(_ *cobra.Command, _ []string) error {
	statter, err := stats.NewS3Statter()
	if err != nil {
		return err
	}
	return stats.Run(stats.Options{Options: root.LoadRootOptions()}, statter)
}
