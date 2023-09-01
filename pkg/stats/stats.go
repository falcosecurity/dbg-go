package stats

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"log/slog"
	"os"
	"strconv"
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
