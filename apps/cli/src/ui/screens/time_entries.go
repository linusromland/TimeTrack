package screens

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"
	"TimeTrack-shared/models"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TimeEntriesScreen(nav *ui.Navigator, ctx *app.AppContext, startDate, endDate time.Time) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	statsView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	statsView.SetBorder(true)
	statsView.SetTitle(" Statistics ")

	table := tview.NewTable().
		SetSelectable(true, false)
	table.SetBorder(true).
		SetTitle(" Time Entries ")

	actionBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	currentPage := 1
	totalPages := 1

	selectionMode := false
	selectedRows := make(map[int]bool)
	entriesCache := []*models.TimeEntry{}
	projectMap := make(map[string]string)

	prettyDuration := func(seconds float64) string {
		d := time.Duration(seconds) * time.Second
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		return fmt.Sprintf("%02dh%02dm", h, m)
	}

	refreshTableHighlight := func() {
		rowCount := table.GetRowCount()
		selectedRow, _ := table.GetSelection()
		for r := 1; r < rowCount; r++ {
			for c := 0; c < table.GetColumnCount(); c++ {
				cell := table.GetCell(r, c)
				if cell == nil {
					continue
				}
				text := strings.TrimPrefix(strings.TrimPrefix(cell.Text, "[green]"), "[white]")
				if r == selectedRow {
					cell.SetText(text).SetTextColor(tcell.ColorBlack).SetBackgroundColor(tcell.ColorYellow)
				} else if selectedRows[r] {
					cell.SetText(text).SetTextColor(tcell.ColorWhite).SetBackgroundColor(tcell.ColorDarkGreen)
				} else {
					cell.SetText("[white]" + text).SetTextColor(tcell.ColorWhite).SetBackgroundColor(tcell.ColorDefault)
				}
			}
		}
	}

	updateTableTitle := func() {
		if selectionMode {
			table.SetTitle("[red] Time Entries (Selection Mode ON) ").SetBorderColor(tcell.ColorRed)
		} else {
			table.SetTitle(" Time Entries ").SetBorderColor(tcell.ColorWhite)
		}
	}

	updateActionBar := func() {
		actionBar.SetText("[yellow](D)[white] Delete   [yellow](R)[white] Report   [yellow](A)[white] Amend   " +
			"[yellow](S)[white] Toggle Selection Mode   [yellow](Space)[white] Select Row   " +
			"[yellow](N)[white] Next Page   [yellow](P)[white] Prev Page   [yellow](Q)[white] Quit")
	}

	showEntryListConfirm := func(title, action string, rows map[int]bool, singleRow int, onConfirm func()) {
		listText := ""
		if rows != nil && len(rows) > 0 {
			for r := range rows {
				if r-1 >= 0 && r-1 < len(entriesCache) {
					entry := entriesCache[r-1]
					project := projectMap[entry.ProjectID]
					if project == "" {
						project = entry.ProjectID
					}
					listText += fmt.Sprintf("• %s (%s - %s)\n",
						project,
						entry.Period.Started.Format("2006-01-02 15:04"),
						entry.Period.Ended.Format("15:04"))
				}
			}
		} else if singleRow > 0 && singleRow-1 < len(entriesCache) {
			entry := entriesCache[singleRow-1]
			project := projectMap[entry.ProjectID]
			if project == "" {
				project = entry.ProjectID
			}
			listText += fmt.Sprintf("• %s (%s - %s)\n",
				project,
				entry.Period.Started.Format("2006-01-02 15:04"),
				entry.Period.Ended.Format("15:04"))
		}

		modal := tview.NewModal().
			SetText(fmt.Sprintf("%s the following entrie(s)?\n\n%s", action, listText)).
			AddButtons([]string{"Confirm", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Confirm" {
					onConfirm()
				}
				nav.Show(TimeEntriesScreen(nav, ctx, startDate, endDate))
			})
		modal.SetTitle(title).SetBorder(true)
		nav.Show(modal)
	}

	var loadData func(page int)

	deleteEntries := func(ids []string) {
		var errs []string
		for _, id := range ids {
			if err := ctx.API.DeleteTimeEntry(id); err != nil {
				errs = append(errs, fmt.Sprintf("Failed to delete %s: %v", id, err))
			}
		}
		if len(errs) > 0 {
			errorText := strings.Join(errs, "\n")
			errorModal := tview.NewModal().
				SetText(fmt.Sprintf("[red]Some deletions failed:\n\n%s", errorText)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					nav.Show(TimeEntriesScreen(nav, ctx, startDate, endDate))
				})
			errorModal.SetTitle("Error").SetBorder(true)
			nav.Show(errorModal)
			return
		}
		loadData(currentPage)
	}

	loadData = func(page int) {
		statsView.Clear()
		table.Clear()
		selectedRows = map[int]bool{}
		updateActionBar()
		updateTableTitle()

		stats, err := ctx.API.GetTimeEntryStatistics(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))
		if err != nil {
			fmt.Fprintf(statsView, "[red]Error loading stats: %v", err)
			return
		}

		projectIDs := []string{}
		for _, p := range stats.EntriesPerProject {
			projectIDs = append(projectIDs, p.ProjectID)
		}
		if len(projectIDs) > 0 {
			projects, _ := ctx.API.GetProjectByIds(projectIDs)
			for _, pr := range projects {
				projectMap[pr.ID] = pr.Name
			}
		}

		fmt.Fprintf(statsView, "[yellow]Total Entries:[white] %d\n", stats.TotalEntries)
		fmt.Fprintf(statsView, "[yellow]Total Time:[white] %s\n", prettyDuration(float64(stats.TotalTime)))

		sort.Slice(stats.EntriesPerDate, func(i, j int) bool {
			return stats.EntriesPerDate[i].TimeFrame < stats.EntriesPerDate[j].TimeFrame
		})

		fmt.Fprintf(statsView, "\n[green]Per Day:[white]\n")
		if len(stats.EntriesPerDate) == 0 {
			fmt.Fprintf(statsView, "  [gray](no data)\n")
		}
		for _, d := range stats.EntriesPerDate {
			fmt.Fprintf(statsView, "  %s: %s\n", d.TimeFrame, prettyDuration(d.TotalTime))
		}

		fmt.Fprintf(statsView, "\n[green]Per Project:[white]\n")
		if len(stats.EntriesPerProject) == 0 {
			fmt.Fprintf(statsView, "  [gray](no data)\n")
		}
		for _, p := range stats.EntriesPerProject {
			name := projectMap[p.ProjectID]
			if name == "" {
				name = p.ProjectID
			}
			fmt.Fprintf(statsView, "  %s: %s\n", name, prettyDuration(p.TotalTime))
		}

		entries, err := ctx.API.GetTimeEntries(startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), page)
		if err != nil {
			table.SetCell(0, 0, tview.NewTableCell(fmt.Sprintf("[red]Error: %v", err)))
			return
		}

		entriesCache = entries

		if stats.TotalEntries > 0 {
			totalPages = int((stats.TotalEntries + 24) / 25)
		} else {
			totalPages = 1
		}

		entryIDs := []string{}
		for _, e := range entries {
			entryIDs = append(entryIDs, e.ProjectID)
		}
		if len(entryIDs) > 0 {
			projects, _ := ctx.API.GetProjectByIds(entryIDs)
			for _, pr := range projects {
				projectMap[pr.ID] = pr.Name
			}
		}

		headers := []string{"Project", "Start", "End", "Duration", "Note", "Reported"}
		for col, h := range headers {
			table.SetCell(0, col, tview.NewTableCell(fmt.Sprintf("[yellow]%s", h)).SetSelectable(false))
		}

		for row, e := range entries {
			project := projectMap[e.ProjectID]
			if project == "" {
				project = e.ProjectID
			}
			start := e.Period.Started.Format("2006-01-02 15:04")
			end := e.Period.Ended.Format("2006-01-02 15:04")
			duration := prettyDuration(e.Period.Ended.Sub(e.Period.Started).Seconds())
			reported := "[red]No"
			if e.Reported != nil && e.Reported.ReportedAt != nil {
				reported = fmt.Sprintf("[green]Yes (%s)", e.Reported.ReportedAt.Format("2006-01-02"))
			}

			values := []string{project, start, end, duration, e.Note, reported}
			for col, val := range values {
				table.SetCell(row+1, col, tview.NewTableCell("[white]"+val))
			}
		}
		refreshTableHighlight()
	}

	table.SetSelectionChangedFunc(func(row, column int) {
		refreshTableHighlight()
	})

	updateActionBar()

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		switch strings.ToLower(string(event.Rune())) {
		case "n":
			if currentPage < totalPages {
				currentPage++
				loadData(currentPage)
			}
		case "p":
			if currentPage > 1 {
				currentPage--
				loadData(currentPage)
			}
		case "s":
			selectionMode = !selectionMode
			updateTableTitle()
			updateActionBar()
			refreshTableHighlight()
		case " ":
			if selectionMode && row > 0 {
				selectedRows[row] = !selectedRows[row]
				refreshTableHighlight()
			}
		case "d":
			if selectionMode && len(selectedRows) > 0 {
				showEntryListConfirm("Bulk Delete", "Delete", selectedRows, -1, func() {
					ids := []string{}
					for r := range selectedRows {
						if r-1 >= 0 && r-1 < len(entriesCache) {
							ids = append(ids, entriesCache[r-1].ID)
						}
					}
					deleteEntries(ids)
				})
			} else if !selectionMode && row > 0 {
				showEntryListConfirm("Delete Entry", "Delete", nil, row, func() {
					if row-1 >= 0 && row-1 < len(entriesCache) {
						deleteEntries([]string{entriesCache[row-1].ID})
					}
				})
			}
		case "r":
			if selectionMode && len(selectedRows) > 0 {
				showEntryListConfirm("Bulk Report", "Report", selectedRows, -1, func() { /* TODO: bulk report */ })
			} else if !selectionMode && row > 0 {
				showEntryListConfirm("Report Entry", "Report", nil, row, func() { /* TODO: report */ })
			}
		case "a":
			if selectionMode && len(selectedRows) > 0 {
				warningModal := tview.NewModal().
					SetText("[red]Bulk amend is not allowed.\n\nPlease amend entries individually.").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						nav.Show(TimeEntriesScreen(nav, ctx, startDate, endDate))
					})
				warningModal.SetTitle("Bulk Amend Blocked").SetBorder(true)
				nav.Show(warningModal)
			}
		case "q":
			nav.Stop()
		}
		return event
	})

	flex.AddItem(statsView, 10, 0, false)
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(actionBar, 1, 0, false)

	loadData(currentPage)

	return flex
}
