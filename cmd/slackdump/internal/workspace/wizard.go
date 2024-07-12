package workspace

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/golang/base"
	"github.com/rusq/slackdump/v3/internal/cache"
	"github.com/rusq/slackdump/v3/logger"
)

func wizSelect(ctx context.Context, cmd *base.Command, args []string) error {
	m, err := cache.NewManager(cfg.CacheDir())
	if err != nil {
		base.SetExitStatus(base.SCacheError)
		return err
	}

	wsps, err := m.List()
	if err != nil {
		base.SetExitStatus(base.SCacheError)
		return err
	}

	if len(wsps) == 0 {
		fmt.Println("No workspaces found")
		return nil
	}

	current, err := m.Current()
	if err != nil {
		base.SetExitStatus(base.SWorkspaceError)
		return fmt.Errorf("error getting the current workspace: %s", err)
	}

	columns := []table.Column{
		{Title: "C", Width: 1},
		{Title: "Name", Width: 14},
		{Title: "Team", Width: 15},
		{Title: "User", Width: 15},
		{Title: "Error", Width: 30},
	}

	var rows []table.Row
	for _, w := range wspInfo(ctx, m, current, wsps) {
		rows = append(rows, table.Row{w[0], w[1], w[4], w[5], w[6]})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		Foreground(cfg.Theme.Focused.NoteTitle.GetForeground()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(cfg.Theme.Focused.Option.GetBackground()).
		Background(cfg.Theme.Focused.SelectedOption.GetForeground()).
		Bold(false)
	t.SetStyles(s)

	mod, err := tea.NewProgram(model{table: t}).Run()
	if err != nil {
		return fmt.Errorf("workspace select wizard error: %w", err)
	}
	if newWsp := mod.(model).selected; newWsp != "" {
		if err := m.Select(newWsp); err != nil {
			base.SetExitStatus(base.SWorkspaceError)
			return fmt.Errorf("error setting the current workspace: %s", err)
		}
		logger.FromContext(ctx).Debugf("selected workspace: %s", newWsp)
	}

	return nil
}

// var baseStyle = lipgloss.NewStyle().
// 	BorderStyle(lipgloss.NormalBorder()).
// 	BorderForeground(lipgloss.Color("240"))

var baseStyle = cfg.Theme.Form

type model struct {
	table    table.Model
	selected string
	finished bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.finished = true
			return m, tea.Quit
		case "enter":
			m.selected = m.table.SelectedRow()[1]
			m.finished = true
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.finished {
		return "" // don't render the table if we've selected a workspace
	}
	return baseStyle.Render(m.table.View()) + "\n\n" + cfg.Theme.Help.Ellipsis.Render("Select the workspace with arrow keys, press [Enter] to confirm, [Esc] to cancel.")
}
