// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generate

import (
	"errors"
	"fmt"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/dbg-go/pkg/validate"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	json "github.com/json-iterator/go"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
)

var (
	testJsonData  []byte
	testCacheData bool
)

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
				if !opts.KernelReleaseFilter(kernelEntry.KernelRelease) {
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

func loadKernelHeadersFromDk(opts Options) ([]string, error) {
	// Try to load kernel headers from driverkit. Don't error out if unable.
	// Just write a config with empty headers.

	// We already received a driverkit target type (not a kernel crawler distro!)
	targetType := opts.Distro
	b, err := builder.Factory(targetType)
	if err != nil {
		return nil, &UnsupportedTargetErr{target: targetType}
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
		var unsupportedTargetError *UnsupportedTargetErr
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

	// Sort kernelurls, so that we always get the same sorting for dbg configs.
	slices.Sort(dkYaml.KernelUrls)

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
