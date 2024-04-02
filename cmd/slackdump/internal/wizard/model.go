package wizard

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/cfg"
	"github.com/rusq/slackdump/v3/cmd/slackdump/internal/workspace"
)

type model struct {
	form *huh.Form
	val  string
}

func newModel(m *menu) model {
	var options []huh.Option[string]
	for i, name := range m.names {
		var text = fmt.Sprintf("%-10s - %s", name, m.items[i].Description)
		if m.items[i].Description == "" {
			text = fmt.Sprintf("%-10s", name)
		}
		options = append(options, huh.NewOption(text, name))
	}
	return model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("selection").
					Title(m.title).
					Description(currentWsp()).
					Options(options...),
			),
		),
	}
}

func (m *model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		m.val = m.form.GetString("selection")
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return m.form.View()
}

func currentWsp() string {
	if current, err := workspace.Current(cfg.CacheDir(), cfg.Workspace); err == nil {
		return fmt.Sprintf("Slack Workspace:  %s", current)
	}
	return ""
}
