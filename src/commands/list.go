package commands

import (
	"TimeTrack/src/calendar"
	"TimeTrack/src/database"
	"TimeTrack/src/utils"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	googleCalendar "google.golang.org/api/calendar/v3"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "list all tasks",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "start",
			Value:   time.Now().Format("2006-01-02"),
			Aliases: []string{"s"},
			Usage:   "Start Date of task. (format: YYYY-MM-DD)",
		},
		&cli.StringFlag{
			Name:    "end",
			Aliases: []string{"e"},
			Usage:   "End date of task. If not set, the start date will be used. (format: YYYY-MM-DD)",
		},
		&cli.StringFlag{
			Name:    "unit",
			Aliases: []string{"u"},
			Usage:   "Unit of time. (options: d (days), w (week), m (month))",
		},
	},
	Action: func(c *cli.Context) error {
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		calendarId := database.GetData(db, "calendarId")
		if calendarId == "" {
			fmt.Println("No calendar selected. Please select a calendar with the selectCalendar command.")
			return nil
		}

		if !utils.IsValidDate(c.String("start")) {
			return cli.Exit("Invalid start date. Please use the following format: YYYY-MM-DD", 1)
		}

		if c.String("end") != "" && !utils.IsValidDate(c.String("end")) {
			return cli.Exit("Invalid end date. Please use the following format: YYYY-MM-DD", 1)
		}

		if c.String("end") == "" {
			c.Set("end", c.String("start"))
		}

		startTime := fmt.Sprintf("%sT00:00:00+02:00", c.String("start"))
		endTime := fmt.Sprintf("%sT23:59:59+02:00", c.String("end"))
		parsedStart, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return nil
		}

		parsedEnd, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			fmt.Printf("Error parsing end time: %v\n", err)
			return nil
		}

		events := calendar.GetEvents(calendarId, startTime, endTime)

		if len(events.Items) == 0 {
			fmt.Println("No events found.")
			return nil
		}

		parsedEvents := googleEventsToEvents(events)

		unit := c.String("unit")
		printEvents(parsedEvents, parsedStart, parsedEnd, unit)

		return nil
	},
}

type Event struct {
	Summary     string
	Description string
	Start       time.Time
	End         time.Time
}

type TicketDuration struct {
	issueKey    string
	duration    int64
	description string
}

type Duration struct {
	period   string
	sortkey  string
	duration int64
}

func googleEventsToEvents(events *googleCalendar.Events) []Event {
	var parsedEvents []Event
	for _, event := range events.Items {
		start, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return nil
		}

		end, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			fmt.Printf("Error parsing end time: %v\n", err)
			return nil
		}

		parsedEvents = append(parsedEvents, Event{
			Summary:     event.Summary,
			Description: event.Description,
			Start:       start,
			End:         end,
		})
	}
	return parsedEvents
}

func getUnit(daysBetween int) string {
	if daysBetween > 60 {
		return "month"
	} else if daysBetween > 14 {
		return "week"
	} else if daysBetween == 0 {
		return "day"
	} else {
		return "days"
	}
}

func verifyUnit(unit string, daysBetween int) (string, error) {
	unitChar := strings.ToLower(string(unit[0]))

	if unitChar == "d" || unit == "w" || unit == "m" {
		if unitChar == "d" && daysBetween == 0 {
			return "day", nil
		}

		unitMap := map[string]string{
			"d": "days",
			"w": "week",
			"m": "month",
		}

		return unitMap[unitChar], nil
	} else {
		return "", fmt.Errorf("invalid unit. Valid units are: d (days), w (week), m (month)")
	}
}

func getDurationInSeconds(start, end time.Time) int64 {
	duration := end.Sub(start)
	return int64(duration.Seconds())
}

func formatTime(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func getTicketDurations(events []Event) []TicketDuration {
	ticketMap := make(map[string]*TicketDuration)

	for _, event := range events {
		duration := getDurationInSeconds(event.Start, event.End)
		issueKey := event.Summary

		ticket, exists := ticketMap[issueKey]
		if !exists {
			ticket = &TicketDuration{
				issueKey: issueKey,
			}
			ticketMap[issueKey] = ticket
		}

		ticket.duration += duration
		if event.Description != "" {
			ticket.description += "\n" + event.Description
		}
	}

	var ticketDurations []TicketDuration
	for _, ticket := range ticketMap {
		ticketDurations = append(ticketDurations, *ticket)
	}

	sort.Slice(ticketDurations, func(i, j int) bool {
		return ticketDurations[i].duration > ticketDurations[j].duration
	})

	return ticketDurations
}

func getTotalTime(events []Event) int64 {
	var totalTime int64
	for _, event := range events {
		totalTime += getDurationInSeconds(event.Start, event.End)
	}
	return totalTime
}

func printEvents(events []Event, startDate, endDate time.Time, unit string) {
	daysBetween := int(endDate.Sub(startDate).Hours() / 24)

	if unit == "" {
		unit = getUnit(daysBetween)
	} else {
		parsedUnit, err := verifyUnit(unit, daysBetween)
		if err != nil {
			fmt.Println(err)
			return
		}
		unit = parsedUnit
	}

	totalTimePerTicket := getTicketDurations(events)
	totalTime := getTotalTime(events)

	// Print ticket durations
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Issue Key", "Time", "Description"})

	for _, ticket := range totalTimePerTicket {
		table.Append([]string{ticket.issueKey, formatTime(ticket.duration), ticket.description})
	}
	table.Render()
	fmt.Println()

	// Print total time per unit
	if unit != "day" {
		unitDurations := getDurationsPerUnit(events, unit)
		unitTable := tablewriter.NewWriter(os.Stdout)
		unitTable.SetHeader([]string{unit, "Time"})

		for _, d := range unitDurations {
			unitTable.Append([]string{d.period, formatTime(d.duration)})
		}
		unitTable.Render()
		fmt.Println()
	}

	// Print total time
	totalTable := tablewriter.NewWriter(os.Stdout)
	totalTable.SetHeader([]string{"Total time"})
	totalTable.Append([]string{formatTime(totalTime)})
	totalTable.Render()
}

func getDurationsPerUnit(events []Event, unit string) []Duration {
	var durations []Duration

	for _, event := range events {
		duration := getDurationInSeconds(event.Start, event.End)

		if unit == "days" {
			period := event.Start.Format("2006-01-02")

			durations = appendIfMissing(durations, Duration{
				period:   period,
				sortkey:  event.Start.Format("20060102"),
				duration: duration,
			})
		} else if unit == "week" {
			_, week := event.Start.ISOWeek()
			period := fmt.Sprintf("Week %d %d", week, event.Start.Year())
			durations = appendIfMissing(durations, Duration{
				period:   period,
				sortkey:  event.Start.Format("20060102"),
				duration: duration,
			})
		} else if unit == "month" {
			period := event.Start.Format("January 2006")
			durations = appendIfMissing(durations, Duration{
				period:   period,
				sortkey:  event.Start.Format("20060102"),
				duration: duration,
			})
		}
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i].sortkey < durations[j].sortkey
	})

	return durations
}

func appendIfMissing(durationArray []Duration, duration Duration) []Duration {
	for i, d := range durationArray {
		if d.period == duration.period {
			durationArray[i].duration += duration.duration
			return durationArray
		}
	}

	return append(durationArray, duration)
}
