package generate

import (
	"github.com/fededp/dbg-go/pkg/generate"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate new dbg configs",
		Long: `
In auto mode, configs will be generated starting from kernel-crawler output. 
In this scenario, target-{distro,kernelrelease,kernelversion} are available to filter to-be-generated configs. Regexes are allowed.
Instead, when auto mode is not enabled, the tool is able to generate a single config (for each driver version).
In this scenario, target-{distro,kernelrelease,kernelversion} CANNOT be regexes but must be exact values.
`,
		RunE: execute,
	}
)

func execute(c *cobra.Command, args []string) error {
	options := generate.Options{
		Options:    root.LoadRootOptions(),
		Auto:       viper.GetBool("auto"),
		DriverName: viper.GetString("driver-name"),
	}
	return generate.Run(options)
}

func init() {
	flags := Cmd.Flags()
	flags.Bool("auto", false, "automatically generate configs from kernel-crawler output")
	flags.String("driver-name", "falco", "driver name to be used")
}
