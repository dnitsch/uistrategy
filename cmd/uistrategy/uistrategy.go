package cmd

import (
	"os"

	"github.com/dnitsch/configmanager"
	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy"
	"github.com/dnitsch/uistrategy/internal/cmdutil"
	"github.com/dnitsch/uistrategy/internal/util"
	"github.com/spf13/cobra"
)

var (
	path    string
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "uistrategy",
		RunE:  runActions,
		Short: "executes a series of actions against a setup config",
		Long:  `executes a series of instructions against a any number of paths under the same host. supports multiple login options - basic/Idp/MFA e.g. `,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.Exit(err)
	}
	util.CleanExit()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVarP(&path, "input", "i", "", "Path to the input file containing the config definition for the UIStrategy")
}

// runActions parses and executes the provided actions
func runActions(cmd *cobra.Command, args []string) error {
	conf := &uistrategy.UiStrategyConf{}
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	cm := &configmanager.ConfigManager{}

	if err := cmdutil.YamlParseInput(conf, f, cm); err != nil {
		return err
	}
	ui := uistrategy.New(conf.Setup).WithLogger(logger(verbose))

	return cmdutil.RunActions(ui, conf)
}

func logger(verbose bool) log.Logger {
	if verbose {
		return log.New(os.Stderr, log.DebugLvl)
	}
	return log.New(os.Stderr, log.ErrorLvl)
}
