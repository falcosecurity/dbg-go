package root

import (
	"github.com/falcosecurity/falcoctl/pkg/output"
	"github.com/pterm/pterm"
	"os"
)

var (
	Printer = output.NewPrinter(pterm.LogLevelInfo, pterm.LogFormatterColorful, os.Stdout)
)
