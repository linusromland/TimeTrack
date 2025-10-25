package screens

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/database"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/components"
	"TimeTrack-cli/src/utils"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DashboardScreen(nav *ui.Navigator, ctx *app.AppContext) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	statusBox := tview.NewTextView().SetDynamicColors(true)
	statusBox.SetBorder(true).SetTitle(" Status ")

	statusBox.Clear()
	_, _ = fmt.Fprintf(statusBox, "Server: %s\n", colorStatus(getServerStatus(ctx)))
	_, _ = fmt.Fprintf(statusBox, "Server URL: %s\n", colorStatus(getServerURL(ctx)))
	_, _ = fmt.Fprintf(statusBox, "User: %s\n", colorStatus(getUserStatus(ctx)))
	_, _ = fmt.Fprintf(statusBox, "Atlassian Integration: %s\n", colorStatus(getAtlassianStatus(ctx)))

	actions := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)
	actions.SetBorder(true).SetTitle(" Actions ")

	_, _ = fmt.Fprintf(actions, "[yellow](E)[-] Edit Server URL  |  [yellow](L)[-] Login  |  [yellow](R)[-] Register  |  [yellow](A)[-] Atlassian Auth  |  [yellow](Q)[-] Quit")

	// Capture key presses
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'e', 'E':
			nav.Show(EditSettingModal(nav, ctx))
		case 'l', 'L':
			nav.Show(LoginModal(nav, ctx, false))
		case 'r', 'R':
			nav.Show(RegisterModal(nav, ctx, false))
		case 'a', 'A':
			doAtlassianAuth(nav, ctx)
		case 'q', 'Q':
			nav.Stop()
		}
		return event
	})

	flex.AddItem(statusBox, 0, 5, false)
	flex.AddItem(actions, 0, 1, true)

	return flex
}

func getServerStatus(ctx *app.AppContext) string {
	health, err := ctx.API.HealthCheck()
	if err != nil {
		return fmt.Sprintf("Unhealthy - %s", err.Error())
	}
	return fmt.Sprintf("Healthy - Version %s", health.Version)
}

func getUserStatus(ctx *app.AppContext) string {
	user, err := ctx.API.GetCurrentUser()
	if err != nil {
		return "Unauthorized / Not logged in"
	}
	return fmt.Sprintf("Logged in as %s", user.Email)
}

func getAtlassianStatus(ctx *app.AppContext) string {
	user, err := ctx.API.GetCurrentUser()
	if err != nil {
		return "Disabled"
	}
	if user.Integration.Atlassian.Enabled {
		return "Enabled"
	}
	return "Disabled"
}

func getServerURL(ctx *app.AppContext) string {
	return ctx.DB.Get(database.ServerURLKey)
}

func colorStatus(text string) string {
	switch {
	case startsWith(text, "Healthy"), startsWith(text, "Logged in"), startsWith(text, "Enabled"):
		return "[green]" + text + "[-]"
	case startsWith(text, "Unhealthy"), startsWith(text, "Unauthorized"), startsWith(text, "Disabled"):
		return "[red]" + text + "[-]"
	default:
		return "[yellow]" + text + "[-]"
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func doAtlassianAuth(nav *ui.Navigator, ctx *app.AppContext) {
	url, err := ctx.API.GetAtlassianAuthURL()
	if err != nil {
		nav.Show(components.StyledModal("Error: "+err.Error(), func() { nav.Show(DashboardScreen(nav, ctx)) }))
		return
	}
	err = utils.OpenBrowser(url)
	if err != nil {
		nav.Show(components.StyledModal("Error: "+err.Error(), func() { nav.Show(DashboardScreen(nav, ctx)) }))
		return
	}
	nav.Show(components.StyledModal("Opened browser for authentication.\nFollow instructions.", func() {
		nav.Stop()
	}))
}
