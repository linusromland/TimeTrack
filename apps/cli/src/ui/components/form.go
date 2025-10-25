package components

import "github.com/rivo/tview"

func StyledForm(title string) *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true).
		SetTitle(" "+title+" ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 2, 2)
	return form
}
