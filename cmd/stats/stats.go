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

func execute(c *cobra.Command, args []string) error {
	switch c.Parent().Name() {
	case "configs":
		return stats.Run(stats.Options{Options: root.LoadRootOptions()}, stats.NewFileStatter())
	case "s3":
		return stats.Run(stats.Options{Options: root.LoadRootOptions()}, stats.NewS3Statter())
	}
	panic("unreachable code.")
}
