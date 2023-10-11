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
	"github.com/falcosecurity/dbg-go/cmd/cleanup"
	"github.com/falcosecurity/dbg-go/cmd/publish"
	"github.com/falcosecurity/dbg-go/cmd/stats"
	"github.com/spf13/cobra"
)

var (
	s3Cmd = &cobra.Command{
		Use:   "drivers",
		Short: "Work with remote drivers bucket",
		Long: `Read only commands will use an S3 client with anonymous credentials.
Write commands will need proper "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY" environment variables set.
`,
	}
)

func init() {
	// Subcommands
	s3Cmd.AddCommand(cleanup.NewCleanupDriversCmd())
	s3Cmd.AddCommand(stats.NewStatsDriversCmd())
	s3Cmd.AddCommand(publish.NewPublishDriversCmd())
}
