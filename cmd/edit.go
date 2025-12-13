package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/codoworks/sshto/internal/ui"
)

var editCmd = &cobra.Command{
	Use:   "edit <server>",
	Short: "Edit an existing server",
	Long:  `Open an interactive form to edit an existing server configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		server, err := App.Config.FindServer(serverName)
		if err != nil {
			return err
		}

		// Make a copy for editing
		serverCopy := *server

		model := ui.NewFormModel(&serverCopy, App.Config.Groups)
		p := tea.NewProgram(model)

		finalModel, err := p.Run()
		if err != nil {
			return err
		}

		m := finalModel.(ui.FormModel)
		if m.Canceled() {
			fmt.Println("Canceled.")
			return nil
		}

		if !m.Done() {
			return nil
		}

		edited := m.Server()
		if err := App.Config.UpdateServer(serverName, *edited); err != nil {
			return err
		}

		if err := App.Save(); err != nil {
			return err
		}

		fmt.Printf("Server %q updated successfully.\n", edited.Name)
		return nil
	},
}
