package autogenerate

import (
	"encoding/json"
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/ompluscator/dynamic-struct"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func loadLastRunDistro() (string, error) {
	lastDistroBytes, err := utils.GetURL(urlLastDistro)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(lastDistroBytes), "\n"), nil
}

func Run(opts Options) error {
	url := fmt.Sprintf(urlArchFmt, opts.Architecture)
	logger.Debug("downloading json data from: ", url)

	// Fetch kernel list json
	jsonData, err := utils.GetURL(url)
	if err != nil {
		return err
	}
	logger.Debug("fetched json")

	return generateConfigs(opts, jsonData)
}

func generateConfigs(opts Options, jsonData []byte) error {
	var err error
	if opts.Target == "" {
		// Fetch last distro kernel-crawler ran against
		opts.Target, err = loadLastRunDistro()
		if err != nil {
			return err
		}
		logger.Debug("loaded last-distro: ", opts.Target)
	}
	if opts.Target != "*" && !slices.Contains(SupportedDistros, opts.Target) {
		return fmt.Errorf("unsupported target distro: %s. Must be one of: %v", opts.Target, SupportedDistros)
	}

	// Generate a dynamic struct with all needed distros
	// NOTE: we might need a single distro when `lastDistro` is != "*";
	// else, we will add all SupportedDistros found in constants.go
	instanceBuilder := dynamicstruct.NewStruct()
	for _, distro := range SupportedDistros {
		if opts.Target == "*" || distro == opts.Target {
			tag := fmt.Sprintf(`json:"%s"`, distro)
			instanceBuilder.AddField(distro, []validate.KernelEntry{}, tag)
		}
	}
	dynamicInstance := instanceBuilder.Build().New()

	// Unmarshal the big json
	err = json.Unmarshal(jsonData, &dynamicInstance)
	if err != nil {
		return err
	}
	logger.Debug("unmarshaled json")
	var errGrp errgroup.Group

	reader := dynamicstruct.NewReader(dynamicInstance)
	for _, f := range reader.GetAllFields() {
		logger.Infof("generating configs for %s\n", f.Name())
		if opts.DryRun {
			logger.Info("skipping because of dry-run.")
			continue
		}
		kernelEntries := f.Interface().([]validate.KernelEntry)
		// A goroutine for each distro
		errGrp.Go(func() error {
			for _, kernelEntry := range kernelEntries {
				driverkitYaml := validate.DriverkitYaml{
					KernelVersion:    kernelEntry.KernelVersion,
					KernelRelease:    kernelEntry.KernelRelease,
					Target:           kernelEntry.Target,
					Architecture:     utils.ToDebArch(opts.Architecture),
					KernelUrls:       kernelEntry.Headers,
					KernelConfigData: string(kernelEntry.KernelConfigData),
				}

				kernelEntryConfName := kernelEntry.ToConfigName()

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
						kernelEntryConfName)

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
			}
			return nil
		})
	}
	return errGrp.Wait()
}
