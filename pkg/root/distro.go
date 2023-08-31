package root

import (
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
	"sort"
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
		builder.TargetTypeTalos:           "Talos",
		builder.TargetTypeUbuntu:          "Ubuntu",
	}
)

func init() {
	SupportedDistroSlice = make([]string, 0)
	for distro, _ := range SupportedDistros {
		SupportedDistroSlice = append(SupportedDistroSlice, string(distro))
	}
	sort.Strings(SupportedDistroSlice)
}
