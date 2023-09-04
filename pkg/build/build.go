package build

import (
	"github.com/falcosecurity/driverkit/cmd"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder"
	"github.com/fededp/dbg-go/pkg/publish"
	"github.com/fededp/dbg-go/pkg/root"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
)

// Used by tests
var testClient *s3utils.Client

func Run(opts Options) error {
	slog.Info("building drivers")
	var (
		client *s3utils.Client
		err    error
	)
	if testClient == nil {
		client, err = s3utils.NewClient(false, opts.AwsProfile)
		if err != nil {
			return err
		}
	} else {
		client = testClient
	}

	return root.LoopConfigsFiltered(opts.Options, "building driver", func(driverVersion, configPath string) error {
		return buildConfig(client, opts, driverVersion, configPath)
	})
}

func buildConfig(client *s3utils.Client, opts Options, driverVersion, configPath string) error {
	logger := slog.With("config", configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var driverkitYaml validate.DriverkitYaml
	err = yaml.Unmarshal(configData, &driverkitYaml)
	if err != nil {
		return errors.WithMessagef(err, "config: %s", configPath)
	}

	ro := cmd.NewRootOptions()

	ro.Architecture = opts.Architecture.String()
	ro.DriverVersion = driverVersion
	ro.KernelVersion = driverkitYaml.KernelVersion
	ro.ModuleDriverName = opts.DriverName
	ro.ModuleDeviceName = opts.DriverName
	ro.KernelRelease = driverkitYaml.KernelRelease
	ro.Target = driverkitYaml.Target
	ro.KernelConfigData = driverkitYaml.KernelConfigData
	ro.KernelUrls = driverkitYaml.KernelUrls
	ro.Output = cmd.OutputOptions{
		Module: driverkitYaml.Output.Module,
		Probe:  driverkitYaml.Output.Probe,
	}

	if opts.SkipExisting {
		if ro.Output.Module != "" {
			moduleName := filepath.Base(ro.Output.Module)
			if client.ObjectExists(opts.Options, driverVersion, moduleName) {
				ro.Output.Module = "" // disable module build
			}
		}
		if ro.Output.Probe != "" {
			probeName := filepath.Base(ro.Output.Probe)
			if client.ObjectExists(opts.Options, driverVersion, probeName) {
				ro.Output.Probe = "" // disable probe build
			}
		}
		if ro.Output.Module == "" && ro.Output.Probe == "" {
			logger.Info("drivers already available on remote, skipping")
			return nil // nothing to do
		}
	}

	err = driverbuilder.NewDockerBuildProcessor(1000, "").Start(ro.ToBuild())
	if err != nil {
		if opts.IgnoreErrors {
			logger.Error(err.Error())
			return nil // do not break the configs loop, just try the next one
		}
		return err
	}

	if opts.Publish {
		err = publish.LoopDriversFiltered(client, opts.Options)
		if err != nil {
			return err
		}
	}
	return nil
}
