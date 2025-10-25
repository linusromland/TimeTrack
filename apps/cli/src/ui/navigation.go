package ui

import "github.com/rivo/tview"

type Navigator struct {
	App *tview.Application
}

func NewNavigator() *Navigator {
	return &Navigator{App: tview.NewApplication()}
}

func (n *Navigator) Show(root tview.Primitive) {
	n.App.SetRoot(root, true)
}

func (n *Navigator) Modal(content tview.Primitive) {
	modal := tview.NewFlex().SetDirection(tview.FlexRow)
	modal.AddItem(nil, 0, 1, false) // top spacer
	modal.AddItem(content, 0, 1, true)
	modal.AddItem(nil, 0, 1, false) // bottom spacer
	n.Show(modal)
}

func (n *Navigator) Run(root tview.Primitive) error {
	n.Show(root)
	return n.App.Run()
}

func (n *Navigator) Stop() {
	n.App.Stop()
}
