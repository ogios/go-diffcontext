package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/ui/comp"
	udiffview "github.com/ogios/merge-repo/ui/src/u-diffview"
	utree "github.com/ogios/merge-repo/ui/src/u-tree"
	uview "github.com/ogios/merge-repo/ui/src/u-view"
)

type HomeDiff struct {
	DiffView *udiffview.DiffViewModel
	HomeCore
}

func newHomeDiff() *HomeDiff {
	w := comp.GlobalUIData.MaxWidth - 2
	h := comp.GlobalUIData.MaxHeight - 2

	modelCount := 3
	modelsHeight := h - 1
	modelsWidth := w - 2*modelCount
	getModelWidth := modelWidthCounter(modelCount, modelsWidth)
	ms := []*childModel{
		newChild([2]int{getModelWidth(0.2), modelsHeight}),
		newChild([2]int{getModelWidth(0.4), modelsHeight}),
		newChild([2]int{getModelWidth(0.4), modelsHeight}),
	}
	ms[0].m = utree.NewTreeModel(comp.TREE_NODE, ms[0].block)
	ms[1].m = uview.NewViewModel(ms[1].block)
	ms[2].m = udiffview.NewDiffViewModel(ms[2].block)

	home := &HomeDiff{
		HomeCore: HomeCore{
			Models: ms,
			Text:   ms[1].m.(*uview.ViewModel),
		},
		DiffView: ms[2].m.(*udiffview.DiffViewModel),
	}

	return home
}

func (m *HomeDiff) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m2 := range m.Models {
		cmds = append(cmds, m2.m.Init())
	}
	return tea.Batch(cmds...)
}

func (m *HomeDiff) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case utree.FileMsg:
		cmds = append(cmds, func() tea.Msg {
			m.DiffView.ViewFile(msg.FileRelPath)
			return nil
		})
	}
	cmds = append(cmds, update(msg, &m.HomeCore))
	return m, tea.Batch(cmds...)
}

func (m *HomeDiff) View() string {
	return view(&m.HomeCore)
}
