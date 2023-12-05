# TimeTrack Changelog

# Test (2023-12-05)

-   Test release

## 0.3.1 (2023-12-03)

-   Fixed bug where `list` command would not use the unit specified by the user.
-   Added `--next` and `--last` flags to `list` command to show the next or last unit of time.

## 0.3.0 (2023-12-03)

-   Added automatic check for updates on startup.
-   Added `update` command to manually check and update.

## 0.2.0 (2023-11-01)

-   Added `change` command to change the current task.
-   Added `info` command to show information about the current task.
-   Fixed issue where timezones were not being handled correctly.
-   Improved logging.
-   Made description optional when ending a task.
-   Fixed type where `selectCalendar` was being called `selectDatabase`.

## 0.1.3 (2023-10-22)

-   Fixed bug where --version flag would always show version as "dev".
-   Added error message when no events are found.

## 0.1.2 (2023-10-22)

-   Updated Google Credentials.

## 0.1.1 (2023-10-22)

-   Added installation scripts for Windows & macOS.
-   Added installation instructions to README.md.
-   Added update instructions to README.md.

## 0.1.0 (2023-10-22)

Initial Release:

-   Introduced CLI tool "TimeTrack".
-   Added task creation support.
-   Implemented start and end task functionality.
-   Feature to list tasks by day, week, and month.
-   View time spent on each task.
-   Google Calendar integration for task tracking.
