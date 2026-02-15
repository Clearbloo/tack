package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/tack/internal/model"
	"github.com/user/tack/internal/store"
	"github.com/user/tack/internal/ui"
)

var staleFlag bool
var staleDays int

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "🗂  Bird's-eye view of all tacks across your machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New()
		if err != nil {
			return err
		}

		if jsonFlag {
			data, err := json.MarshalIndent(s.Tacks, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		dirs := s.AllDirs()
		if len(dirs) == 0 {
			fmt.Println(ui.EmptyStyle.Render("  No tacks anywhere. Your filesystem is a blank slate."))
			return nil
		}

		// Sort dirs alphabetically
		sort.Strings(dirs)

		staleThreshold := time.Duration(staleDays) * 24 * time.Hour

		fmt.Println(ui.TitleStyle.Render("🗂  Tack Board"))
		fmt.Println()

		shown := 0
		for _, dir := range dirs {
			tacks := s.ForDir(dir)

			if staleFlag {
				// Only show dirs with open TODOs older than threshold
				hasStale := false
				for _, t := range tacks {
					if t.Kind == model.KindTodo && !t.Done && time.Since(t.CreatedAt) > staleThreshold {
						hasStale = true
						break
					}
				}
				if !hasStale {
					continue
				}
			}

			fmt.Println(ui.RenderBoard(dir, tacks))
			shown++
		}

		if shown == 0 && staleFlag {
			fmt.Println(ui.SuccessStyle.Render("  ✨ No stale TODOs! You're on top of things."))
		}

		// Summary footer
		total := len(s.Tacks)
		openTodos := 0
		for _, tacks := range s.Tacks {
			for _, t := range tacks {
				if t.Kind == model.KindTodo && !t.Done {
					openTodos++
				}
			}
		}
		fmt.Printf("\n  %s\n",
			ui.CountStyle.Render(fmt.Sprintf("%d tacks across %d directories · %d open todos",
				total, len(dirs), openTodos)),
		)

		return nil
	},
}

func init() {
	boardCmd.Flags().BoolVar(&staleFlag, "stale", false, "Only show directories with stale TODOs")
	boardCmd.Flags().IntVar(&staleDays, "days", 7, "Number of days before a TODO is considered stale")
	rootCmd.AddCommand(boardCmd)
}
