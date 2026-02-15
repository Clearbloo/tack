package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/tack/internal/model"
	"github.com/user/tack/internal/store"
	"github.com/user/tack/internal/ui"
)

var warnCmd = &cobra.Command{
	Use:   "warn [message]",
	Short: "⚠️  Pin a warning to the current directory",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New()
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		t := model.Tack{
			ID:        model.NewID(),
			Kind:      model.KindWarn,
			Message:   strings.Join(args, " "),
			Dir:       dir,
			CreatedAt: time.Now(),
		}

		if err := s.Add(t); err != nil {
			return err
		}

		fmt.Printf("  %s Warning pinned %s\n",
			ui.SuccessStyle.Render("✓"),
			ui.WarnStyle.Render(t.Message),
		)
		fmt.Printf("  %s\n", ui.IDStyle.Render("id: "+t.ID))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(warnCmd)
}
