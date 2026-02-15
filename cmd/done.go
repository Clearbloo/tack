package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/tack/internal/model"
	"github.com/user/tack/internal/store"
	"github.com/user/tack/internal/ui"
)

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "☑ Mark a TODO as done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New()
		if err != nil {
			return err
		}

		id := args[0]
		t := s.FindByID(id)
		if t == nil {
			return fmt.Errorf("no tack with id %q", id)
		}

		if t.Kind != model.KindTodo {
			return fmt.Errorf("tack %s is a %s, not a todo", id, t.Kind)
		}

		if t.Done {
			fmt.Printf("  %s Already done!\n", ui.CountStyle.Render("—"))
			return nil
		}

		now := time.Now()
		t.Done = true
		t.DoneAt = &now

		if err := s.Save(); err != nil {
			return err
		}

		fmt.Printf("  %s Done: %s\n",
			ui.SuccessStyle.Render("✓"),
			ui.DoneStyle.Render(t.Message),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
