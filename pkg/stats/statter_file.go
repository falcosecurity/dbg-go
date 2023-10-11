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

package stats

import (
	"log/slog"
	"os"

	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/dbg-go/pkg/validate"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type fileStatter struct {
	root.Looper
}

func NewFileStatter() Statter {
	return &fileStatter{Looper: root.NewFsLooper(root.BuildConfigPath)}
}

func (f *fileStatter) Info() string {
	return "gathering stats for local config files"
}

func (f *fileStatter) GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error) {
	driverStatsByVersion := make(driverStatsByDriverVersion)
	err := f.LoopFiltered(opts, "computing stats", "config", func(driverVersion, configPath string) error {
		dStats := driverStatsByVersion[driverVersion]
		err := getConfigStats(&dStats, configPath)
		driverStatsByVersion[driverVersion] = dStats
		return err
	})
	return driverStatsByVersion, err
}

func getConfigStats(dStats *driverStats, configPath string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var driverkitYaml validate.DriverkitYaml
	err = yaml.Unmarshal(configData, &driverkitYaml)
	if err != nil {
		return errors.WithMessagef(err, "config: %s", configPath)
	}

	slog.Debug("fetching stats", "parsedConfig", driverkitYaml)

	if driverkitYaml.Output.Probe != "" {
		dStats.NumProbes++
	}
	if driverkitYaml.Output.Module != "" {
		dStats.NumModules++
	}
	return nil
}
