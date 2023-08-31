package stats

import (
	"fmt"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

type DkConfigNamed struct {
	validate.DriverkitYaml
	hasProbe  bool
	hasModule bool
}

func TestStats(t *testing.T) {
	configPath := root.BuildConfigPath(root.Options{
		RepoRoot:     "./test/",
		Architecture: "amd64",
	}, "1.0.0+driver", "")

	fmt.Println(configPath)
	dkConfigs := []DkConfigNamed{
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "1",
				KernelRelease:    "5.10.0",
				Target:           "centos",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasProbe:  true,
			hasModule: true,
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "1",
				KernelRelease:    "5.15.0",
				Target:           "centos",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasProbe:  true,
			hasModule: true,
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "13",
				KernelRelease:    "5.15.0",
				Target:           "ubuntu",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasModule: true,
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "1",
				KernelRelease:    "5.15.25",
				Target:           "bottlerocket",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasProbe:  true,
			hasModule: true,
		},
	}

	err := os.MkdirAll(configPath, 0700)
	t.Cleanup(func() {
		_ = os.RemoveAll("./test")
	})
	assert.NoError(t, err)

	// Create all configs needed by the test
	for _, dkConf := range dkConfigs {
		file, err := os.OpenFile(configPath+dkConf.ToConfigName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		assert.NoError(t, err)

		outputPath := dkConf.ToOutputPath("1.0.0+driver",
			validate.Options{
				DriverName: "falco",
				Options: root.Options{
					Architecture: kernelrelease.Architecture(dkConf.Architecture),
				},
			})
		if dkConf.hasModule {
			dkConf.DriverkitYaml.Output.Module = outputPath + ".ko"
		}
		if dkConf.hasProbe {
			dkConf.DriverkitYaml.Output.Probe = outputPath + ".o"
		}
		enc := yaml.NewEncoder(file)
		err = enc.Encode(dkConf.DriverkitYaml)
		_ = file.Close()
		assert.NoError(t, err)
	}

	tests := map[string]struct {
		opts          Options
		expectedStats driverStats
	}{
		"stats 1.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
			}},
			expectedStats: driverStats{
				NumProbes:  3,
				NumModules: 4,
			},
		},
		"stats 2.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"2.0.0+driver"}, // not present
			}},
			expectedStats: driverStats{
				NumProbes:  0,
				NumModules: 0,
			},
		},
		"stats 1.0.0+driver aarch64": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "arm64", // not present
				DriverVersion: []string{"1.0.0+driver"},
			}},
			expectedStats: driverStats{
				NumProbes:  0,
				NumModules: 0,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by distro": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro: "centos",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  2,
				NumModules: 2,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by distro regex": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro: "cent*",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  2,
				NumModules: 2,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by kernel release": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					KernelRelease: "5.10.*",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  1,
				NumModules: 1,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by kernel version": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					KernelVersion: "1",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  3,
				NumModules: 3,
			},
		},
	}

	// capture output!

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Use logged output to ensure we fetched correct stats
			type MessageJSON struct {
				Message string `json:"msg"`
			}
			var messageJSON MessageJSON
			outputStats := driverStats{}
			startParsing := false
			parsingIdx := 0
			utils.RunTestParsingLogs(t,
				func() error {
					testOutputWriter = log.Default().Writer()
					return Run(test.opts, NewFileStatter())
				},
				&messageJSON,
				func() bool {
					messageJSON.Message = strings.ReplaceAll(messageJSON.Message, " ", "")
					// Example lines:
					//{"time":"2023-08-29T11:38:35.692782942+02:00","level":"INFO","msg":"1.0.0+driver"}
					//{"time":"2023-08-29T11:38:35.692784013+02:00","level":"INFO","msg":""}
					//{"time":"2023-08-29T11:38:35.692784775+02:00","level":"INFO","msg":"|"}
					//{"time":"2023-08-29T11:38:35.692785487+02:00","level":"INFO","msg":""}
					//{"time":"2023-08-29T11:38:35.69279484+02:00","level":"INFO","msg":"4"}
					//{"time":"2023-08-29T11:38:35.692796064+02:00","level":"INFO","msg":""}
					//{"time":"2023-08-29T11:38:35.692797001+02:00","level":"INFO","msg":"|"}
					//{"time":"2023-08-29T11:38:35.692797848+02:00","level":"INFO","msg":""}
					//{"time":"2023-08-29T11:38:35.69280042+02:00","level":"INFO","msg":"3"}
					if startParsing {
						parsingIdx++
						if parsingIdx%4 == 0 {
							n, err := strconv.ParseInt(messageJSON.Message, 10, 64)
							assert.NoError(t, err)
							switch parsingIdx / 4 {
							case 1:
								outputStats.NumModules = n
							case 2:
								outputStats.NumProbes = n
							}
						}
					}
					if messageJSON.Message == "1.0.0+driver" {
						startParsing = true
					} else if parsingIdx == 8 {
						// parsed both numbers
						return false // break out
					}
					return true // continue
				})
			assert.Equal(t, test.expectedStats, outputStats)
		})
	}
}

func TestStatsS3(t *testing.T) {
	keysToBeCreated := []string{
		"driver/1.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
		"driver/1.0.0+driver/x86_64/falco_amazonlinux2022_5.10.96-90.460.amzn2022.x86_64_1.o",
		"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.o",
		"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.ko",
		"driver/2.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
		"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
		"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
	}
	client := utils.S3CreateTestBucket(t, keysToBeCreated)
	statter := s3Statter{client: client}

	tests := map[string]struct {
		opts          Options
		expectedStats driverStatsByDriverVersion
	}{
		"stats 1.0.0+driver, 2.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  2,
					NumModules: 2,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 2.0.0+driver aarch64": {
			opts: Options{Options: root.Options{
				Architecture:  "arm64",
				DriverVersion: []string{"2.0.0+driver"},
			}},
			expectedStats: driverStatsByDriverVersion{
				"2.0.0+driver": {
					NumProbes:  1,
					NumModules: 1,
				},
			},
		},
		"stats 2.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"2.0.0+driver"},
			}},
			expectedStats: driverStatsByDriverVersion{
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 1.0.0+driver, 2.0.0+driver x86_64 filtered by distro": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				Target: root.Target{
					Distro: "almalinux",
				},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 1.0.0+driver, 2.0.0+driver x86_64 filtered by distro regex": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				Target: root.Target{
					Distro: "almali*",
				},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 1.0.0+driver, 2.0.0+driver x86_64 filtered by kernelrelease": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				Target: root.Target{
					KernelRelease: "5.*",
				},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  1,
					NumModules: 1,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dStats, err := statter.GetDriverStats(test.opts.Options)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedStats, dStats)
		})
	}
}
