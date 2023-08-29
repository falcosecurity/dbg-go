package cmd

import (
	"github.com/fededp/dbg-go/cmd/cleanup"
	"github.com/fededp/dbg-go/cmd/generate"
	"github.com/fededp/dbg-go/cmd/stats"
	"github.com/fededp/dbg-go/cmd/validate"
	"github.com/spf13/cobra"
)

var (
	configsCmd = &cobra.Command{
		Use:   "configs",
		Short: "Work with local dbg configs",
	}
)

func init() {
	// Subcommands
	configsCmd.AddCommand(generate.NewGenerateCmd())
	configsCmd.AddCommand(cleanup.NewCleanupCmd())
	configsCmd.AddCommand(validate.NewValidateCmd())
	configsCmd.AddCommand(stats.NewStatsCmd())
}
