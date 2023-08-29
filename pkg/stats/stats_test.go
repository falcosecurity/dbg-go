package stats

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"testing"
)

type DkConfigNamed struct {
	validate.DriverkitYaml
	confPath string
}

func TestStats(t *testing.T) {
	configPath := fmt.Sprintf(root.ConfigPathFmt,
		"./test/",
		"1.0.0+driver",
		"x86_64",
		"")

	dkConfigs := []DkConfigNamed{
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: validate.DriverkitYamlOutputs{
					Module: fmt.Sprintf(validate.OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(validate.OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.10.0",
						"1"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "centos_5.10.0_1.yaml",
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.15.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: validate.DriverkitYamlOutputs{
					Module: fmt.Sprintf(validate.OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.15.0",
						"1"),
					Probe: fmt.Sprintf(validate.OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.15.0",
						"1"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "centos_5.15.0_1.yaml",
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion: "13",
				KernelRelease: "5.15.0",
				Target:        "ubuntu",
				Architecture:  "amd64",
				Output: validate.DriverkitYamlOutputs{
					Module: fmt.Sprintf(validate.OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"ubuntu",
						"5.15.0",
						"13"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "ubuntu_5.15.0_13.yaml",
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.15.25",
				Target:        "bottlerocket",
				Architecture:  "amd64",
				Output: validate.DriverkitYamlOutputs{
					Module: fmt.Sprintf(validate.OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"bottlerocket",
						"5.15.25",
						"1"),
					Probe: fmt.Sprintf(validate.OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"bottlerocket",
						"5.15.25",
						"1"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "bottlerocket_5.15.25_1.yaml",
		},
	}

	err := os.MkdirAll(configPath, 0700)
	t.Cleanup(func() {
		os.RemoveAll("./test")
	})
	assert.NoError(t, err)

	// Create all configs needed by the test
	for _, dkConf := range dkConfigs {
		file, err := os.OpenFile(dkConf.confPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		assert.NoError(t, err)
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
				Architecture:  "x86_64",
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
				Architecture:  "x86_64",
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
				Architecture:  "aarch64", // not present
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
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "CentOS",
					KernelRelease: "",
					KernelVersion: "",
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
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "Cent*",
					KernelRelease: "",
					KernelVersion: "",
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
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "",
					KernelRelease: "5.10.*",
					KernelVersion: "",
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
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "",
					KernelRelease: "",
					KernelVersion: "1",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  3,
				NumModules: 3,
			},
		},
	}

	// Store logged data, will be used by test
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(io.Writer(&buf), nil))
	slog.SetDefault(logger)

	// capture output!
	testOutputWriter = log.Default().Writer()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err = Run(test.opts, NewFileStatter())
			assert.NoError(t, err)

			// Use logged output to ensure we fetched correct stats
			type MessageJSON struct {
				Message string `json:"msg"`
			}
			var messageJSON MessageJSON
			scanner := bufio.NewScanner(&buf)
			outputStats := driverStats{}
			startParsing := false
			parsingIdx := 0
			for scanner.Scan() {
				err = json.Unmarshal(scanner.Bytes(), &messageJSON)
				assert.NoError(t, err)
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
					break
				}
			}
			assert.Equal(t, test.expectedStats, outputStats)
			buf.Reset()
		})
	}
}
