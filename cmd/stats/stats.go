package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/stats"
	"github.com/spf13/cobra"
)

func NewStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Fetch stats about configs",
		RunE:  execute,
	}
	return cmd
}

func execute(_ *cobra.Command, _ []string) error {
	return stats.Run(stats.Options{Options: root.LoadRootOptions()}, stats.NewFileStatter())
}
