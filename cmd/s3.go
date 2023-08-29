package cmd

import (
	"github.com/fededp/dbg-go/cmd/cleanup"
	"github.com/fededp/dbg-go/cmd/stats"
	"github.com/spf13/cobra"
)

var (
	s3Cmd = &cobra.Command{
		Use:   "s3",
		Short: "Work with remote s3 bucket",
	}
)

func init() {
	// Subcommands
	s3Cmd.AddCommand(cleanup.NewCleanupCmd())
	s3Cmd.AddCommand(stats.NewStatsCmd())
}
