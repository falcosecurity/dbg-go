package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/stats"
	"github.com/spf13/cobra"
)

func NewStatsConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Fetch stats about dbg configs",
		RunE:  executeConfigs,
	}
	return cmd
}

func executeConfigs(_ *cobra.Command, _ []string) error {
	return stats.Run(stats.Options{Options: root.LoadRootOptions()}, stats.NewFileStatter())
}
