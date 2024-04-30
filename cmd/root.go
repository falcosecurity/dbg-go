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

package cmd

import (
	"fmt"
	"github.com/falcosecurity/falcoctl/pkg/options"
	"github.com/falcosecurity/falcoctl/pkg/output"
	"github.com/pterm/pterm"
	"os"
	"runtime"
	"strings"

	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel = options.NewLogLevel()
	rootCmd  = &cobra.Command{
		Use:           "dbg-go",
		Short:         "A command line helper tool used by falcosecurity test-infra dbg.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			arch := viper.GetString("architecture")
			driverVersions := viper.GetStringSlice("driver-version")
			if _, present := kernelrelease.SupportedArchs[kernelrelease.Architecture(arch)]; !present {
				return fmt.Errorf("arch %s is not supported", arch)
			}
			if len(driverVersions) == 0 {
				if err := loadDriverVersions(); err != nil {
					return err
				}
			}
			root.Printer = output.NewPrinter(logLevel.ToPtermLogLevel(), pterm.LogFormatterColorful, os.Stdout)
			return nil
		},
	}
)

func loadDriverVersions() error {
	repoRoot := viper.GetString("repo-root")
	configPath := repoRoot + "/driverkit/config/"
	driverVersions := make([]string, 0)
	entries, err := os.ReadDir(configPath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			driverVersions = append(driverVersions, e.Name())
		}
	}

	if len(driverVersions) != 0 {
		viper.Set("driver-version", driverVersions)
		return nil
	}
	return fmt.Errorf("no driver versions found")
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		root.Printer.Logger.Fatal("failed to get working dir.", pterm.DefaultLogger.Args("err", err))
	}

	flags := rootCmd.PersistentFlags()
	flags.Bool("dry-run", false, "enable dry-run mode.")
	flags.VarP(logLevel, "log-level", "l", "set log verbosity "+logLevel.Allowed())
	flags.String("repo-root", cwd, "test-infra repository root path.")
	flags.StringP("architecture", "a", runtime.GOARCH, `architecture to run against. Supported: `+kernelrelease.SupportedArchs.String())
	flags.StringSlice("driver-version", nil, "driver versions to run against.")
	flags.String("driver-name", "falco", "driver name to be used")
	flags.String("target-kernelrelease", "",
		`target kernel release to work against. By default tool will work on any kernel release. Can be a regex.`)
	flags.String("target-kernelversion", "",
		`target kernel version to work against. By default tool will work on any kernel version. Can be a regex.`)
	flags.String("target-distro", "",
		`target distro to work against. By default tool will work on any supported distro. Can be a regex.
Supported: [`+strings.Join(root.SupportedDistroSlice, ",")+"].")

	// Custom completions
	_ = rootCmd.RegisterFlagCompletionFunc("target-distro", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return root.SupportedDistroSlice, cobra.ShellCompDirectiveDefault
	})
	_ = rootCmd.RegisterFlagCompletionFunc("architecture", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return kernelrelease.SupportedArchs.Strings(), cobra.ShellCompDirectiveDefault
	})

	// Subcommands
	rootCmd.AddCommand(configsCmd)
	rootCmd.AddCommand(s3Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
