package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/tack/internal/store"
	"github.com/user/tack/internal/ui"
)

var rmCmd = &cobra.Command{
	Use:     "rm [id]",
	Aliases: []string{"remove", "delete"},
	Short:   "🗑  Remove a tack",
	Args:    cobra.ExactArgs(1),
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

		msg := t.Message
		if err := s.Remove(id); err != nil {
			return err
		}

		fmt.Printf("  %s Removed: %s\n",
			ui.SuccessStyle.Render("✓"),
			ui.CountStyle.Render(msg),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
