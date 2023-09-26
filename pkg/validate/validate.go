package validate

import (
	"encoding/base64"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Run(opts Options) error {
	slog.Info("validate config files")
	looper := root.NewFsLooper(root.BuildConfigPath)
	return looper.LoopFiltered(opts.Options, "validating", "config", func(driverVersion, configPath string) error {
		return validateConfig(configPath, opts, driverVersion)
	})
}

func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func validateConfig(configPath string, opts Options, driverVersion string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var driverkitYaml DriverkitYaml
	err = yaml.Unmarshal(configData, &driverkitYaml)
	if err != nil {
		return errors.WithMessagef(err, "config: %s", configPath)
	}

	slog.Info("validating",
		"config", configPath,
		"target", driverkitYaml.Target,
		"kernelrelease", driverkitYaml.KernelRelease,
		"kernelversion", driverkitYaml.KernelVersion)

	// Check that filename is ok
	expectedFilename := driverkitYaml.ToConfigName()
	configFilename := filepath.Base(configPath)
	if configFilename != expectedFilename {
		return &WrongConfigNameErr{configFilename, expectedFilename}
	}

	// Check that arch is ok
	if driverkitYaml.Architecture != opts.Architecture.String() {
		return &WrongArchInConfigErr{configPath, driverkitYaml.Architecture}
	}

	outputPath := root.BuildOutputPath(opts.Options, driverVersion, driverkitYaml.ToName())
	outputPathFilename := filepath.Base(outputPath)

	kr := kernelrelease.FromString(driverkitYaml.KernelRelease)
	kr.Architecture = opts.Architecture

	// Check output probe if present
	if driverkitYaml.Output.Probe != "" {
		outputProbeFilename := filepath.Base(driverkitYaml.Output.Probe)
		if outputProbeFilename != outputPathFilename+".o" {
			return &WrongOutputProbeNameErr{outputProbeFilename, outputPathFilename}
		}

		if !strings.Contains(driverkitYaml.Output.Probe, opts.Architecture.ToNonDeb()) {
			return &WrongOutputProbeArchErr{driverkitYaml.Output.Probe, opts.Architecture.ToNonDeb()}
		}

		if !kr.SupportsProbe() {
			// Not an error, just throw a warning
			slog.Warn("output probe set on an unsupported kernel release", "kernelrelease", driverkitYaml.KernelRelease)
		}
	}

	// Check output driver if present
	if driverkitYaml.Output.Module != "" {
		outputModuleFilename := filepath.Base(driverkitYaml.Output.Module)
		if outputModuleFilename != outputPathFilename+".ko" {
			return &WrongOutputModuleNameErr{outputModuleFilename, outputPathFilename}
		}

		if !strings.Contains(driverkitYaml.Output.Module, opts.Architecture.ToNonDeb()) {
			return &WrongOutputModuleArchErr{driverkitYaml.Output.Module, opts.Architecture.ToNonDeb()}
		}

		if !kr.SupportsModule() {
			// Not an error, just throw a warning
			slog.Warn("output module set on an unsupported kernel release", "kernelrelease", driverkitYaml.KernelRelease)
		}
	}

	// Kernelconfigdata, if present, must be base64 encoded
	if len(driverkitYaml.KernelConfigData) > 0 && !isBase64(driverkitYaml.KernelConfigData) {
		return &KernelConfigDataNotBase64Err{}
	}

	return nil
}
