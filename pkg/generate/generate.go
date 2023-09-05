package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	testJsonData  []byte
	testCacheData bool
)

func loadLastRunDistro() (string, error) {
	lastDistroBytes, err := getURL(urlLastDistro)
	if err != nil {
		return "", err
	}
	lastDistro := strings.TrimSuffix(string(lastDistroBytes), "\n")
	if lastDistro == "*" {
		// Fix up regex (empty regex -> always true)
		lastDistro = ""
	}
	return lastDistro, nil
}

func Run(opts Options) error {
	slog.Info("generating config files")
	if opts.Auto {
		return autogenerateConfigs(opts)
	} else if opts.IsSet() {
		return generateSingleConfig(opts)
	}
	return fmt.Errorf(`either "auto" or target-{distro,kernelrelease,kernelversion} must be passed`)
}

// This is the only function where opts.Distro gets overridden using KernelCrawler namings
func autogenerateConfigs(opts Options) error {
	url := fmt.Sprintf(urlArchFmt, opts.Architecture.ToNonDeb())
	slog.Debug("downloading json data", "url", url)

	// Fetch kernel list json
	var (
		jsonData []byte
		err      error
	)

	// In case testJsonData is set,
	if testJsonData != nil {
		jsonData = testJsonData
	} else {
		jsonData, err = getURL(url)
		if err != nil {
			return err
		}
	}

	slog.Debug("fetched json")
	if testCacheData {
		testJsonData = jsonData
	}

	// either download latest distro from kernel-crawler
	// or translate the driverkit distro provided by the user to its kernel-crawler naming
	if opts.Distro == "load" {
		// Fetch last distro kernel-crawler ran against
		lastDistro, err := loadLastRunDistro()
		if err != nil {
			return err
		}
		slog.Info("loaded last-distro", "distro", lastDistro)

		// If lastDistro is empty it means we need to run on all supported distros; this is done automatically.
		if lastDistro != "" {
			// Map back the kernel crawler distro to our internal driverkit distro
			opts.Distro = root.ToDriverkitDistro(root.KernelCrawlerDistro(lastDistro))
			if opts.Distro == "" {
				return fmt.Errorf("kernel-crawler last run distro '%s' unsupported.\n", lastDistro)
			}
		} else {
			// This will match all supported distros
			opts.Distro = ""
		}
	}

	slog.SetDefault(slog.With("target-distro", opts.Distro))
	fullJson := map[string][]validate.DriverkitYaml{}

	// Unmarshal the big json
	err = json.Unmarshal(jsonData, &fullJson)
	if err != nil {
		return err
	}
	slog.Debug("unmarshaled json")
	var errGrp errgroup.Group

	for kcDistro, f := range fullJson {
		kernelEntries := f

		dkDistro := root.ToDriverkitDistro(root.KernelCrawlerDistro(kcDistro))

		// Skip unneeded kernel entries
		// optimization for target-distro: skip entire key
		// instead of skipping objects one by one.
		if !opts.DistroFilter(dkDistro.String()) {
			continue
		}

		// A goroutine for each distro
		errGrp.Go(func() error {
			for _, kernelEntry := range kernelEntries {
				if !opts.KernelVersionFilter(kernelEntry.KernelRelease) {
					continue
				}
				if !opts.KernelVersionFilter(kernelEntry.KernelVersion) {
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

type unsupportedTargetErr struct {
	target builder.Type
}

func (err *unsupportedTargetErr) Error() string {
	return fmt.Sprintf("target %s is unsupported by driverkit", err.target.String())
}

func loadKernelHeadersFromDk(opts Options) ([]string, error) {
	// Try to load kernel headers from driverkit. Don't error out if unable.
	// Just write a config with empty headers.

	// We already received a driverkit target type (not a kernel crawler distro!)
	targetType := opts.Distro
	b, err := builder.Factory(targetType)
	if err != nil {
		return nil, &unsupportedTargetErr{target: targetType}
	}

	// Load minimum urls for the builder
	minimumURLs := 1
	if bb, ok := b.(builder.MinimumURLsBuilder); ok {
		minimumURLs = bb.MinimumURLs()
	}

	// Load kernelrelease and architecture
	kr := kernelrelease.FromString(opts.KernelRelease)
	kr.Architecture = opts.Architecture

	// Fetch URLs
	kernelheaders, err := b.URLs(kr)
	if err != nil {
		return nil, err
	}

	// Check actually resolving URLs
	kernelheaders, err = builder.GetResolvingURLs(kernelheaders)
	if err != nil {
		return nil, err
	}
	if len(kernelheaders) < minimumURLs {
		return nil, fmt.Errorf("not enough headers packages found; expected %d, found %d", minimumURLs, len(kernelheaders))
	}
	return kernelheaders, nil
}

func generateSingleConfig(opts Options) error {
	kernelheaders, err := loadKernelHeadersFromDk(opts)
	if err != nil {
		var unsupportedTargetError *unsupportedTargetErr
		if errors.As(err, &unsupportedTargetError) {
			return unsupportedTargetError
		}
		slog.Warn(err.Error())
	}
	driverkitYaml := validate.DriverkitYaml{
		KernelVersion: opts.KernelVersion,
		KernelRelease: opts.KernelRelease,
		Target:        opts.Distro.String(),
		KernelUrls:    kernelheaders,
	}
	return dumpConfig(opts, driverkitYaml)
}

func dumpConfig(opts Options, dkYaml validate.DriverkitYaml) error {
	slog.Info("generating",
		"target", dkYaml.Target,
		"kernelrelease", dkYaml.KernelRelease,
		"kernelversion", dkYaml.KernelVersion)
	if opts.DryRun {
		slog.Info("skipping because of dry-run.")
		return nil
	}

	dkYaml.Architecture = opts.Architecture.String()

	for _, driverVersion := range opts.DriverVersion {
		dkYaml.FillOutputs(driverVersion, opts.Options)
		yamlData, pvtErr := yaml.Marshal(&dkYaml)
		if pvtErr != nil {
			return pvtErr
		}

		configPath := root.BuildConfigPath(opts.Options, driverVersion, dkYaml.ToConfigName())

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
