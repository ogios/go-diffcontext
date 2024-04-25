package main

const (
	code1 = `package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/ui/comp"
	utree "github.com/ogios/merge-repo/ui/src/u-tree"
	uview "github.com/ogios/merge-repo/ui/src/u-view"
)

type HomeDiff struct {
	HomeCore
	DiffView *DiffViewModel
}

func newHomeDiff() *HomeDiff {
	w := comp.GlobalUIData.MaxWidth - 2
	h := comp.GlobalUIData.MaxHeight - 2

	modelCount := 3
	modelsHeight := h - 1
	modelsWidth := w - 2*modelCount
	getModelWidth := modelWidthCounter(modelCount, modelsWidth)
	ms := []tea.Model{
		utree.NewTreeModel(comp.TREE_NODE, [2]int{
			getModelWidth(0.2),
			modelsHeight,
		}),
		uview.NewViewModel([2]int{
			getModelWidth(0.4),
			modelsHeight,
		}),
		NewDiffViewModel([2]int{
			getModelWidth(0.4),
			modelsHeight,
		}),
	}

	home := &HomeDiff{
		HomeCore: HomeCore{
			Models: ms,
			Tree:   ms[0],
			Text:   ms[1].(*uview.ViewModel),
		},
		DiffView: ms[2].(*DiffViewModel),
	}

	return home
}

func (m *HomeDiff) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m2 := range m.Models {
		cmds = append(cmds, m2.Init())
	}
	return tea.Batch(cmds...)
}

func (m *HomeDiff) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	update(msg, &m.HomeCore)
	return m, tea.Batch(cmds...)
}

func (m *HomeDiff) View() string {
	return view(&m.HomeCore)
}
`
	code2 = `package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/ui/comp"
	udiffview "github.com/ogios/merge-repo/ui/src/u-diffview"
	utree "github.com/ogios/merge-repo/ui/src/u-tree"
	uview "github.com/ogios/merge-repo/ui/src/u-view"
)

type HomeDiff struct {
	HomeCore
	DiffView *udiffview.DiffViewModel
}

func newHomeDiff() *HomeDiff {
	w := comp.GlobalUIData.MaxWidth - 2
	h := comp.GlobalUIData.MaxHeight - 2

	modelCount := 3
	modelsHeight := h - 1
	modelsWidth := w - 2*modelCount
	getModelWidth := modelWidthCounter(modelCount, modelsWidth)
	ms := []tea.Model{
		utree.NewTreeModel(comp.TREE_NODE, [2]int{
			getModelWidth(0.2),
			modelsHeight,
		}),
		uview.NewViewModel([2]int{
			getModelWidth(0.4),
			modelsHeight,
		}),
		udiffview.NewDiffViewModel([2]int{
			getModelWidth(0.4),
			modelsHeight,
		}),
	}

	home := &HomeDiff{
		HomeCore: HomeCore{
			Models: ms,
			Tree:   ms[0],
			Text:   ms[1].(*uview.ViewModel),
		},
		DiffView: ms[2].(*udiffview.DiffViewModel),
	}

	return home
}

func (m *HomeDiff) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m2 := range m.Models {
		cmds = append(cmds, m2.Init())
	}
	return tea.Batch(cmds...)
}

func (m *HomeDiff) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	update(msg, &m.HomeCore)
	return m, tea.Batch(cmds...)
}

func (m *HomeDiff) View() string {
	return view(&m.HomeCore)
}
`
)

// const code1 = `
// 	dmp := diffmatchpatch.New()
//
// 	diffs := dmp.DiffMain(text1, text2, false)
//
// 	fmt.Println(dmp.DiffPrettyText(diffs))
// `
// const code2 = `
// 	diffs := dmp.DiffMain(text1, text2, false)
//
// 	dmp := diffmatchpatch.New()
//
// 	fmt.Println(dmp.DiffPrettyText(diffs))
// `

// const code1 = `
// const (
// 	text1 = "Lorem ipsum dolor."
// 	text2 = "Lorem dolor sit amet."
// )
// `
//
// const code2 = `
// const (
// 	text2 = "Lorem sit ipsum dolor amet."
//     text3 = "nothing"
// )
// `

// const code1 = `
// "Lorem ipsum dolor."
// 	"Lorem dolor sit amet."
// `
// const code2 = `
// "Lorem sit ipsum dolor amet."
//     "nothing"
// `
