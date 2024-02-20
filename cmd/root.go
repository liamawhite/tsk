package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/dashboard"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "tsk",
	RunE: func(cmd *cobra.Command, args []string) error {
        f, err := tea.LogToFile("tsk.log", "debug")
        if err != nil {
            return err
        }
        defer f.Close()

	    model := dashboard.NewModel()
        p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
