package build

import (
	"github.com/fededp/dbg-go/pkg/build"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build dbg configs",
		RunE:  execute,
	}
	flags := cmd.Flags()
	flags.String("driver-name", "falco", "driver name to be used")
	flags.Bool("skip-existing", true, "whether to skip the build of drivers existing on S3")
	flags.Bool("publish", false, "whether artifacts must be published on S3")
	flags.String("aws-profile", "", "aws-profile to be used. Mandatory.")

	_ = cmd.MarkFlagRequired("aws-profile")
	return cmd
}

func execute(_ *cobra.Command, _ []string) error {
	options := build.Options{
		Options:      root.LoadRootOptions(),
		DriverName:   viper.GetString("driver-name"),
		SkipExisting: viper.GetBool("skip-existing"),
		Publish:      viper.GetBool("publish"),
		AwsProfile:   viper.GetString("aws-profile"),
	}
	return build.Run(options)
}
