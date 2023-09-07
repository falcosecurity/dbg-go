package cmd

import (
	"github.com/fededp/dbg-go/cmd/cleanup"
	"github.com/fededp/dbg-go/cmd/publish"
	"github.com/fededp/dbg-go/cmd/stats"
	"github.com/spf13/cobra"
)

var (
	s3Cmd = &cobra.Command{
		Use:   "drivers",
		Short: "Work with remote drivers bucket",
	}
)

func init() {
	// Subcommands
	s3Cmd.AddCommand(cleanup.NewCleanupDriversCmd())
	s3Cmd.AddCommand(stats.NewStatsDriversCmd())
	s3Cmd.AddCommand(publish.NewPublishDriversCmd())
}
