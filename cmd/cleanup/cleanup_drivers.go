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

package cleanup

import (
	"github.com/falcosecurity/dbg-go/pkg/cleanup"
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/spf13/cobra"
)

func NewCleanupDriversCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup desired remote drivers",
		RunE:  executeDrivers,
	}
	return cmd
}

func executeDrivers(_ *cobra.Command, _ []string) error {
	cleaner, err := cleanup.NewS3Cleaner()
	if err != nil {
		return err
	}
	return cleanup.Run(cleanup.Options{Options: root.LoadRootOptions()}, cleaner)
}
