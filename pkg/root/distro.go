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

package root

import (
	"sort"

	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
)

var (
	SupportedDistroSlice []string
	// SupportedDistros keeps the list of distros supported by test-infra.
	// We don't want to generate configs for unsupported distros after all.
	// Please add new supported distros here,
	// so that the utility starts building configs for them.
	// Keys must have the same name used by driverkit targets.
	SupportedDistros = map[builder.Type]struct{}{
		builder.TargetTypeAlma:            {},
		builder.TargetTypeAmazonLinux:     {},
		builder.TargetTypeAmazonLinux2:    {},
		builder.TargetTypeAmazonLinux2022: {},
		builder.TargetTypeAmazonLinux2023: {},
		builder.TargetTypeBottlerocket:    {},
		builder.TargetTypeCentos:          {},
		builder.TargetTypeDebian:          {},
		builder.TargetTypeFedora:          {},
		builder.TargetTypeMinikube:        {},
		builder.TargetTypePhoton:          {},
		builder.TargetTypeTalos:           {},
		builder.TargetTypeUbuntu:          {},
	}
)

func init() {
	SupportedDistroSlice = make([]string, 0)
	for distro := range SupportedDistros {
		SupportedDistroSlice = append(SupportedDistroSlice, string(distro))
	}
	sort.Strings(SupportedDistroSlice)
}
