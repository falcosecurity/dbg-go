package cmd

import (
	"fmt"
	"github.com/fededp/dbg-go/cmd/autogenerate"
	"github.com/fededp/dbg-go/cmd/cleanup"
	"github.com/fededp/dbg-go/cmd/validate"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	logger "log/slog"
	"os"
	"runtime"
)

var (
	rootCmd = &cobra.Command{
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
			if !utils.IsArchSupported(arch) {
				return fmt.Errorf("arch %s is not supported", arch)
			}
			if len(driverVersions) == 0 {
				if err := loadDriverVersions(); err != nil {
					return err
				}
			}
			return initLogger(cmd.Name())
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

// NewRootCmd instantiates the root command.
func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	flags := rootCmd.PersistentFlags()
	flags.Bool("dry-run", false, "enable dry-run mode.")
	flags.StringP("log-level", "l", logger.LevelInfo.String(), "set log verbosity.")
	flags.String("repo-root", cwd, "test-infra repository root path.")
	flags.StringP("architecture", "a", utils.FromDebArch(runtime.GOARCH), "architecture to run against.")
	flags.StringSlice("driver-version", nil, "driver versions to generate configs against.")

	// Custom completions
	rootCmd.RegisterFlagCompletionFunc("architecture", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.SupportedArchList(), cobra.ShellCompDirectiveDefault
	})

	// Subcommands
	rootCmd.AddCommand(autogenerate.Cmd)
	rootCmd.AddCommand(cleanup.Cmd)
	rootCmd.AddCommand(validate.Cmd)
}

func initLogger(subcmd string) error {
	var programLevel = new(logger.LevelVar) // Info by default
	h := logger.NewJSONHandler(os.Stderr, &logger.HandlerOptions{Level: programLevel})
	logger.SetDefault(logger.New(h).With("cmd", subcmd))

	logLevel := viper.GetString("log-level")
	return programLevel.UnmarshalText([]byte(logLevel))
}

func Execute() error {
	return rootCmd.Execute()
}
