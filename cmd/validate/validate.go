package validate

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Cmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate dbg configs",
		RunE:  execute,
	}
)

func execute(c *cobra.Command, args []string) error {
	options := validate.Options{
		Options:    root.LoadRootOptions(),
		DriverName: viper.GetString("driver-name"),
	}
	return validate.Run(options)
}

func init() {
	flags := Cmd.Flags()
	flags.String("driver-name", "falco", "driver name to be used")
}
