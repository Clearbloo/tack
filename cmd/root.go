package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/tack/internal/store"
	"github.com/user/tack/internal/ui"
)

var (
	quietFlag bool
	jsonFlag  bool
)

var rootCmd = &cobra.Command{
	Use:   "tack",
	Short: "📌 Sticky context for your terminal",
	Long: `tack — pin notes, TODOs, and warnings to directories.

Your filesystem becomes a spatial notebook. Every time you enter
a directory, see what you left for yourself.

  tack pin "remember: API key rotates monthly"
  tack todo "fix flaky test in auth_test.go"
  tack warn "DO NOT deploy from this branch"
  tack                          show tacks here
  tack board                    see all tacks everywhere`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New()
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		tacks := s.ForDir(dir)

		// JSON output
		if jsonFlag {
			data, err := json.MarshalIndent(tacks, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		// Quiet mode for shell hook
		if quietFlag {
			summary := ui.RenderQuietSummary(tacks)
			if summary != "" {
				fmt.Println(summary)
			}
			return nil
		}

		// Full display
		if len(tacks) == 0 {
			fmt.Println(ui.RenderEmpty())
			return nil
		}

		fmt.Println(ui.TitleStyle.Render(fmt.Sprintf("📋 Tacks in %s", ui.RenderDirHeader(dir))))
		fmt.Println()
		for _, t := range tacks {
			fmt.Println(ui.RenderTack(t))
		}
		fmt.Println()

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("error: "+err.Error()))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output as JSON")
	rootCmd.Flags().BoolVarP(&quietFlag, "quiet", "q", false, "Compact one-line summary (for shell hooks)")
}
