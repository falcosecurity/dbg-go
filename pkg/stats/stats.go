package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"log"
	logger "log/slog"
	"os"
	"strconv"
)

func Run(opts Options) error {
	logger.Info("fetching stats from existing config files")
	driverStatsByVersion := make(map[string]driverStats)
	totalDriverStats := driverStats{}
	err := root.LoopConfigsFiltered(opts.Options, "computing stats", func(driverVersion, configPath string) error {
		dStats := driverStatsByVersion[driverVersion]
		err := getConfigStats(&dStats, configPath)
		driverStatsByVersion[driverVersion] = dStats
		return err
	})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(log.Default().Writer())
	table.SetHeader([]string{"Version", "Modules", "Probes", "Headers", "KernelConfigData"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	data := make([]string, 5)
	for key, stat := range driverStatsByVersion {
		data[0] = key
		data[1] = strconv.FormatInt(stat.NumModules, 10)
		data[2] = strconv.FormatInt(stat.NumProbes, 10)
		data[3] = strconv.FormatInt(stat.NumHeaders, 10)
		data[4] = strconv.FormatInt(stat.NumKernelConfigDatas, 10)
		table.Append(data)

		totalDriverStats.NumModules += stat.NumModules
		totalDriverStats.NumProbes += stat.NumProbes
		totalDriverStats.NumHeaders += stat.NumHeaders
		totalDriverStats.NumKernelConfigDatas += stat.NumKernelConfigDatas
	}
	data[0] = "TOTALS"
	data[1] = strconv.FormatInt(totalDriverStats.NumModules, 10)
	data[2] = strconv.FormatInt(totalDriverStats.NumProbes, 10)
	data[3] = strconv.FormatInt(totalDriverStats.NumHeaders, 10)
	data[4] = strconv.FormatInt(totalDriverStats.NumKernelConfigDatas, 10)
	table.Append(data)
	table.Render() // Send output

	return nil
}

func getConfigStats(dStats *driverStats, configPath string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var driverkitYaml validate.DriverkitYaml
	err = yaml.Unmarshal(configData, &driverkitYaml)
	if err != nil {
		return err
	}

	logger.Debug("fetching stats", "parsedConfig", driverkitYaml)

	if driverkitYaml.Output.Probe != "" {
		dStats.NumProbes++
	}
	if driverkitYaml.Output.Module != "" {
		dStats.NumModules++
	}
	if len(driverkitYaml.KernelUrls) > 0 {
		dStats.NumHeaders++
	}
	if driverkitYaml.KernelConfigData != "" {
		dStats.NumKernelConfigDatas++
	}
	return nil
}
