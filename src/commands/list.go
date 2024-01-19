package commands

import (
	"TimeTrack/src/calendar"
	"TimeTrack/src/database"
	"TimeTrack/src/utils"
	"fmt"
	"os"
	"sort"
	"strconv"
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
			Name:    "next",
			Aliases: []string{"n"},
			Usage:   "Show tasks from the next x units. This will start from the start date. (format: <number><unit>, unit options: d (days), w (week), m (month))",
		},
		&cli.StringFlag{
			Name:    "last",
			Aliases: []string{"l"},
			Usage:   "Show tasks from the last x units. This will start from the start date and go backwards. (format: <number><unit>, unit options: d (days), w (week), m (month))",
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

		start := c.String("start")

		if !utils.IsValidDate(start) {
			return cli.Exit("Invalid start date. Please use the following format: YYYY-MM-DD", 1)
		}

		if c.String("next") != "" && c.String("last") != "" {
			return cli.Exit("Please use either --next or --last, not both.", 1)
		}

		nextDuration, nextUnit, err := parseDurationString(c.String("next"))
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		if nextUnit != "" {
			parsedUnit, err := verifyUnit(nextUnit, 999)
			if err != nil {
				cli.Exit("Invalid unit. Valid units are: d (days), w (week), m (month)", 1)
			}
			nextUnit = parsedUnit
		}

		lastDuration, lastUnit, err := parseDurationString(c.String("last"))
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		if lastUnit != "" {
			parsedUnit, err := verifyUnit(lastUnit, 999)
			if err != nil {
				cli.Exit("Invalid last unit. Valid units are: d (days), w (week), m (month)", 1)
			}
			lastUnit = parsedUnit
		}

		if c.String("end") != "" && (lastUnit != "" || nextUnit != "") {
			return cli.Exit("end date can not be used with --next or --last", 1)
		}

		if c.String("end") != "" && !utils.IsValidDate(c.String("end")) {
			return cli.Exit("Invalid end date. Please use the following format: YYYY-MM-DD", 1)
		}

		startTime := fmt.Sprintf("%sT00:00:00+02:00", start)
		parsedStart, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			return nil
		}

		if c.String("end") == "" {
			if nextUnit != "" {
				c.Set("end", getDateFromDuration(parsedStart, nextDuration, nextUnit, "next"))
			} else if lastUnit != "" {
				startDate := start
				newStartDate := getDateFromDuration(parsedStart, lastDuration, lastUnit, "last")
				startTime = fmt.Sprintf("%sT00:00:00+02:00", newStartDate)
				parsedStart, err = time.Parse(time.RFC3339, startTime)
				if err != nil {
					fmt.Printf("Error parsing start time: %v\n", err)
					return nil
				}

				c.Set("end", startDate)
			} else {
				c.Set("end", start)
			}
		}

		endTime := fmt.Sprintf("%sT23:59:59+02:00", c.String("end"))

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
	difference int64
	days []string
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

func parseDurationString(durationString string) (int, string, error) {
	if durationString == "" {
		return 0, "", nil
	}

	number := durationString[:len(durationString)-1]
	unit := durationString[len(durationString)-1:]

	parsedNumber, err := strconv.Atoi(number)
	if err != nil {
		return 0, "", fmt.Errorf("invalid number: %s", number)
	}

	return parsedNumber, unit, nil
}

func getDateFromDuration(startDate time.Time, duration int, unit string, direction string) string {
	var date string
	if direction == "next" {
		if unit == "days" {
			date = startDate.AddDate(0, 0, duration).Format("2006-01-02")
		} else if unit == "week" {
			date = startDate.AddDate(0, 0, duration*7).Format("2006-01-02")
		} else if unit == "month" {
			date = startDate.AddDate(0, duration, 0).Format("2006-01-02")
		}
	} else if direction == "last" {
		if unit == "days" {
			date = startDate.AddDate(0, 0, -duration).Format("2006-01-02")
		} else if unit == "week" {
			date = startDate.AddDate(0, 0, -duration*7).Format("2006-01-02")
		} else if unit == "month" {
			date = startDate.AddDate(0, -duration, 0).Format("2006-01-02")
		}
	} else {
		fmt.Println("Invalid direction")
		return ""
	}

	return date
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

func formatDifference(seconds int64) string {
	isNegative := seconds < 0

	if isNegative {
		seconds = -seconds
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	operator := ""
	if isNegative {
		operator = "-"
	} else {
		operator = "+"
	}

	if hours == 0 {
		return fmt.Sprintf("%s%dm", operator, minutes)
	}

	return fmt.Sprintf("%s%dh %dm", operator, hours, minutes)
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

	// int64
	totalDifference := int64(0)
	isNotDay := unit != "day";

	// Print total time per unit
	if isNotDay {
		unitDurations := getDurationsPerUnit(events, unit)
		unitTable := tablewriter.NewWriter(os.Stdout)
		unitTable.SetHeader([]string{unit, "Time", "Difference"})

		for _, d := range unitDurations {
			totalDifference += d.difference
			unitTable.Append([]string{d.period, formatTime(d.duration), formatDifference(d.difference)})
		}
		unitTable.Render()
		fmt.Println()
	}

	// Print total time
	totalTable := tablewriter.NewWriter(os.Stdout)

	totalHeaders := []string{"Total time"}
	totalData := []string{formatTime(totalTime)}

	if isNotDay {
		totalHeaders = append(totalHeaders, "Difference")
		totalData = append(totalData, formatDifference(totalDifference))
	}

	totalTable.SetHeader(totalHeaders)
	totalTable.Append(totalData)
	totalTable.Render()
}

func getDurationsPerUnit(events []Event, unit string) []Duration {
	var durations []Duration

	for _, event := range events {
		duration := getDurationInSeconds(event.Start, event.End)
		day := event.Start.Format("2006-01-02")

		if unit == "days" {
			period := event.Start.Format("2006-01-02")
			durations = appendIfMissing(durations, Duration{
				period:   period,
				sortkey:  event.Start.Format("20060102"),
				duration: duration,
			}, day)
		} else if unit == "week" {
			_, week := event.Start.ISOWeek()
			period := fmt.Sprintf("Week %d %d", week, event.Start.Year())
			durations = appendIfMissing(durations, Duration{
				period:   period,
				sortkey:  event.Start.Format("20060102"),
				duration: duration,
			}, day)
		} else if unit == "month" {
			period := event.Start.Format("January 2006")
			durations = appendIfMissing(durations, Duration{
				period:   period,
				sortkey:  event.Start.Format("20060102"),
				duration: duration,
			}, day)
		}
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i].sortkey < durations[j].sortkey
	})

	return durations
}

func appendIfMissing(durationArray []Duration, duration Duration, day string) []Duration {
	averageDay := 8 * 60 * 60; // 8 hours in seconds

	for i, d := range durationArray {
		if d.period == duration.period {
			durationArray[i].duration += duration.duration
			durationArray[i].days = appendIfMissingDays(durationArray[i].days, day)
			
			expectedWorkingTime := averageDay * len(durationArray[i].days)
			durationArray[i].difference = durationArray[i].duration - int64(expectedWorkingTime)

			return durationArray
		}
	}

	return append(durationArray, duration)
}

func appendIfMissingDays(days []string, day string) []string {
	for _, d := range days {
		if d == day {
			return days
		}
	}

	return append(days, day)
}
