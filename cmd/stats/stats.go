package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/stats"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "stats",
		Short: "Fetch stats about configs",
		RunE:  execute,
	}
)

func execute(c *cobra.Command, args []string) error {
	return stats.Run(stats.Options{Options: root.LoadRootOptions()})
}
