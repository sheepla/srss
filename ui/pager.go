package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

const useHighPerformanceRenderer = true

var (
	titleStyle = func() lip.Style {
		b := lip.NormalBorder()
		b.Right = "├"
		return lip.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lip.Style {
		b := lip.NormalBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	title    string
	content  string
	ready    bool
	viewport viewport.Model
}

func NewPager(item *gofeed.Item) (*tea.Program, error) {
	program := tea.NewProgram(
		&model{
			title:   item.Title,
			content: renderContent(item),
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	return program, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
		if msg.String() == "g" {
			m.viewport.GotoTop()
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
		if msg.String() == "G" {
			m.viewport.GotoBottom()
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case tea.WindowSizeMsg:
		headerHeight := lip.Height(m.renderHeader())
		footerHeight := lip.Height(m.renderFooter())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		cmds = append(cmds, viewport.Sync(m.viewport))
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.renderHeader(), m.viewport.View(), m.renderFooter())
}

func (m *model) renderHeader() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("─", larger(0, m.viewport.Width-lip.Width(title)))
	return lip.JoinHorizontal(lip.Center, title, line)
}

func (m *model) renderFooter() string {
	info := infoStyle.Render(scrollPercent(m.viewport.ScrollPercent()))
	line := strings.Repeat("─", larger(0, m.viewport.Width-lip.Width(info)))
	return lip.JoinHorizontal(lip.Center, line, info)
}

func scrollPercent(p float64) string {
	return fmt.Sprintf("%3.f%%", p*100)
}

func larger(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderContent(item *gofeed.Item) string {
	var author string
	if item.Author != nil {
		sprintfIfNotEmpty("by %s ", item.Author.Name)
	}
	return fmt.Sprintf(
		`%s%s %s
──────
%s
%s
──────
%s
`,
		author,
		sprintfIfNotEmpty("published at %s", item.Published),
		sprintfIfNotEmpty("updated at %s", item.Updated),
		sprintfIfNotEmpty("%s", item.Description),
		sprintfIfNotEmpty("%s", item.Content),
		sprintfIfNotEmpty("%s", strings.Join(item.Links, "\n")),
	)
}

func sprintfIfNotEmpty(format string, str string) string {
	if str == "" {
		return ""
	}
	return fmt.Sprintf(format, str)
}
