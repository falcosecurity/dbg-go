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
)

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

	err = driverbuilder.NewDockerBuildProcessor(1000, "").Start( /*ro.ToBuild()*/ toBuild(ro))
	if err != nil {
		if opts.IgnoreErrors {
			logger.Error(err.Error())
			return nil // do not break the configs loop, just try the next one
		}
		return err
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

// copied from driverkit
func toBuild(ro *cmd.RootOptions) *builder.Build {
	kernelConfigData := ro.KernelConfigData
	if len(kernelConfigData) == 0 {
		kernelConfigData = "bm8tZGF0YQ==" // no-data
	}

	build := &builder.Build{
		TargetType:        builder.Type(ro.Target),
		DriverVersion:     ro.DriverVersion,
		KernelVersion:     ro.KernelVersion,
		KernelRelease:     ro.KernelRelease,
		Architecture:      ro.Architecture,
		KernelConfigData:  kernelConfigData,
		ModuleFilePath:    ro.Output.Module,
		ProbeFilePath:     ro.Output.Probe,
		ModuleDriverName:  ro.ModuleDriverName,
		ModuleDeviceName:  ro.ModuleDeviceName,
		GCCVersion:        ro.GCCVersion,
		BuilderImage:      ro.BuilderImage,
		BuilderRepos:      ro.BuilderRepos,
		KernelUrls:        ro.KernelUrls,
		RepoOrg:           ro.Repo.Org,
		RepoName:          ro.Repo.Name,
		Images:            make(builder.ImagesMap),
		RegistryName:      ro.Registry.Name,
		RegistryUser:      ro.Registry.Username,
		RegistryPassword:  ro.Registry.Password,
		RegistryPlainHTTP: ro.Registry.PlainHTTP,
	}

	imageLister, _ := builder.NewRepoImagesLister(ro.BuilderRepos[0], build)
	build.ImagesListers = append(build.ImagesListers, imageLister)

	// attempt the build in case it comes from an invalid config
	kr := build.KernelReleaseFromBuildConfig()
	if len(build.ModuleFilePath) > 0 && !kr.SupportsModule() {
		build.ModuleFilePath = ""
	}
	if len(build.ProbeFilePath) > 0 && !kr.SupportsProbe() {
		build.ProbeFilePath = ""
	}

	return build
}
