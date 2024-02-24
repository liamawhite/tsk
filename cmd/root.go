package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liamawhite/tsk/pkg/models/router"
	"github.com/liamawhite/tsk/pkg/task"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
         	Use: "tsk",
	RunE: func(cmd *cobra.Command, args []string) error {
        f, err := setupLogging(slog.LevelDebug)
        if err != nil {
            return err
        }
        defer f.Close()

        dir, err := persistenceDir()
        if err != nil {
            return err
        }
        tasksDir, err := ensureDir(filepath.Join(dir, "tasks"))
        if err != nil {
            return err
        }
        
        taskClient, err := task.NewClient(tasksDir)
        if err != nil {
            return err
        }

	    model := router.NewModel(taskClient)
        p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func setupLogging(level slog.Level) (*os.File, error) {
    path := "tsk.log"
    f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
    if err != nil {
		return nil, fmt.Errorf("error opening file for logging: %w", err)
	}

    l := new(slog.LevelVar)
    h := slog.NewTextHandler(f, &slog.HandlerOptions{Level: l})
    slog.SetDefault(slog.New(h))
    l.Set(level)
   
    return f, nil
}

func persistenceDir() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return ensureDir(filepath.Join(home, ".tsk"))
}

func ensureDir(dir string) (string, error) {
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        err := os.Mkdir(dir, 0755)
        if err != nil {
            return "", err
        }
    }
    return dir, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
