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

var pinCmd = &cobra.Command{
	Use:   "pin [message]",
	Short: "📌 Pin a note to the current directory",
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
			Kind:      model.KindPin,
			Message:   strings.Join(args, " "),
			Dir:       dir,
			CreatedAt: time.Now(),
		}

		if err := s.Add(t); err != nil {
			return err
		}

		fmt.Printf("  %s Pinned %s\n",
			ui.SuccessStyle.Render("✓"),
			ui.PinStyle.Render(t.Message),
		)
		fmt.Printf("  %s\n", ui.IDStyle.Render("id: "+t.ID))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pinCmd)
}
