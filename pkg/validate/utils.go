package validate

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"strings"
)

func ConfGlobFromDistro(opts root.Target) string {
	// Empty filters fallback at ".*" since we are using a regex match below
	if opts.Distro == "" {
		opts.Distro = "*"
	} else {
		dkDistro, found := root.SupportedDistros[opts.Distro]
		if found {
			// Filenames use driverkit lowercase target, instead of the kernel-crawler naming.
			opts.Distro = dkDistro
		} else {
			// Perhaps a regex? ToLower and pray
			opts.Distro = strings.ToLower(opts.Distro)
		}
	}
	if opts.KernelRelease == "" {
		opts.KernelRelease = "*"
	}
	if opts.KernelVersion == "" {
		opts.KernelVersion = "*"
	}
	return fmt.Sprintf("%s_%s_%s.yaml", opts.Distro, opts.KernelRelease, opts.KernelVersion)
}
