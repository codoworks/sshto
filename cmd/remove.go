package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var removeForce bool

var removeCmd = &cobra.Command{
	Use:     "remove <server>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a server",
	Long:    `Remove a server from the configuration.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		// Check if server exists
		_, err := App.Config.FindServer(serverName)
		if err != nil {
			return err
		}

		// Confirm unless --force is used
		if !removeForce {
			fmt.Printf("Are you sure you want to remove server %q? [y/N]: ", serverName)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("Canceled.")
				return nil
			}
		}

		if err := App.Config.RemoveServer(serverName); err != nil {
			return err
		}

		if err := App.Save(); err != nil {
			return err
		}

		fmt.Printf("Server %q removed.\n", serverName)
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "skip confirmation")
}
