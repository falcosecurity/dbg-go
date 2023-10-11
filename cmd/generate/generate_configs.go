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

package generate

import (
	"github.com/falcosecurity/dbg-go/pkg/generate"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGenerateConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new dbg configs",
		Long: `In auto mode, configs will be generated starting from kernel-crawler output. 
In this scenario, --target-{distro,kernelrelease,kernelversion} are available to filter to-be-generated configs. Regexes are allowed.
Moreover, you can pass special value "load" as target-distro to make the tool automatically fetch latest distro kernel-crawler ran against.
Instead, when auto mode is disabled, the tool is able to generate a single config (for each driver version).
In this scenario, --target-{distro,kernelrelease,kernelversion} CANNOT be regexes but must be exact values.
Also, in non-automatic mode, kernelurls will be retrieved using driverkit libraries.
`,
		RunE: executeConfigs,
	}
	flags := cmd.Flags()
	flags.Bool("auto", false, "automatically generate configs from kernel-crawler output")
	return cmd
}

func executeConfigs(_ *cobra.Command, _ []string) error {
	options := generate.Options{
		Options: root.LoadRootOptions(),
		Auto:    viper.GetBool("auto"),
	}
	return generate.Run(options)
}
