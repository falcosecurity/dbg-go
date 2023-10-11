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

package validate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
)

type Options struct {
	root.Options
}

type DriverkitYamlOutputs struct {
	Module string `yaml:"module"`
	Probe  string `yaml:"probe"`
}

// DriverkitYaml is the driverkit config schema
type DriverkitYaml struct {
	KernelVersion    string               `yaml:"kernelversion" json:"kernelversion"`
	KernelRelease    string               `yaml:"kernelrelease" json:"kernelrelease"`
	Target           string               `yaml:"target" json:"target"`
	Architecture     string               `yaml:"architecture"`
	Output           DriverkitYamlOutputs `yaml:"output"`
	KernelUrls       []string             `yaml:"kernelurls,omitempty" json:"headers"`
	KernelConfigData string               `yaml:"kernelconfigdata,omitempty" json:"kernelconfigdata"`
}

func (dy *DriverkitYaml) ToName() string {
	return fmt.Sprintf("%s_%s_%s", dy.Target, dy.KernelRelease, dy.KernelVersion)
}

func (dy *DriverkitYaml) ToConfigName() string {
	return fmt.Sprintf("%s.yaml", dy.ToName())
}

func (dy *DriverkitYaml) FillOutputs(driverVersion string, opts root.Options) {
	outputPath := root.BuildOutputPath(opts, driverVersion, dy.ToName())
	// Tricky because driverkit configs Outputs assume
	// that the tool is called from the `driverkit` folder of test-infra repo.
	// Only keep last 4 parts, ie: from "output/" onwards
	paths := strings.Split(outputPath, "/")
	configOutputPath := filepath.Join(paths[len(paths)-4:]...)

	kr := kernelrelease.FromString(dy.KernelRelease)
	kr.Architecture = opts.Architecture
	if kr.SupportsModule() {
		dy.Output.Module = configOutputPath + ".ko"
	}
	if kr.SupportsProbe() {
		dy.Output.Probe = configOutputPath + ".o"
	}
}
