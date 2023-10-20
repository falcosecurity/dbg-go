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

type KernelCrawlerDistro string

var (
	SupportedDistroSlice []string
	// SupportedDistros keeps the list of distros supported by test-infra.
	// We don't want to generate configs for unsupported distros after all.
	// Please add new supported build-new-drivers structures here,
	// so that the utility starts building configs for them.
	// Keys must have the same name used by driverkit targets.
	// Values must have the same name used by kernel-crawler json keys.
	SupportedDistros = map[builder.Type]KernelCrawlerDistro{
		builder.TargetTypeAlma:            "AlmaLinux",
		builder.TargetTypeAmazonLinux:     "AmazonLinux",
		builder.TargetTypeAmazonLinux2:    "AmazonLinux2",
		builder.TargetTypeAmazonLinux2022: "AmazonLinux2022",
		builder.TargetTypeAmazonLinux2023: "AmazonLinux2023",
		builder.TargetTypeBottlerocket:    "BottleRocket",
		builder.TargetTypeCentos:          "CentOS",
		builder.TargetTypeDebian:          "Debian",
		builder.TargetTypeFedora:          "Fedora",
		builder.TargetTypeMinikube:        "Minikube",
		builder.TargetTypePhoton:          "PhotonOS",
		builder.TargetTypeTalos:           "Talos",
		builder.TargetTypeUbuntu:          "Ubuntu",
	}
)

func init() {
	SupportedDistroSlice = make([]string, 0)
	for distro := range SupportedDistros {
		SupportedDistroSlice = append(SupportedDistroSlice, string(distro))
	}
	sort.Strings(SupportedDistroSlice)
}

func ToDriverkitDistro(distro KernelCrawlerDistro) builder.Type {
	for key, val := range SupportedDistros {
		if val == distro {
			return key
		}
	}
	return builder.Type("")
}
