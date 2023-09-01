package build

import (
	"github.com/falcosecurity/driverkit/cmd"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
	"github.com/fededp/dbg-go/pkg/root"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
)

func Run(opts Options) error {
	slog.Info("building drivers")
	client, err := s3utils.NewClient(false, opts.AwsProfile)
	if err != nil {
		return err
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

	ro := cmd.RootOptions{
		Architecture:     opts.Architecture.String(),
		DriverVersion:    driverVersion,
		KernelVersion:    driverkitYaml.KernelVersion,
		ModuleDriverName: opts.DriverName,
		ModuleDeviceName: opts.DriverName,
		KernelRelease:    driverkitYaml.KernelRelease,
		Target:           driverkitYaml.Target,
		KernelConfigData: driverkitYaml.KernelConfigData,
		BuilderRepos:     []string{"docker.io/falcosecurity/driverkit-builder"},
		KernelUrls:       driverkitYaml.KernelUrls,
		Repo: cmd.RepoOptions{
			Org:  "falcosecurity",
			Name: "libs",
		},
		Output: cmd.OutputOptions{
			Module: driverkitYaml.Output.Module,
			Probe:  driverkitYaml.Output.Probe,
		},
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
			logger.Info("drivers already available on remote, nothing to do")
			return nil // nothing to do
		}
	}

	logger.Info("building", "config", configPath)

	// Magic to call unexported method
	b := reflect.ValueOf(&ro).MethodByName("toBuild").Call(nil)[0].Interface().(*builder.Build)
	err = driverbuilder.NewDockerBuildProcessor(1000, "").Start( /*ro.ToBuild()*/ b)
	if err != nil {
		logger.Error(err.Error())
		return nil // we don't want to break the builds chain
	}

	if opts.Publish {
		logger.Info("publishing")
		// Publish object!
		if ro.Output.Module != "" {
			f, err := os.Open(ro.Output.Module)
			if err != nil {
				return err
			}
			err = client.PutObject(opts.Options, driverVersion, filepath.Base(ro.Output.Module), f)
			_ = f.Close()
			if err != nil {
				return err
			}
		}
		if ro.Output.Probe != "" {
			f, err := os.Open(ro.Output.Probe)
			if err != nil {
				return err
			}
			err = client.PutObject(opts.Options, driverVersion, filepath.Base(ro.Output.Probe), f)
			_ = f.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
