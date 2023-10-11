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

package build

import (
	"github.com/falcosecurity/dbg-go/pkg/build"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewBuildConfigsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build dbg configs",
		RunE:  executeConfigs,
	}
	flags := cmd.Flags()
	flags.Bool("skip-existing", true, "whether to skip the build of drivers existing on S3")
	flags.Bool("publish", false, "whether artifacts must be published on S3")
	flags.Bool("ignore-errors", false, "whether to ignore build errors and go on looping on config files")
	flags.String("redirect-errors", "", "redirect build errors to the specified file")
	return cmd
}

func executeConfigs(_ *cobra.Command, _ []string) error {
	options := build.Options{
		Options:        root.LoadRootOptions(),
		SkipExisting:   viper.GetBool("skip-existing"),
		Publish:        viper.GetBool("publish"),
		IgnoreErrors:   viper.GetBool("ignore-errors"),
		RedirectErrors: viper.GetString("redirect-errors"),
	}
	return build.Run(options)
}
