package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/codoworks/sshto/internal/config"
	"github.com/codoworks/sshto/internal/ui"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups",
	Long:  `List all configured server groups.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(App.Config.Groups) == 0 {
			fmt.Println("No groups configured. Use 'sshto groups add' to create one.")
			return
		}

		for _, g := range App.Config.Groups {
			tag := ui.GroupTag(g.Name, g.Color)
			serverCount := len(App.Config.ServersByGroup(g.Name))
			fmt.Printf("%s (%d servers)\n", tag, serverCount)
		}
	},
}

var groupsAddCmd = &cobra.Command{
	Use:   "add <name> [color]",
	Short: "Add a new group",
	Long: `Add a new group with an optional color.
Available colors: red, green, yellow, blue, magenta, cyan, white, gray`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		color := ""
		if len(args) > 1 {
			color = args[1]
		}

		group := config.Group{
			Name:  name,
			Color: color,
		}

		if err := App.Config.AddGroup(group); err != nil {
			return err
		}

		if err := App.Save(); err != nil {
			return err
		}

		fmt.Printf("Group %q added.\n", name)
		return nil
	},
}

var groupsRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a group",
	Long:    `Remove a group from the configuration. Servers in this group will not be deleted.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Check if group exists
		_, err := App.Config.FindGroup(name)
		if err != nil {
			return err
		}

		// Warn if servers use this group
		servers := App.Config.ServersByGroup(name)
		if len(servers) > 0 {
			fmt.Printf("Warning: %d server(s) belong to this group.\n", len(servers))
			fmt.Printf("Remove group %q? [y/N]: ", name)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("Canceled.")
				return nil
			}
		}

		if err := App.Config.RemoveGroup(name); err != nil {
			return err
		}

		if err := App.Save(); err != nil {
			return err
		}

		fmt.Printf("Group %q removed.\n", name)
		return nil
	},
}

func init() {
	groupsCmd.AddCommand(groupsAddCmd)
	groupsCmd.AddCommand(groupsRemoveCmd)
}
