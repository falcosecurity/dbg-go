package build

import (
	"github.com/falcosecurity/driverkit/cmd"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder"
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
		client, err = s3utils.NewClient(opts.AwsProfile)
		if err != nil {
			return err
		}
	} else {
		client = testClient
	}
	looper := root.NewFsLooper(root.BuildConfigPath)
	return looper.LoopFiltered(opts.Options, "building driver", "config", func(driverVersion, configPath string) error {
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

	// If Module or Probe are not absolute paths, assume they are relative to the repo-root/driverkit folder.
	if !filepath.IsAbs(driverkitYaml.Output.Module) {
		driverkitYaml.Output.Module = filepath.Join(opts.RepoRoot, "driverkit", driverkitYaml.Output.Module)
	}
	if !filepath.IsAbs(driverkitYaml.Output.Probe) {
		driverkitYaml.Output.Probe = filepath.Join(opts.RepoRoot, "driverkit", driverkitYaml.Output.Probe)
	}
	ro.Output = cmd.OutputOptions{
		Module: driverkitYaml.Output.Module,
		Probe:  driverkitYaml.Output.Probe,
	}

	if opts.SkipExisting {
		if ro.Output.Module != "" {
			moduleName := filepath.Base(ro.Output.Module)
			if client.HeadDriver(opts.Options, driverVersion, moduleName) {
				ro.Output.Module = "" // disable module build
			}
		}
		if ro.Output.Probe != "" {
			probeName := filepath.Base(ro.Output.Probe)
			if client.HeadDriver(opts.Options, driverVersion, probeName) {
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
		if ro.Output.Module != "" {
			err = client.PutDriver(opts.Options, driverVersion, ro.Output.Module)
			if err != nil {
				return err
			}
		}
		if ro.Output.Probe != "" {
			err = client.PutDriver(opts.Options, driverVersion, ro.Output.Probe)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
