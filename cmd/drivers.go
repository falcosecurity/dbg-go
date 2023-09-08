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
		Long: `Read only commands will use an S3 client with anonymous credentials.
Write commands will need proper "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY" environment variables set.
`,
	}
)

func init() {
	// Subcommands
	s3Cmd.AddCommand(cleanup.NewCleanupDriversCmd())
	s3Cmd.AddCommand(stats.NewStatsDriversCmd())
	s3Cmd.AddCommand(publish.NewPublishDriversCmd())
}
