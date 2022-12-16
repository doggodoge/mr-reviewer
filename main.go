package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"mr-reviewer/browser"
	"mr-reviewer/config"
	"mr-reviewer/fetch"
)

var (
	modelStyle = lipgloss.NewStyle().Margin(1, 2)
	conf       = initConf()
)

func initConf() *config.Config {
	conf, err := config.Read()
	if err != nil {
		panic(err)
	}
	return conf
}

const (
	RepositoryView = iota
	MergeRequestView
)

type model struct {
	list        list.Model
	currentView int
	showDraft   bool
	mrResponse  *fetch.MRsResponse
}

func initialModel() model {
	list := list.New(conf.RepositoriesAsItems(), list.NewDefaultDelegate(), 0, 0)
	list.Title = "MR Reviewer"
	return model{
		list:        list,
		currentView: RepositoryView,
		showDraft:   false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.currentView {
			case RepositoryView:
				item, ok := m.list.SelectedItem().(config.Repository)
				if ok {
					m.currentView = MergeRequestView

					res, err := fetch.FetchMRsFromRepo(conf, item.Route)
					if err != nil {
						panic(err)
					}
					m.mrResponse = res

					newMRList := m.mrResponse.ToListItems(m.showDraft)
					cmd := m.list.SetItems(newMRList)
					return m, cmd
				}

			case MergeRequestView:
				item, ok := m.list.SelectedItem().(config.Repository)
				if ok {
					browser.OpenURL(item.Route)
					return m, tea.Quit
				}
			}

		case "d":
			if m.currentView == MergeRequestView {
				m.showDraft = !m.showDraft
				m.list.SetItems(m.mrResponse.ToListItems(m.showDraft))
			}

		case "backspace":
			if m.currentView == MergeRequestView {
				cmd := m.list.SetItems(conf.RepositoriesAsItems())
				m.currentView = RepositoryView
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		h, v := modelStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return modelStyle.Render(m.list.View())
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
