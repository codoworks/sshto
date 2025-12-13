package cmd

import (
	"github.com/spf13/cobra"

	"github.com/codoworks/sshto/internal/ssh"
)

var connectOpts ssh.ConnectOptions

var connectCmd = &cobra.Command{
	Use:     "connect <server>",
	Aliases: []string{"c"},
	Short:   "Connect to a server",
	Long:    `Connect to a server by name. Use flags to override config values.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return App.Connect(args[0], connectOpts)
	},
}

func init() {
	connectCmd.Flags().StringVarP(&connectOpts.User, "user", "u", "", "override user")
	connectCmd.Flags().IntVarP(&connectOpts.Port, "port", "p", 0, "override port")
	connectCmd.Flags().StringVarP(&connectOpts.Key, "key", "k", "", "override key file")

	// Also add these flags to root command for `sshto server --user root` usage
	rootCmd.Flags().StringVarP(&connectOpts.User, "user", "u", "", "override user")
	rootCmd.Flags().IntVarP(&connectOpts.Port, "port", "p", 0, "override port")
	rootCmd.Flags().StringVarP(&connectOpts.Key, "key", "k", "", "override key file")
}
