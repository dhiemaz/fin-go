package cmd

import (
	"fmt"
	"github.com/dhiemaz/fin-go/config"
	"github.com/dhiemaz/fin-go/infrastructure/logger"
	"github.com/spf13/cobra"
)

// CommandEngine is the structure of cli
type Command struct {
	rootCmd *cobra.Command
}

var text = `VISION`

// NewCommandEngine the command line boot loader
func NewCommand() *Command {
	var rootCmd = &cobra.Command{
		Use:   "FIN-GO",
		Short: "Financial Go (Banking) API",
		Long:  "Financial Go (Banking) API",
	}

	return &Command{
		rootCmd: rootCmd,
	}
}

// Run the all command line
func (c *Command) Run() {
	var rootCommands = []*cobra.Command{
		{
			Use:   "serve",
			Short: "Run Financial Go (Banking) API HTTP server",
			Long:  "Run Financial Go (Banking) API HTTP server",
			PreRun: func(cmd *cobra.Command, args []string) {
				// initialize config
				config.InitConfig()

				// Show display text
				fmt.Println(fmt.Sprintf(text))

				logger.WithFields(logger.Fields{"component": "command", "action": "serve with watcher"}).
					Infof("PreRun command done")
			},
			Run: func(cmd *cobra.Command, args []string) {
				//action.RunRest()
			},
			PostRun: func(cmd *cobra.Command, args []string) {
				// close database connection
				defer config.GetConfig().DBPool.Close()
				logger.WithFields(logger.Fields{"component": "command", "action": "serve with watcher"}).
					Infof("PostRun command done")
			},
		},
	}

	for _, command := range rootCommands {
		c.rootCmd.AddCommand(command)
	}

	c.rootCmd.Execute()
}

// GetRoot the command line service
func (c *Command) GetRoot() *cobra.Command {
	return c.rootCmd
}
