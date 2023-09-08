package build

import (
	"fmt"
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
	"sync"
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
		// writable client only if we need to publish
		client, err = s3utils.NewClient(!opts.Publish)
		if err != nil {
			return err
		}
	} else {
		client = testClient
	}
	looper := root.NewFsLooper(root.BuildConfigPath)

	var redirectErrorsF *os.File
	if opts.RedirectErrors != "" {
		redirectErrorsF, err = os.OpenFile(opts.RedirectErrors, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer redirectErrorsF.Close()
	}

	var publishCh chan publishVal
	var wg sync.WaitGroup
	if opts.Publish {
		publishCh = make(chan publishVal, 64)
		wg.Add(1)
		go func() {
			defer wg.Done()
			publishLoop(publishCh, opts.Options, client)
		}()
	}

	err = looper.LoopFiltered(opts.Options, "building driver", "config", func(driverVersion, configPath string) error {
		return buildConfig(client, opts, publishCh, redirectErrorsF, driverVersion, configPath)
	})

	if publishCh != nil {
		close(publishCh)
	}
	wg.Wait()
	return err
}

func buildConfig(client *s3utils.Client, opts Options,
	publishCh chan<- publishVal, redirectErrorsF *os.File,
	driverVersion, configPath string) error {
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
				logger.Info("output module already exists inside S3 bucket - skipping")
				ro.Output.Module = "" // disable module build
			}
		}
		if ro.Output.Probe != "" {
			probeName := filepath.Base(ro.Output.Probe)
			if client.HeadDriver(opts.Options, driverVersion, probeName) {
				logger.Info("output probe already exists inside S3 bucket - skipping")
				ro.Output.Probe = "" // disable probe build
			}
		}
		if ro.Output.Module == "" && ro.Output.Probe == "" {
			logger.Info("drivers already available on S3 bucket, skipping build")
			return nil // nothing to do
		}
	}

	err = driverbuilder.NewDockerBuildProcessor(1000, "").Start(ro.ToBuild())
	if err != nil {
		if redirectErrorsF != nil {
			logLine := fmt.Sprintf("config: %s | error: %s\n", configPath, err.Error())
			_, _ = redirectErrorsF.WriteString(logLine)
		}
		if opts.IgnoreErrors {
			logger.Error(err.Error())
			return nil // do not break the configs loop, just try the next one
		}
		return err
	}

	if publishCh != nil {
		publishCh <- publishVal{
			driverVersion: driverVersion,
			out:           ro.Output,
		}
	}
	return nil
}

func publishLoop(publishCh <-chan publishVal, opts root.Options, client *s3utils.Client) {
	for val := range publishCh {
		if val.out.Module != "" {
			err := client.PutDriver(opts, val.driverVersion, val.out.Module)
			if err != nil {
				slog.Warn("failed to upload module", "path", val.out.Module, "err", err.Error())
			} else {
				slog.Info("published module", "path", val.out.Module)
			}
		}
		if val.out.Probe != "" {
			err := client.PutDriver(opts, val.driverVersion, val.out.Probe)
			if err != nil {
				slog.Warn("failed to upload probe", "path", val.out.Probe, "err", err.Error())
			} else {
				slog.Info("published probe", "path", val.out.Probe)
			}
		}
	}
}
