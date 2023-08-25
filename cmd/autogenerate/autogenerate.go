package autogenerate

import (
	"github.com/fededp/dbg-go/pkg/autogenerate"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Cmd = &cobra.Command{
		Use:   "autogenerate",
		Short: "Fetch updated kernel-crawler lists and generate new dbg configs",
		RunE:  execute,
	}
)

func execute(c *cobra.Command, args []string) error {
	options := autogenerate.Options{
		Options:    root.LoadRootOptions(),
		DriverName: viper.GetString("driver-name"),
	}
	return autogenerate.Run(options)
}

func init() {
	flags := Cmd.Flags()
	flags.String("driver-name", "falco", "driver name to be used")
}
