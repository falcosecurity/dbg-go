package autogenerate

const (
	urlArchFmt    = "https://raw.githubusercontent.com/falcosecurity/kernel-crawler/kernels/%s/list.json"
	urlLastDistro = "https://raw.githubusercontent.com/falcosecurity/kernel-crawler/kernels/last_run_distro.txt"
)

var (
	// Please add new supported build-new-drivers structures here,
	// so that the utility starts building configs for them.
	// Fields must have the same name used in kernel-crawler json keys.
	SupportedDistros = []string{
		"AlmaLinux",
		"AmazonLinux",
		"AmazonLinux2",
		"AmazonLinux2022",
		"AmazonLinux2023",
		"BottleRocket",
		"CentOS",
		"Debian",
		"Fedora",
		"Minikube",
		"Talos",
		"Ubuntu",
	}
)
