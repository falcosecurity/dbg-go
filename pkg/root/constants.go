package root

const ConfigPathFmt = "%s/driverkit/config/%s/%s/%s" // Eg: repo-root/driverkit/config/5.0.1+driver/x86_64/centos_5.14.0-325.el9.x86_64_1.yaml

var (
	// SupportedDistros keeps the list of distros supported by test-infra.
	// We don't want to generate configs for unsupported distros after all.
	// Please add new supported build-new-drivers structures here,
	// so that the utility starts building configs for them.
	// Keys must have the same name used in kernel-crawler json keys.
	// Values must have the same name used by driverkit targets.
	SupportedDistros = map[string]string{
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
