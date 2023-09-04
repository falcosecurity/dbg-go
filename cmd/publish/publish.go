package publish

import (
	"github.com/fededp/dbg-go/pkg/publish"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish local drivers to remote bucket",
		RunE:  execute,
	}
	flags := cmd.Flags()
	flags.String("aws-profile", "", "aws-profile to be used. Mandatory.")

	_ = cmd.MarkFlagRequired("aws-profile")
	return cmd
}

func execute(_ *cobra.Command, _ []string) error {
	options := publish.Options{
		Options:    root.LoadRootOptions(),
		AwsProfile: viper.GetString("aws-profile"),
	}
	return publish.Run(options)
}
