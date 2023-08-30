package root

import "strings"

type KernelCrawlerDistro string
type DriverkitDistro string

var (
	// SupportedDistros keeps the list of distros supported by test-infra.
	// We don't want to generate configs for unsupported distros after all.
	// Please add new supported build-new-drivers structures here,
	// so that the utility starts building configs for them.
	// Keys must have the same name used in kernel-crawler json keys.
	// Values must have the same name used by driverkit targets.
	SupportedDistros = map[KernelCrawlerDistro]DriverkitDistro{
		"AlmaLinux":       "almalinux",
		"AmazonLinux":     "amazonlinux",
		"AmazonLinux2":    "amazonlinux2",
		"AmazonLinux2022": "amazonlinux2022",
		"AmazonLinux2023": "amazonlinux2023",
		"BottleRocket":    "bottlerocket",
		"CentOS":          "centos",
		"Debian":          "debian",
		"Fedora":          "fedora",
		"Minikube":        "minikube",
		"Talos":           "talos",
		"Ubuntu":          "ubuntu",
	}
)

func (kDistro KernelCrawlerDistro) ToDriverkitDistro() DriverkitDistro {
	dkDistro, found := SupportedDistros[kDistro]
	if found {
		return dkDistro
	} else {
		// Perhaps a regex? ToLower and pray
		return DriverkitDistro(strings.ToLower(string(kDistro)))
	}
}
