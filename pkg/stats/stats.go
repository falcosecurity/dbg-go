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

package stats

import (
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

// Used by tests!
// We cannot simply use table = tablewriter.NewWriter(log.Default().Writer())
// as that would completely break tablewriter output.
var testOutputWriter io.Writer

func Run(opts Options, statter Statter) error {
	slog.Info(statter.Info())
	driverStatsByVersion, err := statter.GetDriverStats(opts.Options)
	if err != nil {
		return err
	}

	var table *tablewriter.Table
	if testOutputWriter != nil {
		table = tablewriter.NewWriter(testOutputWriter)
	} else {
		table = tablewriter.NewWriter(os.Stdout)
	}
	table.SetHeader([]string{"Version", "Modules", "Probes"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	totalDriverStats := driverStats{}

	data := make([]string, 3)
	// Keep keys sorted
	// (looping directly on the map {key,value} tuples gives wrong sorting sometimes).
	for _, key := range opts.DriverVersion {
		stat := driverStatsByVersion[key]
		data[0] = key
		data[1] = strconv.FormatInt(stat.NumModules, 10)
		data[2] = strconv.FormatInt(stat.NumProbes, 10)
		table.Append(data)

		totalDriverStats.NumModules += stat.NumModules
		totalDriverStats.NumProbes += stat.NumProbes
	}
	data[0] = "TOTALS"
	data[1] = strconv.FormatInt(totalDriverStats.NumModules, 10)
	data[2] = strconv.FormatInt(totalDriverStats.NumProbes, 10)
	table.Append(data)
	table.Render() // Send output

	return nil
}
