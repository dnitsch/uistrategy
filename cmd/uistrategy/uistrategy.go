package cmd

import (
	"os"

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
		Short: "executes a series of actions against a URL",
		Long:  ``,
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
	ui := uistrategy.New().WithLogger(logger(verbose))

	if err := cmdutil.YamlParseInput(conf, path); err != nil {
		return err
	}

	return cmdutil.RunActions(ui, conf)
}

func logger(verbose bool) log.Logger {
	if verbose {
		return log.New(os.Stderr, log.DebugLvl)
	}
	return log.New(os.Stderr, log.ErrorLvl)
}
