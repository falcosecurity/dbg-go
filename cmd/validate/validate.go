package validate

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/spf13/cobra"
)

func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate dbg configs",
		RunE:  execute,
	}
	return cmd
}

func execute(c *cobra.Command, args []string) error {
	options := validate.Options{
		Options: root.LoadRootOptions(),
	}
	return validate.Run(options)
}
