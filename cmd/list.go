package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/codoworks/sshto/internal/ssh"
	"github.com/codoworks/sshto/internal/ui"
)

var listGroup string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "Interactive server selection",
	Long:    `Open an interactive fuzzy-filterable list of servers to connect to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		servers := App.Config.Servers
		if listGroup != "" {
			servers = ui.FilterByGroup(servers, listGroup)
		}

		if len(servers) == 0 {
			fmt.Println("No servers configured. Use 'sshto add' to add a server.")
			return nil
		}

		model := ui.NewListModel(servers, App.Config.Groups)
		p := tea.NewProgram(model, tea.WithAltScreen())

		finalModel, err := p.Run()
		if err != nil {
			return err
		}

		m := finalModel.(ui.ListModel)
		selected := m.Selected()
		if selected == nil {
			return nil
		}

		fmt.Printf("Connecting to %s...\n", selected.Name)
		return App.Connect(selected.Name, ssh.ConnectOptions{
			User: connectOpts.User,
			Port: connectOpts.Port,
			Key:  connectOpts.Key,
		})
	},
}

func init() {
	listCmd.Flags().StringVarP(&listGroup, "group", "g", "", "filter by group")
}
