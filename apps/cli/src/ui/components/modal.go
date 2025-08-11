package components

import "github.com/rivo/tview"

func StyledModal(message string, onClose func()) *tview.Modal {
	return tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			if onClose != nil {
				onClose()
			}
		})
}
