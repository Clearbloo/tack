package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/tack/internal/ui"
)

const bashHook = `# tack shell hook — add this to your .bashrc or .bash_profile
_tack_cd() {
    builtin cd "$@" && tack --quiet
}
alias cd='_tack_cd'`

const zshHook = `# tack shell hook — add this to your .zshrc
autoload -U add-zsh-hook
_tack_chpwd() {
    tack --quiet
}
add-zsh-hook chpwd _tack_chpwd`

const fishHook = `# tack shell hook — add this to ~/.config/fish/conf.d/tack.fish
function _tack_cd --on-variable PWD
    tack --quiet
end`

var hookCmd = &cobra.Command{
	Use:   "hook [bash|zsh|fish]",
	Short: "🪝 Print shell hook for auto-display on cd",
	Long: `Print a shell hook snippet that shows tack summaries
whenever you cd into a directory.

Add the output to your shell config:
  tack hook bash >> ~/.bashrc
  tack hook zsh  >> ~/.zshrc
  tack hook fish >> ~/.config/fish/conf.d/tack.fish`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish"},
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := args[0]
		switch shell {
		case "bash":
			fmt.Println(bashHook)
		case "zsh":
			fmt.Println(zshHook)
		case "fish":
			fmt.Println(fishHook)
		default:
			return fmt.Errorf("unknown shell %q — use bash, zsh, or fish", shell)
		}

		fmt.Println()
		fmt.Printf("  %s Copy the above into your shell config, then restart your shell.\n",
			ui.SuccessStyle.Render("💡"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
