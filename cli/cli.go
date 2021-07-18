package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"mc-tool/application"
	"mc-tool/config"
	"os"
)

const (
	OptionVerbosity = "verbosity"
	OptionRelease   = "release"
	OptionSnapshot  = "snapshot"
	OptionOldBeta   = "old_beta"
	OptionOldAlpha  = "old_alpha"
	OptionOutput    = "output"
	OptionVersion   = "version"
	OptionJar       = "jar"
	OptionMappings  = "mappings"
	OptionServer    = "server"
	OptionClient    = "client"
)

func Run(application *application.Application) {
	rootCommand := &cobra.Command{
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if verbosity, err := cmd.Flags().GetString(OptionVerbosity); err != nil {
				return err
			} else {
				return setUpLogs(os.Stderr, verbosity)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config.CliRootConfig{}
			application.RunRoot(cfg)
			return nil
		},
	}

	flags := rootCommand.PersistentFlags()
	flags.String(OptionVerbosity, logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	minecraftCommand := &cobra.Command{
		Use: "minecraft",
	}
	rootCommand.AddCommand(minecraftCommand)

	minecraftListCommand := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config.CliMinecraftListConfig{}
			err := GetValues(cmd, cfg)
			if err != nil {
				return err
			}
			application.RunMinecraftList(cfg)
			return nil
		},
	}
	minecraftCommand.AddCommand(minecraftListCommand)
	Configure(minecraftCommand, &config.CliMinecraftListConfig{})

	minecraftDownloadCommand := &cobra.Command{
		Use: "download",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config.CliMinecraftDownloadConfig{}

			err := GetValues(cmd, cfg)
			if err != nil {
				return err
			}

			application.RunMinecraftDownload(cfg)
			return nil
		},
	}
	Configure(minecraftDownloadCommand, &config.CliMinecraftDownloadConfig{})
	minecraftCommand.AddCommand(minecraftDownloadCommand)

	minecraftDecompileCommand := &cobra.Command{
		Use: "decompile",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &config.CliMinecraftDecompileConfig{}

			err := GetValues(cmd, cfg)
			if err != nil {
				return err
			}

			application.RunMinecraftDecompile(cfg)
			return nil
		},
	}
	Configure(minecraftDecompileCommand, &config.CliMinecraftDecompileConfig{})
	minecraftCommand.AddCommand(minecraftDecompileCommand)

	err := rootCommand.Execute()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}
