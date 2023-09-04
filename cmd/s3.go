package cmd

import (
	"github.com/fededp/dbg-go/cmd/cleanup"
	"github.com/fededp/dbg-go/cmd/publish"
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
	s3Cmd.AddCommand(cleanup.NewCleanupS3Cmd())
	s3Cmd.AddCommand(stats.NewStatsS3Cmd())
	s3Cmd.AddCommand(publish.NewPublishCmd())
}
