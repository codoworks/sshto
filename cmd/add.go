package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/codoworks/sshto/internal/ui"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new server",
	Long:  `Open an interactive form to add a new server to the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		model := ui.NewFormModel(nil, App.Config.Groups)
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

		server := m.Server()
		if err := App.Config.AddServer(*server); err != nil {
			return err
		}

		if err := App.Save(); err != nil {
			return err
		}

		fmt.Printf("Server %q added successfully.\n", server.Name)
		return nil
	},
}
