package cleanup

import (
	"github.com/fededp/dbg-go/pkg/cleanup"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCleanupS3Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup desired remote drivers",
		RunE:  executeS3,
	}

	flags := cmd.Flags()
	flags.String("aws-profile", "", "aws-profile to be used. Mandatory.")

	_ = cmd.MarkFlagRequired("aws-profile")
	return cmd
}

func executeS3(_ *cobra.Command, _ []string) error {
	cleaner, err := cleanup.NewS3Cleaner(viper.GetString("aws-profile"))
	if err != nil {
		return err
	}
	return cleanup.Run(cleanup.Options{Options: root.LoadRootOptions()}, cleaner)
}
