package generate

import (
	"encoding/json"
	"fmt"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/fededp/dbg-go/pkg/validate"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	logger "log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	testJsonData  []byte
	testCacheData bool
)

func loadLastRunDistro() (string, error) {
	lastDistroBytes, err := utils.GetURL(urlLastDistro)
	if err != nil {
		return "", err
	}
	lastDistro := strings.TrimSuffix(string(lastDistroBytes), "\n")
	if lastDistro == "*" {
		// Fix up regex (emtpy regex -> always true)
		lastDistro = ""
	}
	return lastDistro, nil
}

func Run(opts Options) error {
	if opts.Auto {
		url := fmt.Sprintf(urlArchFmt, opts.Architecture)
		logger.Debug("downloading json data", "url", url)

		// Fetch kernel list json
		var (
			jsonData []byte
			err      error
		)

		// In case testJsonData is set,
		if testJsonData != nil {
			jsonData = testJsonData
		} else {
			jsonData, err = utils.GetURL(url)
		}
		if err != nil {
			return err
		}
		logger.Debug("fetched json")
		if testCacheData {
			testJsonData = jsonData
		}

		if opts.Distro == "" {
			// Fetch last distro kernel-crawler ran against
			opts.Distro, err = loadLastRunDistro()
			if err != nil {
				return err
			}
			logger.Debug("loaded last-distro")
		}
		return autogenerateConfigs(opts, jsonData)
	} else if opts.IsSet() {
		return generateSingleConfig(opts)
	}
	return fmt.Errorf(`either "auto" or target-{distro,kernelrelease,kernelversion} must be passed`)
}

func autogenerateConfigs(opts Options, jsonData []byte) error {
	logger.SetDefault(logger.With("target-distro", opts.Distro))

	distroFilter := func(distro string) bool {
		matched, _ := regexp.MatchString(opts.Distro, distro)
		return matched
	}

	kernelreleaseFilter := func(kernelrelease string) bool {
		matched, _ := regexp.MatchString(opts.KernelRelease, kernelrelease)
		return matched
	}

	kernelversionFilter := func(kernelversion string) bool {
		matched, _ := regexp.MatchString(opts.KernelVersion, kernelversion)
		return matched
	}

	// Generate a dynamic struct with all needed distros
	// NOTE: we might need a single distro when `lastDistro` is != "*";
	// else, we will add all SupportedDistros found in constants.go
	distroCtr := 0
	instanceBuilder := dynamicstruct.NewStruct()
	for distro, _ := range root.SupportedDistros {
		if distroFilter(distro) {
			tag := fmt.Sprintf(`json:"%s"`, distro)
			instanceBuilder.AddField(distro, []validate.KernelEntry{}, tag)
			distroCtr++
		}
	}
	if distroCtr == 0 {
		return fmt.Errorf("unsupported target distro: %s. Must be one of: %v", opts.Distro, root.SupportedDistros)
	}
	dynamicInstance := instanceBuilder.Build().New()

	// Unmarshal the big json
	err := json.Unmarshal(jsonData, &dynamicInstance)
	if err != nil {
		return err
	}
	logger.Debug("unmarshaled json")
	var errGrp errgroup.Group

	reader := dynamicstruct.NewReader(dynamicInstance)
	for _, f := range reader.GetAllFields() {
		logger.Info("generating configs", "distro", f.Name())
		if opts.DryRun {
			logger.Info("skipping because of dry-run.")
			continue
		}
		kernelEntries := f.Interface().([]validate.KernelEntry)
		// A goroutine for each distro
		errGrp.Go(func() error {
			for _, kernelEntry := range kernelEntries {
				// Skip unneeded kernel entries
				if !kernelreleaseFilter(kernelEntry.KernelRelease) {
					continue
				}
				if !kernelversionFilter(kernelEntry.KernelVersion) {
					continue
				}

				if pvtErr := dumpConfig(opts, kernelEntry); pvtErr != nil {
					return pvtErr
				}
			}
			return nil
		})
	}
	return errGrp.Wait()
}

func generateSingleConfig(opts Options) error {
	// Try to load kernel headers from driverkit
	var kernelheaders []string
	targetType := builder.Type(root.SupportedDistros[opts.Distro])
	b, err := builder.Factory(targetType)
	if err != nil {
		return err
	}
	kr := kernelrelease.FromString(opts.KernelRelease)
	kr.Architecture = kernelrelease.Architecture(utils.ToDebArch(opts.Architecture))
	kernelheaders, err = b.URLs(kr)
	if err != nil {
		return err
	}
	kernelEntry := validate.KernelEntry{
		KernelVersion: opts.KernelVersion,
		KernelRelease: opts.KernelRelease,
		Target:        opts.Distro,
		Headers:       kernelheaders,
	}
	return dumpConfig(opts, kernelEntry)
}

func dumpConfig(opts Options, kernelEntry validate.KernelEntry) error {
	driverkitYaml := validate.DriverkitYaml{
		KernelVersion:    kernelEntry.KernelVersion,
		KernelRelease:    kernelEntry.KernelRelease,
		Target:           kernelEntry.Target,
		Architecture:     utils.ToDebArch(opts.Architecture),
		KernelUrls:       kernelEntry.Headers,
		KernelConfigData: string(kernelEntry.KernelConfigData),
	}

	for _, driverVersion := range opts.DriverVersion {
		outputPath := fmt.Sprintf(validate.OutputPathFmt,
			driverVersion,
			opts.Architecture,
			opts.DriverName,
			kernelEntry.Target,
			kernelEntry.KernelRelease,
			kernelEntry.KernelVersion)
		driverkitYaml.Output = validate.DriverkitYamlOutputs{
			Module: outputPath + ".ko",
			Probe:  outputPath + ".o",
		}
		yamlData, pvtErr := yaml.Marshal(&driverkitYaml)
		if pvtErr != nil {
			return pvtErr
		}

		configPath := fmt.Sprintf(root.ConfigPathFmt,
			opts.RepoRoot,
			driverVersion,
			opts.Architecture,
			kernelEntry.ToConfigName())

		// Make sure folder exists
		pvtErr = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		if pvtErr != nil {
			return pvtErr
		}
		fW, pvtErr := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if pvtErr != nil {
			return pvtErr
		}
		_, _ = fW.Write(yamlData)
		_ = fW.Close()
	}
	return nil
}
