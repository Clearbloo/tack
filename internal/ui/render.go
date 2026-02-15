package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/tack/internal/model"
)

var (
	// Colors
	Subtle  = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	Accent  = lipgloss.AdaptiveColor{Light: "#7B2FBE", Dark: "#BD93F9"}
	Warning = lipgloss.AdaptiveColor{Light: "#CC6600", Dark: "#FFB86C"}
	Success = lipgloss.AdaptiveColor{Light: "#2D8B46", Dark: "#50FA7B"}
	Dimmed  = lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#555555"}

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Accent).
			MarginBottom(1)

	IDStyle = lipgloss.NewStyle().
		Foreground(Subtle).
		Width(6)

	PinStyle = lipgloss.NewStyle().
			Foreground(Accent)

	TodoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#0077CC", Dark: "#8BE9FD"})

	WarnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Warning)

	DoneStyle = lipgloss.NewStyle().
			Foreground(Dimmed).
			Strikethrough(true)

	DirStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	CountStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#FF5555"}).
			Bold(true)

	EmptyStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			Italic(true)

	BoardBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Subtle).
			Padding(0, 1).
			MarginBottom(1)
)

// RenderTack renders a single tack line.
func RenderTack(t model.Tack) string {
	id := IDStyle.Render(t.ID)

	switch {
	case t.Kind == model.KindTodo && t.Done:
		emoji := model.DoneEmoji()
		msg := DoneStyle.Render(t.Message)
		age := timeAgo(t.CreatedAt)
		return fmt.Sprintf("  %s %s %s  %s", id, emoji, msg, CountStyle.Render(age))

	case t.Kind == model.KindWarn:
		emoji := t.Kind.Emoji()
		msg := WarnStyle.Render(t.Message)
		age := timeAgo(t.CreatedAt)
		return fmt.Sprintf("  %s %s %s  %s", id, emoji, msg, CountStyle.Render(age))

	case t.Kind == model.KindTodo:
		emoji := t.Kind.Emoji()
		msg := TodoStyle.Render(t.Message)
		age := timeAgo(t.CreatedAt)
		return fmt.Sprintf("  %s %s %s  %s", id, emoji, msg, CountStyle.Render(age))

	default: // pin
		emoji := t.Kind.Emoji()
		msg := PinStyle.Render(t.Message)
		age := timeAgo(t.CreatedAt)
		return fmt.Sprintf("  %s %s %s  %s", id, emoji, msg, CountStyle.Render(age))
	}
}

// RenderDirHeader renders the header for a directory section.
func RenderDirHeader(dir string) string {
	short := shortenPath(dir)
	return DirStyle.Render("📁 " + short)
}

// RenderQuietSummary renders the compact one-line summary for the shell hook.
func RenderQuietSummary(tacks []model.Tack) string {
	if len(tacks) == 0 {
		return ""
	}

	pins, todos, warns, done := 0, 0, 0, 0
	for _, t := range tacks {
		switch {
		case t.Kind == model.KindTodo && t.Done:
			done++
		case t.Kind == model.KindPin:
			pins++
		case t.Kind == model.KindTodo:
			todos++
		case t.Kind == model.KindWarn:
			warns++
		}
	}

	var parts []string
	if pins > 0 {
		parts = append(parts, fmt.Sprintf("📌 %d note%s", pins, plural(pins)))
	}
	if todos > 0 {
		parts = append(parts, fmt.Sprintf("☐ %d todo%s", todos, plural(todos)))
	}
	if warns > 0 {
		parts = append(parts, WarnStyle.Render(fmt.Sprintf("⚠️  %d warning%s", warns, plural(warns))))
	}
	if done > 0 {
		parts = append(parts, CountStyle.Render(fmt.Sprintf("☑ %d done", done)))
	}

	return CountStyle.Render("📋 ") + strings.Join(parts, CountStyle.Render(" · "))
}

// RenderEmpty renders the "no tacks here" message.
func RenderEmpty() string {
	return EmptyStyle.Render("  No tacks here. Use `tack pin`, `tack todo`, or `tack warn` to add one.")
}

// RenderBoard renders a board view for a directory and its tacks.
func RenderBoard(dir string, tacks []model.Tack) string {
	var lines []string
	lines = append(lines, RenderDirHeader(dir))
	for _, t := range tacks {
		lines = append(lines, RenderTack(t))
	}
	content := strings.Join(lines, "\n")
	return BoardBoxStyle.Render(content)
}

// CountSummary returns a short summary like "2 pins · 1 todo · 1 warning".
func CountSummary(tacks []model.Tack) string {
	pins, todos, warns := 0, 0, 0
	for _, t := range tacks {
		if t.Done {
			continue
		}
		switch t.Kind {
		case model.KindPin:
			pins++
		case model.KindTodo:
			todos++
		case model.KindWarn:
			warns++
		}
	}

	var parts []string
	if pins > 0 {
		parts = append(parts, fmt.Sprintf("%d pin%s", pins, plural(pins)))
	}
	if todos > 0 {
		parts = append(parts, fmt.Sprintf("%d todo%s", todos, plural(todos)))
	}
	if warns > 0 {
		parts = append(parts, fmt.Sprintf("%d warning%s", warns, plural(warns)))
	}
	if len(parts) == 0 {
		return "all clear"
	}
	return strings.Join(parts, " · ")
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func shortenPath(dir string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return dir
	}
	if strings.HasPrefix(dir, home) {
		return "~" + strings.TrimPrefix(dir, home)
	}
	return dir
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		if days < 30 {
			return fmt.Sprintf("%dd ago", days)
		}
		return filepath.Base(t.Format("Jan 2"))
	}
}
