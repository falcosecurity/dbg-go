package validate

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate dbg configs",
		RunE:  execute,
	}
)

func execute(c *cobra.Command, args []string) error {
	return validate.Run(validate.Options{Options: root.LoadRootOptions()})
}
