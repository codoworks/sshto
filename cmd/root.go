package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/codoworks/sshto/internal/app"
	"github.com/codoworks/sshto/internal/config"
)

var (
	cfgFile string
	App     *app.App
)

var rootCmd = &cobra.Command{
	Use:   "sshto [server]",
	Short: "SSH connection manager with interactive menu",
	Long: `sshto is an SSH connection manager that provides an interactive
menu for selecting and connecting to SSH servers.

Run without arguments to open the interactive server selection menu.
Run with a server name to connect directly.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			// Direct connection mode
			return connectCmd.RunE(cmd, args)
		}
		// Interactive mode
		return listCmd.RunE(cmd, args)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show config file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(App.Config.Path())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initApp)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/sshto/config.yaml)")

	// Add subcommands
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(groupsCmd)
}

func initApp() {
	path := cfgFile
	if path == "" {
		path = config.DefaultPath()
	}

	var err error
	App, err = app.New(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
}
