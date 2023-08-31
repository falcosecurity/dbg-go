package validate

import (
	"encoding/base64"
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Run(opts Options) error {
	slog.Info("validate config files")
	return root.LoopConfigsFiltered(opts.Options, "validating", func(driverVersion, configPath string) error {
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
		return fmt.Errorf("config filename is wrong (%s); should be %s", configFilename, expectedFilename)
	}

	// Check that arch is ok
	if driverkitYaml.Architecture != opts.Architecture.String() {
		return fmt.Errorf("wrong architecture in config file %s: %s", configPath, driverkitYaml.Architecture)
	}

	// Either kernelconfigdata or kernelurls must be set
	if driverkitYaml.KernelConfigData == "" && len(driverkitYaml.KernelUrls) == 0 {
		return fmt.Errorf("at least one between `kernelurls` and `kernelconfigdata` must be set: %s", configPath)
	}

	outputPath := driverkitYaml.ToOutputPath(driverVersion, opts)
	outputPathFilename := filepath.Base(outputPath)

	// Check output probe if present
	if len(driverkitYaml.Output.Probe) > 0 {
		outputProbeFilename := filepath.Base(driverkitYaml.Output.Probe)
		if outputProbeFilename != outputPathFilename+".o" {
			return fmt.Errorf("output probe filename is wrong (%s); expected: %s.o", outputProbeFilename, outputPathFilename)
		}

		if !strings.Contains(driverkitYaml.Output.Probe, opts.Architecture.ToNonDeb()) {
			return fmt.Errorf("output probe filename has wrong architecture in its path (%s); expected %s",
				driverkitYaml.Output.Probe, opts.Architecture.ToNonDeb())
		}
	}

	// Check output driver if present
	if len(driverkitYaml.Output.Module) > 0 {
		outputModuleFilename := filepath.Base(driverkitYaml.Output.Module)
		if outputModuleFilename != outputPathFilename+".ko" {
			return fmt.Errorf("output module filename is wrong (%s); expected: %s.ko", outputModuleFilename, outputPathFilename)
		}

		if !strings.Contains(driverkitYaml.Output.Module, opts.Architecture.ToNonDeb()) {
			return fmt.Errorf("output module filename has wrong architecture in its path (%s); expected %s",
				driverkitYaml.Output.Module, opts.Architecture.ToNonDeb())
		}
	}

	// Kernelconfigdata, if present, must be base64 encoded
	if len(driverkitYaml.KernelConfigData) > 0 && !isBase64(driverkitYaml.KernelConfigData) {
		return fmt.Errorf("kernelconfigdata must be a base64 encoded string")
	}

	return nil
}
