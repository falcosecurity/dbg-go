package build

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/build"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewBuildConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build dbg configs",
		RunE:  executeConfigs,
	}
	flags := cmd.Flags()
	flags.Bool("skip-existing", true, "whether to skip the build of drivers existing on S3")
	flags.Bool("publish", false, "whether artifacts must be published on S3")
	flags.Bool("ignore-errors", false, "whether to ignore build errors and go on looping on config files")
	flags.String("aws-profile", "", "aws-profile to be used. Mandatory if publish is enabled.")

	return cmd
}

func executeConfigs(_ *cobra.Command, _ []string) error {
	options := build.Options{
		Options:      root.LoadRootOptions(),
		SkipExisting: viper.GetBool("skip-existing"),
		Publish:      viper.GetBool("publish"),
		IgnoreErrors: viper.GetBool("ignore-errors"),
		AwsProfile:   viper.GetString("aws-profile"),
	}

	if options.Publish && options.AwsProfile == "" {
		return fmt.Errorf("if publish is enabled, aws-profile is mandatory")
	}

	return build.Run(options)
}
